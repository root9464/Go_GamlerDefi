package audiorecorder

import (
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
	mu              sync.RWMutex
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
	}
}

func (ar *AudioRecorder) StartRecordingTrack(trackRemote *webrtc.TrackRemote, roomID, userID string) {
	if trackRemote.Kind() != webrtc.RTPCodecTypeAudio {
		return
	}

	codec := trackRemote.Codec()
	clockRate := codec.ClockRate

	// Для Opus всегда используем 48kHz для лучшего качества
	if codec.MimeType == "audio/opus" {
		clockRate = 48000
	}

	trackKey := fmt.Sprintf("%s_%s_%s", roomID, userID, trackRemote.ID())

	ar.mu.Lock()
	if _, exists := ar.activeRecorders[trackKey]; exists {
		ar.mu.Unlock()
		return
	}

	recorder := &TrackRecorder{
		trackID:         trackRemote.ID(),
		userID:          userID,
		roomID:          roomID,
		segmentDuration: 5 * time.Second,
		segmentIndex:    1,
		// Используем оптимизированные параметры для лучшего качества
		sampleBuilder: samplebuilder.New(10, &codecs.OpusPacket{}, clockRate),
		clockRate:     clockRate,
		stopChan:      make(chan struct{}),
	}

	ar.activeRecorders[trackKey] = recorder
	ar.mu.Unlock()

	go recorder.processTrack(trackRemote, ar.outputDir)
}

func (ar *AudioRecorder) StopRecordingTrack(trackID, roomID, userID string) {
	trackKey := fmt.Sprintf("%s_%s_%s", roomID, userID, trackID)

	ar.mu.Lock()
	recorder, exists := ar.activeRecorders[trackKey]
	if exists {
		delete(ar.activeRecorders, trackKey)
	}
	ar.mu.Unlock()

	if exists {
		close(recorder.stopChan)
	}
}

func (tr *TrackRecorder) processTrack(track *webrtc.TrackRemote, baseOutputDir string) {
	defer tr.cleanup()

	userDir := filepath.Join(baseOutputDir, tr.roomID, tr.userID)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		log.Printf("Ошибка создания папки %s: %v", userDir, err)
		return
	}

	if err := tr.startNewSegment(userDir); err != nil {
		log.Printf("Ошибка создания первого сегмента: %v", err)
		return
	}

	for {
		select {
		case <-tr.stopChan:
			return
		default:
			packet, _, err := track.ReadRTP()
			if err != nil {
				log.Printf("Ошибка чтения RTP: %v", err)
				return
			}

			// Обрабатываем timestamp'ы для правильного времени
			tr.handleTimestamp(packet, uint8(track.PayloadType()))

			// Пушим пакет в SampleBuilder
			tr.sampleBuilder.Push(packet)

			// Обрабатываем готовые сэмплы
			for sample := tr.sampleBuilder.Pop(); sample != nil; sample = tr.sampleBuilder.Pop() {
				if time.Since(tr.segmentStart) >= tr.segmentDuration {
					if err := tr.startNewSegment(userDir); err != nil {
						log.Printf("Ошибка создания нового сегмента: %v", err)
						return
					}
				}

				if tr.currentWriter != nil && len(sample.Data) > 0 {
					// Создаем RTP пакет с правильными параметр��ми для качества
					rtpPacket := &rtp.Packet{
						Header: rtp.Header{
							PayloadType: uint8(track.PayloadType()),
							Timestamp:   sample.PacketTimestamp,
							Marker:      true,
						},
						Payload: sample.Data,
					}

					if err := tr.currentWriter.WriteRTP(rtpPacket); err != nil {
						log.Printf("Ошибка записи сэмпла: %v", err)
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

	// Вычисляем разность timestamp'ов с учетом переполнения
	var timestampDiff uint32
	if packet.Timestamp >= tr.lastTimestamp {
		timestampDiff = packet.Timestamp - tr.lastTimestamp
	} else {
		// Обработка переполнения timestamp
		timestampDiff = (0xFFFFFFFF - tr.lastTimestamp) + packet.Timestamp + 1
	}

	// Если пропуск больше 100ms, заполняем тишиной
	maxGap := tr.clockRate / 10                                   // 100ms
	if timestampDiff > maxGap && timestampDiff < tr.clockRate*2 { // но не больше 2 секунд
		tr.fillSilenceGap(tr.lastTimestamp, packet.Timestamp, payloadType)
	}

	tr.lastTimestamp = packet.Timestamp
}

func (tr *TrackRecorder) fillSilenceGap(startTimestamp, endTimestamp uint32, payloadType uint8) {
	// Opus пакет тишины (DTX - Discontinuous Transmission)
	silencePayload := []byte{0xF8, 0xFF}

	// Шаг зависит от частоты дискретизации
	var timestampStep uint32
	switch tr.clockRate {
	case 48000:
		timestampStep = 960 // 20ms при 48kHz
	case 16000:
		timestampStep = 320 // 20ms при 16kHz
	case 8000:
		timestampStep = 160 // 20ms при 8kHz
	default:
		timestampStep = tr.clockRate / 50 // 20ms для любой частоты
	}

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

	// Используем стерео для лучшего качества, если возможно
	channels := uint16(1) // моно по умолчанию
	if tr.clockRate >= 16000 {
		channels = 2 // стерео для высокого качества
	}

	writer, err := oggwriter.New(fullPath, tr.clockRate, channels)
	if err != nil {
		return err
	}

	tr.currentWriter = writer
	tr.segmentStart = time.Now()
	tr.segmentIndex++

	log.Printf("Начата запись сегмента: %s (%d Hz, %d каналов)", fullPath, tr.clockRate, channels)
	return nil
}

func (tr *TrackRecorder) cleanup() {
	tr.mu.Lock()
	defer tr.mu.Unlock()

	// Принудительно обрабатываем оставшиеся пакеты
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

// Дополнительный метод для остановки всех записей
func (ar *AudioRecorder) StopAll() {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	for _, recorder := range ar.activeRecorders {
		close(recorder.stopChan)
	}
	ar.activeRecorders = make(map[string]*TrackRecorder)
}

func (ar *AudioRecorder) MixAndCleanupRoom(roomID string) error {
	// Останавливаем все записи для этой комнаты
	ar.mu.Lock()
	var roomRecorders []*TrackRecorder
	for key, recorder := range ar.activeRecorders {
		if recorder.roomID == roomID {
			roomRecorders = append(roomRecorders, recorder)
			delete(ar.activeRecorders, key)
		}
	}
	ar.mu.Unlock()

	// Останавливаем все рекордеры комнаты
	for _, recorder := range roomRecorders {
		close(recorder.stopChan)
	}

	// Ждем завершения записи
	time.Sleep(2 * time.Second)

	// Используем production микшер
	mixer := NewProductionRoomMixer(roomID, ar.outputDir)

	// Проверяем доступность FFmpeg
	if err := mixer.CheckFFmpegAvailable(); err != nil {
		return fmt.Errorf("FFmpeg required for production mixing: %w", err)
	}

	// Микшируем аудио
	if err := mixer.MixRoomAudio(); err != nil {
		return fmt.Errorf("failed to mix room audio: %w", err)
	}

	log.Printf("Room %s audio mixed successfully with production quality", roomID)
	return nil
}
