package eventtypes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v17/fastly"
)

func TestClient_EventTypes(t *testing.T) {
	ctx := context.TODO()

	var err error

	var types *Collection
	fastly.Record(t, "list", func(c *fastly.Client) {
		types, err = List(ctx, c, &ListInput{
			ScopeType: fastly.ToPointer("account"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, types)
	require.NotEmpty(t, types.Data)
	for _, et := range types.Data {
		require.NotEmpty(t, et.EventType)
		require.NotEmpty(t, et.DisplayName)
		require.NotEmpty(t, et.ScopeTypes)
	}
}
