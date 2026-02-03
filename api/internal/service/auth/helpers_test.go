package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestParseAccessTokenWithRS256(t *testing.T) {
	username := "testuser"
	actualSigningMethod := jwt.SigningMethodRS256

	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := &privateKey.PublicKey

	kid := "1"
	keys := map[string]any{
		kid: publicKey,
	}

	baseClaims := JWTClaims{
		jwt.RegisteredClaims{
			Issuer:    "wordsjar",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   username,
		},
	}

	t.Run("succesfully returns right username", func(t *testing.T) {
		token := jwt.NewWithClaims(actualSigningMethod, baseClaims)
		token.Header["kid"] = kid
		ss, _ := token.SignedString(privateKey)

		data, err := ParseAccessToken(ss, keys)

		assert.NoError(t, err)
		assert.Equal(t, username, data.Username)
	})

	t.Run("returns token expired error", func(t *testing.T) {
		claims := baseClaims
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-24 * time.Hour))

		token := jwt.NewWithClaims(actualSigningMethod, claims)
		token.Header["kid"] = kid
		ss, _ := token.SignedString(privateKey)

		_, err := ParseAccessToken(ss, keys)

		assert.ErrorIs(t, err, jwt.ErrTokenExpired)
	})

	t.Run("returns invalid key type error when token with different hash lenght", func(t *testing.T) {
		wrongMethod := jwt.SigningMethodRS384
		token := jwt.NewWithClaims(wrongMethod, baseClaims)
		token.Header["kid"] = kid
		ss, _ := token.SignedString(privateKey)

		_, err := ParseAccessToken(ss, keys)

		assert.ErrorIs(t, err, jwt.ErrInvalidKeyType)
	})

	t.Run("returns token malformed error when token with different alg", func(t *testing.T) {
		wrongMethod := jwt.SigningMethodES256

		token := jwt.NewWithClaims(wrongMethod, baseClaims)
		token.Header["kid"] = kid
		ss, _ := token.SignedString(privateKey)

		_, err := ParseAccessToken(ss, keys)

		assert.ErrorIs(t, err, jwt.ErrTokenMalformed)
	})

	t.Run("returns token malformed error when pass broken token", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodES256, baseClaims)
		token.Header["kid"] = kid
		ss, _ := token.SignedString(privateKey)

		invalidToken := "c" + ss + "1"
		_, err := ParseAccessToken(invalidToken, keys)

		assert.ErrorIs(t, err, jwt.ErrTokenMalformed)
	})

	t.Run("returns error when token contains diff issuer", func(t *testing.T) {
		claims := baseClaims
		claims.Subject = "testuser"
		claims.Issuer = "wrongone"

		token := jwt.NewWithClaims(actualSigningMethod, claims)
		token.Header["kid"] = kid
		ss, _ := token.SignedString(privateKey)

		_, err := ParseAccessToken(ss, keys)

		assert.ErrorIs(t, err, jwt.ErrTokenInvalidIssuer)
	})
}
