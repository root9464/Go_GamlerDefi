package jwt_helpers

import (
	"crypto/ecdsa"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	jwt_dto "github.com/root9464/Go_GamlerDefi/src/layers/submodules/jwt/dto"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

var _ IJwtHelper = (*jwtHelper)(nil)

type IJwtHelper interface {
	CreateJwt(claims jwt.Claims, key *ecdsa.PrivateKey) (*string, error)
	VerifyJwt(tokenString string, key *ecdsa.PublicKey) (*jwt.Token, error)
	CheckTokenExpiration(token string, publicKey *ecdsa.PublicKey) (bool, error)
	ParseJwt(tokenString string, key *ecdsa.PublicKey) (*jwt_dto.UserJwtPayload, error)
}

type jwtHelper struct {
	logger    *logger.Logger
	validator *validator.Validate
}

func NewJwtHelper(logger *logger.Logger, validator *validator.Validate) IJwtHelper {
	return &jwtHelper{logger: logger, validator: validator}
}
