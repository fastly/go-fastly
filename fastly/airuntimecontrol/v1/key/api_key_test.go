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
	var keys *VirtualKeys
	fastly.Record(t, "list", func(c *fastly.Client) {
		keys, err = List(ctx, c, &ListInput{
			Provider: fastly.ToPointer("Anthropic"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, keys)

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
