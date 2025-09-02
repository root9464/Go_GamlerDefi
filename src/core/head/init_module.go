package core

import (
	"sync"

	conference_module "github.com/root9464/Go_GamlerDefi/src/modules/hub/conference"
	game_session_module "github.com/root9464/Go_GamlerDefi/src/modules/hub/game_session"
	game_config_module "github.com/root9464/Go_GamlerDefi/src/modules/hub/games/game_config"
	referral_module "github.com/root9464/Go_GamlerDefi/src/modules/referral"
	test_module "github.com/root9464/Go_GamlerDefi/src/modules/test"
	ton_module "github.com/root9464/Go_GamlerDefi/src/modules/ton"
	validation_module "github.com/root9464/Go_GamlerDefi/src/modules/validation"
	file_module "github.com/root9464/Go_GamlerDefi/src/submodules/file"
)

type Modules struct {
	test       *test_module.TestModule
	referral   *referral_module.ReferralModule
	validation *validation_module.ValidationModule
	ton        *ton_module.TonModule

	conference   *conference_module.ConferenceModule
	game_session *game_session_module.GameSessionModule
	game_config  *game_config_module.GameConfigModule

	file_module *file_module.FileModule
}

var monce sync.Once

func (m *Core) init_modules() {
	m.modules = &Modules{}
	monce.Do(func() {
		m.initTestModule()
		m.initReferralModule()
		m.initValidationModule()
		m.initConferenceModule()
		m.initGameConfigModule()
		m.initGameSessionModule()
		m.initTonModule()
	})

	// m.modules = &Modules{
	// 	test:       test_module.NewTestModule(m.logger),
	// 	referral:   referral_module.NewReferralModule(m.config, m.logger, m.validator, m.database, m.ton_client, m.ton_api),
	// 	validation: validation_module.NewValidationModule(m.config, m.logger, m.validator, m.database, m.ton_api),
	// 	conference: conference_module.NewConferencebModule(m.logger),
	//
	// 	game_config:  game_config_module.NewGameConfigModule(m.logger, m.database),
	// 	game_session: game_session_module.NewGameSessionModule(m.logger, m.database, m.postgres, m.modules.game_config.Game_config_repository),
	// 	ton:          ton_module.NewTonModule(m.config, m.logger),
	// }
}

func (m *Core) initTestModule() {
	m.modules.test = test_module.NewTestModule(m.logger)
}

func (m *Core) initReferralModule() {
	m.modules.referral = referral_module.NewReferralModule(m.config, m.logger, m.validator, m.database, m.ton_client, m.ton_api)
}

func (m *Core) initValidationModule() {
	m.modules.validation = validation_module.NewValidationModule(m.config, m.logger, m.validator, m.database, m.ton_api)
}

func (m *Core) initConferenceModule() {

	m.logger.Warnf("Conf ws = %v", m.modules.conference)

	m.modules.conference = conference_module.NewConferencebModule(m.logger)
	m.modules.conference.Init()

	m.logger.Warnf("Conf ws = %v", m.modules.conference.Conference_ws)

}

func (m *Core) initGameConfigModule() {
	m.modules.game_config = game_config_module.NewGameConfigModule(m.logger, m.database)
}

func (m *Core) initGameSessionModule() {
	m.modules.game_session = game_session_module.NewGameSessionModule(m.logger, m.database, m.postgres, m.modules.game_config.Game_config_repository, m.modules.conference.Conference_ws, m.config)
}

func (m *Core) initTonModule() {
	m.modules.ton = ton_module.NewTonModule(m.config, m.logger)
}
