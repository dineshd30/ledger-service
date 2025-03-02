package ledger

import "github.com/google/uuid"

type UUIDGenerator interface {
	Generate() string
}

type uuidGenerator struct{}

func NewUUIDGenerator() UUIDGenerator {
	return &uuidGenerator{}
}

func (ug *uuidGenerator) Generate() string {
	return uuid.New().String()
}
