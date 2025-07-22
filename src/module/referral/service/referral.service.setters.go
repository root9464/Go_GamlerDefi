package referral_service

import (
	"context"

	referral_adapters "github.com/root9464/Go_GamlerDefi/src/module/referral/adapters"
	errors "github.com/root9464/Go_GamlerDefi/src/packages/lib/error"
	"github.com/shopspring/decimal"
)

func (s *ReferralService) CalculateAuthorDebt(ctx context.Context, authorID int) (decimal.Decimal, error) {
	s.logger.Infof("calculating author debt for author ID: %d", authorID)

	s.logger.Infof("getting payment orders by author ID: %d", authorID)
	paymentOrders, err := s.referral_repository.GetPaymentOrdersByAuthorID(ctx, authorID)
	if err != nil {
		return decimal.NewFromInt(0), errors.NewError(500, "failed to get payment orders by author ID")
	}

	s.logger.Infof("converting payment orders to DTO")
	paymentOrdersDto, err := referral_adapters.CreatePaymentOrderFromModelList(paymentOrders)
	if err != nil {
		return decimal.NewFromInt(0), errors.NewError(500, "failed to convert payment orders to DTO")
	}
	s.logger.Infof("payment orders converted to DTO: %+v", paymentOrdersDto)

	s.logger.Infof("calculating total debt")
	totalDebt := decimal.NewFromInt(0)
	for _, paymentOrder := range paymentOrdersDto {
		totalDebt = totalDebt.Add(paymentOrder.TotalAmount)
	}
	s.logger.Infof("total debt: %s", totalDebt.String())
	return totalDebt, nil
}
