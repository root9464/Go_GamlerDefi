package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/root9464/Go_GamlerDefi/database"
	"github.com/root9464/Go_GamlerDefi/packages/lib/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	referral_model "github.com/root9464/Go_GamlerDefi/module/referral/model"
	referral_repository "github.com/root9464/Go_GamlerDefi/module/referral/repository"
)

const (
	db_url  = "mongodb://root:example@localhost:27017"
	db_name = "gamer_defi_test"
)

type ReferralRepositoryTestSuite struct {
	suite.Suite
	db         *mongo.Database
	repository *referral_repository.ReferralRepository
	logger     *logger.Logger
}

func (s *ReferralRepositoryTestSuite) SetupSuite() {
	s.logger = logger.GetLogger()
	_, database, err := database.ConnectDatabase(db_url, s.logger, db_name)
	require.NoError(s.T(), err, "Failed to connect to database")
	s.db = database
	repo, ok := referral_repository.NewReferralRepository(s.logger, s.db).(*referral_repository.ReferralRepository)
	require.True(s.T(), ok, "Failed to type assert repository")
	s.repository = repo
}

func (s *ReferralRepositoryTestSuite) TestCreatePaymentOrder() {
	order := referral_model.PaymentOrder{
		AuthorID:    1,
		ReferrerID:  1,
		ReferralID:  2,
		TotalAmount: 150,
		TicketCount: 750,
		CreatedAt:   time.Now().Unix(),
		Levels: []referral_model.Level{
			{LevelNumber: 0, Rate: 0.2, Amount: 150, Address: "0QC3PUCoxBdLfOmO8xFQ84TGFPQUatxvvRsSAODKEvjbb4OS"},
		},
	}

	err := s.repository.CreatePaymentOrder(context.Background(), order)
	assert.NoError(s.T(), err, "Failed to create payment order")
}

func (s *ReferralRepositoryTestSuite) TestGetPaymentOrdersByAuthorID() {
	orders, err := s.repository.GetPaymentOrdersByAuthorID(context.Background(), 1)
	s.logger.Infof("Payment orders: %+v", orders)
	assert.NoError(s.T(), err, "Failed to get payment orders")
}

func (s *ReferralRepositoryTestSuite) TestGetPaymentOrdersByAuthorID_Empty() {
	orders, err := s.repository.GetPaymentOrdersByAuthorID(context.Background(), 999)
	assert.NoError(s.T(), err, "Failed to get payment orders")
	assert.Empty(s.T(), orders, "Expected no payment orders")
}

func (s *ReferralRepositoryTestSuite) TestGetAllPaymentOrders() {
	orders, err := s.repository.GetAllPaymentOrders(context.Background())
	assert.NoError(s.T(), err, "Failed to get all payment orders")
	assert.NotEmpty(s.T(), orders, "No payment orders found")
}

func (s *ReferralRepositoryTestSuite) TestDeletePaymentOrder() {
	orderIDStr := "682092314a00a7558247b21f"
	orderID, err := bson.ObjectIDFromHex(orderIDStr)
	if err != nil {
		s.logger.Fatalf("Invalid ObjectID string: %v", err)
	}

	err = s.repository.DeletePaymentOrder(context.Background(), orderID)
	assert.NoError(s.T(), err, "Failed to delete payment order")
}

func (s *ReferralRepositoryTestSuite) TestGetDebtFromAuthorToReferrer() {
	orders, err := s.repository.GetDebtFromAuthorToReferrer(context.Background(), 3, 1)
	s.logger.Infof("Debt from author to referrer: %+v", orders)
	assert.NoError(s.T(), err, "Failed to get debt from author to referrer")
	assert.NotEmpty(s.T(), orders, "No debt found")
}

func (s *ReferralRepositoryTestSuite) TestGetPaymentOrderByID() {
	orderIDStr := "6823dc5bcb80d8ea88f9b32b"
	orderID, err := bson.ObjectIDFromHex(orderIDStr)
	if err != nil {
		s.logger.Fatalf("Invalid ObjectID string: %v", err)
	}

	order, err := s.repository.GetPaymentOrderByID(context.Background(), orderID)
	assert.NoError(s.T(), err, "Failed to get payment order by ID")
	assert.NotNil(s.T(), order, "Payment order not found")
}

func TestReferralRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ReferralRepositoryTestSuite))
}

// func (s *ReferralRepositoryTestSuite) TearDownSuite() {
// 	err := s.db.Drop(context.Background())
// 	assert.NoError(s.T(), err, "Failed to drop test database")
// }
