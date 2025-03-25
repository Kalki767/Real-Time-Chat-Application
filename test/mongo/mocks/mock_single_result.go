package mocks

import (
	"Real-Time-Chat-Application/repository"

	"github.com/stretchr/testify/mock"
)

type MockSingleResult struct {
	mock.Mock
}


func (m *MockSingleResult) Decode(v interface{}) error {
	args := m.Called(v)
	return args.Error(0)
}

var _ repository.SingleResultInterface = &MockSingleResult{}