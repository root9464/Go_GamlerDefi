package jwt_utils

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"fmt"
)

func HexToKeys(privateKeyHex, publicKeyHex string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding private key hex: %w", err)
	}

	pubKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding public key hex: %w", err)
	}

	privateKey, err := x509.ParseECPrivateKey(privKeyBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing ECDSA private key: %w", err)
	}

	pubKeyInterface, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing ECDSA public key: %w", err)
	}

	publicKey, ok := pubKeyInterface.(*ecdsa.PublicKey)
	if !ok {
		return nil, nil, fmt.Errorf("parsed public key is not ECDSA")
	}

	return privateKey, publicKey, nil
}
