package middleware

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the RepoInterface
type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) IncrBase(key string) (int64, error) {
	args := m.Called(key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepo) SetExBase(key string, duration time.Duration) (bool, error) {
	args := m.Called(key, duration)
	return args.Bool(0), args.Error(1)
}

func (m *MockRepo) GetBase(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockRepo) SetBase(key string, value interface{}) (string, error) {
	args := m.Called(key, value)
	return args.String(0), args.Error(1)
}

func TestIncrementRequest(t *testing.T) {
	mockRepo := new(MockRepo)
	key := "test_key"
	duration := time.Minute

	mockRepo.On("IncrBase", key).Return(int64(1), nil)
	mockRepo.On("SetExBase", key, duration).Return(true, nil)

	IncrementRequest(mockRepo, key, duration)

	mockRepo.AssertExpectations(t)
}

func TestIncrementRequestWithEmptyKey(t *testing.T) {
	mockRepo := new(MockRepo)
	IncrementRequest(mockRepo, "", time.Minute)
	mockRepo.AssertNotCalled(t, "IncrBase")
	mockRepo.AssertNotCalled(t, "SetExBase")
}

func TestExceededLimit(t *testing.T) {
	mockRepo := new(MockRepo)
	key := "test_key"
	limit := 10

	mockRepo.On("GetBase", key).Return("10", nil)

	exceeded := ExceededLimit(mockRepo, key, limit)
	assert.True(t, exceeded)

	mockRepo.AssertExpectations(t)
}

func TestExceededLimitNotExceeded(t *testing.T) {
	mockRepo := new(MockRepo)
	key := "test_key"
	limit := 10

	mockRepo.On("GetBase", key).Return("9", nil)

	exceeded := ExceededLimit(mockRepo, key, limit)
	assert.False(t, exceeded)

	mockRepo.AssertExpectations(t)
}

func TestExceededLimitWithError(t *testing.T) {
	mockRepo := new(MockRepo)
	key := "test_key"
	limit := 10

	mockRepo.On("GetBase", key).Return("", errors.New("some error"))

	exceeded := ExceededLimit(mockRepo, key, limit)
	assert.False(t, exceeded)

	mockRepo.AssertExpectations(t)
}

func TestGetRateLimitConfigs(t *testing.T) {
	limit, duration, err := GetRateLimitConfigs("10", "1m")
	assert.NoError(t, err)
	assert.Equal(t, 10, limit)
	assert.Equal(t, time.Minute, duration)
}

func TestGetRateLimitConfigsInvalidLimit(t *testing.T) {
	_, _, err := GetRateLimitConfigs("invalid", "1m")
	assert.Error(t, err)
}

func TestGetRateLimitConfigsInvalidDuration(t *testing.T) {
	_, _, err := GetRateLimitConfigs("10", "invalid")
	assert.Error(t, err)
}
