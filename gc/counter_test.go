package gc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newCounter(test *testing.T) {
	counterInstance := newCounter(20, 0.25)

	wantCounterInstance := counter{maxIteratedCount: 20, minExpiredPercent: 0.25}
	assert.Equal(test, wantCounterInstance, counterInstance)
}

func Test_counter_stopIterate(test *testing.T) {
	type fields struct {
		maxIteratedCount int

		iteratedCount int
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   assert.BoolAssertionFunc
	}{
		{
			name:   "without iterations",
			fields: fields{maxIteratedCount: 20, iteratedCount: 0},
			want:   assert.False,
		},
		{
			name:   "with iteration count less than maximum",
			fields: fields{maxIteratedCount: 20, iteratedCount: 10},
			want:   assert.False,
		},
		{
			name:   "with iteration count greater than maximum",
			fields: fields{maxIteratedCount: 20, iteratedCount: 30},
			want:   assert.True,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			counter := counter{ // nolint: vetshadow
				maxIteratedCount: data.fields.maxIteratedCount,

				iteratedCount: data.fields.iteratedCount,
			}
			got := counter.stopIterate()

			data.want(test, got)
		})
	}
}

func Test_counter_stopClean(test *testing.T) {
	type fields struct {
		minExpiredPercent float64

		iteratedCount int
		expiredCount  int
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   assert.BoolAssertionFunc
	}{
		{
			name:   "without iterations",
			fields: fields{minExpiredPercent: 0.25, iteratedCount: 0},
			want:   assert.True,
		},
		{
			name:   "with expired percent less than minimum",
			fields: fields{minExpiredPercent: 0.25, iteratedCount: 15, expiredCount: 3},
			want:   assert.True,
		},
		{
			name:   "with expired percent greater than minimum",
			fields: fields{minExpiredPercent: 0.25, iteratedCount: 15, expiredCount: 5},
			want:   assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			counter := counter{ // nolint: vetshadow
				minExpiredPercent: data.fields.minExpiredPercent,

				iteratedCount: data.fields.iteratedCount,
				expiredCount:  data.fields.expiredCount,
			}
			got := counter.stopClean()

			data.want(test, got)
		})
	}
}
