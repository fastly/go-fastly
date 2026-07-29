package key

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v17/fastly"
)

func TestClient_Create_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := Create(ctx, fastly.TestClient, &CreateInput{
		Model:    fastly.ToPointer("opus-4.6"),
		Provider: fastly.ToPointer("Anthropic"),
		UserID:   fastly.ToPointer("user-1"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingName)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name:     fastly.ToPointer("my-key"),
		Provider: fastly.ToPointer("Anthropic"),
		UserID:   fastly.ToPointer("user-1"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingModel)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name:   fastly.ToPointer("my-key"),
		Model:  fastly.ToPointer("opus-4.6"),
		UserID: fastly.ToPointer("user-1"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingProvider)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name:     fastly.ToPointer("my-key"),
		Model:    fastly.ToPointer("opus-4.6"),
		Provider: fastly.ToPointer("Anthropic"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingUserID)
}

func TestClient_Get_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := Get(ctx, fastly.TestClient, &GetInput{KeyID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingKeyID)
}

func TestClient_Update_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := Update(ctx, fastly.TestClient, &UpdateInput{KeyID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingKeyID)
}

func TestClient_Delete_validation(t *testing.T) {
	ctx := context.TODO()

	err := Delete(ctx, fastly.TestClient, &DeleteInput{KeyID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingKeyID)
}

func TestClient_Rotate_validation(t *testing.T) {
	ctx := context.TODO()

	expiresAt := time.Now()
	_, err := Rotate(ctx, fastly.TestClient, &RotateInput{
		KeyID:     nil,
		ExpiresAt: &expiresAt,
	})
	require.ErrorIs(t, err, fastly.ErrMissingKeyID)

	_, err = Rotate(ctx, fastly.TestClient, &RotateInput{
		KeyID:     fastly.ToPointer("key-1"),
		ExpiresAt: nil,
	})
	require.ErrorIs(t, err, fastly.ErrMissingExpiresAt)
}
