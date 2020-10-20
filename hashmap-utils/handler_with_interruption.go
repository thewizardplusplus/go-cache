package hashmaputils

import (
	"context"

	hashmap "github.com/thewizardplusplus/go-hashmap"
)

// HandlerWithInterruption ...
func HandlerWithInterruption(
	ctx context.Context,
	handler hashmap.Handler,
) hashmap.Handler {
	return func(key hashmap.Key, value interface{}) bool {
		select {
		case <-ctx.Done():
			return false
		default:
			return handler(key, value)
		}
	}
}
