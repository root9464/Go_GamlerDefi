package conference_service

import (
	"sync/atomic"
	"time"

	"github.com/gofiber/contrib/socketio"
	"github.com/pion/webrtc/v4"
	conference_entity "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/entity"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/util"
)

func (u *ConferenceService) CreateConnection(roomID string, pc *webrtc.PeerConnection, kws *socketio.Websocket) *conference_entity.Connection {
	conn := &conference_entity.Connection{
		Pc:             pc,
		Kws:            kws,
		Tracks:         make(map[string]*webrtc.TrackLocalStaticRTP),
		RoomID:         roomID,
		Closed:         make(chan struct{}),
		SignalDebounce: time.Millisecond * 500,
	}

	return conn
}

func (s *ConferenceService) GetOrCreateRoom(roomID string, requestID string, conn *conference_entity.Connection) *conference_entity.Room {
	s.roomsLock.Lock()
	r, exists := s.rooms[roomID]
	if !exists {
		r = &conference_entity.Room{
			Connections: make([]*conference_entity.Connection, 0),
			TrackLocals: make(map[string]*webrtc.TrackLocalStaticRTP),
			TrackCount:  0,
		}
		s.rooms[roomID] = r
		s.logger.Infof("New room created, room_id: %s, request_id: %s", roomID, requestID)
	}
	s.roomsLock.Unlock()

	if conn != nil {
		r.Lock.Lock()
		r.Connections = append(r.Connections, conn)
		r.Lock.Unlock()
	}

	return r
}

func (s *ConferenceService) Disconect(ep *socketio.EventPayload) {
	requestID := conference_utils.GenerateRequestID()
	roomID := ep.Kws.GetStringAttribute("room_id")
	if roomID == "" {
		s.logger.Infof("Disconnected without room, uuid: %s, request_id: %s", ep.Kws.GetUUID(), requestID)
		return
	}

	s.roomsLock.RLock()
	r, exists := s.rooms[roomID]
	s.roomsLock.RUnlock()

	if !exists {
		s.logger.Infof("Disconnected from non-existing room, uuid: %s, request_id: %s", ep.Kws.GetUUID(), requestID)
		return
	}

	r.Lock.Lock()
	for i, conn := range r.Connections {
		if conn.Kws.GetUUID() == ep.Kws.GetUUID() {
			// Останавливаем запись для этого пользователя
			// if u.audioRecorder != nil {
			// 	for trackID := range conn.Tracks {
			// 		u.audioRecorder.StopRecordingTrack(trackID, conn.RoomID, conn.Kws.GetUUID())
			// 	}
			// }

			for trackID := range conn.Tracks {
				if _, ok := r.TrackLocals[trackID]; ok {
					delete(r.TrackLocals, trackID)
					atomic.AddInt64(&r.TrackCount, -1)
				}
			}
			close(conn.Closed)
			if err := conn.Pc.Close(); err != nil {
				s.logger.Errorf("Failed to close PeerConnection, error: %v, uuid: %s, request_id: %s", err, conn.Kws.GetUUID(), requestID)
			}
			r.Connections = append(r.Connections[:i], r.Connections[i+1:]...)
			break
		}
	}

	// МИКШИРОВАНИЕ ТОЛЬКО ПРИ УДАЛЕНИИ КОМНАТЫ
	if len(r.Connections) == 0 {
		// Микшируем аудио перед удалением комнаты
		// if u.audioRecorder != nil {
		// 	go func() {
		// 		if err := u.audioRecorder.MixAndCleanupRoom(roomID); err != nil {
		// 			u.logger.Error("Failed to mix room audio",
		// 				"room_id", roomID,
		// 				"error", err,
		// 				"request_id", requestID,
		// 			)
		// 		} else {
		// 			u.logger.Info("Room audio mixed successfully",
		// 				"room_id", roomID,
		// 				"request_id", requestID,
		// 			)
		// 		}
		// 	}()
		// }

		s.roomsLock.Lock()
		delete(s.rooms, roomID)
		s.roomsLock.Unlock()
		s.logger.Infof("Room deleted, room_id: %s, request_id: %s", roomID, requestID)
	} else {
		r.Lock.Unlock()
		s.SignalPeerConnections(requestID, roomID)
		r.Lock.Lock()
	}

	r.Lock.Unlock()
	s.logger.Infof("Disconnected from room, uuid: %s, room_id: %s, request_id: %s, remaining_tracks: %d", ep.Kws.GetUUID(), roomID, requestID, r.TrackCount)
}
