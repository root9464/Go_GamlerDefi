package conference_service

import (
	"encoding/json"
	"io"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
	conference_entity "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/entity"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/util"
)

func (s *ConferenceService) SetubWebRTC(conn *conference_entity.Connection, r *conference_entity.Room, requestID string) {
	pc := conn.Pc
	kws := conn.Kws

	pc.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}
		candidateString, err := json.Marshal(i.ToJSON())
		if err != nil {
			s.logger.Errorf("Failed to marshal candidate, error: %v, request_id: %s", err, requestID)
			return
		}

		err = conference_utils.WriteJSON(kws, &conn.Lock, &conference_utils.WebsocketMessage{
			Event: "candidate",
			Data:  string(candidateString),
		})

		if err != nil {
			s.logger.Errorf("Failed to send candidate, error: %v, request_id: %s", err, requestID)
		}
	})

	pc.OnConnectionStateChange(func(p webrtc.PeerConnectionState) {
		s.logger.Infof("PeerConnection state changed, state: %s, uuid: %s, request_id: %s", p.String(), kws.GetUUID(), requestID)
		switch p {
		case webrtc.PeerConnectionStateFailed:
			if err := pc.Close(); err != nil {
				s.logger.Errorf("Failed to close PeerConnection, error: %v, request_id: %s", err, requestID)
			}
		case webrtc.PeerConnectionStateClosed:
			s.SignalPeerConnections(requestID, conn.RoomID)
		}
	})

	pc.OnTrack(func(t *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		if conn.Pc.ConnectionState() != webrtc.PeerConnectionStateConnected {
			s.logger.Warnf("track received but connection not stable, track_id: %s, uuid: %s, ssrc: %d", t.ID(), kws.GetUUID(), uint32(t.SSRC()))
			return
		}

		s.logger.Infof("track received, track_id: %s, uuid: %s, ssrc: %d, kind: %s",
			t.ID(), kws.GetUUID(), uint32(t.SSRC()), t.Kind().String())

		trackLocal := s.AddTrack(conn, t)
		if trackLocal == nil {
			return
		}

		go func() {
			defer func() {
				s.RemoveTrack(trackLocal, conn.RoomID)
				s.logger.Infof("track removed, track_id: %s, uuid: %s, ssrc: %d", t.ID(), kws.GetUUID(), uint32(t.SSRC()))
			}()

			rtpPkt := new(rtp.Packet)

			for {
				select {
				case <-conn.Closed:
					s.logger.Infof("track removed, track_id: %s, uuid: %s, ssrc: %d", t.ID(), kws.GetUUID(), uint32(t.SSRC()))
					return
				default:
					bufPtr := s.bufferPool.Get().(*[]byte)
					buf := *bufPtr

					i, _, err := t.Read(buf)
					if err != nil {
						s.bufferPool.Put(bufPtr)
						if err == io.EOF {
							s.logger.Infof("track removed, track_id: %s, uuid: %s, ssrc: %d", t.ID(), kws.GetUUID(), uint32(t.SSRC()))
							return
						}
						s.logger.Errorf("failed to read RTP packet, error: %v, request_id: %s", err, requestID)
						return
					}

					if err = rtpPkt.Unmarshal(buf[:i]); err != nil {
						s.bufferPool.Put(bufPtr)
						s.logger.Infof("track removed, track_id: %s, uuid: %s, ssrc: %d", t.ID(), kws.GetUUID(), uint32(t.SSRC()))
						return
					}

					rtpPkt.Extension = false
					rtpPkt.Extensions = nil

					if err = trackLocal.WriteRTP(rtpPkt); err != nil {
						s.bufferPool.Put(bufPtr)
						s.logger.Infof("track removed, track_id: %s, uuid: %s, ssrc: %d", t.ID(), kws.GetUUID(), uint32(t.SSRC()))
						return
					}

					s.bufferPool.Put(bufPtr)
				}
			}
		}()
	})

	pc.OnICEConnectionStateChange(func(is webrtc.ICEConnectionState) {
		s.logger.Infof("ICE connection state changed, state: %s, uuid: %s, request_id: %s", is.String(), kws.GetUUID(), requestID)
	})
}
