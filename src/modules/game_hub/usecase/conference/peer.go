package conference_usecase

import (
	"encoding/json"
	"sync/atomic"
	"time"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v4"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/utils/conference"
)

func (u *ConferenceUsecase) SignalPeerConnections(requestID string, roomID string) {
	u.roomsLock.RLock()
	room, exists := u.rooms[roomID]
	u.roomsLock.RUnlock()
	if !exists {
		u.logger.Info("Room not found, skipping signaling",
			"room_id", roomID,
			"request_id", requestID)
		return
	}
	room.Lock.Lock()
	defer room.Lock.Unlock()

	for _, conn := range room.Connections {
		if time.Since(conn.LastSignal) < conn.SignalDebounce {
			u.logger.Debug("Skipping signaling due to debounce",
				"uuid", conn.Kws.GetUUID(),
				"request_id", requestID,
			)
			continue
		}

		if conn.Pc.ConnectionState() == webrtc.PeerConnectionStateClosed {
			continue
		}

		if conn.Pc.SignalingState() != webrtc.SignalingStateStable && conn.Pc.SignalingState() != webrtc.SignalingStateHaveRemoteOffer {
			u.logger.Debug("Skipping signaling due to non-stable signaling state",
				"uuid", conn.Kws.GetUUID(),
				"state", conn.Pc.SignalingState().String(),
				"request_id", requestID,
			)
			continue
		}

		existingSenders := map[string]bool{}
		for _, sender := range conn.Pc.GetSenders() {
			if sender.Track() == nil {
				continue
			}
			existingSenders[sender.Track().ID()] = true
			if _, ok := room.TrackLocals[sender.Track().ID()]; !ok {
				if err := conn.Pc.RemoveTrack(sender); err != nil {
					u.logger.Error("Failed to remove track",
						"error", err,
						"uuid", conn.Kws.GetUUID(),
						"request_id", requestID,
					)
					continue
				}
			}
		}

		for _, receiver := range conn.Pc.GetReceivers() {
			if receiver.Track() == nil {
				continue
			}
			existingSenders[receiver.Track().ID()] = true
		}

		for trackID := range room.TrackLocals {
			if _, ok := existingSenders[trackID]; !ok {
				_, err := conn.Pc.AddTransceiverFromTrack(room.TrackLocals[trackID], webrtc.RTPTransceiverInit{
					Direction: webrtc.RTPTransceiverDirectionSendonly,
				})
				if err != nil {
					u.logger.Error("Failed to add track",
						"error", err,
						"track_id", trackID,
						"uuid", conn.Kws.GetUUID(),
						"request_id", requestID,
					)
					continue
				}
			}
		}

		offer, err := conn.Pc.CreateOffer(nil)
		if err != nil {
			u.logger.Error("Failed to create offer",
				"error", err,
				"uuid", conn.Kws.GetUUID(),
				"request_id", requestID,
			)
			continue
		}

		if err = conn.Pc.SetLocalDescription(offer); err != nil {
			u.logger.Error("Failed to set local description",
				"error", err,
				"uuid", conn.Kws.GetUUID(),
				"request_id", requestID,
			)
			continue
		}

		offerString, err := json.Marshal(offer)
		if err != nil {
			u.logger.Error("Failed to marshal offer",
				"error", err,
				"uuid", conn.Kws.GetUUID(),
				"request_id", requestID,
			)
			continue
		}

		if err := conference_utils.WriteJSON(conn.Kws, &conn.Lock, &conference_utils.WebsocketMessage{
			Event: "offer",
			Data:  string(offerString),
		}); err != nil {
			u.logger.Error("Failed to send offer",
				"error", err,
				"uuid", conn.Kws.GetUUID(),
				"request_id", requestID,
			)
			continue
		}

		conn.LastSignal = time.Now()
		u.logger.Info("Offer sent",
			"uuid", conn.Kws.GetUUID(),
			"request_id", requestID,
		)
	}

	if len(room.Connections) > 0 {
		go func() {
			time.Sleep(time.Second * 5)
			if atomic.LoadUint32(&u.serverRunning) == 1 {
				u.SignalPeerConnections(conference_utils.GenerateRequestID(), roomID)
			}
		}()
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
			u.logger.Info("Dispatching key frame",
				"room_id", roomID,
				"track_count", r.TrackCount)
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
