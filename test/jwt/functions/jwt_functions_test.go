package jwt_functions_test

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	jwt_dto "github.com/root9464/Go_GamlerDefi/src/modules/jwt/dto"
	jwt_functions "github.com/root9464/Go_GamlerDefi/src/modules/jwt/functions"
	jwt_helpers "github.com/root9464/Go_GamlerDefi/src/modules/jwt/helpers"
	"github.com/root9464/Go_GamlerDefi/src/packages/lib/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	privateKeyStr = "MHcCAQEEIPemCJai8w+gAm+3N30cyqvIuZqmudIulBf6soXQD+iooAoGCCqGSM49AwEHoUQDQgAEQ/TqYy3uYp8JyM2Yoh7PkEXZ9bF8CJk+yKaHImsB/nhBixuA4W7PEknvUWch25is4IPyiVPge6LjAUWUP+tq+w=="
	publicKeyStr  = "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEQ/TqYy3uYp8JyM2Yoh7PkEXZ9bF8CJk+yKaHImsB/nhBixuA4W7PEknvUWch25is4IPyiVPge6LjAUWUP+tq+w=="
)

type JwtFuncsTestSuite struct {
	suite.Suite
	logger     *logger.Logger
	validator  *validator.Validate
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	helpers    jwt_helpers.IJwtHelper
	jwtFuncs   *jwt_functions.JwtFuncs
}

func (s *JwtFuncsTestSuite) SetupSuite() {
	s.logger = logger.GetLogger()
	s.validator = validator.New()

	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyStr)
	require.NoError(s.T(), err, "Failed to decode private key")
	privateKey, err := x509.ParseECPrivateKey(privateKeyBytes)
	require.NoError(s.T(), err, "Failed to parse private key")
	s.privateKey = privateKey

	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyStr)
	require.NoError(s.T(), err, "Failed to decode public key")
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	require.NoError(s.T(), err, "Failed to parse public key")
	s.publicKey = publicKey.(*ecdsa.PublicKey)

	s.helpers = jwt_helpers.NewJwtHelper(s.logger, s.validator)
	s.jwtFuncs = jwt_functions.NewJwtFuncs(s.logger, s.validator, s.privateKey, s.publicKey, s.helpers).(*jwt_functions.JwtFuncs)
}

func (s *JwtFuncsTestSuite) TestGenerateKeyPair_Success() {
	userData := jwt_dto.UserData{
		ID: 5187512201,
	}
	accessToken, refreshToken, err := s.jwtFuncs.GenerateKeyPair(userData)
	assert.NoError(s.T(), err, "Failed to generate key pair")
	assert.NotNil(s.T(), accessToken, "Access token should not be nil")
	assert.NotNil(s.T(), refreshToken, "Refresh token should not be nil")

	s.logger.Infof("Access token: %s", *accessToken)
	s.logger.Infof("Refresh token: %s", *refreshToken)

	parsedAccessToken, err := s.helpers.VerifyJwt(*accessToken, s.publicKey)
	assert.NoError(s.T(), err, "Failed to verify access token")
	assert.True(s.T(), parsedAccessToken.Valid, "Access token should be valid")

	parsedRefreshToken, err := s.helpers.VerifyJwt(*refreshToken, s.publicKey)
	assert.NoError(s.T(), err, "Failed to verify refresh token")
	assert.True(s.T(), parsedRefreshToken.Valid, "Refresh token should be valid")

	s.logger.Infof("Parsed access token: %v", parsedAccessToken)
	s.logger.Infof("Parsed refresh token: %v", parsedRefreshToken)

	accessClaims, ok := parsedAccessToken.Claims.(jwt.MapClaims)
	assert.True(s.T(), ok, "Access token claims should be valid")
	assert.Equal(s.T(), float64(5187512201), accessClaims["sub"], "Access token subject should match user ID")
	assert.Equal(s.T(), "GamlerDefi::admin", accessClaims["iss"], "Access token issuer should match")

	refreshClaims, ok := parsedRefreshToken.Claims.(jwt.MapClaims)
	assert.True(s.T(), ok, "Refresh token claims should be valid")
	assert.Equal(s.T(), float64(5187512201), refreshClaims["sub"], "Refresh token subject should match user ID")
	assert.Equal(s.T(), "GamlerDefi::admin", refreshClaims["iss"], "Refresh token issuer should match")

	s.logger.Infof("Access claims: %v", accessClaims)
	s.logger.Infof("Refresh claims: %v", refreshClaims)
}

