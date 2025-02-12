package types

import "time"

type LimiterInfo struct {
	Count          int
	BlockTimestamp time.Time
}
