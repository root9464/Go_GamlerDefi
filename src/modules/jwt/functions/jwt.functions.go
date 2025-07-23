package jwt_functions

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	jwt_dto "github.com/root9464/Go_GamlerDefi/src/modules/jwt/dto"
)

const (
	AccessTokenExpiry  = 15 * time.Minute
	RefreshTokenExpiry = 24 * time.Hour
	Issuer             = "GamlerDefi::admin"
)

var (
	Now = time.Now().Unix()
)

func (f *JwtFuncs) GenerateKeyPair(userData jwt_dto.UserData) (*string, *string, error) {
	f.logger.Info("Validating user data...")
	if err := f.validator.Struct(userData); err != nil {
		f.logger.Warnf("validate error: %s", err.Error())
		return nil, nil, &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		}
	}
	f.logger.Info("Validating success...")

	userRaw := fmt.Sprintf("%d", userData.ID)

	hash := sha256.New()
	hash.Write([]byte(userRaw))
	refinedHash := hex.EncodeToString(hash.Sum(nil))

	accessClaims := jwt.MapClaims{
		"iss":       Issuer,
		"sub":       userData.ID,
		"iat":       Now,
		"exp":       time.Now().Add(AccessTokenExpiry).Unix(),
		"user_hash": refinedHash,
	}

	f.logger.Infof("access claims: %+v", accessClaims)

	refreshClaims := jwt.MapClaims{
		"iss":       Issuer,
		"sub":       userData.ID,
		"iat":       Now,
		"exp":       time.Now().Add(RefreshTokenExpiry).Unix(),
		"user_hash": refinedHash,
	}

	f.logger.Infof("refresh claims: %+v", refreshClaims)

	if f.privateKey == nil {
		return nil, nil, &fiber.Error{
			Code:    500,
			Message: "privateKey is not initialized",
		}
	}
	if f.helpers == nil {
		return nil, nil, &fiber.Error{
			Code:    500,
			Message: "helpers is not initialized",
		}
	}

	accessToken, err := f.helpers.CreateJwt(accessClaims, f.privateKey)
	if err != nil {
		f.logger.Warnf("create access token error: %s", err.Error())
		return nil, nil, &fiber.Error{
			Code:    500,
			Message: err.Error(),
		}
	}

	refreshToken, err := f.helpers.CreateJwt(refreshClaims, f.privateKey)
	if err != nil {
		f.logger.Warnf("create refresh token error: %s", err.Error())
		return nil, nil, &fiber.Error{
			Code:    500,
			Message: err.Error(),
		}
	}

	return accessToken, refreshToken, nil
}

func (f *JwtFuncs) RefreshAccessToken(refreshToken string, publicKey *ecdsa.PublicKey, privateKey *ecdsa.PrivateKey) (*string, error) {
	parsedToken, err := jwt.ParseWithClaims(refreshToken, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		f.logger.Warnf("invalid refresh token: %s", err.Error())
		return nil, fmt.Errorf("invalid refresh token: %s", err.Error())
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		f.logger.Warn("invalid claims in refresh token")
		return nil, fmt.Errorf("invalid claims in refresh token")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		f.logger.Warn("refresh token missing expiration")
		return nil, fmt.Errorf("refresh token missing expiration")
	}

	expTime := time.Unix(int64(exp), 0)
	if time.Now().After(expTime) {
		f.logger.Warn("refresh token has expired")
		return nil, fmt.Errorf("refresh token has expired")
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		f.logger.Warn("refresh token missing user ID")
		return nil, fmt.Errorf("refresh token missing user ID")
	}

	userHash, ok := claims["user_hash"].(string)
	if !ok {
		f.logger.Warn("refresh token missing user hash")
		return nil, fmt.Errorf("refresh token missing user hash")
	}

	accessClaims := jwt.MapClaims{
		"iss":       Issuer,
		"sub":       int64(userID),
		"iat":       Now,
		"exp":       time.Now().Add(AccessTokenExpiry).Unix(),
		"user_hash": userHash,
	}

	accessToken, err := f.helpers.CreateJwt(accessClaims, privateKey)
	if err != nil {
		f.logger.Warnf("create access token error: %s", err.Error())
		return nil, &fiber.Error{
			Code:    500,
			Message: err.Error(),
		}
	}

	return accessToken, nil
}

func (f *JwtFuncs) GenerateAdminToken(userID jwt_dto.UserData) (string, error) {
	if f.privateKey == nil {
		return "", fmt.Errorf("private key is not initialized")
	}

	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%d", userID.ID)))
	refinedHash := hex.EncodeToString(hash.Sum(nil))

	adminClaims := jwt.MapClaims{
		"iss":       Issuer,
		"sub":       userID.ID,
		"iat":       time.Now().Unix(),
		"exp":       time.Now().Add(100 * 365 * 24 * time.Hour).Unix(),
		"role":      "admin",
		"user_hash": refinedHash,
	}

	token, err := f.helpers.CreateJwt(adminClaims, f.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to create admin token: %w", err)
	}

	return *token, nil
}
