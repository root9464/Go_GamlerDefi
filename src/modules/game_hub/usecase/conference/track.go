package conference_usecase

import (
	"fmt"

	"github.com/pion/webrtc/v4"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/utils/conference"
)

func (u *ConferenceUsecase) AddTrack(peer *conference_utils.PeerConnection, t *webrtc.TrackRemote) (*webrtc.TrackLocalStaticRTP, error) {
	u.mu.Lock()
	defer func() {
		u.mu.Unlock()
		u.SignalPeers()
	}()

	trackLocal, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(), t.StreamID())
	if err != nil {
		return nil, fmt.Errorf("create track: %w", err)
	}

	u.trackLocals[t.ID()] = trackLocal
	u.logger.Infof("Added track %s from peer %p", trackLocal.ID(), peer)
	return trackLocal, nil
}

func (u *ConferenceUsecase) RemoveTrack(t *webrtc.TrackLocalStaticRTP) {
	u.mu.Lock()
	defer func() {
		u.mu.Unlock()
		u.SignalPeers()
	}()
	delete(u.trackLocals, t.ID())
}

func (u *ConferenceUsecase) UpdatePeerTracks(peer *conference_utils.PeerConnection) error {
	senders := peer.PC.GetSenders()
	for _, sender := range senders {
		if sender.Track() == nil {
			continue
		}

		if _, ok := u.trackLocals[sender.Track().ID()]; !ok {
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

	for _, receiver := range peer.PC.GetReceivers() {
		if receiver.Track() == nil {
			continue
		}
		sendingTracks[receiver.Track().ID()] = true
	}

	for trackID, track := range u.trackLocals {
		if _, ok := sendingTracks[trackID]; !ok {
			if _, err := peer.PC.AddTrack(track); err != nil {
				return fmt.Errorf("add track: %w", err)
			}
		}
	}
	return nil
}
