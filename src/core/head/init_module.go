package core

import (
	game_hub_module "github.com/root9464/Go_GamlerDefi/src/modules/game_hub"
	referral_module "github.com/root9464/Go_GamlerDefi/src/modules/referral"
	test_module "github.com/root9464/Go_GamlerDefi/src/modules/test"
	ton_module "github.com/root9464/Go_GamlerDefi/src/modules/ton"
	validation_module "github.com/root9464/Go_GamlerDefi/src/modules/validation"
)

type Modules struct {
	test       *test_module.TestModule
	referral   *referral_module.ReferralModule
	validation *validation_module.ValidationModule
	ton        *ton_module.TonModule
	game_hub   *game_hub_module.GameHubModule
	// conference *conference_module.ConferenceModule
	// jwt        *jwt_module.JwtModule
}

func (m *Core) init_modules() {
	m.modules = &Modules{
		test:       test_module.NewTestModule(m.logger),
		referral:   referral_module.NewReferralModule(m.config, m.logger, m.validator, m.database, m.ton_client, m.ton_api),
		validation: validation_module.NewValidationModule(m.config, m.logger, m.validator, m.database, m.ton_api),
		game_hub:   game_hub_module.NewGameHubModule(m.logger),
		ton:        ton_module.NewTonModule(m.config, m.logger),
		// jwt:        jwt_module.NewJwtModule(m.logger, m.validator, m.config.PrivateKey, m.config.PublicKey),
	}
}
