package cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/go-cache/gc"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

func Test_newConfigWithGC(test *testing.T) {
	type args struct {
		options []OptionWithGC
	}

	for _, data := range []struct {
		name          string
		args          args
		wantCtx       context.Context
		wantStorage   hashmap.Storage
		wantClockTime time.Time
		wantGCType    gc.GC
		wantGCPeriod  time.Duration
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			got := newConfigWithGC(data.args.options)

			for _, value := range []interface{}{got.ctx, got.storage} {
				_, ok := value.(interface {
					AssertExpectations(assert.TestingT) bool // nolint: staticcheck
				})
				if ok {
					mock.AssertExpectationsForObjects(test, value)
				}
			}
			assert.Equal(test, data.wantCtx, got.ctx)
			assert.Equal(test, data.wantStorage, got.storage)
			assert.Equal(test, data.wantGCPeriod, got.gcPeriod)

			// don't use the reflect.Value.Pointer() method for checks below;
			// see details:
			// * https://golang.org/pkg/reflect/#Value.Pointer
			// * https://stackoverflow.com/a/9644797

			require.NotNil(test, got.clock)
			assert.WithinDuration(test, data.wantClockTime, got.clock(), time.Hour)

			require.NotNil(test, got.gcFactory)
			assert.IsType(test, data.wantGCType, got.gcFactory(got.storage, got.clock))
		})
	}
}
