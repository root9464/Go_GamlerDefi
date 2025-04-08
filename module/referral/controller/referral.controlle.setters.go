package referral_controller

import (
	"github.com/gofiber/fiber/v2"
	referral_dto "github.com/root9464/Go_GamlerDefi/module/referral/dto"
	errors "github.com/root9464/Go_GamlerDefi/packages/lib/error"
)

func (c *ReferralController) ReferralProcessPlatform(ctx *fiber.Ctx) error {
	var dto referral_dto.ReferralProcessRequest
	if err := ctx.BodyParser(&dto); err != nil {
		c.logger.Errorf("Error: %v", err)
		return errors.NewError(400, err.Error())
	}

	if err := c.validator.Struct(dto); err != nil {
		c.logger.Warnf("validate error: %s", err.Error())
		return errors.NewError(400, err.Error())
	}

	c.logger.Infof("Cleaned User ID: %d", dto.ReferrerID)

	err := c.referralService.CalculateReferralBonuses(ctx.Context(), dto)
	if err != nil {
		c.logger.Errorf("Error: %v", err)
		return errors.NewError(404, err.Error())
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Referral process completed successfully",
	})
}
