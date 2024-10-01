package httpclient

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithQueryParameters(t *testing.T) {
	tests := []struct {
		name    string
		params  map[string]string
		params2 map[string]string
		want    url.Values
	}{
		{
			name: "set query parameters with a non-nil map",
			params: map[string]string{
				"q1": "2",
			},
			want: url.Values{
				"q1": []string{"2"},
			},
		},
		{
			name:   "query parameters are nil",
			params: nil,
			want:   url.Values{},
		},
		{
			name: "query parameters are already set",
			params: map[string]string{
				"q1": "2",
				"q2": "3",
			},
			params2: map[string]string{
				"q1": "4",
				"q4": "5",
			},
			want: url.Values{
				"q1": []string{"4"},
				"q2": []string{"3"},
				"q4": []string{"5"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qp := []RequestParameter{WithQueryParameters(tt.params)}
			if tt.params2 != nil {
				qp = append(qp, WithQueryParameters(tt.params2))
			}
			opts := NewRequestParameters(qp...)

			assert.Equal(t, tt.want, opts.QueryParameters())
		})
	}
}
