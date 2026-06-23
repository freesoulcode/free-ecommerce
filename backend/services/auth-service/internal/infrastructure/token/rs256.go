package token

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	applicationcredential "github.com/freesoulcode/free-ecommerce/backend/services/auth-service/internal/application/credential"
	"github.com/golang-jwt/jwt/v5"
)

type RS256Signer struct {
	privateKey *rsa.PrivateKey
}

func NewRS256Signer(privateKeyPEM string) (*RS256Signer, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("decode rsa private key pem: invalid pem")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		parsedKey, parseErr := x509.ParsePKCS8PrivateKey(block.Bytes)
		if parseErr != nil {
			return nil, fmt.Errorf("parse rsa private key: %w", err)
		}

		rsaKey, ok := parsedKey.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("parse rsa private key: invalid key type")
		}
		privateKey = rsaKey
	}

	return &RS256Signer{privateKey: privateKey}, nil
}

func (s *RS256Signer) SignAccessToken(input applicationcredential.AccessTokenClaims) (string, error) {
	claims := jwt.MapClaims{
		"sub": input.Subject,
		"iss": input.Issuer,
		"aud": input.Audience,
		"uid": input.UserID,
		"sid": input.SessionID,
		"typ": input.TokenType,
		"idt": input.Identifier,
		"iat": input.IssuedAt.Unix(),
		"exp": input.ExpiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.privateKey)
}

type RandomTokenGenerator struct{}

func NewRandomTokenGenerator() *RandomTokenGenerator {
	return &RandomTokenGenerator{}
}

func (g *RandomTokenGenerator) Generate() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate random token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}
