package conference_usecase

import (
	"sync/atomic"
	"time"

	"github.com/gofiber/contrib/socketio"
	"github.com/pion/webrtc/v4"
	conference_entity "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/entity"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/utils"
)

func (u *ConferenceUsecase) CreateConnection(roomID string, pc *webrtc.PeerConnection, kws *socketio.Websocket) *conference_entity.Connection {
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

func (u *ConferenceUsecase) GetOrCreateRoom(roomID string, requestID string, conn *conference_entity.Connection) *conference_entity.Room {
	u.roomsLock.Lock()
	r, exists := u.rooms[roomID]
	if !exists {
		r = &conference_entity.Room{
			Connections: make([]*conference_entity.Connection, 0),
			TrackLocals: make(map[string]*webrtc.TrackLocalStaticRTP),
			TrackCount:  0,
		}
		u.rooms[roomID] = r
		u.logger.Infof("new room created, room id: %s, request id: %s", roomID, requestID)
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
		u.logger.Infof("disconnected without room uuid: %s, request id: %s", ep.Kws.GetUUID(), requestID)
		return
	}

	u.roomsLock.RLock()
	r, exists := u.rooms[roomID]
	u.roomsLock.RUnlock()

	if !exists {
		u.logger.Infof("disconnected from non-existing room id: %s, request id: %s", roomID, requestID)
		return
	}

	r.Lock.Lock()

	var disconnectedConn *conference_entity.Connection
	for i, conn := range r.Connections {
		if conn.Kws.GetUUID() == ep.Kws.GetUUID() {
			disconnectedConn = conn
			u.logger.Infof("stop recording, uuid: %s, request id: %s", conn.Kws.GetUUID(), requestID)
			if u.audioRecorder != nil {
				for trackID := range conn.Tracks {
					u.audioRecorder.StopRecordingTrack(trackID, conn.RoomID, conn.Kws.GetUUID())
				}
			}
			u.logger.Infof("disconnected from room, uuid: %s, request id: %s", conn.Kws.GetUUID(), requestID)
			for trackID := range conn.Tracks {
				if _, ok := r.TrackLocals[trackID]; ok {
					delete(r.TrackLocals, trackID)
					atomic.AddInt64(&r.TrackCount, -1)
				}
			}
			close(conn.Closed)
			if err := conn.Pc.Close(); err != nil {
				u.logger.Errorf("failed to close peer connection error: %v, uuid: %s, request id: %s", err, conn.Kws.GetUUID(), requestID)
			}
			r.Connections = append(r.Connections[:i], r.Connections[i+1:]...)
			break
		}
	}

	shouldMixAudio := len(r.Connections) == 0 && u.audioRecorder != nil && disconnectedConn != nil
	r.Lock.Unlock()

	if shouldMixAudio {
		u.logger.Info("mixing only when deleting room")
		go func() {
			if err := u.audioRecorder.MixAndCleanupRoom(roomID); err != nil {
				u.logger.Errorf("failed to mix room audio room id: %s, error: %v, request id: %s", roomID, err, requestID)
			} else {
				u.logger.Infof("room audio mixed successfully room id: %s, request id: %s", roomID, requestID)
			}

			u.roomsLock.Lock()
			delete(u.rooms, roomID)
			u.roomsLock.Unlock()
			u.logger.Infof("room deleted with audio mixing room id: %s, request id: %s", roomID, requestID)
		}()
	} else if len(r.Connections) > 0 {
		u.SignalPeerConnections(requestID, roomID)
	}

	u.logger.Infof("disconnected from room uuid: %s, room id: %s, request id: %s", ep.Kws.GetUUID(), roomID, requestID)
}
