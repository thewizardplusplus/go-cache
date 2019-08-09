package gc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

func Test_withInterruption(test *testing.T) {
	type wrapperArgs struct {
		ctx     context.Context
		handler Handler
	}
	type handlerArgs struct {
		key   hashmap.Key
		value interface{}
	}

	for _, data := range []struct {
		name        string
		wrapperArgs wrapperArgs
		handlerArgs handlerArgs
		want        assert.BoolAssertionFunc
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			handler := withInterruption(
				data.wrapperArgs.ctx,
				data.wrapperArgs.handler.Handle,
			)
			require.NotNil(test, handler)

			got := handler(data.handlerArgs.key, data.handlerArgs.value)

			mock.AssertExpectationsForObjects(
				test,
				data.wrapperArgs.handler,
				data.handlerArgs.key,
			)
			data.want(test, got)
		})
	}
}
