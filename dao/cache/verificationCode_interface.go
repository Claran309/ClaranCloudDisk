package cache

import "context"

type VerificationCodeCache interface {
	SaveVerificationCode(ctx context.Context, email, code string) error
	GetVerificationCode(ctx context.Context, email string) (string, error)
	DeleteVerificationCode(ctx context.Context, email string) error
	SetRateLimit(ctx context.Context, email string) error
	CheckRateLimit(ctx context.Context, email string) (bool, error)
}
