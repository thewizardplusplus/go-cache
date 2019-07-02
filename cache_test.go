package cache

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewCache(test *testing.T) {
	storage := new(MockStorage)
	cache := NewCache(storage, time.Now)

	mock.AssertExpectationsForObjects(test, storage)
	assert.Equal(test, storage, cache.storage)
	assert.Equal(test, getPointer(time.Now), getPointer(cache.clock))
}

func getPointer(value interface{}) uintptr {
	return reflect.ValueOf(value).Pointer()
}
