package conference_service

import (
	"encoding/json"
	"sync/atomic"
	"time"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v4"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/util"
)

func (s *ConferenceService) SignalPeerConnections(requestID string, roomID string) {
	s.roomsLock.RLock()
	room, exists := s.rooms[roomID]
	s.roomsLock.RUnlock()
	if !exists {
		s.logger.Infof("Room not found, skipping signaling, room_id: %s, request_id: %s", roomID, requestID)
		return
	}
	room.Lock.Lock()
	defer room.Lock.Unlock()

	for _, conn := range room.Connections {
		if time.Since(conn.LastSignal) < conn.SignalDebounce {
			s.logger.Debugf("Skipping signaling due to debounce uuid: %s, request_id: %s", conn.Kws.GetUUID(), requestID)
			continue
		}

		if conn.Pc.ConnectionState() == webrtc.PeerConnectionStateClosed {
			continue
		}

		if conn.Pc.SignalingState() != webrtc.SignalingStateStable && conn.Pc.SignalingState() != webrtc.SignalingStateHaveRemoteOffer {
			s.logger.Debugf("Skipping signaling due to non-stable signaling state uuid: %s, state: %s, request_id: %s", conn.Kws.GetUUID(), conn.Pc.SignalingState().String(), requestID)
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
					s.logger.Errorf("Failed to remove track, error: %v, uuid: %s, request_id: %s", err, conn.Kws.GetUUID(), requestID)
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
					s.logger.Errorf("Failed to add track, error: %v, uuid: %s, request_id: %s", err, conn.Kws.GetUUID(), requestID)
					continue
				}
			}
		}

		offer, err := conn.Pc.CreateOffer(nil)
		if err != nil {
			s.logger.Errorf("Failed to create offer, error: %v, uuid: %s, request_id: %s", err, conn.Kws.GetUUID(), requestID)
			continue
		}

		if err = conn.Pc.SetLocalDescription(offer); err != nil {
			s.logger.Errorf("Failed to set local description, error: %v, uuid: %s, request_id: %s", err, conn.Kws.GetUUID(), requestID)
			continue
		}

		offerString, err := json.Marshal(offer)
		if err != nil {
			s.logger.Errorf("Failed to marshal offer, error: %v, uuid: %s, request_id: %s", err, conn.Kws.GetUUID(), requestID)
			continue
		}

		if err := conference_utils.WriteJSON(conn.Kws, &conn.Lock, &conference_utils.WebsocketMessage{
			Event: "offer",
			Data:  string(offerString),
		}); err != nil {
			s.logger.Errorf("Failed to send offer, error: %v, uuid: %s, request_id: %s", err, conn.Kws.GetUUID(), requestID)
			continue
		}

		conn.LastSignal = time.Now()

		s.logger.Infof("Offer sent, uuid: %s, request_id: %s", conn.Kws.GetUUID(), requestID)
	}

	if len(room.Connections) > 0 {
		go func() {
			time.Sleep(time.Second * 5)
			if atomic.LoadUint32(&s.serverRunning) == 1 {
				s.SignalPeerConnections(conference_utils.GenerateRequestID(), roomID)
			}
		}()
	}
}

func (s *ConferenceService) StartKeyFrameDispatcher() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	for range ticker.C {
		if atomic.LoadUint32(&s.serverRunning) == 0 {
			break
		}
		s.roomsLock.RLock()
		for roomID, r := range s.rooms {
			s.dispatchKeyFrame(roomID)
			s.logger.Infof("Dispatching key frame, room_id: %s, track_count: %d", roomID, r.TrackCount)
		}
		s.roomsLock.RUnlock()
	}
}

func (s *ConferenceService) dispatchKeyFrame(roomID string) {
	room := s.rooms[roomID]

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
