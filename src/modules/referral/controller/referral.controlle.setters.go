package referral_controller

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	referral_dto "github.com/root9464/Go_GamlerDefi/src/modules/referral/dto"
	errors "github.com/root9464/Go_GamlerDefi/src/packages/lib/error"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ReferralProcessPlatform handles referral bonus calculation
// @Summary Process of awarding referral bonuses to the user
// @Description Calculate and distribute referral bonuses between users
// @Tags Referrals
// @Accept json
// @Produce json
// @Param request body referral_dto.ReferralProcessRequest true "Referral processing data"
// @Success 200 {object} fiber.Map "Success response"
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

	err := c.referral_service.ReferralProcess(ctx.Context(), dto)
	if err != nil {
		c.logger.Errorf("error calculating referral bonuses: %v", err)
		return errors.NewError(500, err.Error())
	}

	return ctx.Status(201).JSON(fiber.Map{
		"message": "Referral process completed successfully",
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
	paramOrderID := ctx.Query("order_id")
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

// @Summary Delete all payment orders
// @Description Delete all payment orders
// @Tags Referrals
// @Accept json
// @Produce json
// @Success 200 {object} fiber.Map "Success response"
// @Failure 500 {object} fiber.Map "Internal server error"
// @Router /api/referrals/delete/all [delete]
func (c *ReferralController) DeleteAllPaymentOrders(ctx *fiber.Ctx) error {
	paramAuthorID := ctx.Query("author_id")
	c.logger.Infof("author ID: %s", paramAuthorID)
	if paramAuthorID == "" {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Author ID is required",
		})
	}

	authorID, err := strconv.Atoi(paramAuthorID)
	if err != nil {
		c.logger.Errorf("error converting author ID to int: %v", err)
		return errors.NewError(400, err.Error())
	}

	err = c.referral_repository.DeleteAllPaymentOrders(ctx.Context(), authorID)
	if err != nil {
		c.logger.Errorf("error deleting all payment orders: %v", err)
		return errors.NewError(500, err.Error())
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "All payment orders deleted successfully",
	})
}

func (c *ReferralController) AddTrHashToPaymentOrder(ctx *fiber.Ctx) error {
	var dto referral_dto.AddTrHashToPaymentOrderRequest
	if err := ctx.BodyParser(&dto); err != nil {
		c.logger.Errorf("error parsing request body: %v", err)
		return errors.NewError(400, err.Error())
	}

	orderID, err := bson.ObjectIDFromHex(dto.OrderID)
	if err != nil {
		c.logger.Fatalf("Invalid ObjectID string: %v", err)
	}

	err = c.referral_repository.AddTrHashToPaymentOrder(ctx.Context(), orderID, dto.TrHash)
	if err != nil {
		c.logger.Errorf("error adding tr hash to payment order: %v", err)
		return errors.NewError(500, err.Error())
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Tr hash added to payment order successfully",
	})
}
