package token

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"sync"
)

type SecretRepository interface {
	PublicKeys(ctx context.Context) (map[string]any, error)
	GetKID(ctx context.Context) (string, error)
	SignData(ctx context.Context, data string) (string, error)
}

type InMemorySecretRepository struct {
	currentKID string
	publicKeys sync.Map
	signKeys   sync.Map
}

func NewInMemory() *InMemorySecretRepository {
	kid := "1"

	repo := InMemorySecretRepository{}
	repo.currentKID = kid

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	publicKey := &privateKey.PublicKey

	repo.signKeys.Store(kid, privateKey)
	repo.publicKeys.Store(kid, publicKey)

	return &repo
}

func (r *InMemorySecretRepository) PublicKeys(ctx context.Context) (map[string]any, error) {
	keys := make(map[string]any)

	v, ok := r.publicKeys.Load(r.currentKID)
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	pk, ok := v.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to key with right format from publicKeys")
	}

	keys[r.currentKID] = pk
	return keys, nil
}

func (r *InMemorySecretRepository) GetKID(ctx context.Context) (string, error) {
	return r.currentKID, nil
}

func (r *InMemorySecretRepository) SignData(ctx context.Context, data string) (string, error) {
	privateKey, ok := r.signKeys.Load(r.currentKID)
	if !ok {
		return "", fmt.Errorf("key for signing data is not stored")
	}

	msgHash := sha256.New()
	_, err := msgHash.Write([]byte(data))
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	msgHashSum := msgHash.Sum(nil)

	signatureRaw, err := rsa.SignPKCS1v15(rand.Reader, privateKey.(*rsa.PrivateKey), crypto.SHA256, msgHashSum)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	encodedSignature := base64.RawURLEncoding.EncodeToString(signatureRaw)
	return fmt.Sprintf("%s.%s", data, encodedSignature), nil

}
