package conference_usecase

import (
	"sync/atomic"

	"github.com/pion/webrtc/v4"
	conference_entity "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/entity"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/utils"
)

func (u *ConferenceUsecase) AddTrack(conn *conference_entity.Connection, t *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP {
	room := u.rooms[conn.RoomID]
	room.Lock.Lock()
	defer func() {
		room.Lock.Unlock()
		u.SignalPeerConnections(conference_utils.GenerateRequestID(), conn.RoomID)
	}()

	trackLocal, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(), t.StreamID())
	if err != nil {
		u.logger.Errorf("failed to create track local error: %v, track id: %s, uuid: %s, request id: %s", err, t.ID(), conn.Kws.GetUUID(), conference_utils.GenerateRequestID())
		return nil
	}

	if oldTrack, exists := room.TrackLocals[t.ID()]; exists {
		u.logger.Warnf("replacing existing track track id: %s, uuid: %s, ssrc: %d", t.ID(), conn.Kws.GetUUID(), uint32(t.SSRC()))
		u.RemoveTrack(oldTrack, conn.RoomID)
	}

	room.TrackLocals[t.ID()] = trackLocal
	conn.Tracks[t.ID()] = trackLocal
	atomic.AddInt64(&room.TrackCount, 1)
	u.logger.Infof("track added track id: %s, uuid: %s, track count: %d, ssrc: %d", t.ID(), conn.Kws.GetUUID(), room.TrackCount, uint32(t.SSRC()))
	return trackLocal
}

func (u *ConferenceUsecase) RemoveTrack(t *webrtc.TrackLocalStaticRTP, roomID string) {
	room := u.rooms[roomID]

	room.Lock.Lock()
	defer func() {
		room.Lock.Unlock()
		u.SignalPeerConnections(conference_utils.GenerateRequestID(), roomID)
	}()
	delete(room.TrackLocals, t.ID())
	atomic.AddInt64(&room.TrackCount, -1)
	u.logger.Infof("track removed track id: %s, track count: %d", t.ID(), room.TrackCount)

}
