package utils

import (
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"go-pos/internal/domain"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtCustomClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID uuid.UUID, username, role string, secretPrivateKey *rsa.PrivateKey) (string, error) {
	now := time.Now()
	accessClaims := jwtCustomClaims{
		UserID:   userID.String(),
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			Issuer:    "go-pos-auth",
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(getEnvAsUint32("JWT_ACCESS_TOKEN_EXPIRED_SECONDS", 3600)) * time.Second)),
		},
	}

	accessJwt := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)

	return accessJwt.SignedString(secretPrivateKey)
}

func GenerateRefreshToken(userID uuid.UUID, secretPrivateKey *rsa.PrivateKey) (string, error) {
	now := time.Now()
	accessClaims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		Issuer:    "go-pos-auth",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(getEnvAsUint32("JWT_REFRESH_TOKEN_EXPIRED_SECONDS", 604800)) * time.Second)),
	}

	accessJwt := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)

	return accessJwt.SignedString(secretPrivateKey)
}

func LoadRSAPrivateKey() (*rsa.PrivateKey, error) {
	b64Key := os.Getenv("RSA_PRIVATE_KEY_B64")
	if b64Key == "" {
		return nil, errors.New("RSA_PRIVATE_KEY_B64 not set")
	}

	pemBytes, err := base64.StdEncoding.DecodeString(b64Key)
	if err != nil {
		return nil, errors.New("Base64 RSA_PRIVATE_KEY_B64 invalid")
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pemBytes)
	if err != nil {
		return nil, errors.New("RSA_PRIVATE_KEY_B64 invalid value")
	}

	return privateKey, nil
}

func LoadRSAPublicKey() (*rsa.PublicKey, error) {
	b64Key := os.Getenv("RSA_PUBLIC_KEY_B64")
	if b64Key == "" {
		return nil, errors.New("RSA_PUBLIC_KEY_B64 not set")
	}

	pemBytes, err := base64.StdEncoding.DecodeString(b64Key)
	if err != nil {
		return nil, errors.New("Base64 RSA_PUBLIC_KEY_B64 invalid")
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(pemBytes)
	if err != nil {
		return nil, errors.New("RSA_PUBLIC_KEY_B64 invalid value")
	}

	return publicKey, nil
}

func GetClaimsFromAccessToken(token string, secretPublicKey *rsa.PublicKey) (string, string, error) {
	accessToken, err := jwt.ParseWithClaims(token, &jwtCustomClaims{}, func(t *jwt.Token) (any, error) {
		return secretPublicKey, nil
	})
	if err != nil {
		return "", "", err
	}

	accessTokenClaims, ok := accessToken.Claims.(*jwtCustomClaims)
	if !ok {
		return "", "", domain.ErrUnauthorized
	}

	if accessToken.Valid {
		return accessTokenClaims.UserID, accessTokenClaims.Role, nil
	}

	return "", "", domain.ErrUnauthorized
}
