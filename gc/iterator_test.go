package gc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	cache "github.com/thewizardplusplus/go-cache"
	hashmap "github.com/thewizardplusplus/go-hashmap"
)

func Test_newIterator(test *testing.T) {
	storage := new(MockStorage)
	iterator := newIterator(storage, time.Now)

	mock.AssertExpectationsForObjects(test, storage)
	assert.Zero(test, iterator.counter)
	assert.Equal(test, storage, iterator.storage)
	assert.Equal(test, getPointer(time.Now), getPointer(iterator.clock))
}

func Test_iterator_handleIteration(test *testing.T) {
	type fields struct {
		counter counter
		storage Storage
		clock   cache.Clock
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
		// TODO: add test cases
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
