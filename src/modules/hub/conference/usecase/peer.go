package conference_usecase

import (
	"encoding/json"
	"sync/atomic"
	"time"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v4"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/utils"
)

func (u *ConferenceUsecase) SignalPeerConnections(requestID string, roomID string) {
	u.roomsLock.RLock()
	room, exists := u.rooms[roomID]
	u.roomsLock.RUnlock()
	if !exists {
		u.logger.Infof("room not found, skipping signaling, room id: %s, request id: %s", roomID, requestID)
		return
	}
	room.Lock.Lock()
	defer room.Lock.Unlock()

	for _, conn := range room.Connections {
		if conn.Pc.SignalingState() == webrtc.SignalingStateClosed ||
			conn.Pc.ConnectionState() == webrtc.PeerConnectionStateClosed {
			continue
		}

		// создаем оффер только если стабильное состояние
		if conn.Pc.SignalingState() != webrtc.SignalingStateStable {
			continue
		}

		existingSenders := map[string]bool{}
		for _, sender := range conn.Pc.GetSenders() {
			if sender.Track() == nil {
				continue
			}
			existingSenders[sender.Track().ID()] = true
			if _, ok := room.TrackLocals[sender.Track().ID()]; !ok {
				_ = conn.Pc.RemoveTrack(sender)
			}
		}

		for trackID := range room.TrackLocals {
			if _, ok := existingSenders[trackID]; !ok {
				_, _ = conn.Pc.AddTransceiverFromTrack(room.TrackLocals[trackID], webrtc.RTPTransceiverInit{
					Direction: webrtc.RTPTransceiverDirectionSendonly,
				})
			}
		}

		offer, err := conn.Pc.CreateOffer(nil)
		if err != nil {
			u.logger.Errorf("failed to create offer, error: %v, uuid: %s, request id: %s", err, conn.Kws.GetUUID(), requestID)
			continue
		}

		if err = conn.Pc.SetLocalDescription(offer); err != nil {
			u.logger.Errorf("failed to set local description, error: %v, uuid: %s, request id: %s", err, conn.Kws.GetUUID(), requestID)
			continue
		}

		offerString, _ := json.Marshal(offer)
		_ = conference_utils.WriteJSON(conn.Kws, &conn.Mu, &conference_utils.WebsocketMessage{
			Event: "offer",
			Data:  string(offerString),
		})
		u.logger.Infof("offer sent, uuid: %s, request id: %s", conn.Kws.GetUUID(), requestID)
	}
}

func (u *ConferenceUsecase) StartKeyFrameDispatcher() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for range ticker.C {
		if atomic.LoadUint32(&u.serverRunning) == 0 {
			break
		}
		u.roomsLock.RLock()
		for roomID, r := range u.rooms {
			u.dispatchKeyFrame(roomID)
			u.logger.Infof("dispatching key frame, room id: %s, track count: %d", roomID, r.TrackCount)
		}
		u.roomsLock.RUnlock()
	}
}

func (u *ConferenceUsecase) dispatchKeyFrame(roomID string) {
	room := u.rooms[roomID]

	room.Lock.RLock()
	defer room.Lock.RUnlock()

	for i := range room.Connections {
		for _, receiver := range room.Connections[i].Pc.GetReceivers() {
			if receiver.Track() == nil {
				continue
			}
			_ = room.Connections[i].Pc.WriteRTCP([]rtcp.Packet{
				&rtcp.PictureLossIndication{
					MediaSSRC: uint32(receiver.Track().SSRC()),
				},
			})
		}
	}
}
