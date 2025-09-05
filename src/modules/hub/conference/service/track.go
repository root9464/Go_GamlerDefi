package conference_service

import (
	"sync/atomic"

	"github.com/pion/webrtc/v4"
	conference_entity "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/entity"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/util"
)

func (s *ConferenceService) AddTrack(conn *conference_entity.Connection, t *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP {
	s.roomsLock.RLock()
	room, exists := s.rooms[conn.RoomID]
	s.roomsLock.RUnlock()

	if !exists {
		s.logger.Warnf("Room not found for track addition, room_id: %s, track_id: %s", conn.RoomID, t.ID())
		return nil
	}

	room.Lock.Lock()
	defer func() {
		room.Lock.Unlock()
		s.SignalPeerConnections(conference_utils.GenerateRequestID(), conn.RoomID)
	}()

	trackLocal, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(), t.StreamID())
	if err != nil {
		s.logger.Errorf("Failed to create track local, error: %v, track_id: %s, uuid: %s, request_id: %s", err, t.ID(), conn.Kws.GetUUID(), conference_utils.GenerateRequestID())
		return nil
	}

	if oldTrack, exists := room.TrackLocals[t.ID()]; exists {
		s.logger.Warnf("Replacing existing track, track_id: %s, uuid: %s, ssrc: %d", t.ID(), conn.Kws.GetUUID(), uint32(t.SSRC()))
		s.RemoveTrack(oldTrack, conn.RoomID)
	}

	room.TrackLocals[t.ID()] = trackLocal
	conn.Tracks[t.ID()] = trackLocal
	atomic.AddInt64(&room.TrackCount, 1)

	if conn.User == nil {
		conn.User = &conference_entity.User{
			StreamID: t.StreamID(),
			TrackID:  t.ID(),
			UserID:   conn.Kws.GetUUID(),
			Username: "",
		}
	} else {
		conn.User.StreamID = t.StreamID()
		conn.User.TrackID = t.ID()
	}
	s.logger.Infof("Track added, track_id: %s, uuid: %s, ssrc: %d", t.ID(), conn.Kws.GetUUID(), uint32(t.SSRC()))

	s.broadcastRoomEvent(room, "participant:updated", conn.User)
	return trackLocal
}

func (s *ConferenceService) RemoveTrack(t *webrtc.TrackLocalStaticRTP, roomID string) {
	room := s.rooms[roomID]

	room.Lock.Lock()
	defer func() {
		room.Lock.Unlock()
		s.SignalPeerConnections(conference_utils.GenerateRequestID(), roomID)
	}()
	var ownerUser *conference_entity.User
	for _, c := range room.Connections {
		if _, ok := c.Tracks[t.ID()]; ok {
			delete(c.Tracks, t.ID())
			ownerUser = c.User
			break
		}
	}
	delete(room.TrackLocals, t.ID())
	atomic.AddInt64(&room.TrackCount, -1)
	if ownerUser != nil {
		ownerUser.TrackID = ""
		ownerUser.StreamID = ""
		s.broadcastRoomEvent(room, "participant:updated", ownerUser)
		s.logger.Infof("Track removed, track_id: %s, user: %s", t.ID(), ownerUser.UserID)
		return
	}
	if len(room.Connections) > 0 {
		s.logger.Infof("Track removed, track_id: %s, uuid: %s", t.ID(), room.Connections[0].Kws.GetUUID())
	}
}
