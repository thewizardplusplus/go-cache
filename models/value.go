package models

import (
	"time"
)

// Value ...
type Value struct {
	Data           interface{}
	ExpirationTime time.Time
}

// IsExpired ...
func (value Value) IsExpired(clock Clock) bool {
	return !value.ExpirationTime.IsZero() && clock().After(value.ExpirationTime)
}
