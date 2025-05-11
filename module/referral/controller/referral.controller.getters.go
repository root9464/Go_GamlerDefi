package referral_controller

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	referral_dto "github.com/root9464/Go_GamlerDefi/module/referral/dto"
	errors "github.com/root9464/Go_GamlerDefi/packages/lib/error"
	"github.com/root9464/Go_GamlerDefi/packages/utils"
)

const (
	url = "https://serv.gamler.atma-dev.ru/referral"
)

// @Summary checking the referrer
// @Description Сheck if the user is a referrer
// @Tags referral
// @Accept json
// @Produce json
// @Success 200 {object} referral_dto.ReferrerResponse "Success response"
// @Failure 400 {object} errors.MapError "Validation error"
// @Failure 404 {object} errors.MapError "Not found"
// @Failure 500 {object} errors.MapError "Internal server error"
// @Router /api/referral/precheckout/{user_id} [get]
func (c *ReferralController) PrecheckoutReferrer(ctx *fiber.Ctx) error {
	paramUserID := ctx.Params("user_id")
	c.logger.Infof("User ID: %s", paramUserID)

	if paramUserID == "" {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "User ID is required",
		})
	}

	userID := strings.TrimPrefix(paramUserID, "user_id=")
	c.logger.Infof("сleaned User ID: %s", userID)

	c.logger.Infof("get referral URL: %s", fmt.Sprintf("%s/referrer/%s", url, userID))
	referral, err := utils.Get[referral_dto.ReferrerResponse](fmt.Sprintf("%s/referrer/%s", url, userID))
	if err != nil {
		c.logger.Errorf("error fetching referrer data: %v", err)
		return errors.NewError(404, err.Error())
	}

	c.logger.Infof("referral: %+v", referral)

	return ctx.Status(200).JSON(fiber.Map{
		"message": "The referrer has been confirmed",
	})
}

// @Summary getting debt from author to referrer
// @Description Getting debt from author to referrer
// @Tags referral
// @Accept json
// @Produce json
// @Success 200 {object} referral_dto.PaymentOrder "Success response"
// @Failure 400 {object} errors.MapError "Validation error"
// @Failure 404 {object} errors.MapError "Not found"
// @Failure 500 {object} errors.MapError "Internal server error"
// @Router /api/referral/debt/{author_id}/{referrer_id} [get]
func (c *ReferralController) GetDebtAuthor(ctx *fiber.Ctx) error {
	paramAuthorID := ctx.Params("author_id")
	c.logger.Infof("author ID: %s", paramAuthorID)

	if paramAuthorID == "" {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Author ID is required",
		})
	}

	authorID, err := strconv.Atoi(paramAuthorID)
	if err != nil {
		c.logger.Errorf("error converting author ID: %v", err)
		return errors.NewError(400, err.Error())
	}

	c.logger.Infof("author ID to int: %d", authorID)
	authorOrders, err := c.referral_repository.GetPaymentOrdersByAuthorID(ctx.Context(), authorID)
	if err != nil {
		c.logger.Errorf("error getting author orders: %v", err)
		return errors.NewError(500, err.Error())
	}

	c.logger.Infof("author orders: %+v", authorOrders)

	return ctx.Status(200).JSON(authorOrders)
}
