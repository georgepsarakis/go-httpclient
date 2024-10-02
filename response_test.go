package httpclient

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeserializeJSON(t *testing.T) {
	type args struct {
		resp   *http.Response
		target any
	}
	tests := []struct {
		name           string
		args           args
		wantErrMessage string
		want           map[string]any
	}{
		{
			name: "returns error from JSON marshalling",
			args: args{
				resp: &http.Response{
					Body: io.NopCloser(strings.NewReader("{")),
				},
				target: &map[string]any{},
			},
			wantErrMessage: "unexpected end of JSON input",
		},
		{
			name: "returns error when not passing a pointer",
			args: args{
				resp: &http.Response{
					Body: io.NopCloser(strings.NewReader("{}")),
				},
				target: map[string]any{},
			},
			wantErrMessage: "pointer required, got map[string]interface {}",
		},
		{
			name: "unmarshals the JSON payload to the passed pointer",
			args: args{
				resp: &http.Response{
					Body: io.NopCloser(strings.NewReader(`{"hello": "world"}`)),
				},
				target: &map[string]any{},
			},
			want: map[string]any{"hello": "world"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DeserializeJSON(tt.args.resp, tt.args.target)
			if tt.wantErrMessage != "" {
				assert.ErrorContains(t, err, tt.wantErrMessage)
			} else {
				assert.Equal(t, tt.want, *tt.args.target.(*map[string]any))
			}
		})
	}
}
