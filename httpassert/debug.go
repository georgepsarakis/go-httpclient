package httpassert

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// PrintJSON is a debugging test helper that will pretty-print any JSON-marshallable value.
func PrintJSON(t *testing.T, v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	require.NoError(t, err)
	t.Log(string(b))
}
