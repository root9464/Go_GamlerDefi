package audiorecorder

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type TimeSyncRoomMixer struct {
	roomID    string
	outputDir string
	tempDir   string
}

func NewTimeSyncRoomMixer(roomID, outputDir string) *TimeSyncRoomMixer {
	tempDir := filepath.Join(outputDir, "temp_sync_mixing")

	return &TimeSyncRoomMixer{
		roomID:    roomID,
		outputDir: outputDir,
		tempDir:   tempDir,
	}
}

func (tsm *TimeSyncRoomMixer) CheckFFmpegAvailable() error {
	cmd := exec.Command("ffmpeg", "-version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("FFmpeg not available: %v", err)
	}
	return nil
}

func (tsm *TimeSyncRoomMixer) MixRoomAudio(roomMeta *RoomMetadata) error {
	log.Printf("Starting time-synchronized mixing for room: %s with %d users", tsm.roomID, len(roomMeta.Users))

	if len(roomMeta.Users) == 0 {
		return fmt.Errorf("no users found for room %s", tsm.roomID)
	}

	// Создаем временную директорию
	if err := os.MkdirAll(tsm.tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Логируем все временные метки для отладки
	roomMeta.mu.RLock()
	for userID, userRecord := range roomMeta.Users {
		relativeStart := userRecord.JoinTime.Sub(roomMeta.StartTime).Seconds()
		duration := userRecord.LeaveTime.Sub(userRecord.JoinTime).Seconds()
		log.Printf("User %s: join=%.2fs, duration=%.2fs", userID, relativeStart, duration)
	}
	roomMeta.mu.RUnlock()

	// Проверяем существование директории комнаты
	roomDir := filepath.Join(tsm.outputDir, tsm.roomID)
	if _, err := os.Stat(roomDir); err != nil {
		return fmt.Errorf("room directory does not exist: %s", roomDir)
	}

	// Ждем готовности файлов
	if err := tsm.waitForFilesReady(roomDir); err != nil {
		log.Printf("Warning: %v", err)
	}

	// Этап 1: Конкатенируем сегменты каждого пользователя
	userFiles, err := tsm.concatenateAllUserSegments(roomMeta)
	if err != nil {
		return fmt.Errorf("failed to concatenate user segments: %w", err)
	}

	if len(userFiles) == 0 {
		return fmt.Errorf("no user audio files found for room %s", tsm.roomID)
	}

	log.Printf("Successfully concatenated %d user files", len(userFiles))

	// Этап 2: Микшируем с временной синхронизацией
	mixedFile := filepath.Join(tsm.outputDir, fmt.Sprintf("%s_mixed.ogg", tsm.roomID))
	if err := tsm.mixWithTimeSync(userFiles, roomMeta, mixedFile); err != nil {
		return fmt.Errorf("failed to mix with time sync: %w", err)
	}

	tsm.cleanup()
	log.Printf("Time-synchronized mixing completed for room %s", tsm.roomID)
	return nil
}

func (tsm *TimeSyncRoomMixer) concatenateAllUserSegments(roomMeta *RoomMetadata) ([]string, error) {
	roomDir := filepath.Join(tsm.outputDir, tsm.roomID)
	var userFiles []string

	roomMeta.mu.RLock()
	defer roomMeta.mu.RUnlock()

	for userID := range roomMeta.Users {
		userPath := filepath.Join(roomDir, userID)
		segments, err := tsm.collectUserSegments(userPath)
		if err != nil || len(segments) == 0 {
			log.Printf("No segments found for user %s", userID)
			continue
		}

		concatenatedFile := filepath.Join(tsm.tempDir, fmt.Sprintf("user_%s.ogg", userID))
		if err := tsm.concatenateWithFFmpeg(segments, concatenatedFile); err != nil {
			log.Printf("Failed to concatenate segments for user %s: %v", userID, err)
			continue
		}

		userFiles = append(userFiles, concatenatedFile)
		log.Printf("Concatenated %d segments for user %s", len(segments), userID)
	}

	return userFiles, nil
}

func (tsm *TimeSyncRoomMixer) collectUserSegments(userPath string) ([]string, error) {
	files, err := os.ReadDir(userPath)
	if err != nil {
		return nil, err
	}

	var segments []string
	for _, file := range files {
		if !file.Type().IsRegular() || filepath.Ext(file.Name()) != ".ogg" {
			continue
		}
		segments = append(segments, filepath.Join(userPath, file.Name()))
	}

	// ПРАВИЛЬНАЯ сортировка по номеру сегмента
	sort.Slice(segments, func(i, j int) bool {
		name1 := filepath.Base(segments[i])
		name2 := filepath.Base(segments[j])

		var num1, num2 int
		fmt.Sscanf(name1, "segment_%d.ogg", &num1)
		fmt.Sscanf(name2, "segment_%d.ogg", &num2)

		return num1 < num2
	})

	return segments, nil
}

func (tsm *TimeSyncRoomMixer) concatenateWithFFmpeg(segments []string, outputFile string) error {
	if len(segments) == 0 {
		return fmt.Errorf("no segments to concatenate")
	}

	if len(segments) == 1 {
		return tsm.copyFile(segments[0], outputFile)
	}

	listFile := filepath.Join(tsm.tempDir, fmt.Sprintf("concat_%d.txt", time.Now().UnixNano()))
	f, err := os.Create(listFile)
	if err != nil {
		return err
	}

	for _, segment := range segments {
		absPath, _ := filepath.Abs(segment)
		fmt.Fprintf(f, "file '%s'\n", absPath)
	}
	f.Close()

	// Правильная конкатенация с перекодированием
	cmd := exec.Command("ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", listFile,
		"-c:a", "libopus",
		"-b:a", "128k",
		"-ar", "48000",
		"-ac", "1",
		"-avoid_negative_ts", "make_zero",
		"-fflags", "+genpts",
		"-y",
		outputFile)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ffmpeg concat failed: %v, output: %s", err, string(output))
	}

	os.Remove(listFile)
	return nil
}

