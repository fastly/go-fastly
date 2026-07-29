package providerconnection

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v17/fastly"
)

func TestClient_ProviderConnections(t *testing.T) {
	ctx := context.TODO()

	var err error

	// Create a provider connection.
	var created *ProviderConnection
	fastly.Record(t, "create", func(c *fastly.Client) {
		created, err = Create(ctx, c, &CreateInput{
			Name:    fastly.ToPointer("go-fastly-test-connection"),
			Models:  []string{"claude-opus-4-7"},
			BaseURL: fastly.ToPointer("https://api.anthropic.com"),
			// `sk-go-fastly-test` is a placeholder and will need to be replaced with a real Anthropic API key
			// in order to re-run these tests.
			APIKey: fastly.ToPointer("sk-go-fastly-test"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotEmpty(t, created.ID)

	connID := created.ID

	defer func() {
		fastly.Record(t, "delete", func(c *fastly.Client) {
			_ = Delete(ctx, c, &DeleteInput{
				ID: fastly.ToPointer(connID),
			})
		})
	}()

	// Get the provider connection.
	var fetched *ProviderConnection
	fastly.Record(t, "get", func(c *fastly.Client) {
		fetched, err = Get(ctx, c, &GetInput{
			ID: fastly.ToPointer(connID),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, fetched)
	require.Equal(t, connID, fetched.ID)

	// List provider connections.
	var conns *ProviderConnections
	fastly.Record(t, "list", func(c *fastly.Client) {
		conns, err = List(ctx, c, &ListInput{})
	})
	require.NoError(t, err)
	require.NotNil(t, conns)

	// Update the provider connection.
	var updated *ProviderConnection
	fastly.Record(t, "update", func(c *fastly.Client) {
		updated, err = Update(ctx, c, &UpdateInput{
			ID:     fastly.ToPointer(connID),
			Models: []string{"claude-opus-4-6", "claude-sonnet-4-6"},
		})
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, connID, updated.ID)
}
