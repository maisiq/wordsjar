package auth

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/maisiq/go-words-jar/internal/errors"
	"github.com/maisiq/go-words-jar/internal/logger"
	"github.com/maisiq/go-words-jar/internal/models"
)

var (
	ErrAutheticationFailed = fmt.Errorf("bad username-password pair")
)

type Repository interface {
	AddUser(ctx context.Context, user models.User) error
	User(ctx context.Context, username string) (models.User, error)
}

type SecretRepository interface {
	PublicKeys(ctx context.Context) (map[string]any, error)
	GetKID(ctx context.Context) (string, error)
	SignData(ctx context.Context, data string) (signedData string, err error)
}

type tokenStorage interface {
	Store(ctx context.Context, key string, value any, ttl int64) error
	Get(ctx context.Context, key string) (any, error)
}

type Tokens struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}

type UserInfo struct {
	Username string
}

type JWTClaims struct {
	jwt.RegisteredClaims
}

type AuthService struct {
	repo    Repository
	tokens  tokenStorage
	secrets SecretRepository
	log     logger.Logger
}

func NewAuthService(logger logger.Logger, repo Repository, tokenStorage tokenStorage, secrets SecretRepository) *AuthService {
	return &AuthService{
		repo:    repo,
		tokens:  tokenStorage,
		log:     logger,
		secrets: secrets,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, username, password string) error {
	hashed, err := hashPassword(password)
	if err != nil {
		return err
	}
	name := normalizeUsername(username)
	user := models.User{
		ID:             uuid.NewString(),
		Username:       name,
		HashedPassword: hashed,
	}

	repoErr := s.repo.AddUser(ctx, user)
	if repoErr != nil {
		switch repoErr {
		case errors.ErrUserAlreadyExists:
			return repoErr
		default:
			s.log.Errorw("failed to add user to repo", "error", repoErr)
			return errors.ErrInternal
		}
	}
	return nil
}

func (s *AuthService) Authenticate(ctx context.Context, username, plainPassword string) (Tokens, error) {
	username = normalizeUsername(username)

	user, err := s.repo.User(ctx, username)
	if err != nil {
		return Tokens{}, err
	}
	ok, err := verifyPassword(plainPassword, user.HashedPassword)
	if err != nil {
		s.log.Errorw("failed to verify password", "error", err)
		return Tokens{}, errors.ErrInternal
	}

	if !ok {
		return Tokens{}, ErrAutheticationFailed
	}

	kid, kidErr := s.secrets.GetKID(ctx)
	if kidErr != nil {
		s.log.Errorw("failed to get kid", "error", kidErr)
		return Tokens{}, errors.ErrInternal
	}

	unsignedToken, err := CreateUnsignedAccessToken(AccessTokenInput{
		KID:      kid,
		Username: username,
	})
	if err != nil {
		s.log.Errorw("failed to generate access token", "error", err)
		return Tokens{}, errors.ErrInternal
	}

	ss, signErr := s.secrets.SignData(ctx, unsignedToken)
	if signErr != nil {
		s.log.Errorw("failed to sign access token", "error", signErr)
		return Tokens{}, errors.ErrInternal
	}

	refesh, err := GenerateRefreshToken()
	if err != nil {
		s.log.Errorw("failed to generate refresh token", "error", err)
		return Tokens{}, errors.ErrInternal
	}

	return Tokens{
		Access:  ss,
		Refresh: refesh,
	}, nil
}

func (s *AuthService) UserInfo(ctx context.Context, username string) (UserInfo, error) {
	user, err := s.repo.User(ctx, username)
	if err != nil {
		return UserInfo{}, err
	}
	return UserInfo{
		Username: user.Username,
	}, nil
}

func (s *AuthService) GetUser(ctx context.Context, username string) (models.User, error) {
	user, err := s.repo.User(ctx, username)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (s *AuthService) Logout(ctx context.Context, tokens Tokens) error {
	if err := s.tokens.Store(ctx, tokens.Access, "true", 10); err != nil {
		return err
	}
	if err := s.tokens.Store(ctx, tokens.Refresh, "true", 10); err != nil {
		return err
	}
	return nil
}
