package fastly

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetLogInsights_validation(t *testing.T) {
	assert := assert.New(t)

	_, err := TestClient.GetLogInsights(context.TODO(), &GetLogInsightsInput{
		Start:         "2026-08-12T15:00:00Z",
		End:           "2026-08-13T15:00:00Z",
		Visualization: LogInsightsVisualizationTopURLByRequests,
	})
	assert.ErrorIs(err, ErrMissingServiceID)

	_, err = TestClient.GetLogInsights(context.TODO(), &GetLogInsightsInput{
		ServiceID:     TestDeliveryServiceID,
		End:           "2026-08-13T15:00:00Z",
		Visualization: LogInsightsVisualizationTopURLByRequests,
	})
	assert.ErrorIs(err, ErrMissingStart)

	_, err = TestClient.GetLogInsights(context.TODO(), &GetLogInsightsInput{
		ServiceID:     TestDeliveryServiceID,
		Start:         "2026-08-12T15:00:00Z",
		Visualization: LogInsightsVisualizationTopURLByRequests,
	})
	assert.ErrorIs(err, ErrMissingEnd)

	_, err = TestClient.GetLogInsights(context.TODO(), &GetLogInsightsInput{
		ServiceID: TestDeliveryServiceID,
		Start:     "2026-08-12T15:00:00Z",
		End:       "2026-08-13T15:00:00Z",
	})
	assert.ErrorIs(err, ErrMissingVisualization)
}

func TestClient_GetLogInsights(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	t.Parallel()

	// NOTE: Update these to a recent time range when regenerating the test
	// fixture, otherwise the data may be outside Log Insights retention.
	const (
		start  = "2026-08-12T15:00:00Z"
		end    = "2026-08-13T15:00:00Z"
		domain = "example.com"
	)
	const limit = 10

	var (
		result *LogInsightsResponse
		err    error
	)
	Record(t, "observability_log_insights/get", func(c *Client) {
		result, err = c.GetLogInsights(context.TODO(), &GetLogInsightsInput{
			ServiceID:        TestDeliveryServiceID,
			Start:            start,
			End:              end,
			Domain:           ToPointer(domain),
			DomainExactMatch: ToPointer(true),
			Limit:            ToPointer(limit),
			POPs:             []string{"IAD", "DFW"},
			Visualization:    LogInsightsVisualizationTopURLByRequests,
		})
	})
	require.NoError(err)
	require.NotNil(result)
	require.NotNil(result.Meta)
	require.NotNil(result.Meta.Filters)
	filters := result.Meta.Filters
	require.NotNil(filters.Domain)
	require.NotNil(filters.DomainExactMatch)
	require.NotNil(filters.End)
	require.NotNil(filters.Limit)
	require.NotNil(filters.ServiceID)
	require.NotNil(filters.Start)

	assert.Equal(domain, *filters.Domain)
	assert.True(*filters.DomainExactMatch)
	assertRFC3339TimeEqual(t, end, *filters.End)
	assert.Equal(limit, *filters.Limit)
	assert.ElementsMatch([]string{"IAD", "DFW"}, filters.POPs)
	assert.Equal(TestDeliveryServiceID, *filters.ServiceID)
	assertRFC3339TimeEqual(t, start, *filters.Start)

	require.NotNil(result.Data)
	assert.Empty(result.Data)
}
