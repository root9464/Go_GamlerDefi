package audiorecorder

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pion/rtp"
	"github.com/pion/rtp/codecs"
	"github.com/pion/webrtc/v4"
	"github.com/pion/webrtc/v4/pkg/media/oggwriter"
	"github.com/pion/webrtc/v4/pkg/media/samplebuilder"
)

type AudioRecorder struct {
	outputDir       string
	activeRecorders map[string]*TrackRecorder
	roomMetadata    map[string]*RoomMetadata
	mu              sync.RWMutex
}

type RoomMetadata struct {
	RoomID    string                 `json:"room_id"`
	StartTime time.Time              `json:"start_time"`
	Users     map[string]*UserRecord `json:"users"`
	mu        sync.RWMutex
}

type UserRecord struct {
	UserID    string    `json:"user_id"`
	JoinTime  time.Time `json:"join_time"`
	LeaveTime time.Time `json:"leave_time,omitempty"`
}

type TrackRecorder struct {
	trackID         string
	userID          string
	roomID          string
	currentWriter   *oggwriter.OggWriter
	segmentDuration time.Duration
	segmentStart    time.Time
	segmentIndex    int
	sampleBuilder   *samplebuilder.SampleBuilder
	clockRate       uint32
	stopChan        chan struct{}
	mu              sync.Mutex
	lastTimestamp   uint32
	timestampInit   bool
}

func NewAudioRecorder(outputDir string) *AudioRecorder {
	return &AudioRecorder{
		outputDir:       outputDir,
		activeRecorders: make(map[string]*TrackRecorder),
		roomMetadata:    make(map[string]*RoomMetadata),
	}
}

func (ar *AudioRecorder) StartRecordingTrack(trackRemote *webrtc.TrackRemote, roomID, userID string) {
	if trackRemote.Kind() != webrtc.RTPCodecTypeAudio {
		return
	}

	clockRate := uint32(48000)
	trackKey := fmt.Sprintf("%s_%s_%s", roomID, userID, trackRemote.ID())
	currentTime := time.Now()

	ar.mu.Lock()
	defer ar.mu.Unlock()

	// Проверяем, что рекордер еще не существует
	if _, exists := ar.activeRecorders[trackKey]; exists {
		return
	}

	// Инициализируем метаданные комнаты если нужно
	if _, exists := ar.roomMetadata[roomID]; !exists {
		ar.roomMetadata[roomID] = &RoomMetadata{
			RoomID:    roomID,
			StartTime: currentTime,
			Users:     make(map[string]*UserRecord),
		}
		log.Printf("Room %s recording started at %v", roomID, currentTime)
	}

	// Добавляем пользователя в метаданные комнаты
	roomMeta := ar.roomMetadata[roomID]
	roomMeta.mu.Lock()
	if _, exists := roomMeta.Users[userID]; !exists {
		roomMeta.Users[userID] = &UserRecord{
			UserID:   userID,
			JoinTime: currentTime,
		}
		log.Printf("User %s joined room %s at %v (%.2fs from room start)",
			userID, roomID, currentTime, currentTime.Sub(roomMeta.StartTime).Seconds())
	}
	roomMeta.mu.Unlock()

	// Создаем рекордер
	recorder := &TrackRecorder{
		trackID:         trackRemote.ID(),
		userID:          userID,
		roomID:          roomID,
		segmentDuration: 5 * time.Second,
		segmentIndex:    1,
		sampleBuilder:   samplebuilder.New(50, &codecs.OpusPacket{}, clockRate),
		clockRate:       clockRate,
		stopChan:        make(chan struct{}),
	}

	ar.activeRecorders[trackKey] = recorder
	go recorder.processTrack(trackRemote, ar.outputDir)
}

func (ar *AudioRecorder) StopRecordingTrack(trackID, roomID, userID string) {
	trackKey := fmt.Sprintf("%s_%s_%s", roomID, userID, trackID)

	ar.mu.Lock()
	recorder, exists := ar.activeRecorders[trackKey]
	if exists {
		delete(ar.activeRecorders, trackKey)
	}

	// Обновляем время выхода пользователя
	if roomMeta, roomExists := ar.roomMetadata[roomID]; roomExists {
		roomMeta.mu.Lock()
		if userRecord, userExists := roomMeta.Users[userID]; userExists {
			userRecord.LeaveTime = time.Now()
		}
		roomMeta.mu.Unlock()
	}
	ar.mu.Unlock()

	if exists {
		close(recorder.stopChan)
		log.Printf("User %s left room %s", userID, roomID)
	}
}

func (tr *TrackRecorder) processTrack(track *webrtc.TrackRemote, baseOutputDir string) {
	defer tr.cleanup()

	userDir := filepath.Join(baseOutputDir, tr.roomID, tr.userID)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		log.Printf("Error creating directory %s: %v", userDir, err)
		return
	}

	if err := tr.startNewSegment(userDir); err != nil {
		log.Printf("Error creating first segment: %v", err)
		return
	}

	for {
		select {
		case <-tr.stopChan:
			return
		default:
			packet, _, err := track.ReadRTP()
			if err != nil {
				log.Printf("Error reading RTP: %v", err)
				return
			}

			tr.handleTimestamp(packet, uint8(track.PayloadType()))
			tr.sampleBuilder.Push(packet)

			for sample := tr.sampleBuilder.Pop(); sample != nil; sample = tr.sampleBuilder.Pop() {
				if time.Since(tr.segmentStart) >= tr.segmentDuration {
					if err := tr.startNewSegment(userDir); err != nil {
						log.Printf("Error creating new segment: %v", err)
						return
					}
				}

				if tr.currentWriter != nil && len(sample.Data) > 0 {
					rtpPacket := &rtp.Packet{
						Header: rtp.Header{
							PayloadType: uint8(track.PayloadType()),
							Timestamp:   sample.PacketTimestamp,
							Marker:      true,
						},
						Payload: sample.Data,
					}

					if err := tr.currentWriter.WriteRTP(rtpPacket); err != nil {
						log.Printf("Error writing sample: %v", err)
					}
				}
			}
		}
	}
}

