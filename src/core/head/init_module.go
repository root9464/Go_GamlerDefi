package core

import (
	referral_module "github.com/root9464/Go_GamlerDefi/src/module/referral"
	test_module "github.com/root9464/Go_GamlerDefi/src/module/test"
	ton_module "github.com/root9464/Go_GamlerDefi/src/module/ton"
	validation_module "github.com/root9464/Go_GamlerDefi/src/module/validation"
)

type Modules struct {
	test       *test_module.TestModule
	referral   *referral_module.ReferralModule
	validation *validation_module.ValidationModule
	ton        *ton_module.TonModule
	// jwt        *jwt_module.JwtModule
}

func (m *Core) init_modules() {
	m.modules = &Modules{
		test:       test_module.NewTestModule(m.logger),
		referral:   referral_module.NewReferralModule(m.config, m.logger, m.validator, m.database, m.ton_client, m.ton_api),
		validation: validation_module.NewValidationModule(m.config, m.logger, m.validator, m.database, m.ton_api),
		ton:        ton_module.NewTonModule(m.config, m.logger),
		// jwt:        jwt_module.NewJwtModule(m.logger, m.validator, m.config.PrivateKey, m.config.PublicKey),
	}
}
