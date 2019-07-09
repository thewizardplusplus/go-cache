package gc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	cache "github.com/thewizardplusplus/go-cache"
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
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			gc := PartialGC{data.fields.storage, data.fields.clock}
			gc.Clean()

			mock.AssertExpectationsForObjects(test, data.fields.storage)
		})
	}
}
