package conference_usecase

import (
	"fmt"

	"github.com/pion/webrtc/v4"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/utils/conference"
)

func (u *ConferenceUsecase) AddTrack(peer *conference_utils.PeerConnection, t *webrtc.TrackRemote) (*webrtc.TrackLocalStaticRTP, error) {
	trackLocal, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(), t.StreamID())
	if err != nil {
		return nil, fmt.Errorf("create track: %w", err)
	}
	u.mu.Lock()
	u.trackLocals[t.ID()] = trackLocal
	u.trackOwners[t.ID()] = peer
	u.mu.Unlock()
	u.logger.Infof("Added track %s from peer %p", trackLocal.ID(), peer)
	return trackLocal, nil
}

func (u *ConferenceUsecase) RemoveTrack(track *webrtc.TrackLocalStaticRTP) {
	u.mu.Lock()
	delete(u.trackLocals, track.ID())
	delete(u.trackOwners, track.ID())
	u.mu.Unlock()
}

func (u *ConferenceUsecase) UpdatePeerTracks(peer *conference_utils.PeerConnection) error {
	senders := peer.PC.GetSenders()
	for _, sender := range senders {
		if sender.Track() == nil {
			continue
		}
		trackID := sender.Track().ID()
		u.mu.RLock()
		_, ok := u.trackLocals[trackID]
		u.mu.RUnlock()
		if !ok {
			if err := peer.PC.RemoveTrack(sender); err != nil {
				return fmt.Errorf("remove track: %w", err)
			}
		}
	}

	sendingTracks := make(map[string]bool)
	for _, sender := range peer.PC.GetSenders() {
		if track := sender.Track(); track != nil {
			sendingTracks[track.ID()] = true
		}
	}

	u.mu.RLock()
	defer u.mu.RUnlock()
	for trackID, track := range u.trackLocals {
		if u.trackOwners[trackID] != peer && !sendingTracks[trackID] {
			if _, err := peer.PC.AddTrack(track); err != nil {
				return fmt.Errorf("add track: %w", err)
			}
		}
	}
	return nil
}
