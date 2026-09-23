package eventmappings

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v17/fastly"
)

func TestClient_EventMappings(t *testing.T) {
	ctx := context.TODO()

	var err error

	// Create an event mapping.
	var created *EventMapping
	fastly.Record(t, "create", func(c *fastly.Client) {
		created, err = Create(ctx, c, &CreateInput{
			Name:           new("go-fastly-test-mapping"),
			Description:    new("Sends a notification when any user logs in"),
			ScopeType:      new(ScopeTypeAccount),
			EventTypes:     []string{"user.login"},
			IntegrationIDs: []string{"7znp3LzS0yF0jiNU9FuQxW"},
		})
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotEmpty(t, created.ID)
	require.Equal(t, MappingStatusActive, created.MappingStatus)

	mappingID := created.ID

	defer func() {
		fastly.Record(t, "delete", func(c *fastly.Client) {
			_ = Delete(ctx, c, &DeleteInput{
				MappingID: new(mappingID),
			})
		})
	}()

	// Get the event mapping.
	var fetched *EventMapping
	fastly.Record(t, "get", func(c *fastly.Client) {
		fetched, err = Get(ctx, c, &GetInput{
			MappingID: new(mappingID),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, fetched)
	require.Equal(t, mappingID, fetched.ID)

	// List event mappings.
	var mappings []EventMapping
	fastly.Record(t, "list", func(c *fastly.Client) {
		mappings, err = List(ctx, c, &ListInput{
			ScopeType: new(ScopeTypeAccount),
		})
	})
	require.NoError(t, err)
	require.NotEmpty(t, mappings)

	// Update the event mapping.
	var updated *EventMapping
	fastly.Record(t, "update", func(c *fastly.Client) {
		updated, err = Update(ctx, c, &UpdateInput{
			MappingID:      new(mappingID),
			Name:           new("go-fastly-test-mapping-updated"),
			Description:    new("Updated description"),
			ScopeType:      new(ScopeTypeAccount),
			EventTypes:     []string{"user.login", "user.create"},
			IntegrationIDs: []string{"7znp3LzS0yF0jiNU9FuQxW"},
		})
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, mappingID, updated.ID)
	require.Equal(t, MappingStatusActive, updated.MappingStatus)
}

func TestClient_Create_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := Create(ctx, fastly.TestClient, &CreateInput{
		ScopeType:      new(ScopeTypeAccount),
		EventTypes:     []string{"user.login"},
		IntegrationIDs: []string{"integration-1"},
	})
	require.ErrorIs(t, err, fastly.ErrMissingName)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name:           new("my-mapping"),
		EventTypes:     []string{"user.login"},
		IntegrationIDs: []string{"integration-1"},
	})
	require.ErrorIs(t, err, fastly.ErrMissingScopeType)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name:           new("my-mapping"),
		ScopeType:      new(ScopeTypeAccount),
		IntegrationIDs: []string{"integration-1"},
	})
	require.ErrorIs(t, err, fastly.ErrMissingEventTypes)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name:       new("my-mapping"),
		ScopeType:  new(ScopeTypeAccount),
		EventTypes: []string{"user.login"},
	})
	require.ErrorIs(t, err, fastly.ErrMissingIntegrationIDs)
}

func TestClient_Get_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := Get(ctx, fastly.TestClient, &GetInput{MappingID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingMappingID)
}

func TestClient_Update_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := Update(ctx, fastly.TestClient, &UpdateInput{
		ScopeType:      new(ScopeTypeAccount),
		EventTypes:     []string{"user.login"},
		IntegrationIDs: []string{"integration-1"},
	})
	require.ErrorIs(t, err, fastly.ErrMissingMappingID)

	_, err = Update(ctx, fastly.TestClient, &UpdateInput{
		MappingID:      new("mapping-1"),
		EventTypes:     []string{"user.login"},
		IntegrationIDs: []string{"integration-1"},
	})
	require.ErrorIs(t, err, fastly.ErrMissingName)

	_, err = Update(ctx, fastly.TestClient, &UpdateInput{
		MappingID:      new("mapping-1"),
		Name:           new("my-mapping"),
		IntegrationIDs: []string{"integration-1"},
	})
	require.ErrorIs(t, err, fastly.ErrMissingScopeType)

	_, err = Update(ctx, fastly.TestClient, &UpdateInput{
		MappingID:      new("mapping-1"),
		Name:           new("my-mapping"),
		ScopeType:      new(ScopeTypeAccount),
		IntegrationIDs: []string{"integration-1"},
	})
	require.ErrorIs(t, err, fastly.ErrMissingEventTypes)

	_, err = Update(ctx, fastly.TestClient, &UpdateInput{
		MappingID:  new("mapping-1"),
		Name:       new("my-mapping"),
		ScopeType:  new(ScopeTypeAccount),
		EventTypes: []string{"user.login"},
	})
	require.ErrorIs(t, err, fastly.ErrMissingIntegrationIDs)
}

func TestClient_Delete_validation(t *testing.T) {
	ctx := context.TODO()

	err := Delete(ctx, fastly.TestClient, &DeleteInput{MappingID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingMappingID)
}
