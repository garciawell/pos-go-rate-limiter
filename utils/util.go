package utils

import (
	"time"

	"github.com/garciawell/pos-go-rate-limiter/types"
)

func ValidateMiddleware(data types.LimiterInfo, duration time.Duration) {
	if data.BlockTimestamp.IsZero() {
		limiterInfo := data
		limiterInfo.BlockTimestamp = time.Now()
		data = limiterInfo
	}
	if time.Since(data.BlockTimestamp) > duration {
		limiterInfo := data
		limiterInfo.Count = 0
		limiterInfo.BlockTimestamp = time.Time{}
		data = limiterInfo
	}
}
