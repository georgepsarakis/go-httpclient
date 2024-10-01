package httpassert

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func PrintJSON(t *testing.T, v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	require.NoError(t, err)
	t.Log(string(b))
}
