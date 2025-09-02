package file_service

import (
	"fmt"
	"path/filepath"

	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
	file_entity "github.com/root9464/Go_GamlerDefi/src/submodules/file/entity"
)

type IFileStorage interface {
	Upload(name *string, data []byte) error
	Delete(name string) error
	Get(name string) ([]byte, error)
}

type FileService struct {
	fileStorage IFileStorage
	logger      *logger.Logger
}

func NewFileService(fileStorage IFileStorage, logger *logger.Logger) *FileService {
	return &FileService{
		fileStorage: fileStorage,
		logger:      logger,
	}
}

func (s *FileService) Upload(data []byte) (string, error) {
	fileEntity := &file_entity.File{
		Data: data,
	}
	if err := fileEntity.GenerateName(); err != nil {
		s.logger.Errorf("failed to generate file name: %v", err)
		return "", err
	}

	fmt.Printf("File service: %+v \n", s)
	if err := s.fileStorage.Upload(&fileEntity.Name, fileEntity.Data); err != nil {
		s.logger.Errorf("failed to upload file: %v", err)
		return "", err
	}

	return fileEntity.Name, nil
}

func (s *FileService) Get(name string) (*file_entity.File, error) {
	file, err := s.fileStorage.Get(name)
	if err != nil {
		return nil, err
	}

	return &file_entity.File{
		Name: name,
		Data: file,
	}, nil
}

func (s *FileService) Delete(name string) error {
	return s.fileStorage.Delete(name)
}

func (s *FileService) GetContentType(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".pdf":
		return "application/pdf"
	case ".txt":
		return "text/plain"
	case ".html":
		return "text/html"
	default:
		return "application/octet-stream"
	}
}
