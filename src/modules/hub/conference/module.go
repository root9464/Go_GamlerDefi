package conference_module

import (
	"github.com/gofiber/fiber/v2"
	conference_ws "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/delivery/ws"
	conference_usecase "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/usecase"
	logger "github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type ConferenceModule struct {
	conference_usecase conference_usecase.IConferenceUsecase
	conference_ws      *conference_ws.WSHandler
	logger             *logger.Logger
}

func NewConferencebModule(logger *logger.Logger) *ConferenceModule {
	return &ConferenceModule{
		logger: logger,
	}
}

func (m *ConferenceModule) init() {
	m.conference_usecase = conference_usecase.NewConferenceUsecase(m.logger)
	m.conference_ws = conference_ws.NewWSHanler(m.logger, m.conference_usecase)
}

func (m *ConferenceModule) InitDelivery(router fiber.Router) {
	m.init()
	go m.conference_usecase.StartKeyFrameDispatcher()

	m.conference_ws.RegisterRoutes(router)
}
