package gc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	cache "github.com/thewizardplusplus/go-cache"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

func TestNewPartialGC(test *testing.T) {
	storage := new(MockStorage)
	gc := NewPartialGC(storage, time.Now)

	mock.AssertExpectationsForObjects(test, storage)
	assert.Equal(test, storage, gc.storage)
	assert.Equal(test, getPointer(time.Now), getPointer(gc.clock))
}

func TestPartialGC_Clean(test *testing.T) {
	type fields struct {
		storage Storage
		clock   cache.Clock
	}

	for _, data := range []struct {
		name   string
		fields fields
	}{
		{
			name: "without values",
			fields: fields{
				storage: func() Storage {
					storage := new(MockStorage)
					storage.
						On("Iterate", mock.MatchedBy(func(handler hashmap.Handler) bool {
							return handler != nil
						})).
						Return(true)

					return storage
				}(),
				clock: clock,
			},
		},
		{
			name: "with value and its expiration time less than the current one",
			fields: fields{
				storage: func() Storage {
					var expiredCount int
					storage := new(MockStorage)
					storage.
						On("Iterate", mock.MatchedBy(func(handler hashmap.Handler) bool {
							return handler != nil
						})).
						Return(func(handler hashmap.Handler) bool {
							expiredCount++
							if expiredCount > 1 {
								return true
							}

							return handler(NewMockKeyWithID(23), cache.Value{
								Data:           "data",
								ExpirationTime: clock().Add(-time.Second),
							})
						})
					storage.On("Delete", NewMockKeyWithID(23))

					return storage
				}(),
				clock: clock,
			},
		},
		{
			name: "with value and its expiration time greater than the current one",
			fields: fields{
				storage: func() Storage {
					storage := new(MockStorage)
					storage.
						On("Iterate", mock.MatchedBy(func(handler hashmap.Handler) bool {
							return handler != nil
						})).
						Return(func(handler hashmap.Handler) bool {
							return handler(NewMockKeyWithID(23), cache.Value{
								Data:           "data",
								ExpirationTime: clock().Add(time.Second),
							})
						})

					return storage
				}(),
				clock: clock,
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gc := PartialGC{data.fields.storage, data.fields.clock}
			gc.Clean()

			mock.AssertExpectationsForObjects(test, data.fields.storage)
		})
	}
}
