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

	status := game_session_models.StatusScheduled
	switch entity.Status {
	case game_session_entity.StatusActive:
		status = game_session_models.StatusActive
	case game_session_entity.StatusFinished:
		status = game_session_models.StatusFinished
	}

	players := []game_session_models.Player{}
	for _, player := range entity.Participants {
		players = append(players, game_session_models.Player{
			ID:   player.ID,
			Name: player.Name,
		})
	}

	return &game_session_models.GameSession{
		ID:       objectID,
		HostID:   entity.HostID,
		GameName: entity.GameName,

		Participants: players,
		Status:       status,
		TimeToStart:  entity.TimeToStart,
		Price:        entity.Price,
	}
}

func (c *Converter) ModelToEntity(model *game_session_models.GameSession) *game_session_entity.GameSession {
	status := game_session_entity.StatusScheduled
	switch model.Status {
	case game_session_models.StatusActive:
		status = game_session_entity.StatusActive
	case game_session_models.StatusFinished:
		status = game_session_entity.StatusFinished
	}

	players := []game_session_entity.Player{}
	for _, player := range model.Participants {
		players = append(players, game_session_entity.Player{
			ID:   player.ID,
			Name: player.Name,
		})
	}

	return &game_session_entity.GameSession{
		ID:           bson.ObjectID(model.ID).Hex(),
		HostID:       model.HostID,
		GameName:     model.GameName,
		Participants: players,
		Status:       status,
		TimeToStart:  model.TimeToStart,
		Price:        model.Price,
	}
}

func (c *Converter) StringToID(id string) (bson.ObjectID, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return objectID, err
	}
	return objectID, nil
}
