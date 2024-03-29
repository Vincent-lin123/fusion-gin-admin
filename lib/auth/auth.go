package auth

import (
	"context"
	"github.com/pkg/errors"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type TokenInfo interface {
	GetAccessToken() string
	GetTokenType() string
	GetExpiresAt() int64
	EncodeToJSON() ([]byte, error)
}

type Auther interface {
	GenerateToken(ctx context.Context, userID string) (TokenInfo, error)

	DestroyToken(ctx context.Context, accessToken string) error

	ParseUserID(ctx context.Context, accessToken string) (string, error)

	Release() error
}
