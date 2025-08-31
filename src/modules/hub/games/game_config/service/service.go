package game_config_service

import (
	"context"
	"fmt"

	"github.com/mitchellh/mapstructure"
	game_session_registry "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/usecase/registry"
	game_config_dto "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/game_config/dto"
	game_config_model "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/game_config/model"
	game_config_repository "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/game_config/repository"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type GameConfigService struct {
	logger     *logger.Logger
	repository *game_config_repository.GameConfigRepository
}

func NewGameConfigService(logger *logger.Logger, repository *game_config_repository.GameConfigRepository) *GameConfigService {
	return &GameConfigService{logger: logger, repository: repository}
}

func (s *GameConfigService) CreateOrUpdate(ctx context.Context, dto *game_config_dto.CreateGameConfigDTO) error {
	if !game_session_registry.IsGameRegistered(dto.GameName) {
		return fmt.Errorf("game %s not registered", dto.GameName)
	}

	settingsModel, _ := game_session_registry.GetGameSettingsModel(dto.GameName)

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:      settingsModel,
		ErrorUnused: true,
	})

	if err != nil {
		return err
	}

	if err := decoder.Decode(dto.Settings); err != nil {
		return fmt.Errorf("setting structure is not valid: %v", err)
	}

	configModel := &game_config_model.GameConfigModel{
		ID:          dto.GameName,
		Title:       dto.Title,
		Description: dto.Description,
		Settings:    dto.Settings,
	}
	return s.repository.Upsert(ctx, configModel)
}
