package admin_middleware

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"
	jwt_helpers "github.com/root9464/Go_GamlerDefi/src/submodules/jwt/helpers"
)

type Middleware struct {
	logger     *logger.Logger
	jwtHelpers jwt_helpers.IJwtHelper
	publicKey  string
}

func NewMiddleware(
	logger *logger.Logger,
	jwtHelpers jwt_helpers.IJwtHelper,
	publicKey string,
) *Middleware {
	return &Middleware{
		logger:     logger,
		jwtHelpers: jwtHelpers,
		publicKey:  publicKey,
	}
}

func (m *Middleware) AdminOnly() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		tokenString := ctx.Get("Authorization")
		if tokenString == "" {
			m.logger.Warn("Missing Authorization header")
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing Authorization header",
			})
		}
		m.logger.Infof("tokenString: %s", tokenString)

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		publicKeyBytes, err := base64.StdEncoding.DecodeString(m.publicKey)
		if err != nil {
			m.logger.Warnf("Failed to decode base64 public key: %s", err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid public key",
			})
		}
		m.logger.Infof("publicKey: %s", m.publicKey)

		parsedKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
		if err != nil {
			m.logger.Warnf("Failed to parse public key: %s", err.Error())
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid public key",
			})
		}
		m.logger.Infof("parsedKey: %v", parsedKey)
		ecdsaPublicKey, ok := parsedKey.(*ecdsa.PublicKey)
		if !ok {
			m.logger.Warn("Public key is not ECDSA")
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Invalid public key type",
			})
		}

		m.logger.Infof("ecdsaPublicKey: %v", ecdsaPublicKey)
		payload, err := m.jwtHelpers.ParseJwt(tokenString, ecdsaPublicKey)
		if err != nil {
			m.logger.Warnf("Invalid JWT token: %s", err.Error())
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid JWT token",
			})
		}
		m.logger.Infof("payload: %v", payload)
		token, err := m.jwtHelpers.VerifyJwt(tokenString, ecdsaPublicKey)
		if err != nil {
			m.logger.Warnf("Token verification failed: %s", err.Error())
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid JWT token",
			})
		}
		m.logger.Infof("token: %v", token)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["role"] != "admin" {
			m.logger.Warn("Token does not have admin role")
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Admin access required",
			})
		}
		m.logger.Infof("claims: %v", claims)
		isValid, err := m.jwtHelpers.CheckTokenExpiration(tokenString, ecdsaPublicKey)
		if err != nil {
			m.logger.Warnf("Token expiration check failed: %s", err.Error())
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		if !isValid {
			m.logger.Warn("Token has expired")
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token has expired",
			})
		}
		m.logger.Infof("isValid: %v", isValid)
		ctx.Locals("user", payload)
		return ctx.Next()
	}
}
