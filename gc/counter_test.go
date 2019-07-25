package gc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_counter_stopIterate(test *testing.T) {
	type fields struct {
		iterated int
		expired  int
	}

	for _, data := range []struct {
		name   string
		fields fields
		want   assert.BoolAssertionFunc
	}{
		// TODO: add test cases
	} {
		test.Run(data.name, func(test *testing.T) {
			var counter counter // nolint: vetshadow
			counter.iterated = data.fields.iterated
			counter.expired = data.fields.expired

			got := counter.stopIterate()

			data.want(test, got)
		})
	}
}