func (s *JwtFuncsTestSuite) TestGenerateAdminToken_Success() {
	userData := jwt_dto.UserData{
		ID: 5187512201,
	}
	adminToken, err := s.jwtFuncs.GenerateAdminToken(userData)
	assert.NoError(s.T(), err, "Failed to generate admin token")
	assert.NotNil(s.T(), adminToken, "Admin token should not be nil")

	s.logger.Infof("Admin token: %s", adminToken)

	parsedAccessToken, err := s.helpers.VerifyJwt(adminToken, s.publicKey)
	assert.NoError(s.T(), err, "Failed to verify admin token")
	assert.True(s.T(), parsedAccessToken.Valid, "Admin token should be valid")

	s.logger.Infof("Parsed admin token: %v", parsedAccessToken)

	adminClaims, ok := parsedAccessToken.Claims.(jwt.MapClaims)
	assert.True(s.T(), ok, "Admin token claims should be valid")
	assert.Equal(s.T(), float64(5187512201), adminClaims["sub"], "Admin token subject should match user ID")
	assert.Equal(s.T(), "GamlerDefi::admin", adminClaims["iss"], "Admin token issuer should match")
	assert.Equal(s.T(), "admin", adminClaims["role"], "Admin token role should be admin")

	s.logger.Infof("Admin claims: %v", adminClaims)
}

func (s *JwtFuncsTestSuite) TestGenerateKeyPair_InvalidUserData() {
	userData := jwt_dto.UserData{
		ID: 0,
	}
	accessToken, refreshToken, err := s.jwtFuncs.GenerateKeyPair(userData)
	assert.Error(s.T(), err, "Expected error for invalid user data")
	assert.IsType(s.T(), &fiber.Error{}, err, "Error should be of type fiber.Error")
	assert.Equal(s.T(), fiber.StatusBadRequest, err.(*fiber.Error).Code, "Expected bad request status")
	assert.Nil(s.T(), accessToken, "Access token should be nil")
	assert.Nil(s.T(), refreshToken, "Refresh token should be nil")
}

func (s *JwtFuncsTestSuite) TestGenerateKeyPair_NilPrivateKey() {
	jwtFuncs := jwt_functions.NewJwtFuncs(s.logger, s.validator, nil, s.publicKey, s.helpers).(*jwt_functions.JwtFuncs)
	userData := jwt_dto.UserData{
		ID: 123,
	}
	accessToken, refreshToken, err := jwtFuncs.GenerateKeyPair(userData)
	assert.Error(s.T(), err, "Expected error for nil private key")
	assert.IsType(s.T(), &fiber.Error{}, err, "Error should be of type fiber.Error")
	assert.Equal(s.T(), 500, err.(*fiber.Error).Code, "Expected internal server error status")
	assert.Nil(s.T(), accessToken, "Access token should be nil")
	assert.Nil(s.T(), refreshToken, "Refresh token should be nil")
}

func (s *JwtFuncsTestSuite) TestGenerateKeyPair_NilHelpers() {
	jwtFuncs := jwt_functions.NewJwtFuncs(s.logger, s.validator, s.privateKey, s.publicKey, nil).(*jwt_functions.JwtFuncs)
	userData := jwt_dto.UserData{
		ID: 123,
	}
	accessToken, refreshToken, err := jwtFuncs.GenerateKeyPair(userData)
	assert.Error(s.T(), err, "Expected error for nil helpers")
	assert.IsType(s.T(), &fiber.Error{}, err, "Error should be of type fiber.Error")
	assert.Equal(s.T(), 500, err.(*fiber.Error).Code, "Expected internal server error status")
	assert.Nil(s.T(), accessToken, "Access token should be nil")
	assert.Nil(s.T(), refreshToken, "Refresh token should be nil")
}

