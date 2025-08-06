package conference_usecase

import (
	"encoding/json"
	"io"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
	hub_entity "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/entity"
	conference_utils "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/utils/conference"
)

func (u *ConferenceUsecase) SetubWebRTC(conn *hub_entity.Connection, r *hub_entity.Room, requestID string) {
	pc := conn.Pc
	kws := conn.Kws

	pc.OnICECandidate(func(i *webrtc.ICECandidate) {
		if i == nil {
			return
		}
		candidateString, err := json.Marshal(i.ToJSON())
		if err != nil {
			u.logger.Error("Failed to marshal candidate",
				"error", err,
				"request_id", requestID,
			)
			return
		}
		if err := conference_utils.WriteJSON(kws, &conn.Lock, &conference_utils.WebsocketMessage{
			Event: "candidate",
			Data:  string(candidateString),
		}); err != nil {
			u.logger.Error("Failed to send candidate",
				"error", err,
				"request_id", requestID,
			)
		}
	})

	pc.OnConnectionStateChange(func(p webrtc.PeerConnectionState) {
		u.logger.Info("PeerConnection state changed",
			"state", p.String(),
			"uuid", kws.GetUUID(),
			"request_id", requestID,
		)
		switch p {
		case webrtc.PeerConnectionStateFailed:
			if err := pc.Close(); err != nil {
				u.logger.Error("Failed to close PeerConnection",
					"error", err,
					"request_id", requestID,
				)
			}
		case webrtc.PeerConnectionStateClosed:
			u.SignalPeerConnections(requestID, conn.RoomID)
		}
	})

	pc.OnTrack(func(t *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
		if conn.Pc.ConnectionState() != webrtc.PeerConnectionStateConnected {
			u.logger.Info("Track received but connection not stable",
				"track_id", t.ID(),
				"uuid", kws.GetUUID(),
				"request_id", requestID,
				"ssrc", uint32(t.SSRC()),
			)
			return
		}
		u.logger.Info("Track received",
			"track_id", t.ID(),
			"uuid", kws.GetUUID(),
			"request_id", requestID,
			"ssrc", uint32(t.SSRC()),
		)

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
				u.logger.Info("Track removed",
					"track_id", t.ID(),
					"uuid", kws.GetUUID(),
					"request_id", requestID,
				)
			}()

			rtpPkt := &rtp.Packet{}
			for {
				select {
				case <-conn.Closed:
					u.logger.Info("Track processing stopped due to connection closed",
						"track_id", t.ID(),
						"request_id", requestID,
					)
					return
				default:
					bufPtr := u.bufferPool.Get().(*[]byte)
					buf := *bufPtr
					i, _, err := t.Read(buf)
					if err != nil {
						u.bufferPool.Put(bufPtr)
						if err == io.EOF {
							u.logger.Info("Track closed",
								"track_id", t.ID(),
								"request_id", requestID,
							)
							return
						}
						u.logger.Error("Failed to read RTP packet",
							"error", err,
							"track_id", t.ID(),
							"uuid", kws.GetUUID(),
							"request_id", requestID,
						)
						return
					}
					if err = rtpPkt.Unmarshal(buf[:i]); err != nil {
						u.bufferPool.Put(bufPtr)
						u.logger.Error("Failed to unmarshal RTP packet",
							"error", err,
							"track_id", t.ID(),
							"request_id", requestID,
						)
						return
					}
					rtpPkt.Extension = false
					rtpPkt.Extensions = nil
					if err = trackLocal.WriteRTP(rtpPkt); err != nil {
						u.bufferPool.Put(bufPtr)
						u.logger.Error("Failed to write RTP packet",
							"error", err,
							"track_id", t.ID(),
							"request_id", requestID,
						)
						return
					}
					u.bufferPool.Put(bufPtr)
				}
			}
		}()
	})

	pc.OnICEConnectionStateChange(func(is webrtc.ICEConnectionState) {
		u.logger.Info("ICE connection state changed",
			"state", is.String(),
			"uuid", kws.GetUUID(),
			"request_id", requestID,
		)
	})
}
