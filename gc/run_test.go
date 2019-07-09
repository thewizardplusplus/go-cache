package gc

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

func TestRun(test *testing.T) {
	var waiter sync.WaitGroup
	waiter.Add(1)

	gc := new(MockGC)
	gc.On("Clean")

	const period = 100 * time.Millisecond
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer waiter.Done()
		Run(ctx, gc, period)
	}()

	time.Sleep(period * 2)
	cancel()
	waiter.Wait()

	mock.AssertExpectationsForObjects(test, gc)
}
