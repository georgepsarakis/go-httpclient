package httpassert

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseEqual(t *testing.T) {
	type args struct {
		t        *testing.T
		actual   *http.Response
		expected *http.Response
	}
	tests := []struct {
		name string
		args args
		want assert.BoolAssertionFunc
	}{
		{
			name: "status code does not match",
			args: args{
				t: &testing.T{},
				actual: &http.Response{
					StatusCode: http.StatusBadRequest,
				},
				expected: &http.Response{
					StatusCode: http.StatusOK,
				},
			},
			want: assert.True,
		},
		{
			name: "status code matches",
			args: args{
				t: &testing.T{},
				actual: &http.Response{
					StatusCode: http.StatusBadRequest,
				},
				expected: &http.Response{
					StatusCode: http.StatusBadRequest,
				},
			},
			want: assert.False,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ResponseEqual(tt.args.t, tt.args.actual, tt.args.expected)
			tt.want(t, tt.args.t.Failed())
		})
	}
}
