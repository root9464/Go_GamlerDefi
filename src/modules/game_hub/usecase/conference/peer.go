package conference_usecase

import (
	"encoding/json"
	"fmt"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v4"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/utils/conference"
)

func (u *ConferenceUsecase) AddPeer(wpc *webrtc.PeerConnection, pc *conference_utils.PeerConnection) {
	hub := u.hubs[pc.HubID]
	hub.mu.Lock()
	defer hub.mu.Unlock()

	hub.peers = append(hub.peers, pc)
}

func (u *ConferenceUsecase) SignalPeers(pc *conference_utils.PeerConnection) error {
	hub := u.hubs[pc.HubID]
	hub.mu.Lock()
	defer func() {
		hub.mu.Unlock()
		u.DispatchKeyFrames(pc.HubID)
	}()

	peers := u.activePeers(pc.HubID)
	for _, peer := range peers {
		if err := u.UpdatePeerTracks(peer); err != nil {
			return err
		}
		if peer.PC.SignalingState() == webrtc.SignalingStateHaveLocalOffer {
			u.logger.Info("offer already exist")
			continue
		}
		if err := u.sendOffer(peer); err != nil {
			return err
		}
	}
	return nil
}

func (u *ConferenceUsecase) activePeers(hubID string) []*conference_utils.PeerConnection {
	hub := u.hubs[hubID]
	var active []*conference_utils.PeerConnection
	for _, peer := range hub.peers {
		if peer.PC.ConnectionState() != webrtc.PeerConnectionStateClosed {
			active = append(active, peer)
		}
	}
	hub.peers = active
	return active
}

func (u *ConferenceUsecase) DispatchKeyFrames(hubID string) {
	hub := u.hubs[hubID]
	hub.mu.RLock()
	defer hub.mu.RUnlock()
	for _, peer := range hub.peers {
		for _, receiver := range peer.PC.GetReceivers() {
			if track := receiver.Track(); track != nil && track.Kind() == webrtc.RTPCodecTypeVideo {
				_ = peer.PC.WriteRTCP([]rtcp.Packet{
					&rtcp.PictureLossIndication{MediaSSRC: uint32(track.SSRC())},
				})
			}
		}
	}
}

func (u *ConferenceUsecase) sendOffer(peer *conference_utils.PeerConnection) error {
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
	u.logger.Infof("Sending offer")

	return peer.Writer.WriteJSON(conference_utils.WebsocketMessage{
		Event: "offer",
		Data:  string(offerData),
	})
}
