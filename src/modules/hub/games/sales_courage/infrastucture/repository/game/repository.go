package sales_courage_repository

import (
	games_constant "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/constant"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IGameRepository interface{}

type GameRepository struct {
	logger *logger.Logger
	db     *mongo.Database
}

func NewGameRepository(logger *logger.Logger, db *mongo.Database) IGameRepository {
	return &GameRepository{logger: logger, db: db}
}

func (r *GameRepository) GetGameCollection() *mongo.Collection {
	return r.db.Collection(games_constant.GamesCollection)
}

// func (r *GameRepository) CreateGame(ctx context.Context, game *sales_courage_models.Game) error {
// 	r.logger.Info("creating game...")
// 	collection := r.GetGameCollection()
// 	return err
// }
