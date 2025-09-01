package file_module

import (
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
	file_filesystem "github.com/root9464/Go_GamlerDefi/src/submodules/file/infrastucture/filesystem"
	file_service "github.com/root9464/Go_GamlerDefi/src/submodules/file/service"
)

type FileModule struct {
	file_storage *file_filesystem.LocalFileStorage
	File_service *file_service.FileService

	logger  *logger.Logger
	fileDir string
}

func NewFileModule(logger *logger.Logger, fileDir string) *FileModule {
	return &FileModule{
		logger:  logger,
		fileDir: fileDir,
	}
}

func (m *FileModule) Init() {
	m.file_storage = file_filesystem.NewLocalFileStorage(m.logger, m.fileDir)
	m.File_service = file_service.NewFileService(m.file_storage, m.logger)
}
