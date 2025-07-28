package conference_usecase

import (
	"encoding/json"
	"fmt"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v4"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/utils/conference"
)

func (u *ConferenceUsecase) AddPeer(pc *conference_utils.PeerConnection) {
	hub := u.hubs[pc.HubID]
	hub.mu.Lock()
	defer hub.mu.Unlock()

	hub.peers[conference_utils.GetPeerID(pc)] = pc
	u.logger.Infof("Added new peer %s to hub %s", pc.UserID, pc.HubID)
}

func (u *ConferenceUsecase) SignalPeers(pc *conference_utils.PeerConnection) error {
	hub := u.hubs[pc.HubID]
	hub.mu.Lock()
	defer func() {
		hub.mu.Unlock()
		u.logger.Debugf("Dispatching keyframes for hub %s", pc.HubID)
		u.DispatchKeyFrames(pc.HubID)
	}()

	peers := u.activePeers(pc.HubID)
	u.logger.Debugf("Signaling %d peers in hub %s", len(peers), pc.HubID)

	for _, peer := range peers {
		u.logger.Debugf("Updating tracks for peer %s", peer.UserID)
		if err := u.UpdatePeerTracks(peer); err != nil {
			u.logger.Errorf("Failed to update tracks for peer %s: %v", peer.UserID, err)
			continue // Продолжить с другими пирами
		}

		if peer.PC.ConnectionState() == webrtc.PeerConnectionStateClosed {
			u.logger.Warnf("Peer %s connection is closed — skipping", peer.UserID)
			continue
		}

		// Проверить signaling state
		if peer.PC.SignalingState() != webrtc.SignalingStateStable {
			u.logger.Warnf("Peer %s not in stable state (%s) — skipping",
				peer.UserID, peer.PC.SignalingState().String())
			continue
		}

		u.logger.Debugf("Sending offer to peer %s", peer.UserID)
		if err := u.sendOffer(peer); err != nil {
			u.logger.Errorf("Failed to send offer to peer %s: %v", peer.UserID, err)
			// Решить: return err или continue
		}
	}
	return nil
}

func (u *ConferenceUsecase) activePeers(hubID string) map[string]*conference_utils.PeerConnection {
	hub := u.hubs[hubID]
	active := make(map[string]*conference_utils.PeerConnection)

	for peerID, peer := range hub.peers {
		if peer.PC.ConnectionState() != webrtc.PeerConnectionStateClosed {
			active[peerID] = peer
		} else {
			u.logger.Warnf("Peer %s in hub %s is closed — removing", peer.UserID, hubID)
		}
	}

	if len(hub.peers) != len(active) {
		u.logger.Infof("Filtered out %d inactive peers in hub %s", len(hub.peers)-len(active), hubID)
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
			track := receiver.Track()
			if track != nil && track.Kind() == webrtc.RTPCodecTypeVideo {
				err := peer.PC.WriteRTCP([]rtcp.Packet{
					&rtcp.PictureLossIndication{MediaSSRC: uint32(track.SSRC())},
				})
				if err != nil {
					u.logger.Errorf("Failed to send PLI to peer %s: %v", peer.UserID, err)
				} else {
					u.logger.Debugf("Sent PLI to peer %s for track %s", peer.UserID, track.ID())
				}
			}
		}
	}
}

func (u *ConferenceUsecase) sendOffer(peer *conference_utils.PeerConnection) error {
	u.logger.Debugf("Creating offer for peer %s", peer.UserID)
	offer, err := peer.PC.CreateOffer(nil)
	if err != nil {
		u.logger.Errorf("Failed to create offer for peer %s: %v", peer.UserID, err)
		return fmt.Errorf("create offer: %w", err)
	}

	err = peer.PC.SetLocalDescription(offer)
	if err != nil {
		u.logger.Errorf("Failed to set local description for peer %s: %v", peer.UserID, err)
		return fmt.Errorf("set local description: %w", err)
	}

	offerData, err := json.Marshal(offer)
	if err != nil {
		u.logger.Errorf("Failed to marshal offer for peer %s: %v", peer.UserID, err)
		return fmt.Errorf("marshal offer: %w", err)
	}

	u.logger.Infof("Sending offer to peer %s", peer.UserID)
	return peer.Writer.WriteJSON(conference_utils.WebsocketMessage{
		Event: "offer",
		Data:  string(offerData),
	})
}
