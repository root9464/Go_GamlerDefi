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

type ProductionRoomMixer struct {
	roomID    string
	outputDir string
	tempDir   string
}

func NewProductionRoomMixer(roomID, outputDir string) *ProductionRoomMixer {
	tempDir := filepath.Join(outputDir, "temp_mixing")
	os.MkdirAll(tempDir, 0755)

	return &ProductionRoomMixer{
		roomID:    roomID,
		outputDir: outputDir,
		tempDir:   tempDir,
	}
}

func (prm *ProductionRoomMixer) MixRoomAudio() error {
	log.Printf("Starting production audio mixing for room: %s", prm.roomID)

	// Ждем, пока все файлы будут готовы
	roomDir := filepath.Join(prm.outputDir, prm.roomID)
	if err := prm.waitForFilesReady(roomDir); err != nil {
		log.Printf("Warning: %v", err)
	}

	// Этап 1: Конкатенация сегментов каждого пользователя
	userFiles, err := prm.concatenateAllUserSegments()
	if err != nil {
		return fmt.Errorf("failed to concatenate user segments: %w", err)
	}

	if len(userFiles) == 0 {
		return fmt.Errorf("no user audio files found for room %s", prm.roomID)
	}

	// Этап 2: Реальное микширование через FFmpeg
	mixedFile := filepath.Join(prm.outputDir, fmt.Sprintf("%s_mixed.ogg", prm.roomID))
	if err := prm.mixWithFFmpeg(userFiles, mixedFile); err != nil {
		return fmt.Errorf("failed to mix audio with FFmpeg: %w", err)
	}

	// Очистка временных файлов
	prm.cleanup()

	log.Printf("Production mixing completed for room %s", prm.roomID)
	return nil
}
func (prm *ProductionRoomMixer) concatenateAllUserSegments() ([]string, error) {
	roomDir := filepath.Join(prm.outputDir, prm.roomID)

	userDirs, err := os.ReadDir(roomDir)
	if err != nil {
		return nil, err
	}

	var userFiles []string

	for _, userDir := range userDirs {
		if !userDir.IsDir() {
			continue
		}

		userID := userDir.Name()
		userPath := filepath.Join(roomDir, userID)

		// Собираем все сегменты пользователя
		segments, err := prm.collectUserSegments(userPath)
		if err != nil || len(segments) == 0 {
			log.Printf("No segments found for user %s", userID)
			continue
		}

		// Конкатенируем через FFmpeg для лучшего качества
		concatenatedFile := filepath.Join(prm.tempDir, fmt.Sprintf("user_%s.ogg", userID))
		if err := prm.concatenateWithFFmpeg(segments, concatenatedFile); err != nil {
			log.Printf("Failed to concatenate segments for user %s: %v", userID, err)
			continue
		}

		userFiles = append(userFiles, concatenatedFile)
		log.Printf("Concatenated %d segments for user %s", len(segments), userID)
	}

	return userFiles, nil
}

func (prm *ProductionRoomMixer) mixWithFFmpeg(userFiles []string, outputFile string) error {
	if len(userFiles) == 0 {
		return fmt.Errorf("no user files to mix")
	}

	if len(userFiles) == 1 {
		// Только один пользователь - просто копируем
		return prm.copyFile(userFiles[0], outputFile)
	}

	// Строим команду FFmpeg для микширования
	args := []string{"-y"} // Перезаписывать выходной файл

	// Добавляем входные файлы
	for _, file := range userFiles {
		args = append(args, "-i", file)
	}

	// Создаем фильтр для микширования
	filterInputs := make([]string, len(userFiles))
	for i := range userFiles {
		filterInputs[i] = fmt.Sprintf("[%d:a]", i)
	}

	// Фильтр амикс для качественного микширования аудио
	mixFilter := fmt.Sprintf("%samix=inputs=%d:duration=longest:dropout_transition=2",
		strings.Join(filterInputs, ""), len(userFiles))

	args = append(args,
		"-filter_complex", mixFilter,
		"-c:a", "libopus", // Кодек Opus
		"-b:a", "128k", // Битрейт 128kbps
		"-ar", "48000", // Частота дискретизации 48kHz
		"-ac", "1", // Моно
		"-application", "voip", // Оптимизация для голоса
		"-frame_duration", "20", // 20ms фреймы
		"-packet_loss", "1", // Устойчивость к потерям
		outputFile)

	cmd := exec.Command("ffmpeg", args...)

	log.Printf("Running FFmpeg command: %s", cmd.String())

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("ffmpeg mixing failed: %v, output: %s", err, string(output))
	}

	log.Printf("Successfully mixed %d audio tracks", len(userFiles))
	return nil
}

