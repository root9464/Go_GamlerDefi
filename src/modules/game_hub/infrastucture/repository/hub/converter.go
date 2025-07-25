package hub_repository

import (
	hub_entity "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/entity"
	hub_model "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/infrastucture/repository/model"
)

type Converter struct{}

func (c *Converter) ToModel(hub *hub_entity.Hub) *hub_model.Hub {
	return &hub_model.Hub{
		ID:        hub.ID,
		GameID:    hub.GameID,
		HostID:    hub.HostID,
		Status:    string(hub.Status),
		StartTime: hub.StartTime,
		EndTime:   hub.EndTime,
		EntryFee:  hub.EntryFee,
		Currency:  string(hub.Currency),
	}
}

func (c *Converter) ToEntity(hub *hub_model.Hub) *hub_entity.Hub {
	return &hub_entity.Hub{
		ID:        hub.ID,
		GameID:    hub.GameID,
		HostID:    hub.HostID,
		Status:    hub_entity.HubStatus(hub.Status),
		StartTime: hub.StartTime,
		EndTime:   hub.EndTime,
		EntryFee:  hub.EntryFee,
		Currency:  hub_entity.Currency(hub.Currency),
	}
}
