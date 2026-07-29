package providerconnection

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v17/fastly"
)

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
