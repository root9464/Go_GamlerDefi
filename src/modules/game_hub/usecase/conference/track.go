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
		u.logger.Debugf("Signaling peers after adding track %s", t.ID())
		u.SignalPeers(pc)
	}()

	u.logger.Debugf("Creating local static RTP track for %s (stream: %s)", t.ID(), t.StreamID())
	trackLocal, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(), t.StreamID())
	if err != nil {
		u.logger.Errorf("Failed to create local track: %v", err)
		return nil, fmt.Errorf("create track: %w", err)
	}

	hub.trackLocals[t.ID()] = trackLocal
	u.logger.Infof("Added local track %s to hub %s from peer %p", trackLocal.ID(), pc.HubID, pc)

	return trackLocal, nil
}

func (u *ConferenceUsecase) RemoveTrack(t *webrtc.TrackLocalStaticRTP, pc *conference_utils.PeerConnection) {
	hub := u.hubs[pc.HubID]
	hub.mu.Lock()
	defer func() {
		hub.mu.Unlock()
		u.logger.Debugf("Signaling peers after removing track %s", t.ID())
		u.SignalPeers(pc)
	}()

	u.logger.Infof("Removing track %s from all peers in hub %s", t.ID(), pc.HubID)

	delete(hub.trackLocals, t.ID())
	u.logger.Infof("Track %s removed from hub %s trackLocals", t.ID(), pc.HubID)
}

func (u *ConferenceUsecase) UpdatePeerTracks(pc *conference_utils.PeerConnection) error {
	hub := u.hubs[pc.HubID]
	u.logger.Debugf("Updating peer tracks for peer %p in hub %s", pc, pc.HubID)

	// Проверка состояния соединения
	if pc.PC.ConnectionState() == webrtc.PeerConnectionStateClosed {
		return fmt.Errorf("peer connection is closed")
	}

	// Удаление устаревших senders
	senders := pc.PC.GetSenders()
	for _, sender := range senders {
		if sender.Track() == nil {
			continue
		}
		trackID := sender.Track().ID()

		if _, ok := hub.trackLocals[trackID]; !ok {
			u.logger.Debugf("Track %s not found in hub.trackLocals — removing sender", trackID)
			if err := pc.PC.RemoveTrack(sender); err != nil {
				u.logger.Errorf("Failed to remove sender track %s: %v", trackID, err)
				return fmt.Errorf("remove track: %w", err)
			}
			u.logger.Infof("Sender track %s removed from peer %p", trackID, pc)
		}
	}

	// Проверка уже отправляемых треков (только senders)
	sendingTracks := make(map[string]bool)
	for _, sender := range pc.PC.GetSenders() {
		if track := sender.Track(); track != nil {
			sendingTracks[track.ID()] = true
		}
	}

	// Добавление недостающих треков
	for trackID, track := range hub.trackLocals {
		if _, ok := sendingTracks[trackID]; !ok {
			u.logger.Debugf("Track %s not being sent — adding to peer %p", trackID, pc)
			if _, err := pc.PC.AddTrack(track); err != nil {
				u.logger.Errorf("Failed to add track %s to peer %p: %v", trackID, pc, err)
				return fmt.Errorf("add track: %w", err)
			}
			u.logger.Infof("Track %s added to peer %p", trackID, pc)
		} else {
			u.logger.Debugf("Track %s already sending — skipping", trackID)
		}
	}
	return nil
}
