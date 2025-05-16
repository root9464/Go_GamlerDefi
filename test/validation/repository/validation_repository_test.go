package validation_repository_test

import (
	"testing"
	"time"

	"github.com/root9464/Go_GamlerDefi/database"
	validation_model "github.com/root9464/Go_GamlerDefi/module/validation/model"
	validation_repository "github.com/root9464/Go_GamlerDefi/module/validation/repository"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	db_url  = "mongodb://root:example@localhost:27017"
	db_name = "gamer_defi_test"
)

type ValidationRepositoryTestSuite struct {
	suite.Suite
	db         *mongo.Database
	repository validation_repository.IValidationRepository
	logger     *logger.Logger
}

func (s *ValidationRepositoryTestSuite) SetupSuite() {
	s.logger = logger.GetLogger()
	_, database, err := database.ConnectDatabase(db_url, s.logger, db_name)
	require.NoError(s.T(), err, "Failed to connect to database")
	s.db = database
	repo := validation_repository.NewValidationRepository(s.logger, s.db)
	s.repository = repo
}

func (s *ValidationRepositoryTestSuite) TestCreateTransactionObserver() {

	paymentOrderId, err := bson.ObjectIDFromHex("6823e92b5d53ea679cbd4426")
	require.NoError(s.T(), err, "Failed to convert payment order id to bson.ObjectID")
	transaction := validation_model.WorkerTransaction{
		ID:                 bson.NewObjectID(),
		TxHash:             "105f7620bf78d534941ebcf97dda0dbe8e79c134a8ab346843787c71fe3308d5",
		TxQueryID:          1747000636,
		TargetJettonSymbol: "FROGE",
		TargetJettonMaster: "kQAE0xZ5bHIOdDBCGxNEwoJCzptm5bpcs5KtVIqQbl3-CL0N",
		TargetAddress:      "0QANsjLvOX2MERlT4oyv2bSPEVc9lunSPIs5a1kPthCXydUX",
		PaymentOrderId:     paymentOrderId,
		Status:             validation_model.WorkerStatusPending,
		CreatedAt:          time.Now().Unix(),
		UpdatedAt:          time.Now().Unix(),
	}
	transaction, err = s.repository.CreateTransactionObserver(transaction)
	require.NoError(s.T(), err, "Failed to create transaction observer")
	require.NotNil(s.T(), transaction, "Transaction observer should not be nil")
}

func (s *ValidationRepositoryTestSuite) TestGetTransactionObserver() {
	observerId, err := bson.ObjectIDFromHex("none")
	require.NoError(s.T(), err, "Failed to convert observer id to bson.ObjectID")
	transaction, err := s.repository.GetTransactionObserver(observerId)
	require.NoError(s.T(), err, "Failed to get transaction observer")
	require.NotNil(s.T(), transaction, "Transaction observer should not be nil")
}

func (s *ValidationRepositoryTestSuite) TestUpdateStatus() {
	observerId, err := bson.ObjectIDFromHex("none")
	require.NoError(s.T(), err, "Failed to convert observer id to bson.ObjectID")
	err = s.repository.UpdateStatus(observerId, validation_model.WorkerStatusFailed)
	require.NoError(s.T(), err, "Failed to update status")
}

func (s *ValidationRepositoryTestSuite) TestPrecheckoutTransaction() {
	observerId, err := bson.ObjectIDFromHex("none")
	require.NoError(s.T(), err, "Failed to convert observer id to bson.ObjectID")
	err = s.repository.PrecheckoutTransaction(observerId)
	require.NoError(s.T(), err, "Failed to precheckout transaction")
}

func (s *ValidationRepositoryTestSuite) TestDeleteTransactionObserver() {
	observerId, err := bson.ObjectIDFromHex("none")
	require.NoError(s.T(), err, "Failed to convert observer id to bson.ObjectID")
	err = s.repository.DeleteTransactionObserver(observerId)
	require.NoError(s.T(), err, "Failed to delete transaction observer")
}

func TestValidationRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ValidationRepositoryTestSuite))
}
