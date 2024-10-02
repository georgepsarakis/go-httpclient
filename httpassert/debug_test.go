package httpassert

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestPrintJSON(t *testing.T) {
	t.Parallel()

	t.Run("fails if the input cannot be marshalled to JSON", func(t *testing.T) {
		grp := errgroup.Group{}
		stubTest := &testing.T{}
		grp.Go(func() error {
			PrintJSON(stubTest, func() {})
			return nil
		})
		require.NoError(t, grp.Wait())
		require.True(t, stubTest.Failed())
	})
}
