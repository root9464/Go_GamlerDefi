package core

import (
	referral_module "github.com/root9464/Go_GamlerDefi/module/referral"
	test_module "github.com/root9464/Go_GamlerDefi/module/test"
)

type Modules struct {
	test     *test_module.TestModule
	referral *referral_module.ReferralModule
}

func (m *Core) init_modules() {
	m.modules = &Modules{
		test:     test_module.NewTestModule(m.logger),
		referral: referral_module.NewReferralModule(m.logger, m.validator, m.ton_client, m.ton_api),
	}
}
