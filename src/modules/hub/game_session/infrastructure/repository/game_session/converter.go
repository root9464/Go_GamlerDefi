package game_session_repository

import (
	game_session_entity "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/entity"
	game_session_models "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/infrastructure/repository/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Converter struct{}

func NewConverter() *Converter {
	return &Converter{}
}

func (c *Converter) EntityToModel(entity *game_session_entity.GameSession) *game_session_models.GameSession {
	return &game_session_models.GameSession{
		ID:       entity.ID,
		HostID:   entity.HostID,
		GameName: entity.GameName,

		TimeToStart: entity.TimeToStart,
	}
}

func (c *Converter) ModelToEntity(model *game_session_models.GameSession) *game_session_entity.GameSession {
	return &game_session_entity.GameSession{
		ID:          model.ID,
		HostID:      model.HostID,
		GameName:    model.GameName,
		TimeToStart: model.TimeToStart,
	}
}

func (c *Converter) StringToID(id string) (bson.ObjectID, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return objectID, err
	}
	return objectID, nil
}
