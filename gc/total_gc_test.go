package gc

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	cache "github.com/thewizardplusplus/go-cache"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

type MockKeyWithID struct {
	MockKey

	ID int
}

func NewMockKeyWithID(id int) *MockKeyWithID {
	return &MockKeyWithID{ID: id}
}

func TestNewTotalGC(test *testing.T) {
	storage := new(MockStorage)
	gc := NewTotalGC(storage, time.Now)

	mock.AssertExpectationsForObjects(test, storage)
	assert.Equal(test, storage, gc.storage)
	assert.Equal(test, getPointer(time.Now), getPointer(gc.clock))
}

func TestTotalGC_Clean(test *testing.T) {
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
					storage := new(MockStorage)
					storage.
						On("Iterate", mock.MatchedBy(func(handler hashmap.Handler) bool {
							return handler != nil
						})).
						Return(func(handler hashmap.Handler) bool {
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
			gc := TotalGC{data.fields.storage, data.fields.clock}
			gc.Clean()

			mock.AssertExpectationsForObjects(test, data.fields.storage)
		})
	}
}

func getPointer(value interface{}) uintptr {
	return reflect.ValueOf(value).Pointer()
}

func clock() time.Time {
	return time.Date(
		2006, time.January, 2, // year, month, day
		15, 4, 5, // hour, minute, second
		0,        // nanosecond
		time.UTC, // location
	)
}
