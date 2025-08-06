package conference_usecase

import (
	"sync/atomic"
	"time"

	"github.com/gofiber/contrib/socketio"
	"github.com/pion/webrtc/v4"
	hub_entity "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/entity"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/utils/conference"
)

func (u *ConferenceUsecase) CreateConnection(roomID string, pc *webrtc.PeerConnection, kws *socketio.Websocket) *hub_entity.Connection {
	conn := &hub_entity.Connection{
		Pc:             pc,
		Kws:            kws,
		Tracks:         make(map[string]*webrtc.TrackLocalStaticRTP),
		RoomID:         roomID,
		Closed:         make(chan struct{}),
		SignalDebounce: time.Millisecond * 500,
	}

	return conn
}

func (u *ConferenceUsecase) GetOrCreateRoom(roomID string, requestID string, conn *hub_entity.Connection) *hub_entity.Room {
	u.roomsLock.Lock()
	r, exists := u.rooms[roomID]
	if !exists {
		r = &hub_entity.Room{
			Connections: make([]*hub_entity.Connection, 0),
			TrackLocals: make(map[string]*webrtc.TrackLocalStaticRTP),
			TrackCount:  0,
		}
		u.rooms[roomID] = r
		u.logger.Info("New room created",
			"room_id", roomID,
			"request_id", requestID)
	}
	u.roomsLock.Unlock()

	if conn != nil {
		r.Lock.Lock()
		r.Connections = append(r.Connections, conn)
		r.Lock.Unlock()
	}

	return r
}

func (u *ConferenceUsecase) Disconect(ep *socketio.EventPayload) {
	requestID := conference_utils.GenerateRequestID()
	roomID := ep.Kws.GetStringAttribute("room_id")
	if roomID == "" {
		u.logger.Info("Disconnected without room",
			"uuid", ep.Kws.GetUUID(),
			"request_id", requestID,
		)
		return
	}

	u.roomsLock.RLock()
	r, exists := u.rooms[roomID]
	u.roomsLock.RUnlock()

	if !exists {
		u.logger.Info("Disconnected from non-existing room",
			"room_id", roomID,
			"request_id", requestID,
		)
		return
	}

	r.Lock.Lock()
	for i, conn := range r.Connections {
		if conn.Kws.GetUUID() == ep.Kws.GetUUID() {
			// Останавливаем запись для этого пользователя
			if u.audioRecorder != nil {
				for trackID := range conn.Tracks {
					u.audioRecorder.StopRecordingTrack(trackID, conn.RoomID, conn.Kws.GetUUID())
				}
			}

			for trackID := range conn.Tracks {
				if _, ok := r.TrackLocals[trackID]; ok {
					delete(r.TrackLocals, trackID)
					atomic.AddInt64(&r.TrackCount, -1)
				}
			}
			close(conn.Closed)
			if err := conn.Pc.Close(); err != nil {
				u.logger.Error("Failed to close PeerConnection",
					"error", err,
					"uuid", conn.Kws.GetUUID(),
					"request_id", requestID,
				)
			}
			r.Connections = append(r.Connections[:i], r.Connections[i+1:]...)
			break
		}
	}

	// МИКШИРОВАНИЕ ТОЛЬКО ПРИ УДАЛЕНИИ КОМНАТЫ
	if len(r.Connections) == 0 {
		// Микшируем аудио перед удалением комнаты
		if u.audioRecorder != nil {
			go func() {
				if err := u.audioRecorder.MixAndCleanupRoom(roomID); err != nil {
					u.logger.Error("Failed to mix room audio",
						"room_id", roomID,
						"error", err,
						"request_id", requestID,
					)
				} else {
					u.logger.Info("Room audio mixed successfully",
						"room_id", roomID,
						"request_id", requestID,
					)
				}
			}()
		}

		u.roomsLock.Lock()
		delete(u.rooms, roomID)
		u.roomsLock.Unlock()
		u.logger.Info("Room deleted with audio mixing",
			"room_id", roomID,
			"request_id", requestID,
		)
	} else {
		r.Lock.Unlock()
		u.SignalPeerConnections(requestID, roomID)
		r.Lock.Lock()
	}

	r.Lock.Unlock()
	u.logger.Info("Disconnected from room",
		"uuid", ep.Kws.GetUUID(),
		"room_id", roomID,
		"request_id", requestID,
		"remaining_tracks", r.TrackCount,
	)
}
