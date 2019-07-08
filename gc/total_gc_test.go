package gc

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	cache "github.com/thewizardplusplus/go-cache"
)

func TestNewTotalGC(test *testing.T) {
	storage := new(MockStorage)
	gc := NewTotalGC(time.Second, storage, time.Now)

	mock.AssertExpectationsForObjects(test, storage)
	assert.Equal(test, time.Second, gc.period)
	assert.Equal(test, storage, gc.storage)
	assert.Equal(test, getPointer(time.Now), getPointer(gc.clock))
}

func TestTotalGC_Clean(test *testing.T) {
	type fields struct {
		period  time.Duration
		storage Storage
		clock   cache.Clock
	}

	for _, data := range []struct {
		name   string
		fields fields
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			gc := TotalGC{data.fields.period, data.fields.storage, data.fields.clock}
			gc.Clean()

			mock.AssertExpectationsForObjects(test, data.fields.storage)
		})
	}
}

func getPointer(value interface{}) uintptr {
	return reflect.ValueOf(value).Pointer()
}
