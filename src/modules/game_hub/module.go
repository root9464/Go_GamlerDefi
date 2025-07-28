package game_hub_module

import (
	"github.com/gofiber/fiber/v2"
	conference_ws "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/delivery/ws"
	conference_usecase "github.com/root9464/Go_GamlerDefi/src/modules/game_hub/usecase/conference"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type GameHubModule struct {
	conference_usecase conference_usecase.IConferenceUsecase
	conference_ws      *conference_ws.WSHandler
	logger             *logger.Logger
}

func NewGameHubModule(logger *logger.Logger) *GameHubModule {
	return &GameHubModule{
		logger: logger,
	}
}

func (m *GameHubModule) init() {
	m.conference_usecase = conference_usecase.NewConferenceUsecase(m.logger)
	m.conference_ws = conference_ws.NewWSHanler(m.logger, m.conference_usecase)
}

func (m *GameHubModule) InitDelivery(router fiber.Router) {
	m.init()
	// Регистрация маршрутов, например:
	// router.Get("/websocket", m.conference_ws.ConferenceWebsocketHandler)
	m.conference_ws.RegisterRoutes(router) // Предполагая, что у вас есть такой метод
}
