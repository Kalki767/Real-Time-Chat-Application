package mocks

import (
	"context"
	"Real-Time-Chat-Application/repository"

	"github.com/stretchr/testify/mock"
)

type MockCursor struct {
	mock.Mock
}

func (m *MockCursor) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCursor) Next(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

func (m *MockCursor) Decode(val interface{}) error {
	args := m.Called(val)
	return args.Error(0)
}

func (m *MockCursor) Err() error {
	args := m.Called()
	return args.Error(0)
}

// Ensure MockCursor implements CursorInterface
var _ repository.CursorInterface = &MockCursor{}
