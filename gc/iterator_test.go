package gc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_newIterator(test *testing.T) {
	storage := new(MockStorage)
	iterator := newIterator(storage, time.Now)

	mock.AssertExpectationsForObjects(test, storage)
	assert.Zero(test, iterator.counter)
	assert.Equal(test, storage, iterator.storage)
	assert.Equal(test, getPointer(time.Now), getPointer(iterator.clock))
}
