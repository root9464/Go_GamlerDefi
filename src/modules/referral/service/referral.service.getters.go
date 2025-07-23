package referral_service

import (
	"context"

	referral_adapters "github.com/root9464/Go_GamlerDefi/src/modules/referral/adapters"
	errors "github.com/root9464/Go_GamlerDefi/src/packages/lib/error"
	"github.com/shopspring/decimal"
)

func (s *ReferralService) AssessInvitationAbility(ctx context.Context, authorID int) (bool, error) {
	s.logger.Infof("assessing invitation ability for author_id: %d", authorID)
	s.logger.Infof("fetching payment orders in database by author_id: %d", authorID)
	paymentOrders, err := s.referral_repository.GetPaymentOrdersByAuthorID(ctx, authorID)
	if err != nil {
		s.logger.Errorf("failed to get payment orders: %v", err)
		return false, errors.NewError(500, "failed to get payment orders")
	}

	s.logger.Infof("payment orders fetched successfully: %+v", paymentOrders)

	s.logger.Infof("converting payment order to DTO")
	paymentOrderDTO, err := referral_adapters.CreatePaymentOrderFromModelList(paymentOrders)
	if err != nil {
		s.logger.Errorf("failed to convert payment order to DTO: %v", err)
		return false, errors.NewError(500, "failed to convert payment order to DTO")
	}

	s.logger.Infof("converted payment order to DTO: %+v", paymentOrderDTO)

	s.logger.Infof("calculating total amount of payment orders")
	totalAmount := decimal.NewFromFloat(0)
	for _, paymentOrder := range paymentOrderDTO {
		totalAmount = totalAmount.Add(paymentOrder.TotalAmount)
	}

	s.logger.Infof("total amount of payment orders: %s", totalAmount.String())
	if totalAmount.GreaterThan(decimal.NewFromInt(int64(maxDebt))) {
		s.logger.Infof("insufficient funds on the balance sheet to pay the debt: %s", totalAmount.String())
		return false, errors.NewError(402, "insufficient funds on the balance sheet to pay the debt")
	}

	s.logger.Infof("the invitation has been approved")
	return true, nil
}
