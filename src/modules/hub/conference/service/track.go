package conference_service

import (
	"sync/atomic"

	"github.com/pion/webrtc/v4"
	conference_entity "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/entity"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/util"
)

func (s *ConferenceService) AddTrack(conn *conference_entity.Connection, t *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP {
	room := s.rooms[conn.RoomID]
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
	s.logger.Infof("Track added, track_id: %s, uuid: %s, ssrc: %d", t.ID(), conn.Kws.GetUUID(), uint32(t.SSRC()))
	return trackLocal
}

func (s *ConferenceService) RemoveTrack(t *webrtc.TrackLocalStaticRTP, roomID string) {
	room := s.rooms[roomID]

	room.Lock.Lock()
	defer func() {
		room.Lock.Unlock()
		s.SignalPeerConnections(conference_utils.GenerateRequestID(), roomID)
	}()
	delete(room.TrackLocals, t.ID())
	atomic.AddInt64(&room.TrackCount, -1)
	s.logger.Infof("Track removed, track_id: %s, uuid: %s", t.ID(), room.Connections[0].Kws.GetUUID())
}