func (s *JwtFuncsTestSuite) TestRefreshAccessToken_Success() {
	userData := jwt_dto.UserData{
		ID: 123,
	}
	_, refreshToken, err := s.jwtFuncs.GenerateKeyPair(userData)
	require.NoError(s.T(), err, "Failed to generate refresh token")

	accessToken, err := s.jwtFuncs.RefreshAccessToken(*refreshToken, s.publicKey, s.privateKey)
	assert.NoError(s.T(), err, "Failed to refresh access token")
	assert.NotNil(s.T(), accessToken, "Access token should not be nil")

	parsedAccessToken, err := s.helpers.VerifyJwt(*accessToken, s.publicKey)
	assert.NoError(s.T(), err, "Failed to verify new access token")
	assert.True(s.T(), parsedAccessToken.Valid, "New access token should be valid")

	accessClaims, ok := parsedAccessToken.Claims.(jwt.MapClaims)
	assert.True(s.T(), ok, "Access token claims should be valid")
	assert.Equal(s.T(), float64(123), accessClaims["sub"], "Access token subject should match user ID")
	assert.Equal(s.T(), "GamlerDefi::admin", accessClaims["iss"], "Access token issuer should match")
	assert.Equal(s.T(), "admin", accessClaims["role"], "Access token role should be admin")
}

func (s *JwtFuncsTestSuite) TestRefreshAccessToken_InvalidToken() {
	invalidToken := "invalid.token.string"
	accessToken, err := s.jwtFuncs.RefreshAccessToken(invalidToken, s.publicKey, s.privateKey)
	assert.Error(s.T(), err, "Expected error for invalid refresh token")
	assert.Contains(s.T(), err.Error(), "invalid refresh token", "Error message should indicate invalid token")
	assert.Nil(s.T(), accessToken, "Access token should be nil")
}

func (s *JwtFuncsTestSuite) TestRefreshAccessToken_ExpiredToken() {
	userData := jwt_dto.UserData{
		ID: 123,
	}
	expiredClaims := jwt.MapClaims{
		"iss":       "GamlerDefi::admin",
		"sub":       userData.ID,
		"iat":       time.Now().Add(-25 * time.Hour).Unix(),
		"exp":       time.Now().Add(-24 * time.Hour).Unix(),
		"user_hash": "somehash",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, expiredClaims)
	expiredRefreshToken, err := token.SignedString(s.privateKey)
	require.NoError(s.T(), err, "Failed to create expired refresh token")

	accessToken, err := s.jwtFuncs.RefreshAccessToken(expiredRefreshToken, s.publicKey, s.privateKey)
	assert.Error(s.T(), err, "Expected error for expired refresh token")
	assert.Contains(s.T(), err.Error(), "refresh token has expired", "Error message should indicate expired token")
	assert.Nil(s.T(), accessToken, "Access token should be nil")
}

func (s *JwtFuncsTestSuite) TestRefreshAccessToken_MissingClaims() {
	invalidClaims := jwt.MapClaims{
		"iss": "GamlerDefi::admin",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, invalidClaims)
	invalidRefreshToken, err := token.SignedString(s.privateKey)
	require.NoError(s.T(), err, "Failed to create invalid refresh token")

	accessToken, err := s.jwtFuncs.RefreshAccessToken(invalidRefreshToken, s.publicKey, s.privateKey)
	assert.Error(s.T(), err, "Expected error for missing claims")
	assert.Contains(s.T(), err.Error(), "refresh token missing user ID", "Error message should indicate missing user ID")
	assert.Nil(s.T(), accessToken, "Access token should be nil")
}

func TestJwtFuncsTestSuite(t *testing.T) {
	suite.Run(t, new(JwtFuncsTestSuite))
}
