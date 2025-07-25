package hub_repository

import (
	"context"

	hub_entity "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/entity"
	hub_model "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/infrastucture/repository/model"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
	"gorm.io/gorm"
)

type HubRepository struct {
	logger   *logger.Logger
	db       *gorm.DB
	converer Converter
}

func NewHubRepository(logger *logger.Logger, db *gorm.DB) *HubRepository {
	return &HubRepository{
		logger:   logger,
		db:       db,
		converer: Converter{},
	}
}

func (r *HubRepository) Create(ctx context.Context, hub *hub_entity.Hub) error {
	r.logger.Infof("creating hub: %v", hub)
	hubModel := r.converer.ToModel(hub)
	if err := r.db.WithContext(ctx).Create(hubModel).Error; err != nil {
		r.logger.Errorf("error creating hub: %v", err)
		return err
	}
	hub.ID = hubModel.ID
	r.logger.Info("hub created successfully")
	return nil
}

func (r *HubRepository) Update(ctx context.Context, hub *hub_entity.Hub) error {
	r.logger.Infof("updating hub: %v", hub)
	hubModel := r.converer.ToModel(hub)
	if err := r.db.WithContext(ctx).Save(hubModel).Error; err != nil {
		r.logger.Errorf("error updating hub: %v", err)
		return err
	}
	r.logger.Info("hub updated successfully")
	return nil
}

func (r *HubRepository) Delete(ctx context.Context, id string) error {
	r.logger.Infof("deleting hub by id: %v", id)
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&hub_model.Hub{}).Error; err != nil {
		r.logger.Errorf("error deleting hub by id: %v", err)
		return err
	}
	r.logger.Info("hub deleted successfully")
	return nil
}

func (r *HubRepository) GtByID(ctx context.Context, id string) (*hub_entity.Hub, error) {
	r.logger.Infof("getting hub by id: %v", id)
	hubModel := new(hub_model.Hub)
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(hubModel).Error; err != nil {
		r.logger.Errorf("error getting hub by id: %v", err)
		return nil, err
	}
	hub := r.converer.ToEntity(hubModel)
	r.logger.Info("hub got successfully")
	return hub, nil
}

func (r *HubRepository) GetAll(ctx context.Context) ([]hub_entity.Hub, error) {
	r.logger.Info("getting all hubs")
	hubsModel := make([]hub_model.Hub, 0)
	if err := r.db.WithContext(ctx).Find(&hubsModel).Error; err != nil {
		r.logger.Errorf("error getting all hubs: %v", err)
		return nil, err
	}
	hubs := make([]hub_entity.Hub, 0)
	for _, hubModel := range hubsModel {
		hubs = append(hubs, *r.converer.ToEntity(&hubModel))
	}
	r.logger.Info("hubs got successfully")
	return hubs, nil
}
