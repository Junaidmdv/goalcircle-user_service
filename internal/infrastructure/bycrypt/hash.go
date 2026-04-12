package bycrypt

import "golang.org/x/crypto/bcrypt"

type PasswordHasher interface {
	HashPassword(string) (string, error)
	ComparePassword(string, string) error
}

type bycriptHasher struct {
	cost int
}

func NewBycriptHasher(cost int) *bycriptHasher {
	return &bycriptHasher{
		cost: cost,
	}
}

func (b *bycriptHasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	return string(bytes), err
}

func (b *bycriptHasher) ComparePassword(hashedpassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedpassword), []byte(password))
	return err
}
