package hub_usecase

import (
	"context"

	hub_entity "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/entity"
)

func (u *HubUsecase) Create(ctx context.Context, hub *hub_entity.Hub) error {
	if err := u.repository.Create(ctx, hub); err != nil {
		return err
	}
	return nil
}

func (u *HubUsecase) Join(ctx context.Context) error {
	return nil
}
