package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
)

const JWTIssuer = "wordsjar"

type AccessTokenData struct {
	Username string
}

type AccessTokenInput struct {
	KID      string
	Username string
}

func hashPassword(plain string) (string, error) {
	return argon2id.CreateHash(plain, argon2id.DefaultParams)
}

func verifyPassword(plain, hashed string) (bool, error) {
	return argon2id.ComparePasswordAndHash(plain, hashed)
}

func normalizeUsername(s string) string {
	return strings.Trim(s, " ")
}

func CreateUnsignedAccessToken(in AccessTokenInput) (string, error) {
	claims := JWTClaims{
		jwt.RegisteredClaims{
			Issuer:    JWTIssuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   in.Username, //TODO: swap with user.ID
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = in.KID
	return token.SigningString()
}

func ParseAccessToken(token string, keys map[string]any) (AccessTokenData, error) {
	var claims JWTClaims
	t, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		if t.Header["alg"] == "RS256" {
			kid, ok := t.Header["kid"]
			if !ok {
				return "", jwt.ErrInvalidKeyType
			}
			key, ok := keys[kid.(string)]
			if !ok {
				return "", jwt.ErrInvalidKeyType
			}
			rsaKey, ok := key.(*rsa.PublicKey)
			if !ok {
				return "", fmt.Errorf("failed to convert key to rsa.PublicKey")
			}
			return rsaKey, nil
		}
		return "", jwt.ErrInvalidKeyType
	})
	if err != nil {
		return AccessTokenData{}, err
	}

	if !t.Valid {
		return AccessTokenData{}, fmt.Errorf("invalid token")
	}

	if claims.Issuer != JWTIssuer {
		return AccessTokenData{}, jwt.ErrTokenInvalidIssuer
	}

	return AccessTokenData{
		Username: claims.Subject,
	}, nil

}

func GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
