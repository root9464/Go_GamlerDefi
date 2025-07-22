package jwt_module

import (
	"github.com/go-playground/validator/v10"
	jwt_functions "github.com/root9464/Go_GamlerDefi/src/layers/submodules/jwt/functions"
	jwt_helpers "github.com/root9464/Go_GamlerDefi/src/layers/submodules/jwt/helpers"
	jwt_utils "github.com/root9464/Go_GamlerDefi/src/layers/submodules/jwt/utils"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
)

type JwtModule struct {
	jwtFuncs   jwt_functions.IJwtFuncs
	jwtHelpers jwt_helpers.IJwtHelper

	privateKey string
	publicKey  string

	logger    *logger.Logger
	validator *validator.Validate
}

func (m *JwtModule) JwtHelpers() jwt_helpers.IJwtHelper {
	if m.jwtHelpers == nil {
		m.jwtHelpers = jwt_helpers.NewJwtHelper(m.logger, m.validator)
	}
	return m.jwtHelpers
}

func (m *JwtModule) JwtFuncs() jwt_functions.IJwtFuncs {
	privateKey, publicKey, err := jwt_utils.HexToKeys(m.privateKey, m.publicKey)
	if err != nil {
		panic(err)
	}

	if m.jwtFuncs == nil {
		m.jwtFuncs = jwt_functions.NewJwtFuncs(m.logger, m.validator, privateKey, publicKey, m.JwtHelpers())
	}
	return m.jwtFuncs
}

func NewJwtModule(logger *logger.Logger, validator *validator.Validate, privateKey string, publicKey string) *JwtModule {
	return &JwtModule{
		logger:     logger,
		validator:  validator,
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}
