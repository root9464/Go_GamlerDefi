package referral_service

import (
	"context"

	referral_adapters "github.com/root9464/Go_GamlerDefi/module/referral/adapters"
	errors "github.com/root9464/Go_GamlerDefi/packages/lib/error"
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
	paymentOrderDTO := referral_adapters.CreatePaymentOrderFromModelList(ctx, paymentOrders)

	s.logger.Infof("converted payment order to DTO: %+v", paymentOrderDTO)

	s.logger.Infof("calculating total amount of payment orders")
	totalAmount := 0.0
	for _, paymentOrder := range paymentOrderDTO {
		totalAmount += paymentOrder.TotalAmount
	}

	s.logger.Infof("total amount of payment orders: %f", totalAmount)
	if totalAmount > maxDebt {
		s.logger.Infof("insufficient funds on the balance sheet to pay the debt: %f", totalAmount)
		return false, errors.NewError(402, "insufficient funds on the balance sheet to pay the debt")
	}

	s.logger.Infof("the invitation has been approved")
	return true, nil
}
