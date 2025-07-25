package hub_usecase

import (
	"context"

	hub_entity "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/entity"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type IHubUsecase interface {
}

type IHubRepository interface {
	Create(ctx context.Context, hub *hub_entity.Hub) error
	Update(ctx context.Context, hub *hub_entity.Hub) error
	Delete(ctx context.Context, id string) error
	GtByID(ctx context.Context, id string) (*hub_entity.Hub, error)
	GetAll(ctx context.Context) ([]hub_entity.Hub, error)
}

type HubUsecase struct {
	logger     *logger.Logger
	repository IHubRepository
}

func NewHubUsecase(logger *logger.Logger, repository IHubRepository) IHubUsecase {
	return &HubUsecase{
		logger:     logger,
		repository: repository,
	}
}