func (prm *ProductionRoomMixer) copyFile(src, dst string) error {
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

func (prm *ProductionRoomMixer) cleanup() {
	os.RemoveAll(prm.tempDir)
}

// Проверка доступности FFmpeg
func (prm *ProductionRoomMixer) CheckFFmpegAvailable() error {
	cmd := exec.Command("ffmpeg", "-version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("FFmpeg not available: %v", err)
	}
	return nil
}

func (prm *ProductionRoomMixer) concatenateWithFFmpeg(segments []string, outputFile string) error {
	if len(segments) == 0 {
		return fmt.Errorf("no segments to concatenate")
	}

	if len(segments) == 1 {
		return prm.copyFile(segments[0], outputFile)
	}

	// Убеждаемся, что временная директория существует
	if err := os.MkdirAll(prm.tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Создаем уникальное имя для списка файлов
	listFileName := fmt.Sprintf("concat_list_%d.txt", time.Now().UnixNano())
	listFile := filepath.Join(prm.tempDir, listFileName)

	log.Printf("Creating concat list file: %s", listFile)

	f, err := os.Create(listFile)
	if err != nil {
		return fmt.Errorf("failed to create concat list file %s: %w", listFile, err)
	}

	// Записываем АБСОЛЮТНЫЕ пути
	validSegments := 0
	for _, segment := range segments {
		absPath, err := filepath.Abs(segment)
		if err != nil {
			log.Printf("Failed to get absolute path for %s: %v", segment, err)
			continue
		}

		// Проверяем существование и размер файла
		if info, err := os.Stat(absPath); err != nil {
			log.Printf("Segment file does not exist: %s", absPath)
			continue
		} else if info.Size() == 0 {
			log.Printf("Segment file is empty: %s", absPath)
			continue
		}

		fmt.Fprintf(f, "file '%s'\n", absPath)
		validSegments++
	}
	f.Close()

	if validSegments == 0 {
		os.Remove(listFile)
		return fmt.Errorf("no valid segments found")
	}

	// Проверяем, что файл действительно создался
	if _, err := os.Stat(listFile); err != nil {
		return fmt.Errorf("concat list file was not created: %w", err)
	}

	// Логируем содержимое файла для отладки
	if content, err := os.ReadFile(listFile); err == nil {
		log.Printf("Concat list file %s content (%d bytes):\n%s", listFile, len(content), string(content))
	} else {
		log.Printf("Failed to read concat list file: %v", err)
	}

	// Получаем абсолютный путь к выходному файлу
	absOutputFile, err := filepath.Abs(outputFile)
	if err != nil {
		os.Remove(listFile)
		return fmt.Errorf("failed to get absolute path for output file: %w", err)
	}

	// Конкатенируем через FFmpeg
	cmd := exec.Command("ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", listFile,
		"-c", "copy",
		"-y",
		absOutputFile)

	// НЕ устанавливаем рабочую директорию - используем абсолютные пути
	log.Printf("Running FFmpeg command: %s", cmd.String())

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Сохраняем файл списка для отладки при ошибке
		log.Printf("FFmpeg concat failed. List file preserved at: %s", listFile)
		log.Printf("FFmpeg output: %s", string(output))
		return fmt.Errorf("ffmpeg concat failed: %v", err)
	}

	// Удаляем файл списка только при успехе
	if err := os.Remove(listFile); err != nil {
		log.Printf("Failed to remove concat list file: %v", err)
	}

	log.Printf("Successfully concatenated %d segments to %s", validSegments, absOutputFile)
	return nil
}

func (prm *ProductionRoomMixer) collectUserSegments(userPath string) ([]string, error) {
	log.Printf("Collecting segments from user path: %s", userPath)

	// Проверяем, что директория существует
	if _, err := os.Stat(userPath); err != nil {
		return nil, fmt.Errorf("user path does not exist: %w", err)
	}

	files, err := os.ReadDir(userPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read user directory: %w", err)
	}

	log.Printf("Found %d files in user directory", len(files))

	var segments []string
	for _, file := range files {
		if !file.Type().IsRegular() {
			log.Printf("Skipping non-regular file: %s", file.Name())
			continue
		}

		fileName := file.Name()
		filePath := filepath.Join(userPath, fileName)

		log.Printf("Checking file: %s", filePath)

		// Проверяем расширение
		if filepath.Ext(fileName) != ".ogg" {
			log.Printf("Skipping non-ogg file: %s", fileName)
			continue
		}

		// Проверяем существование и размер файла
		if info, err := os.Stat(filePath); err != nil {
			log.Printf("Cannot stat file %s: %v", filePath, err)
			continue
		} else if info.Size() == 0 {
			log.Printf("Skipping empty file: %s", filePath)
			continue
		} else {
			log.Printf("Adding segment: %s (size: %d bytes)", filePath, info.Size())
			segments = append(segments, filePath)
		}
	}

	// Сортируем по времени создания
	sort.Slice(segments, func(i, j int) bool {
		info1, _ := os.Stat(segments[i])
		info2, _ := os.Stat(segments[j])
		return info1.ModTime().Before(info2.ModTime())
	})

	log.Printf("Found %d valid OGG segments for user path %s", len(segments), userPath)
	return segments, nil
}

// Добавляем метод для ожидания завершения записи файлов
func (prm *ProductionRoomMixer) waitForFilesReady(roomDir string) error {
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
