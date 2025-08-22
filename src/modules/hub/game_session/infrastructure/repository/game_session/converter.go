package game_session_repository

import (
	"fmt"

	game_session_entity "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/entity"
	game_session_models "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session/infrastructure/repository/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Converter struct{}

func NewConverter() *Converter {
	return &Converter{}
}

func (c *Converter) EntityToModel(entity *game_session_entity.GameSession) *game_session_models.GameSession {
	objectID, err := bson.ObjectIDFromHex(entity.ID)
	if err != nil {
		fmt.Println(err)
	}

	return &game_session_models.GameSession{
		ID:     objectID,
		HostID: entity.HostID,
		GameID: entity.GameID,
	}
}

func (c *Converter) ModelToEntity(model *game_session_models.GameSession) *game_session_entity.GameSession {
	return &game_session_entity.GameSession{
		ID:     bson.ObjectID(model.ID).Hex(),
		HostID: model.HostID,
		GameID: model.GameID,
	}
}
