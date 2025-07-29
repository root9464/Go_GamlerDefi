package conference_usecase

import (
	"sync/atomic"

	"github.com/pion/webrtc/v4"
	hub_entity "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/entity"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/utils/conference"
)

func (u *ConferenceUsecase) AddTrack(conn *hub_entity.Connection, t *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP {
	room := u.rooms[conn.RoomID]
	room.Lock.Lock()
	defer func() {
		room.Lock.Unlock()
		u.SignalPeerConnections(conference_utils.GenerateRequestID(), conn.RoomID)
	}()

	trackLocal, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(), t.StreamID())
	if err != nil {
		u.logger.Error("Failed to create track local",
			"error", err,
			"track_id", t.ID(),
			"uuid", conn.Kws.GetUUID(),
			"request_id", conference_utils.GenerateRequestID(),
		)
		return nil
	}

	if oldTrack, exists := room.TrackLocals[t.ID()]; exists {
		u.logger.Warn("Replacing existing track",
			"track_id", t.ID(),
			"uuid", conn.Kws.GetUUID(),
			"ssrc", uint32(t.SSRC()),
		)
		u.RemoveTrack(oldTrack, conn.RoomID)
	}

	room.TrackLocals[t.ID()] = trackLocal
	conn.Tracks[t.ID()] = trackLocal
	atomic.AddInt64(&room.TrackCount, 1)
	u.logger.Info("Track added",
		"track_id", t.ID(),
		"uuid", conn.Kws.GetUUID(),
		"track_count", room.TrackCount,
		"ssrc", uint32(t.SSRC()),
	)
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
	u.logger.Info("Track removed",
		"track_id", t.ID(),
		"track_count", room.TrackCount,
	)
}
