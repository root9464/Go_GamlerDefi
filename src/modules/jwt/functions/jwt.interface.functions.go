package jwt_functions

import (
	"crypto/ecdsa"

	"github.com/go-playground/validator/v10"
	jwt_dto "github.com/root9464/Go_GamlerDefi/src/modules/jwt/dto"
	jwt_helpers "github.com/root9464/Go_GamlerDefi/src/modules/jwt/helpers"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

var _ IJwtFuncs = (*JwtFuncs)(nil)

type IJwtFuncs interface {
	GenerateKeyPair(userData jwt_dto.UserData) (*string, *string, error)
	RefreshAccessToken(refreshToken string, publicKey *ecdsa.PublicKey, privateKey *ecdsa.PrivateKey) (*string, error)
	GenerateAdminToken(userID jwt_dto.UserData) (string, error)
}

type JwtFuncs struct {
	logger    *logger.Logger
	validator *validator.Validate

	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey

	helpers jwt_helpers.IJwtHelper
}

func NewJwtFuncs(logger *logger.Logger, validator *validator.Validate, privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey, helpers jwt_helpers.IJwtHelper) IJwtFuncs {
	return &JwtFuncs{
		logger:     logger,
		validator:  validator,
		privateKey: privateKey,
		publicKey:  publicKey,
		helpers:    helpers,
	}
}
