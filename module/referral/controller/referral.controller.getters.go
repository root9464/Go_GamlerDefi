package referral_controller

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	referral_adapters "github.com/root9464/Go_GamlerDefi/module/referral/adapters"
	referral_dto "github.com/root9464/Go_GamlerDefi/module/referral/dto"
	errors "github.com/root9464/Go_GamlerDefi/packages/lib/error"
	"github.com/root9464/Go_GamlerDefi/packages/utils"
)

const (
	url = "https://serv.gamler.atma-dev.ru/referral"
)

// @Summary checking the referrer
// @Description Сheck if the user is a referrer
// @Tags Referrals
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
// @Tags Referrals
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
	authorOrdersDTO, err := referral_adapters.CreatePaymentOrderFromModelList(authorOrders)
	if err != nil {
		c.logger.Errorf("error converting author orders to DTO: %v", err)
		return errors.NewError(500, err.Error())
	}

	return ctx.Status(200).JSON(authorOrdersDTO)
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
// @Router /api/referral/payment-orders/pay [get]
func (c *ReferralController) PayDebtAuthor(ctx *fiber.Ctx) error {
	paramOrderID := ctx.Query("order_id")
	walletAddress := ctx.Get("Wallet-Address")
	c.logger.Infof("order ID: %s", paramOrderID)
	c.logger.Infof("wallet address: %s", walletAddress)

	if paramOrderID == "" || walletAddress == "" {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "Wallet address is required",
		})
	}

	cell, err := c.referral_service.PayPaymentOrder(ctx.Context(), paramOrderID, walletAddress)
	if err != nil {
		c.logger.Errorf("error paying payment order: %v", err)
		return errors.NewError(500, err.Error())
	}

	return ctx.Status(200).JSON(referral_dto.CellResponse{
		Cell: cell,
	})
}

// @Summary Pay all debt from author to referrer
// @Description Paying all debt from author to referrer
// @Tags Referrals
// @Accept json
// @Produce json
// @Param author_id query int true "Author ID"
// @Success 200 {object} referral_dto.CellResponse
// @Failure 400 {object} errors.MapError
// @Failure 404 {object} errors.MapError
// @Failure 500 {object} errors.MapError
// @Router /api/referral/payment-orders/all [get]
func (c *ReferralController) PayAllDebtAuthor(ctx *fiber.Ctx) error {
	paramAuthorID := ctx.Query("author_id")
	walletAddress := ctx.Get("Wallet-Address")

	c.logger.Infof("author ID: %s", paramAuthorID)
	c.logger.Infof("wallet address: %s", walletAddress)

	if paramAuthorID == "" || walletAddress == "" {
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
	cell, err := c.referral_service.PayAllPaymentOrders(ctx.Context(), authorID, walletAddress)
	if err != nil {
		c.logger.Errorf("error paying all payment orders: %v", err)
		return errors.NewError(500, err.Error())
	}

	return ctx.Status(200).JSON(referral_dto.CellResponse{
		Cell: cell,
	})
}

// @Summary Validate invitation conditions
// @Description Validate invitation conditions
// @Tags Referrals
// @Accept json
// @Produce json
// @Param author_id query int true "Author ID"
// @Success 200 {object} referral_dto.CellResponse
// @Failure 400 {object} errors.MapError
// @Failure 404 {object} errors.MapError
// @Failure 500 {object} errors.MapError
// @Router /api/referral/validate-invitation-conditions [get]
func (c *ReferralController) ValidateInvitationConditions(ctx *fiber.Ctx) error {
	paramAuthorID := ctx.Query("author_id")
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
	valid, err := c.referral_service.AssessInvitationAbility(ctx.Context(), authorID)
	if err != nil {
		c.logger.Errorf("error assessing invitation ability: %v", err)
		return errors.NewError(500, err.Error())
	}

	return ctx.Status(200).JSON(referral_dto.ValidateInvitationConditionsResponse{
		Valid: valid,
	})
}

func (c *ReferralController) GetCalculateAuthorDebt(ctx *fiber.Ctx) error {
	paramAuthorID := ctx.Query("author_id")
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

	debt, err := c.referral_service.CalculateAuthorDebt(ctx.Context(), authorID)
	if err != nil {
		c.logger.Errorf("error calculating author debt: %v", err)
		return errors.NewError(500, err.Error())
	}

	return ctx.Status(200).JSON(debt)
}
