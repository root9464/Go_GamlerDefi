package file_filesystem

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type LocalFileStorage struct {
	logger  *logger.Logger
	fileDir string
}

func NewLocalFileStorage(
	logger *logger.Logger,
	fileDir string,
) *LocalFileStorage {
	return &LocalFileStorage{
		logger:  logger,
		fileDir: fileDir,
	}
}

func (s *LocalFileStorage) Upload(name *string, data []byte) error {
	outputPath := s.fileDir + *name
	reader := bytes.NewReader(data)

	// if img, err := webp.Decode(reader); err == nil {
	// 	return s.saveAsWebP(img, outputPath)
	// }

	reader.Seek(0, io.SeekStart)
	img, _, err := image.Decode(reader)
	if err == nil {
		newPath := s.replaceExtToWebP(outputPath)
		*name = *name + ".webp"
		return s.saveAsWebP(img, newPath)
	}

	if err := os.MkdirAll(s.fileDir, 0755); err != nil {
		s.logger.Errorf("failed to create directory: %s", s.fileDir)
		return err
	}

	if os.WriteFile(outputPath, data, 0644) != nil {
		s.logger.Errorf("failed to write file: %s", outputPath)
		return err
	}

	return nil
}

func (s *LocalFileStorage) Delete(name string) error {
	filePath := s.fileDir + name
	if err := os.Remove(filePath); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete file %s: %w", name, err)
		}
	}
	return nil
}

func (s *LocalFileStorage) Get(name string) ([]byte, error) {
	filePath := s.fileDir + name
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s not found", name)
	}
	return os.ReadFile(filePath)
}

func (s *LocalFileStorage) replaceExtToWebP(path string) string {
	ext := filepath.Ext(path)
	return strings.TrimSuffix(path, ext) + ".webp"
}

func (s *LocalFileStorage) saveAsWebP(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// options := webp.Options{Quality: 80, Lossless: false}
	// if err := webp.Encode(file, img, &options); err != nil {
	// 	return err
	// }

	return nil
}
