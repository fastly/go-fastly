package provider

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v17/fastly"
)

func TestClient_Providers(t *testing.T) {
	ctx := context.TODO()

	var err error

	var providers *Providers
	fastly.Record(t, "list", func(c *fastly.Client) {
		providers, err = List(ctx, c)
	})
	require.NoError(t, err)
	require.NotNil(t, providers)
	require.NotEmpty(t, providers.Data)

	providerID := providers.Data[0].ID

	var models *Models
	fastly.Record(t, "list_models", func(c *fastly.Client) {
		models, err = ListModels(ctx, c, &ListModelsInput{
			ProviderID: fastly.ToPointer(providerID),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, models)
}

func TestClient_ListModels_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := ListModels(ctx, fastly.TestClient, &ListModelsInput{ProviderID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingProviderID)
}
