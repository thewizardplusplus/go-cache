package cache

import (
	"time"

	"github.com/thewizardplusplus/go-cache/models"
)

// Value ...
type Value struct {
	Data           interface{}
	ExpirationTime time.Time
}

// IsExpired ...
func (value Value) IsExpired(clock models.Clock) bool {
	return !value.ExpirationTime.IsZero() && clock().After(value.ExpirationTime)
}
