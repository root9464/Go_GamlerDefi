package conference_usecase

import (
	"fmt"

	"github.com/pion/webrtc/v4"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/utils/conference"
)

func (u *ConferenceUsecase) AddTrack(pc *conference_utils.PeerConnection, t *webrtc.TrackRemote) (*webrtc.TrackLocalStaticRTP, error) {
	hub := u.hubs[pc.HubID]
	hub.mu.Lock()
	defer func() {
		hub.mu.Unlock()
		u.SignalPeers(pc)
	}()

	trackLocal, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(), t.StreamID())
	if err != nil {
		return nil, fmt.Errorf("create track: %w", err)
	}

	hub.trackLocals[t.ID()] = trackLocal
	u.logger.Infof("Added track %s from peer %p", trackLocal.ID(), pc)

	return trackLocal, nil
}

func (u *ConferenceUsecase) RemoveTrack(t *webrtc.TrackLocalStaticRTP, pc *conference_utils.PeerConnection) {
	hub := u.hubs[pc.HubID]
	hub.mu.Lock()
	defer func() {
		hub.mu.Unlock()
		u.SignalPeers(pc)
	}()

	for _, peer := range hub.peers {
		for _, sender := range peer.PC.GetSenders() {
			if sender.Track() != nil && sender.Track().ID() == t.ID() {
				_ = peer.PC.RemoveTrack(sender)
			}
		}
	}
	delete(hub.trackLocals, t.ID())
}

func (u *ConferenceUsecase) UpdatePeerTracks(pc *conference_utils.PeerConnection) error {
	hub := u.hubs[pc.HubID]
	senders := pc.PC.GetSenders()
	for _, sender := range senders {
		if sender.Track() == nil {
			continue
		}

		if _, ok := hub.trackLocals[sender.Track().ID()]; !ok {
			if err := pc.PC.RemoveTrack(sender); err != nil {
				return fmt.Errorf("remove track: %w", err)
			}
		}
	}

	sendingTracks := make(map[string]bool)
	for _, sender := range pc.PC.GetSenders() {
		if track := sender.Track(); track != nil {
			sendingTracks[track.ID()] = true
		}
	}

	for _, receiver := range pc.PC.GetReceivers() {
		if receiver.Track() == nil {
			continue
		}
		sendingTracks[receiver.Track().ID()] = true
	}

	for trackID, track := range hub.trackLocals {
		if _, ok := sendingTracks[trackID]; !ok {
			if _, err := pc.PC.AddTrack(track); err != nil {
				return fmt.Errorf("add track: %w", err)
			}
		}
	}
	return nil
}
