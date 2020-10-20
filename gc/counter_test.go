package gc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_counter_stopIterate(test *testing.T) {
	type fields struct {
		maxIteratedCount int

		iterated int
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   assert.BoolAssertionFunc
	}{
		{
			name:   "without iterations",
			fields: fields{maxIteratedCount: 20, iterated: 0},
			want:   assert.False,
		},
		{
			name:   "with iteration count less than maximum",
			fields: fields{maxIteratedCount: 20, iterated: 10},
			want:   assert.False,
		},
		{
			name:   "with iteration count greater than maximum",
			fields: fields{maxIteratedCount: 20, iterated: 30},
			want:   assert.True,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			counter := counter{ // nolint: vetshadow
				maxIteratedCount: data.fields.maxIteratedCount,

				iterated: data.fields.iterated,
			}
			got := counter.stopIterate()

			data.want(test, got)
		})
	}
}

func Test_counter_stopClean(test *testing.T) {
	type fields struct {
		iterated int
		expired  int
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   assert.BoolAssertionFunc
	}{
		{
			name:   "without iterations",
			fields: fields{iterated: 0},
			want:   assert.True,
		},
		{
			name:   "with expired percent less than minimum",
			fields: fields{iterated: 15, expired: 3},
			want:   assert.True,
		},
		{
			name:   "with expired percent greater than minimum",
			fields: fields{iterated: 15, expired: 5},
			want:   assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			var counter counter // nolint: vetshadow
			counter.iterated = data.fields.iterated
			counter.expired = data.fields.expired

			got := counter.stopClean()

			data.want(test, got)
		})
	}
}
