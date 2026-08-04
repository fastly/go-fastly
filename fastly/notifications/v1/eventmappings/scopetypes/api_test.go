package scopetypes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v17/fastly"
)

func TestClient_ScopeTypes(t *testing.T) {
	ctx := context.TODO()

	var err error

	var types *Collection
	fastly.Record(t, "list", func(c *fastly.Client) {
		types, err = List(ctx, c, &ListInput{})
	})
	require.NoError(t, err)
	require.NotNil(t, types)
	require.NotEmpty(t, types.Data)
	for _, st := range types.Data {
		require.NotEmpty(t, st.ScopeType)
	}
}
