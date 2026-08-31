package session

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v17/fastly"
)

func TestClient_Sessions(t *testing.T) {
	ctx := context.TODO()

	var err error

	var sessions []Session
	fastly.Record(t, "list", func(c *fastly.Client) {
		sessions, err = List(ctx, c, &ListInput{
			Limit: fastly.ToPointer(2),
		})
	})
	require.NoError(t, err)
	require.Len(t, sessions, 3)

	var foundSecondPageItem bool
	for _, s := range sessions {
		if s.ID == "sess_546" {
			foundSecondPageItem = true
			break
		}
	}
	require.True(t, foundSecondPageItem, "expected an item from the second page")
}
