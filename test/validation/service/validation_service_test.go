package validation_service_test

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/root9464/Go_GamlerDefi/database"
	validation_dto "github.com/root9464/Go_GamlerDefi/module/validation/dto"
	validation_repository "github.com/root9464/Go_GamlerDefi/module/validation/repository"
	validation_service "github.com/root9464/Go_GamlerDefi/module/validation/service"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tonkeeper/tonapi-go"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ValidationServiceTestSuite struct {
	suite.Suite
	logger     *logger.Logger
	validator  *validator.Validate
	ton_api    *tonapi.Client
	database   *mongo.Database
	service    validation_service.IValidationService
	repository validation_repository.IValidationRepository
}

const (
	db_url  = "mongodb://root:example@localhost:27017"
	db_name = "gamer_defi_test"
)

func (s *ValidationServiceTestSuite) SetupSuite() {
	s.logger = logger.GetLogger()
	validator := validator.New()
	client, err := tonapi.NewClient(tonapi.TestnetTonApiURL, &tonapi.Security{})
	require.NoError(s.T(), err)

	_, database, err := database.ConnectDatabase(db_url, s.logger, db_name)
	require.NoError(s.T(), err)

	s.validator = validator
	s.ton_api = client
	s.database = database

	s.repository = validation_repository.NewValidationRepository(s.logger, s.database)
	s.service = validation_service.NewValidationService(s.logger, s.ton_api, s.repository)
}

func (s *ValidationServiceTestSuite) MockTransaction() validation_dto.WorkerTransactionDTO {
	paymentOrderId, err := bson.ObjectIDFromHex("6826ac79ff2f0eb00db5fa1d")
	require.NoError(s.T(), err, "failed to convert payment order id to bson.ObjectID")
	transaction := validation_dto.WorkerTransactionDTO{
		TxHash:         "105f7620bf78d534941ebcf97dda0dbe8e79c134a8ab346843787c71fe3308d5",
		TxQueryID:      1747000636,
		TargetAddress:  "0QANsjLvOX2MERlT4oyv2bSPEVc9lunSPIs5a1kPthCXydUX",
		PaymentOrderId: paymentOrderId.Hex(),
		Status:         validation_dto.WorkerStatusPending,
		CreatedAt:      time.Now().Unix(),
		UpdatedAt:      time.Now().Unix(),
	}
	return transaction
}

func (s *ValidationServiceTestSuite) TestRunnerTransaction_Success() {
	transaction := s.MockTransaction()
	tr, state, err := s.service.RunnerTransaction(&transaction)
	require.NoError(s.T(), err, "transaction should be success")
	require.True(s.T(), state, "transaction should be success")
	require.NotNil(s.T(), tr, "transaction should not be nil")
}

func (s *ValidationServiceTestSuite) TestSubWorkerTransaction_Success() {
	transaction := s.MockTransaction()
	tr, err := s.service.SubWorkerTransaction(&transaction)
	require.NoError(s.T(), err, "transaction should be success")
	require.True(s.T(), tr, "transaction should be success")
}

func (s *ValidationServiceTestSuite) TestWorkerTransaction_Success() {
	transaction := s.MockTransaction()
	tr, err := s.service.WorkerTransaction(&transaction)
	require.NoError(s.T(), err, "transaction should be success")
	require.True(s.T(), tr, "transaction should be success")
}

func TestValidationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ValidationServiceTestSuite))
}
