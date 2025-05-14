package referral_controller

import (
	"github.com/gofiber/fiber/v2"
	referral_dto "github.com/root9464/Go_GamlerDefi/module/referral/dto"
	errors "github.com/root9464/Go_GamlerDefi/packages/lib/error"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ReferralProcessPlatform handles referral bonus calculation
// @Summary Process of awarding referral bonuses to the user
// @Description Calculate and distribute referral bonuses between users
// @Tags Referrals
// @Accept json
// @Produce json
// @Param request body referral_dto.ReferralProcessRequest true "Referral processing data"
// @Success 200 {object} referral_dto.CellResponse "Success response"
// @Failure 400 {object} errors.MapError "Validation error"
// @Failure 500 {object} errors.MapError "Internal server error"
// @Router /api/referrals/process [post]
func (c *ReferralController) ReferralProcessPlatform(ctx *fiber.Ctx) error {
	var dto referral_dto.ReferralProcessRequest
	if err := ctx.BodyParser(&dto); err != nil {
		c.logger.Errorf("error parsing request body: %v", err)
		return errors.NewError(400, err.Error())
	}
	if err := c.validator.Struct(dto); err != nil {
		c.logger.Errorf("validation error: %s", err.Error())
		return errors.NewError(400, err.Error())
	}

	c.logger.Infof("processing referral for referrer ID: %d", dto.ReferrerID)

	cell, err := c.referral_service.ReferralProcess(ctx.Context(), dto)
	if err != nil {
		c.logger.Errorf("error calculating referral bonuses: %v", err)
		return errors.NewError(500, err.Error())
	}

	return ctx.Status(200).JSON(referral_dto.CellResponse{
		Cell: cell,
	})
}

// @Summary Delete payment order
// @Description Delete a payment order by ID
// @Tags Referrals
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Success 200 {object} fiber.Map "Success response"
// @Failure 400 {object} errors.MapError "Validation error"
// @Failure 500 {object} errors.MapError "Internal server error"
// @Router /api/referrals/delete/{order_id} [delete]
func (c *ReferralController) DeletePaymentOrder(ctx *fiber.Ctx) error {
	paramOrderID := ctx.Params("order_id")
	c.logger.Infof("order ID: %s", paramOrderID)

	orderID, err := bson.ObjectIDFromHex(paramOrderID)
	if err != nil {
		c.logger.Fatalf("Invalid ObjectID string: %v", err)
	}

	err = c.referral_repository.DeletePaymentOrder(ctx.Context(), orderID)
	if err != nil {
		c.logger.Errorf("error deleting payment order: %v", err)
		return errors.NewError(500, err.Error())
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Payment order deleted successfully",
	})
}

// @Summary Pay payment order
// @Description Pay a payment order by ID
// @Tags Referrals
// @Accept json
// @Produce json
// @Param order_id path string true "Order ID"
// @Success 200 {object} referral_dto.CellResponse "Success response"
// @Failure 400 {object} errors.MapError "Validation error"
// @Failure 500 {object} errors.MapError "Internal server error"
// @Router /api/referrals/pay/{order_id} [post]
func (c *ReferralController) PayDebtAuthor(ctx *fiber.Ctx) error {
	paramOrderID := ctx.Query("order_id")
	c.logger.Infof("order ID: %s", paramOrderID)

	cell, err := c.referral_service.PayPaymentOrder(ctx.Context(), paramOrderID)
	if err != nil {
		c.logger.Errorf("error paying payment order: %v", err)
		return errors.NewError(500, err.Error())
	}

	return ctx.Status(200).JSON(referral_dto.CellResponse{
		Cell: cell,
	})
}
