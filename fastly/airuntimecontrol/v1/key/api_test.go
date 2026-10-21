package key

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v18/fastly"
)

// arcSecretFields are the response fields that carry live secrets or PII and must never
// be written to a fixture verbatim.
var arcSecretFields = []string{"access_token", "user_id", "customer_id", "user_name", "created_by"}

func TestClient_Keys(t *testing.T) {
	ctx := context.TODO()

	var err error

	expiresAt := time.Now().Add(24 * time.Hour).UTC().Truncate(time.Second)

	// Create a virtual key.
	var created *VirtualKeyWithToken
	fastly.RecordRedacted(t, "create", arcSecretFields, func(c *fastly.Client) {
		created, err = Create(ctx, c, &CreateInput{
			Name:            new("go-fastly-test-key"),
			Model:           new("claude-sonnet-4-20250514"),
			Provider:        new("Anthropic"),
			ExpiresAt:       &expiresAt,
			SecurityEnabled: new(true),
			SecurityAction:  new("log"),
			RPMLimit:        new(60),
			TPMLimit:        new(10000),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotEmpty(t, created.ID)
	require.NotEmpty(t, created.AccessToken)
	require.NotNil(t, created.RPMLimit)
	require.Equal(t, 60, *created.RPMLimit)
	require.NotNil(t, created.TPMLimit)
	require.Equal(t, 10000, *created.TPMLimit)

	keyID := created.ID

	defer func() {
		fastly.Record(t, "delete", func(c *fastly.Client) {
			_ = Delete(ctx, c, &DeleteInput{
				KeyID: new(keyID),
			})
		})
	}()

	// Get the virtual key.
	var fetched *VirtualKeyListItem
	fastly.RecordRedacted(t, "get", arcSecretFields, func(c *fastly.Client) {
		fetched, err = Get(ctx, c, &GetInput{
			KeyID: new(keyID),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, fetched)
	require.Equal(t, keyID, fetched.ID)
	require.NotNil(t, fetched.SecurityEnabled)
	require.True(t, *fetched.SecurityEnabled)
	require.NotNil(t, fetched.SecurityAction)
	require.Equal(t, "log", *fetched.SecurityAction)
	require.NotNil(t, fetched.RPMLimit)
	require.Equal(t, 60, *fetched.RPMLimit)

	// List virtual keys.
	var keys []VirtualKeyListItem
	fastly.RecordRedacted(t, "list", arcSecretFields, func(c *fastly.Client) {
		keys, err = List(ctx, c, &ListInput{
			Limit: new(1),
		})
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(keys), 2, "expected at least two pages worth of keys")

	var foundCreated bool
	for _, k := range keys {
		if k.ID == keyID {
			foundCreated = true
			break
		}
	}
	require.True(t, foundCreated, "expected the created key in the paginated results")

	// Update the virtual key.
	var updated *VirtualKey
	fastly.RecordRedacted(t, "update", arcSecretFields, func(c *fastly.Client) {
		updated, err = Update(ctx, c, &UpdateInput{
			KeyID:             new(keyID),
			Name:              new("go-fastly-test-key-updated"),
			RPMLimit:          fastly.NewNullable(120),
			TPMLimit:          fastly.NullValue[int](),
			SpendLimitEnabled: new(true),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, keyID, updated.ID)
	require.NotNil(t, updated.RPMLimit)
	require.Equal(t, 120, *updated.RPMLimit)
	require.Nil(t, updated.TPMLimit)
	require.True(t, updated.SpendLimitEnabled)
	require.NotNil(t, updated.SpendLimitEnabledOn)

	// Get the failover policy.
	var failover *FailoverPolicy
	fastly.Record(t, "get_failover_policy", func(c *fastly.Client) {
		failover, err = GetFailoverPolicy(ctx, c, &GetFailoverPolicyInput{
			KeyID: new(keyID),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, failover)

	// Update the failover policy.
	fastly.Record(t, "update_failover_policy", func(c *fastly.Client) {
		failover, err = UpdateFailoverPolicy(ctx, c, &UpdateFailoverPolicyInput{
			KeyID:                 new(keyID),
			TimeoutEnabled:        new(true),
			RateLimit429Enabled:   new(true),
			ServerError5xxEnabled: new(false),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, failover)
	require.True(t, failover.TimeoutEnabled)
	require.True(t, failover.RateLimit429Enabled)
	require.False(t, failover.ServerError5xxEnabled)

	// Rotate the virtual key.
	newExpiresAt := time.Now().Add(48 * time.Hour).UTC().Truncate(time.Second)
	var rotated *VirtualKeyWithToken
	fastly.RecordRedacted(t, "rotate", arcSecretFields, func(c *fastly.Client) {
		rotated, err = Rotate(ctx, c, &RotateInput{
			KeyID:     new(keyID),
			ExpiresAt: &newExpiresAt,
		})
	})
	require.NoError(t, err)
	require.NotNil(t, rotated)
	require.NotEmpty(t, rotated.AccessToken)

	// Rotate the virtual key without an expiration.
	fastly.RecordRedacted(t, "rotate_without_expires_at", arcSecretFields, func(c *fastly.Client) {
		rotated, err = Rotate(ctx, c, &RotateInput{
			KeyID: new(keyID),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, rotated)
	require.NotEmpty(t, rotated.AccessToken)
	require.Nil(t, rotated.ExpiresAt)
}

func TestClient_Create_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := Create(ctx, fastly.TestClient, &CreateInput{
		Model:    new("opus-4.6"),
		Provider: new("Anthropic"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingName)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name:     new("my-key"),
		Provider: new("Anthropic"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingModel)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name:  new("my-key"),
		Model: new("opus-4.6"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingProvider)
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

	_, err := Rotate(ctx, fastly.TestClient, &RotateInput{KeyID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingKeyID)
}

func TestClient_GetFailoverPolicy_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := GetFailoverPolicy(ctx, fastly.TestClient, &GetFailoverPolicyInput{KeyID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingKeyID)
}

func TestClient_UpdateFailoverPolicy_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := UpdateFailoverPolicy(ctx, fastly.TestClient, &UpdateFailoverPolicyInput{KeyID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingKeyID)
}
