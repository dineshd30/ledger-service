package mock

import (
	"fmt"

	"github.com/stretchr/testify/mock"
)

type UUIDGenerator struct {
	mock.Mock
}

func (s *UUIDGenerator) Generate() string {
	fmt.Println("Called mocked Generate function")
	args := s.Called()
	return args.Get(0).(string)
}
