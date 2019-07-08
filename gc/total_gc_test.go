package gc

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewTotalGC(test *testing.T) {
	storage := new(MockStorage)
	gc := NewTotalGC(time.Second, storage, time.Now)

	mock.AssertExpectationsForObjects(test, storage)
	assert.Equal(test, time.Second, gc.period)
	assert.Equal(test, storage, gc.storage)
	assert.Equal(test, getPointer(time.Now), getPointer(gc.clock))
}

func getPointer(value interface{}) uintptr {
	return reflect.ValueOf(value).Pointer()
}
