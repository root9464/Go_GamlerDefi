package core

import test_module "github.com/root9464/Go_GamlerDefi/module/test"

type Modules struct {
	test *test_module.TestModule
}

func (m *Core) init_modules() {
	m.modules = &Modules{
		test: test_module.NewTestModule(m.logger),
	}
}