func (tsm *TimeSyncRoomMixer) mixWithTimeSync(userFiles []string, roomMeta *RoomMetadata, outputFile string) error {
	if len(userFiles) == 0 {
		return fmt.Errorf("no user files to mix")
	}

	if len(userFiles) == 1 {
		return tsm.copyFile(userFiles[0], outputFile)
	}

	args := []string{"-y"}

	// Добавляем входные файлы
	for _, file := range userFiles {
		args = append(args, "-i", file)
	}

	// Создаем фильтр с правильными задержками
	var filterParts []string
	roomMeta.mu.RLock()
	for i, file := range userFiles {
		userID := tsm.extractUserIDFromFilename(file)
		var delayMs int

		// Находим задержку для этого пользователя
		if userRecord, exists := roomMeta.Users[userID]; exists {
			delaySeconds := userRecord.JoinTime.Sub(roomMeta.StartTime).Seconds()
			delayMs = int(delaySeconds * 1000) // Конвертируем в миллисекунды
		}

		if delayMs > 0 {
			filterParts = append(filterParts, fmt.Sprintf("[%d:a]adelay=%d[a%d]", i, delayMs, i))
		} else {
			filterParts = append(filterParts, fmt.Sprintf("[%d:a]acopy[a%d]", i, i))
		}

		log.Printf("User %s will have delay of %dms", userID, delayMs)
	}
	roomMeta.mu.RUnlock()

	// Создаем финальный микс
	var mixInputs []string
	for i := range userFiles {
		mixInputs = append(mixInputs, fmt.Sprintf("[a%d]", i))
	}

	filterComplex := strings.Join(filterParts, ";") + ";" +
		strings.Join(mixInputs, "") +
		fmt.Sprintf("amix=inputs=%d:duration=longest:dropout_transition=0:normalize=0", len(userFiles))

	args = append(args,
		"-filter_complex", filterComplex,
		"-c:a", "libopus",
		"-b:a", "128k",
		"-ar", "48000",
		"-ac", "1",
		"-application", "voip",
		"-frame_duration", "20",
		"-avoid_negative_ts", "make_zero",
		"-fflags", "+genpts",
		outputFile)

	cmd := exec.Command("ffmpeg", args...)
	log.Printf("Running FFmpeg mix command: %s", cmd.String())

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ffmpeg mixing failed: %v, output: %s", err, string(output))
	}

	log.Printf("Successfully mixed %d tracks with time synchronization", len(userFiles))
	return nil
}

func (tsm *TimeSyncRoomMixer) extractUserIDFromFilename(filename string) string {
	base := filepath.Base(filename)
	if strings.HasPrefix(base, "user_") && strings.HasSuffix(base, ".ogg") {
		return strings.TrimSuffix(strings.TrimPrefix(base, "user_"), ".ogg")
	}
	return ""
}

func (tsm *TimeSyncRoomMixer) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func (tsm *TimeSyncRoomMixer) waitForFilesReady(roomDir string) error {
	maxWait := 10 * time.Second
	checkInterval := 500 * time.Millisecond

	start := time.Now()
	for time.Since(start) < maxWait {
		allReady := true

		err := filepath.Walk(roomDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if filepath.Ext(path) == ".ogg" && info.Size() == 0 {
				allReady = false
				return filepath.SkipDir
			}

			return nil
		})

		if err != nil {
			return err
		}

		if allReady {
			return nil
		}

		time.Sleep(checkInterval)
	}

	return fmt.Errorf("timeout waiting for files to be ready")
}

func (tsm *TimeSyncRoomMixer) cleanup() {
	os.RemoveAll(tsm.tempDir)
}
