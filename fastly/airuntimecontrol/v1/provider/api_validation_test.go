package provider

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v17/fastly"
)

func TestClient_ListModels_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := ListModels(ctx, fastly.TestClient, &ListModelsInput{ProviderID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingProviderID)
}
