package key

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v17/fastly"
)

func TestClient_Keys(t *testing.T) {
	ctx := context.TODO()

	var err error

	expiresAt := time.Now().Add(24 * time.Hour).UTC().Truncate(time.Second)

	// Create a virtual key.
	var created *VirtualKeyWithToken
	fastly.Record(t, "create", func(c *fastly.Client) {
		created, err = Create(ctx, c, &CreateInput{
			Name:      fastly.ToPointer("go-fastly-test-key"),
			Model:     fastly.ToPointer("claude-sonnet-4-20250514"),
			Provider:  fastly.ToPointer("Anthropic"),
			UserID:    fastly.ToPointer("go-fastly-test-user"),
			ExpiresAt: &expiresAt,
		})
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotEmpty(t, created.ID)
	require.NotEmpty(t, created.AccessToken)

	keyID := created.ID

	defer func() {
		fastly.Record(t, "delete", func(c *fastly.Client) {
			_ = Delete(ctx, c, &DeleteInput{
				KeyID: fastly.ToPointer(keyID),
			})
		})
	}()

	// Get the virtual key.
	var fetched *VirtualKeyListItem
	fastly.Record(t, "get", func(c *fastly.Client) {
		fetched, err = Get(ctx, c, &GetInput{
			KeyID: fastly.ToPointer(keyID),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, fetched)
	require.Equal(t, keyID, fetched.ID)

	// List virtual keys.
	var keys []VirtualKeyListItem
	fastly.Record(t, "list", func(c *fastly.Client) {
		keys, err = List(ctx, c, &ListInput{
			Provider: fastly.ToPointer("Anthropic"),
		})
	})
	require.NoError(t, err)
	require.NotEmpty(t, keys)

	// Update the virtual key.
	var updated *VirtualKey
	fastly.Record(t, "update", func(c *fastly.Client) {
		updated, err = Update(ctx, c, &UpdateInput{
			KeyID: fastly.ToPointer(keyID),
			Name:  fastly.ToPointer("go-fastly-test-key-updated"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, keyID, updated.ID)

	// Rotate the virtual key.
	newExpiresAt := time.Now().Add(48 * time.Hour).UTC().Truncate(time.Second)
	var rotated *VirtualKeyWithToken
	fastly.Record(t, "rotate", func(c *fastly.Client) {
		rotated, err = Rotate(ctx, c, &RotateInput{
			KeyID:     fastly.ToPointer(keyID),
			ExpiresAt: &newExpiresAt,
		})
	})
	require.NoError(t, err)
	require.NotNil(t, rotated)
	require.NotEmpty(t, rotated.AccessToken)
}

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
