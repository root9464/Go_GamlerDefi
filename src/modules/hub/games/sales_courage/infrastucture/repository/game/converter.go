package sales_courage_repository

import (
	sales_courage_entity "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/sales_courage/entity"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type converter struct {
	logger *logger.Logger
}

func newGameConverter(logger *logger.Logger) *converter {
	return &converter{logger: logger}
}

func (c *converter) EntityToModel(entity *sales_courage_entity.Game)
