package tsigkeys

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v14/fastly"
)

func TestTSIGKeys(t *testing.T) {
	ctx := context.TODO()

	var err error

	// Create TSIG key.
	var key *TSIGKey
	fastly.Record(t, "create_tsig_key", func(c *fastly.Client) {
		key, err = Create(ctx, c, &CreateInput{
			Name:        fastly.ToPointer("go-fastly-test-key"),
			Algorithm:   fastly.ToPointer("hmac-sha256"),
			Secret:      fastly.ToPointer("dGVzdHNlY3JldA=="),
			Description: fastly.ToPointer("go-fastly test TSIG key"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, key)
	require.NotNil(t, key.ID)
	require.Equal(t, "go-fastly-test-key", *key.Name)
	require.Equal(t, "hmac-sha256", *key.Algorithm)
	require.Equal(t, "go-fastly test TSIG key", *key.Description)

	defer func() {
		fastly.Record(t, "delete_tsig_key", func(c *fastly.Client) {
			err = Delete(ctx, c, &DeleteInput{
				TSIGKeyID: key.ID,
			})
		})
		require.NoError(t, err)
	}()

	// Get TSIG key.
	var got *TSIGKey
	fastly.Record(t, "get_tsig_key", func(c *fastly.Client) {
		got, err = Get(ctx, c, &GetInput{
			TSIGKeyID: key.ID,
		})
	})
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, *key.ID, *got.ID)
	require.Equal(t, *key.Name, *got.Name)

	// Update TSIG key.
	var updated *TSIGKey
	fastly.Record(t, "update_tsig_key", func(c *fastly.Client) {
		updated, err = Update(ctx, c, &UpdateInput{
			TSIGKeyID:   key.ID,
			Description: fastly.NewNullable("updated description"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, *key.ID, *updated.ID)
	require.Equal(t, "updated description", *updated.Description)

	// Update TSIG key, unsetting the description.
	fastly.Record(t, "update_tsig_key_null_description", func(c *fastly.Client) {
		updated, err = Update(ctx, c, &UpdateInput{
			TSIGKeyID:   key.ID,
			Description: fastly.NullValue[string](),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Nil(t, updated.Description)

	// Create a second TSIG key for pagination testing.
	var key2 *TSIGKey
	fastly.Record(t, "create_tsig_key_2", func(c *fastly.Client) {
		key2, err = Create(ctx, c, &CreateInput{
			Name:      fastly.ToPointer("go-fastly-test-key-2"),
			Algorithm: fastly.ToPointer("hmac-sha256"),
			Secret:    fastly.ToPointer("dGVzdHNlY3JldA=="),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, key2)

	defer func() {
		fastly.Record(t, "delete_tsig_key_2", func(c *fastly.Client) {
			err = Delete(ctx, c, &DeleteInput{
				TSIGKeyID: key2.ID,
			})
		})
		require.NoError(t, err)
	}()

	// List TSIG keys — should return both keys via auto-pagination.
	var keys []TSIGKey
	fastly.Record(t, "list_tsig_keys", func(c *fastly.Client) {
		keys, err = List(ctx, c, &ListInput{})
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(keys), 2, "expected at least both test keys")
}

func TestTSIGKeys_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := Get(ctx, fastly.TestClient, &GetInput{TSIGKeyID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingID)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Algorithm: fastly.ToPointer("hmac-sha256"),
		Secret:    fastly.ToPointer("dGVzdHNlY3JldA=="),
	})
	require.ErrorIs(t, err, fastly.ErrMissingName)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name:   fastly.ToPointer("go-fastly-test-key"),
		Secret: fastly.ToPointer("dGVzdHNlY3JldA=="),
	})
	require.ErrorIs(t, err, fastly.ErrMissingAlgorithm)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name:      fastly.ToPointer("go-fastly-test-key"),
		Algorithm: fastly.ToPointer("hmac-sha256"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingSecret)

	_, err = Update(ctx, fastly.TestClient, &UpdateInput{TSIGKeyID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingID)

	err = Delete(ctx, fastly.TestClient, &DeleteInput{TSIGKeyID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingID)
}
