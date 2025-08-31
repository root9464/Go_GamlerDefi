package conference_module

import (
	"github.com/gofiber/fiber/v2"
	conference_ws "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/handler"
	conference_service "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference/service"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type ConferenceModule struct {
	conference_service conference_service.IConferenceService
	Conference_ws      *conference_ws.WSHandler
	logger             *logger.Logger
}

func NewConferencebModule(logger *logger.Logger) *ConferenceModule {
	return &ConferenceModule{
		logger: logger,
	}
}

func (m *ConferenceModule) Init() {
	m.conference_service = conference_service.NewConferenceUsecase(m.logger)
	m.Conference_ws = conference_ws.NewWSHanler(m.logger, m.conference_service)
}

func (m *ConferenceModule) InitDelivery(router fiber.Router) {
	go m.conference_service.StartKeyFrameDispatcher()
}
