package gc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/go-cache/models"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

func Test_newIterator(test *testing.T) {
	storage := new(MockStorage)
	iterator := newIterator(storage, time.Now, 20, 0.25)

	wantCounter := counter{maxIteratedCount: 20, minExpiredPercent: 0.25}
	assert.Equal(test, wantCounter, iterator.counter)

	mock.AssertExpectationsForObjects(test, storage)
	assert.Equal(test, storage, iterator.storage)

	// don't use the reflect.Value.Pointer() method for this check; see details:
	// * https://golang.org/pkg/reflect/#Value.Pointer
	// * https://stackoverflow.com/a/9644797
	require.NotNil(test, iterator.clock)
	assert.WithinDuration(test, time.Now(), iterator.clock(), time.Hour)
}

func Test_iterator_handleIteration(test *testing.T) {
	type fields struct {
		counter counter

		storage hashmap.Storage
		clock   models.Clock
	}
	type args struct {
		key   hashmap.Key
		value interface{}
	}

	for _, data := range []struct {
		name        string
		fields      fields
		args        args
		wantCounter counter
		wantOk      assert.BoolAssertionFunc
	}{
		{
			name: "with a not expired value " +
				"and with iteration count less than maximum",
			fields: fields{
				counter: counter{
					maxIteratedCount:  20,
					minExpiredPercent: 0.25,

					iteratedCount: 15,
					expiredCount:  3,
				},

				storage: new(MockStorage),
				clock:   clock,
			},
			args: args{
				key:   NewMockKeyWithID(23),
				value: models.Value{Data: "data", ExpirationTime: clock().Add(time.Second)},
			},
			wantCounter: counter{
				maxIteratedCount:  20,
				minExpiredPercent: 0.25,

				iteratedCount: 16,
				expiredCount:  3,
			},
			wantOk: assert.True,
		},
		{
			name: "with a not expired value " +
				"and with iteration count greater than maximum",
			fields: fields{
				counter: counter{
					maxIteratedCount:  20,
					minExpiredPercent: 0.25,

					iteratedCount: 25,
					expiredCount:  3,
				},

				storage: new(MockStorage),
				clock:   clock,
			},
			args: args{
				key:   NewMockKeyWithID(23),
				value: models.Value{Data: "data", ExpirationTime: clock().Add(time.Second)},
			},
			wantCounter: counter{
				maxIteratedCount:  20,
				minExpiredPercent: 0.25,

				iteratedCount: 26,
				expiredCount:  3,
			},
			wantOk: assert.False,
		},
		{
			name: "with an expired value",
			fields: fields{
				counter: counter{
					maxIteratedCount:  20,
					minExpiredPercent: 0.25,

					iteratedCount: 15,
					expiredCount:  3,
				},

				storage: func() hashmap.Storage {
					storage := new(MockStorage)
					storage.On("Delete", NewMockKeyWithID(23))

					return storage
				}(),
				clock: clock,
			},
			args: args{
				key: NewMockKeyWithID(23),
				value: models.Value{
					Data:           "data",
					ExpirationTime: clock().Add(-time.Second),
				},
			},
			wantCounter: counter{
				maxIteratedCount:  20,
				minExpiredPercent: 0.25,

				iteratedCount: 16,
				expiredCount:  4,
			},
			wantOk: assert.True,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			iterator := &iterator{
				counter: data.fields.counter,

				storage: data.fields.storage,
				clock:   data.fields.clock,
			}
			gotOk := iterator.handleIteration(data.args.key, data.args.value)

			mock.AssertExpectationsForObjects(test, data.fields.storage, data.args.key)
			assert.Equal(test, data.wantCounter, iterator.counter)
			data.wantOk(test, gotOk)
		})
	}
}
