package conference_usecase

import (
	"encoding/json"
	"fmt"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v4"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/utils/conference"
)

func (u *ConferenceUsecase) AddPeer(pc *webrtc.PeerConnection, conn *conference_utils.ThreadSafeWriter) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.peers = append(u.peers, &conference_utils.PeerConnection{PC: pc, Conn: conn})
}

func (u *ConferenceUsecase) SignalPeers() error {
	peers := u.activePeers()
	for _, peer := range peers {
		if err := u.UpdatePeerTracks(peer); err != nil {
			return err
		}
		if err := u.sendOffer(peer); err != nil {
			return err
		}
	}
	u.DispatchKeyFrames()
	return nil
}

func (u *ConferenceUsecase) activePeers() []*conference_utils.PeerConnection {
	u.mu.Lock()
	defer u.mu.Unlock()
	var active []*conference_utils.PeerConnection
	for _, peer := range u.peers {
		if peer.PC.ConnectionState() != webrtc.PeerConnectionStateClosed {
			active = append(active, peer)
		}
	}
	u.peers = active
	return active
}

func (u *ConferenceUsecase) DispatchKeyFrames() {
	u.mu.RLock()
	defer u.mu.RUnlock()
	for _, peer := range u.peers {
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
	u.logger.Infof("Sending offer: %s", string(offerData))
	return peer.Conn.WriteJSON(conference_utils.WebsocketMessage{
		Event: "offer",
		Data:  string(offerData),
	})
}
