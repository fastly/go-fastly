package dnszones

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v14/fastly"
	"github.com/fastly/go-fastly/v14/fastly/dns/v1/tsigkeys"
)

func TestZones(t *testing.T) {
	ctx := context.TODO()

	var err error

	// Create a TSIG key to associate with the zone.
	var key *tsigkeys.TSIGKey
	fastly.Record(t, "create_tsig_key", func(c *fastly.Client) {
		key, err = tsigkeys.Create(ctx, c, &tsigkeys.CreateInput{
			Name:      fastly.ToPointer("go-fastly-test-key"),
			Algorithm: fastly.ToPointer("hmac-sha256"),
			Secret:    fastly.ToPointer("dGVzdHNlY3JldA=="),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, key)
	require.NotNil(t, key.ID)

	defer func() {
		fastly.Record(t, "delete_tsig_key", func(c *fastly.Client) {
			_ = tsigkeys.Delete(ctx, c, &tsigkeys.DeleteInput{
				TSIGKeyID: key.ID,
			})
		})
	}()

	// Create dns zone with all fields populated.
	var zone *Zone
	fastly.Record(t, "create_zone", func(c *fastly.Client) {
		zone, err = Create(ctx, c, &CreateInput{
			Name:        fastly.ToPointer("go-fastly-test.com"),
			Type:        fastly.ToPointer("secondary"),
			Description: fastly.ToPointer("go-fastly test zone"),
			XfrConfigInbound: &XfrConfigInboundInput{
				Primaries: []Primary{
					{
						Address:     fastly.ToPointer("1.2.3.4"),
						Description: fastly.ToPointer("primary DNS server"),
					},
				},
				InboundTSIGKeyID: fastly.NewNullable(*key.ID),
			},
		})
	})
	require.NoError(t, err)
	require.NotNil(t, zone)
	require.NotNil(t, zone.ID)
	// The API normalizes zone names to fully-qualified domain names (trailing dot).
	require.Equal(t, "go-fastly-test.com.", *zone.Name)
	require.Equal(t, "secondary", *zone.Type)
	require.Equal(t, "go-fastly test zone", *zone.Description)
	require.NotNil(t, zone.XfrConfigInbound)
	require.Equal(t, *key.ID, *zone.XfrConfigInbound.InboundTSIGKeyID)

	defer func() {
		fastly.Record(t, "delete_zone", func(c *fastly.Client) {
			_ = Delete(ctx, c, &DeleteInput{
				ZoneID: zone.ID,
			})
		})
	}()

	// Get dns zone.
	var got *Zone
	fastly.Record(t, "get_zone", func(c *fastly.Client) {
		got, err = Get(ctx, c, &GetInput{
			ZoneID: zone.ID,
		})
	})
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, *zone.ID, *got.ID)
	require.Equal(t, *zone.Name, *got.Name)

	// Update zone description.
	var updated *Zone
	fastly.Record(t, "update_zone", func(c *fastly.Client) {
		updated, err = Update(ctx, c, &UpdateInput{
			ZoneID:      zone.ID,
			Description: fastly.ToPointer("updated description"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, *zone.ID, *updated.ID)
	require.Equal(t, "updated description", *updated.Description)

	// Update zone, unsetting the inbound TSIG key ID.
	fastly.Record(t, "update_zone_null_tsig_key", func(c *fastly.Client) {
		updated, err = Update(ctx, c, &UpdateInput{
			ZoneID: zone.ID,
			XfrConfigInbound: &XfrConfigInboundInput{
				InboundTSIGKeyID: fastly.NullValue[string](),
			},
		})
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Nil(t, updated.XfrConfigInbound.InboundTSIGKeyID)

	// Create a second dns zone for pagination testing.
	var zone2 *Zone
	fastly.Record(t, "create_zone_2", func(c *fastly.Client) {
		zone2, err = Create(ctx, c, &CreateInput{
			Name:        fastly.ToPointer("go-fastly-test-2.com"),
			Type:        fastly.ToPointer("secondary"),
			Description: fastly.ToPointer("go-fastly test zone 2"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, zone2)

	defer func() {
		fastly.Record(t, "delete_zone_2", func(c *fastly.Client) {
			_ = Delete(ctx, c, &DeleteInput{
				ZoneID: zone2.ID,
			})
		})
	}()

	// List dns zones — page 1.
	var page1 *Zones
	fastly.Record(t, "list_zones_page_1", func(c *fastly.Client) {
		page1, err = List(ctx, c, &ListInput{
			Limit: fastly.ToPointer(1),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, page1)
	require.Len(t, page1.Data, 1)
	require.NotNil(t, page1.Meta.NextCursor, "expected a next_cursor for page 2")

	// List dns zones — page 2, using cursor from page 1.
	var page2 *Zones
	fastly.Record(t, "list_zones_page_2", func(c *fastly.Client) {
		page2, err = List(ctx, c, &ListInput{
			Limit:  fastly.ToPointer(1),
			Cursor: page1.Meta.NextCursor,
		})
	})
	require.NoError(t, err)
	require.NotNil(t, page2)
	require.Len(t, page2.Data, 1)
	require.NotEqual(t, *page1.Data[0].ID, *page2.Data[0].ID, "pages should return different zones")
}

func TestZones_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := Get(ctx, fastly.TestClient, &GetInput{ZoneID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingID)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Type: fastly.ToPointer("primary"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingName)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name: fastly.ToPointer("example.com"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingType)

	_, err = Update(ctx, fastly.TestClient, &UpdateInput{ZoneID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingID)

	err = Delete(ctx, fastly.TestClient, &DeleteInput{ZoneID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingID)
}
