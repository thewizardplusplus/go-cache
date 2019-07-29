package gc

import (
	cache "github.com/thewizardplusplus/go-cache"
)

type iterator struct {
	counter

	storage Storage
	clock   cache.Clock
}
