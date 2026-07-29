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

	var sessions *Sessions
	fastly.Record(t, "list", func(c *fastly.Client) {
		sessions, err = List(ctx, c, &ListInput{
			Limit: fastly.ToPointer(10),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, sessions)
}
