package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	referral_model "github.com/root9464/Go_GamlerDefi/module/referral/model"
	referral_repository "github.com/root9464/Go_GamlerDefi/module/referral/repository"
)

type ReferralRepositoryTestSuite struct {
	suite.Suite
	db         *mongo.Database
	repository *referral_repository.ReferralRepository
	logger     *logger.Logger
}

func (s *ReferralRepositoryTestSuite) SetupSuite() {
	s.logger = logger.GetLogger()
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://root:example@localhost:27017"))
	require.NoError(s.T(), err, "Failed to connect to MongoDB")
	s.db = client.Database("test_referral")
	s.repository = referral_repository.NewReferralRepository(client, s.logger)
}

func (s *ReferralRepositoryTestSuite) TearDownSuite() {
	err := s.db.Drop(context.Background())
	assert.NoError(s.T(), err, "Failed to drop test database")
}

func (s *ReferralRepositoryTestSuite) TestCreatePaymentOrder() {
	order := referral_model.PaymentOrder{
		AuthorID:    1,
		ReferrerID:  2,
		ReferralID:  3,
		TotalAmount: 100,
		TicketCount: 10,
		CreatedAt:   time.Now(),
		Levels: []referral_model.Level{
			{LevelNumber: 1, Rate: 0.1, Amount: 10},
		},
	}

	err := s.repository.CreatePaymentOrder(context.Background(), order)
	assert.NoError(s.T(), err, "Failed to create payment order")
}

func (s *ReferralRepositoryTestSuite) TestGetPaymentOrdersByAuthorID() {
	order := referral_model.PaymentOrder{
		AuthorID:    1,
		ReferrerID:  2,
		ReferralID:  3,
		TotalAmount: 100,
		TicketCount: 10,
		CreatedAt:   time.Now(),
		Levels: []referral_model.Level{
			{LevelNumber: 1, Rate: 0.1, Amount: 10},
		},
	}

	err := s.repository.CreatePaymentOrder(context.Background(), order)
	require.NoError(s.T(), err, "Failed to create payment order")

	orders, err := s.repository.GetPaymentOrdersByAuthorID(context.Background(), 1)
	assert.NoError(s.T(), err, "Failed to get payment orders")
	assert.NotEmpty(s.T(), orders, "No payment orders found")
}

func (s *ReferralRepositoryTestSuite) TestGetPaymentOrdersByAuthorID_Empty() {
	orders, err := s.repository.GetPaymentOrdersByAuthorID(context.Background(), 999)
	assert.NoError(s.T(), err, "Failed to get payment orders")
	assert.Empty(s.T(), orders, "Expected no payment orders")
}

func TestReferralRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ReferralRepositoryTestSuite))
}
