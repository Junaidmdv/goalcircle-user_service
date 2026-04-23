package bycrypt

import (
	"errors"
	"fmt"

	"github.com/Junaidmdv/goalcircle-user_service/internal/domain"
	"github.com/Junaidmdv/goalcircle-user_service/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	HashPassword(string) (string, error)
	ComparePassword(string, string) error
}

type bycriptHasher struct {
	cost   int
	logger logger.Logger
}

func NewBycriptHasher(cost int, logger logger.Logger) *bycriptHasher {
	return &bycriptHasher{
		cost:   cost,
		logger: logger,
	}
}

func (b *bycriptHasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	return string(bytes), err
}

// compare the login password with actual password entered while registeration
func (b *bycriptHasher) ComparePassword(hashedpassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedpassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			b.logger.Warn("invalid password", err, fmt.Errorf("invalid password entered %v", err))
			return domain.NewUnAuthenticatedError("Invalid password. Please try again.")
		}
		b.logger.Error("internal error", "error", err)
		return domain.NewInternalError("Something went wrong. Please try again later.", err)
	}
	return err
}
