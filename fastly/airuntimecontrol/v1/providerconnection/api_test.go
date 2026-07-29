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

func TestClient_Create_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := Create(ctx, fastly.TestClient, &CreateInput{
		Models:  []string{"gpt-4"},
		BaseURL: fastly.ToPointer("https://api.openai.com/v1"),
		APIKey:  fastly.ToPointer("secret"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingName)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name:    fastly.ToPointer("OpenAI"),
		BaseURL: fastly.ToPointer("https://api.openai.com/v1"),
		APIKey:  fastly.ToPointer("secret"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingModels)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name:   fastly.ToPointer("OpenAI"),
		Models: []string{"gpt-4"},
		APIKey: fastly.ToPointer("secret"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingBaseURL)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name:    fastly.ToPointer("OpenAI"),
		Models:  []string{"gpt-4"},
		BaseURL: fastly.ToPointer("https://api.openai.com/v1"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingAPIKey)
}

func TestClient_Get_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := Get(ctx, fastly.TestClient, &GetInput{ID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingID)
}

func TestClient_Update_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := Update(ctx, fastly.TestClient, &UpdateInput{ID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingID)
}

func TestClient_Delete_validation(t *testing.T) {
	ctx := context.TODO()

	err := Delete(ctx, fastly.TestClient, &DeleteInput{ID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingID)
}
