package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/go-cache/gc"
	"github.com/thewizardplusplus/go-cache/models"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

func Test_newConfigWithGC(test *testing.T) {
	type args struct {
		options []OptionWithGC
	}

	for _, data := range []struct {
		name          string
		args          args
		wantStorage   hashmap.Storage
		wantClockTime time.Time
		wantGCType    gc.GC
		wantGCPeriod  time.Duration
	}{
		{
			name: "with the default config",
			args: args{
				options: nil,
			},
			wantStorage:   hashmap.NewConcurrentHashMap(),
			wantClockTime: time.Now(),
			wantGCType:    gc.PartialGC{},
			wantGCPeriod:  100 * time.Millisecond,
		},
		{
			name: "with the set storage",
			args: args{
				options: []OptionWithGC{WithGCAndStorage(new(MockStorage))},
			},
			wantStorage:   new(MockStorage),
			wantClockTime: time.Now(),
			wantGCType:    gc.PartialGC{},
			wantGCPeriod:  100 * time.Millisecond,
		},
		{
			name: "with the set clock",
			args: args{
				options: []OptionWithGC{WithGCAndClock(clock)},
			},
			wantStorage:   hashmap.NewConcurrentHashMap(),
			wantClockTime: clock(),
			wantGCType:    gc.PartialGC{},
			wantGCPeriod:  100 * time.Millisecond,
		},
		{
			name: "with the set GC factory",
			args: args{
				options: []OptionWithGC{
					WithGCAndGCFactory(
						func(storage hashmap.Storage, clock models.Clock) gc.GC {
							return new(MockGC)
						},
					),
				},
			},
			wantStorage:   hashmap.NewConcurrentHashMap(),
			wantClockTime: time.Now(),
			wantGCType:    new(MockGC),
			wantGCPeriod:  100 * time.Millisecond,
		},
		{
			name: "with the set GC period",
			args: args{
				options: []OptionWithGC{WithGCAndGCPeriod(23 * time.Second)},
			},
			wantStorage:   hashmap.NewConcurrentHashMap(),
			wantClockTime: time.Now(),
			wantGCType:    gc.PartialGC{},
			wantGCPeriod:  23 * time.Second,
		},
		{
			name: "with the set config",
			args: args{
				options: []OptionWithGC{
					WithGCAndStorage(new(MockStorage)),
					WithGCAndClock(clock),
					WithGCAndGCFactory(
						func(storage hashmap.Storage, clock models.Clock) gc.GC {
							return new(MockGC)
						},
					),
					WithGCAndGCPeriod(23 * time.Second),
				},
			},
			wantStorage:   new(MockStorage),
			wantClockTime: clock(),
			wantGCType:    new(MockGC),
			wantGCPeriod:  23 * time.Second,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := newConfigWithGC(data.args.options)

			_, ok := got.storage.(interface {
				AssertExpectations(assert.TestingT) bool // nolint: staticcheck
			})
			if ok {
				mock.AssertExpectationsForObjects(test, got.storage)
			}
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
