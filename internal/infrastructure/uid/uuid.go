package uid

import "github.com/google/uuid"

type UuidGenerater interface {
	Generate() string
}

type uuidGenerater struct{}

func NewUUIDGenerater() *uuidGenerater {
	return &uuidGenerater{}
}

func (u *uuidGenerater) Generate() string {
	return uuid.New().String()
}
