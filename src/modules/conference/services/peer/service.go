package conference_peer_service

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v4"
	"github.com/root9464/Go_GamlerDefi/src/modules/conference/utils"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type ITrackService interface {
	UpdatePeerTracks(peer *conference_utils.PeerConnection) error
}

type IPeerService interface {
	AddPeer(pc *webrtc.PeerConnection, conn *conference_utils.ThreadSafeWriter)
	SignalPeers() error
	DispatchKeyFrames()
}

type PeerService struct {
	logger       *logger.Logger
	trackLocals  map[string]*webrtc.TrackLocalStaticRTP
	trackOwners  map[string]*conference_utils.PeerConnection
	peers        []*conference_utils.PeerConnection
	trackService ITrackService

	mu sync.RWMutex
}

func NewPeerService(
	logger *logger.Logger,
	trackLocals map[string]*webrtc.TrackLocalStaticRTP,
	trackOwners map[string]*conference_utils.PeerConnection,
	trackService ITrackService,
) IPeerService {
	return &PeerService{
		logger:       logger,
		trackLocals:  trackLocals,
		trackOwners:  trackOwners,
		trackService: trackService,
		mu:           sync.RWMutex{},
	}
}

func (s *PeerService) AddPeer(pc *webrtc.PeerConnection, conn *conference_utils.ThreadSafeWriter) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.peers = append(s.peers, &conference_utils.PeerConnection{PC: pc, Conn: conn})
}

func (s *PeerService) SignalPeers() error {
	peers := s.activePeers()
	for _, peer := range peers {
		if err := s.trackService.UpdatePeerTracks(peer); err != nil {
			return err
		}
		if err := s.sendOffer(peer); err != nil {
			return err
		}
	}
	s.DispatchKeyFrames()
	return nil
}

func (s *PeerService) activePeers() []*conference_utils.PeerConnection {
	s.mu.Lock()
	defer s.mu.Unlock()
	var active []*conference_utils.PeerConnection
	for _, peer := range s.peers {
		if peer.PC.ConnectionState() != webrtc.PeerConnectionStateClosed {
			active = append(active, peer)
		}
	}
	s.peers = active
	return active
}

func (s *PeerService) DispatchKeyFrames() {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, peer := range s.peers {
		for _, receiver := range peer.PC.GetReceivers() {
			if track := receiver.Track(); track != nil && track.Kind() == webrtc.RTPCodecTypeVideo {
				_ = peer.PC.WriteRTCP([]rtcp.Packet{
					&rtcp.PictureLossIndication{MediaSSRC: uint32(track.SSRC())},
				})
			}
		}
	}
}

func (s *PeerService) sendOffer(peer *conference_utils.PeerConnection) error {
	offer, err := peer.PC.CreateOffer(nil)
	if err != nil {
		return fmt.Errorf("create offer: %w", err)
	}
	if err := peer.PC.SetLocalDescription(offer); err != nil {
		return fmt.Errorf("set local description: %w", err)
	}
	offerData, err := json.Marshal(offer)
	if err != nil {
		return fmt.Errorf("marshal offer: %w", err)
	}
	s.logger.Infof("Sending offer: %s", string(offerData))
	return peer.Conn.WriteJSON(conference_utils.WebsocketMessage{
		Event: "offer",
		Data:  string(offerData),
	})
}
