package utils

import (
	"testing"
	"time"

	"github.com/garciawell/pos-go-rate-limiter/types"
)

func TestValidateMiddleware(t *testing.T) {
	tests := []struct {
		name     string
		data     types.LimiterInfo
		duration time.Duration
		expected types.LimiterInfo
	}{
		{
			name: "BlockTimestamp is zero",
			data: types.LimiterInfo{
				BlockTimestamp: time.Time{},
				Count:          5,
			},
			duration: 10 * time.Second,
			expected: types.LimiterInfo{
				BlockTimestamp: time.Time{},
				Count:          5,
			},
		},
		{
			name: "BlockTimestamp is not zero and duration not exceeded",
			data: types.LimiterInfo{
				BlockTimestamp: time.Now().Add(-5 * time.Second),
				Count:          5,
			},
			duration: 10 * time.Second,
			expected: types.LimiterInfo{
				BlockTimestamp: time.Now().Add(-5 * time.Second),
				Count:          5,
			},
		},
		{
			name: "BlockTimestamp is not zero and duration exceeded",
			data: types.LimiterInfo{
				BlockTimestamp: time.Now().Add(-15 * time.Second),
				Count:          5,
			},
			duration: 10 * time.Second,
			expected: types.LimiterInfo{
				BlockTimestamp: time.Time{},
				Count:          5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ValidateMiddleware(tt.data, tt.duration)
			if tt.data.Count != tt.expected.Count {
				t.Errorf("expected Count %d, got %d", tt.expected.Count, tt.data.Count)
			}
			if tt.data.BlockTimestamp.IsZero() && !tt.expected.BlockTimestamp.IsZero() {
				t.Errorf("expected BlockTimestamp to be non-zero, got zero")
			}
		})
	}
}
