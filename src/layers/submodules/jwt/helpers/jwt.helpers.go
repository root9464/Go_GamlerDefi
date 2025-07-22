package jwt_helpers

import (
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	jwt_dto "github.com/root9464/Go_GamlerDefi/src/layers/submodules/jwt/dto"
)

func (h *jwtHelper) CreateJwt(claims jwt.Claims, key *ecdsa.PrivateKey) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	signedToken, err := token.SignedString(key)
	if err != nil {
		return nil, err
	}
	return &signedToken, nil
}

func (h *jwtHelper) VerifyJwt(tokenString string, key *ecdsa.PublicKey) (*jwt.Token, error) {
	h.logger.Info("Verifying JWT token...")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	h.logger.Info("JWT token verified successfully")
	return token, nil
}

func (h *jwtHelper) CheckTokenExpiration(tokenString string, publicKey *ecdsa.PublicKey) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return false, fmt.Errorf("failed to parse token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		exp, ok := claims["exp"].(float64)
		if !ok {
			return false, fmt.Errorf("exp field not found in token")
		}

		expirationTime := time.Unix(int64(exp), 0)
		if expirationTime.Before(time.Now()) {
			return false, nil
		}
		return true, nil
	}

	return false, fmt.Errorf("invalid token claims")
}

func (h *jwtHelper) ParseJwt(tokenString string, key *ecdsa.PublicKey) (*jwt_dto.UserJwtPayload, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return key, nil
	})

	h.logger.Info("Parsing JWT token...")
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return &jwt_dto.UserJwtPayload{
			Iss:  claims["iss"].(string),
			Sub:  int64(claims["sub"].(float64)),
			Iat:  int64(claims["iat"].(float64)),
			Exp:  int64(claims["exp"].(float64)),
			Hash: claims["user_hash"].(string),
		}, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}
