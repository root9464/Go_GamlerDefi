package validation_service_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/root9464/Go_GamlerDefi/src/database"
	validation_dto "github.com/root9464/Go_GamlerDefi/src/modules/validation/dto"
	validation_repository "github.com/root9464/Go_GamlerDefi/src/modules/validation/repository"
	validation_service "github.com/root9464/Go_GamlerDefi/src/modules/validation/service"

	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
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
	ctx := context.Background()
	tr, state, err := s.service.RunnerTransaction(ctx, &transaction)
	require.NoError(s.T(), err, "transaction should be success")
	require.True(s.T(), state, "transaction should be success")
	require.NotNil(s.T(), tr, "transaction should not be nil")
}

func (s *ValidationServiceTestSuite) TestSubWorkerTransaction_Success() {
	transaction := s.MockTransaction()
	ctx := context.Background()
	tr, state, err := s.service.SubWorkerTransaction(ctx, &transaction)
	require.NoError(s.T(), err, "transaction should be success")
	require.True(s.T(), state, "transaction should be success")
	require.NotNil(s.T(), tr, "transaction should not be nil")
}

func (s *ValidationServiceTestSuite) TestWorkerTransaction_Success() {
	transaction := s.MockTransaction()
	ctx := context.Background()
	tr, status, err := s.service.WorkerTransaction(ctx, &transaction)
	s.logger.Infof("transaction: %+v", tr)
	require.NoError(s.T(), err, "transaction should be success")
	require.True(s.T(), status, "transaction should be success")
	require.NotNil(s.T(), tr, "transaction should not be nil")
}

func TestValidationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ValidationServiceTestSuite))
}