func (tr *TrackRecorder) handleTimestamp(packet *rtp.Packet, payloadType uint8) {
	if !tr.timestampInit {
		tr.lastTimestamp = packet.Timestamp
		tr.timestampInit = true
		return
	}

	var timestampDiff uint32
	if packet.Timestamp >= tr.lastTimestamp {
		timestampDiff = packet.Timestamp - tr.lastTimestamp
	} else {
		timestampDiff = (0xFFFFFFFF - tr.lastTimestamp) + packet.Timestamp + 1
	}

	maxGap := tr.clockRate / 10
	if timestampDiff > maxGap && timestampDiff < tr.clockRate*2 {
		tr.fillSilenceGap(tr.lastTimestamp, packet.Timestamp, payloadType)
	}

	tr.lastTimestamp = packet.Timestamp
}

func (tr *TrackRecorder) fillSilenceGap(startTimestamp, endTimestamp uint32, payloadType uint8) {
	silencePayload := []byte{0xF8, 0xFF}
	timestampStep := uint32(960)

	for ts := startTimestamp + timestampStep; ts < endTimestamp; ts += timestampStep {
		silencePacket := &rtp.Packet{
			Header: rtp.Header{
				PayloadType: payloadType,
				Timestamp:   ts,
				Marker:      false,
			},
			Payload: silencePayload,
		}
		tr.sampleBuilder.Push(silencePacket)
	}
}

func (tr *TrackRecorder) startNewSegment(outputDir string) error {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	if tr.currentWriter != nil {
		tr.currentWriter.Close()
	}

	filename := fmt.Sprintf("segment_%03d.ogg", tr.segmentIndex)
	fullPath := filepath.Join(outputDir, filename)

	writer, err := oggwriter.New(fullPath, tr.clockRate, 1)
	if err != nil {
		return err
	}

	tr.currentWriter = writer
	tr.segmentStart = time.Now()
	tr.segmentIndex++

	return nil
}

func (tr *TrackRecorder) cleanup() {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	tr.sampleBuilder.Flush()
	for sample := tr.sampleBuilder.Pop(); sample != nil; sample = tr.sampleBuilder.Pop() {
		if tr.currentWriter != nil && len(sample.Data) > 0 {
			rtpPacket := &rtp.Packet{
				Header: rtp.Header{
					Timestamp: sample.PacketTimestamp,
					Marker:    true,
				},
				Payload: sample.Data,
			}
			tr.currentWriter.WriteRTP(rtpPacket)
		}
	}

	if tr.currentWriter != nil {
		tr.currentWriter.Close()
		tr.currentWriter = nil
	}
}

func (ar *AudioRecorder) MixAndCleanupRoom(roomID string) error {
	log.Printf("Starting room cleanup and mixing for room: %s", roomID)

	ar.mu.Lock()

	// Получаем метаданные комнаты
	roomMeta, roomExists := ar.roomMetadata[roomID]
	if !roomExists {
		ar.mu.Unlock()
		log.Printf("Room %s metadata not found", roomID)
		return fmt.Errorf("room %s not found", roomID)
	}

	// Собираем все активные рекордеры для этой комнаты
	var roomRecorders []*TrackRecorder
	for key, recorder := range ar.activeRecorders {
		if recorder.roomID == roomID {
			roomRecorders = append(roomRecorders, recorder)
			delete(ar.activeRecorders, key)
		}
	}

	// Удаляем метаданные комнаты
	delete(ar.roomMetadata, roomID)
	ar.mu.Unlock()

	// Обновляем время выхода для всех пользователей
	currentTime := time.Now()
	roomMeta.mu.Lock()
	for _, userRecord := range roomMeta.Users {
		if userRecord.LeaveTime.IsZero() {
			userRecord.LeaveTime = currentTime
		}
	}
	roomMeta.mu.Unlock()

	log.Printf("Found %d active recorders for room %s", len(roomRecorders), roomID)
	log.Printf("Room has %d users total", len(roomMeta.Users))

	// Останавливаем все рекордеры
	for _, recorder := range roomRecorders {
		close(recorder.stopChan)
	}

	// Ждем завершения записи
	time.Sleep(3 * time.Second)

	// Создаем выходную директорию если не существует
	if err := os.MkdirAll(ar.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Сохраняем метаданные времени
	if err := ar.saveRoomMetadata(roomID, roomMeta); err != nil {
		log.Printf("Failed to save room metadata: %v", err)
		return err
	}

	// Микшируем аудио
	mixer := NewTimeSyncRoomMixer(roomID, ar.outputDir)
	if err := mixer.CheckFFmpegAvailable(); err != nil {
		return fmt.Errorf("FFmpeg required: %w", err)
	}

	if err := mixer.MixRoomAudio(roomMeta); err != nil {
		return fmt.Errorf("failed to mix room audio: %w", err)
	}

	log.Printf("Room %s mixed successfully", roomID)
	return nil
}

func (ar *AudioRecorder) saveRoomMetadata(roomID string, roomMeta *RoomMetadata) error {
	metadataFile := filepath.Join(ar.outputDir, fmt.Sprintf("%s_timing.json", roomID))

	data, err := json.MarshalIndent(roomMeta, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal room metadata: %w", err)
	}

	if err := os.WriteFile(metadataFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write room metadata: %w", err)
	}

	log.Printf("Saved room metadata to %s", metadataFile)
	return nil
}
