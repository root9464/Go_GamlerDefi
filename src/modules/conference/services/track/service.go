package track_service

import (
	"fmt"
	"sync"

	"github.com/pion/webrtc/v4"
	"github.com/root9464/Go_GamlerDefi/src/modules/conference/utils"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type TrackService struct {
	logger      *logger.Logger
	mu          sync.RWMutex
	trackLocals map[string]*webrtc.TrackLocalStaticRTP
	trackOwners map[string]*utils.PeerConnection
}

func NewTrackService(
	logger *logger.Logger,
	trackLocals map[string]*webrtc.TrackLocalStaticRTP,
	trackOwners map[string]*utils.PeerConnection,
) *TrackService {
	return &TrackService{
		logger:      logger,
		trackLocals: trackLocals,
		trackOwners: trackOwners,
		mu:          sync.RWMutex{},
	}
}

func (s *TrackService) AddTrack(peer *utils.PeerConnection, t *webrtc.TrackRemote) (*webrtc.TrackLocalStaticRTP, error) {
	trackLocal, err := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(), t.StreamID())
	if err != nil {
		return nil, fmt.Errorf("create track: %w", err)
	}
	s.mu.Lock()
	s.trackLocals[t.ID()] = trackLocal
	s.trackOwners[t.ID()] = peer
	s.mu.Unlock()
	s.logger.Infof("Added track %s from peer %p", trackLocal.ID(), peer)
	return trackLocal, nil
}

func (s *TrackService) RemoveTrack(track *webrtc.TrackLocalStaticRTP) {
	s.mu.Lock()
	delete(s.trackLocals, track.ID())
	delete(s.trackOwners, track.ID())
	s.mu.Unlock()
}

func (s *TrackService) UpdatePeerTracks(peer *utils.PeerConnection) error {
	senders := peer.PC.GetSenders()
	for _, sender := range senders {
		if sender.Track() == nil {
			continue
		}
		trackID := sender.Track().ID()
		s.mu.RLock()
		_, ok := s.trackLocals[trackID]
		s.mu.RUnlock()
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

	s.mu.RLock()
	defer s.mu.RUnlock()
	for trackID, track := range s.trackLocals {
		if s.trackOwners[trackID] != peer && !sendingTracks[trackID] {
			if _, err := peer.PC.AddTrack(track); err != nil {
				return fmt.Errorf("add track: %w", err)
			}
		}
	}
	return nil
}
