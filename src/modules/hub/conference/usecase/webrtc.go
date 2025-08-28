package conference_usecase

import (
	"encoding/json"
	"io"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
	conference_entity "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/entity"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/utils"
)

func (u *ConferenceUsecase) SetubWebRTC(conn *conference_entity.Connection, r *conference_entity.Room, requestID string) {
	pc := conn.Pc
	kws := conn.Kws

	pc.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}
		candidateString, err := json.Marshal(i.ToJSON())
		if err != nil {
			u.logger.Errorf("failed to marshal candidate error: %v, request id: %s", err, requestID)
			return
		}
		if err := conference_utils.WriteJSON(kws, &conn.Lock, &conference_utils.WebsocketMessage{
			Event: "candidate",
			Data:  string(candidateString),
		}); err != nil {
			u.logger.Errorf("failed to send candidate error: %v, request id: %s", err, requestID)

		}
	})

	pc.OnConnectionStateChange(func(p webrtc.PeerConnectionState) {
		u.logger.Infof("peer connection state changed state: %s, uuid: %s, request id: %s", p.String(), kws.GetUUID(), requestID)
		switch p {
		case webrtc.PeerConnectionStateFailed:
			if err := pc.Close(); err != nil {
				u.logger.Errorf("failed to close peer connection error: %v, request id: %s", err, requestID)
			}
		case webrtc.PeerConnectionStateClosed:
			u.logger.Infof("peer connection closed, uuid: %s, request id: %s", kws.GetUUID(), requestID)
			u.SignalPeerConnections(requestID, conn.RoomID)
		}
	})

	pc.OnTrack(func(t *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		if conn.Pc.ConnectionState() != webrtc.PeerConnectionStateConnected {
			u.logger.Infof("track received but connection not stable track id: %s, uuid: %s, request id: %s, ssrc: %d", t.ID(), kws.GetUUID(), requestID, uint32(t.SSRC()))
			return
		}
		u.logger.Infof("track received track id: %s, uuid: %s, request id: %s, ssrc: %d", t.ID(), kws.GetUUID(), requestID, uint32(t.SSRC()))

		if u.audioRecorder != nil && t.Kind() == webrtc.RTPCodecTypeAudio {
			u.audioRecorder.StartRecordingTrack(t, conn.RoomID, kws.GetUUID())
		}

		trackLocal := u.AddTrack(conn, t)
		if trackLocal == nil {
			return
		}

		go func() {
			defer func() {
				if u.audioRecorder != nil && t.Kind() == webrtc.RTPCodecTypeAudio {
					u.audioRecorder.StopRecordingTrack(t.ID(), conn.RoomID, kws.GetUUID())
				}

				u.RemoveTrack(trackLocal, conn.RoomID)
				u.logger.Infof("track removed track id: %s, uuid: %s, request id: %s", t.ID(), kws.GetUUID(), requestID)
			}()

			rtpPkt := &rtp.Packet{}
			for {
				select {
				case <-conn.Closed:
					u.logger.Infof("track processing stopped due to connection closed track id: %s, request id: %s", t.ID(), requestID)
					return
				default:
					bufPtr := u.bufferPool.Get().(*[]byte)
					buf := *bufPtr
					i, _, err := t.Read(buf)
					if err != nil {
						u.bufferPool.Put(bufPtr)
						if err == io.EOF {
							u.logger.Infof("track closed track id: %s, request id: %s", t.ID(), requestID)
							return
						}
						u.logger.Errorf("failed to read rtp packet error: %v, track id: %s, uuid: %s, request id: %s", err, t.ID(), kws.GetUUID(), requestID)
						return
					}
					if err = rtpPkt.Unmarshal(buf[:i]); err != nil {
						u.bufferPool.Put(bufPtr)
						u.logger.Errorf("failed to unmarshal rtp packet error: %v, track id: %s, request id: %s", err, t.ID(), requestID)
						return
					}
					rtpPkt.Extension = false
					rtpPkt.Extensions = nil
					if err = trackLocal.WriteRTP(rtpPkt); err != nil {
						u.bufferPool.Put(bufPtr)
						u.logger.Errorf("failed to write rtp packet error: %v, track id: %s, request id: %s", err, t.ID(), requestID)
						return
					}
					u.bufferPool.Put(bufPtr)
				}
			}
		}()
	})

	pc.OnICEConnectionStateChange(func(is webrtc.ICEConnectionState) {
		u.logger.Infof("ice connection state changed state: %s, uuid: %s, request id: %s", is.String(), kws.GetUUID(), requestID)
	})
}
