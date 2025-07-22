package test_service

import "github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"

type ITestService interface {
	Ping() string
}

type testService struct {
	logger *logger.Logger
}

func NewTestService(logger *logger.Logger) ITestService {
	return &testService{logger: logger}
}

func (s *testService) Ping() string {
	return "pong"
}
