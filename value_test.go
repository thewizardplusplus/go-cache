package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValue_IsExpired(test *testing.T) {
	type fields struct {
		Data           interface{}
		ExpirationTime time.Time
	}
	type args struct {
		clock Clock
	}

	for _, data := range []struct {
		name   string
		fields fields
		args   args
		want   assert.BoolAssertionFunc
	}{
		{
			name: "zero expiration time",
			fields: fields{
				Data:           "data",
				ExpirationTime: time.Time{},
			},
			args: args{
				clock: clock,
			},
			want: assert.False,
		},
		{
			name: "expiration time less than the current one",
			fields: fields{
				Data:           "data",
				ExpirationTime: clock().Add(-time.Second),
			},
			args: args{
				clock: clock,
			},
			want: assert.True,
		},
		{
			name: "expiration time equal to the current one",
			fields: fields{
				Data:           "data",
				ExpirationTime: clock(),
			},
			args: args{
				clock: clock,
			},
			want: assert.False,
		},
		{
			name: "expiration time greater than the current one",
			fields: fields{
				Data:           "data",
				ExpirationTime: clock().Add(time.Second),
			},
			args: args{
				clock: clock,
			},
			want: assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			value := Value{data.fields.Data, data.fields.ExpirationTime}
			got := value.IsExpired(data.args.clock)

			data.want(test, got)
		})
	}
}
