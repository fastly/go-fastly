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
		start = "2026-08-12T15:00:00Z"
		end   = "2026-08-13T15:00:00Z"
	)
	const limit = 10

	var (
		result *LogInsightsResponse
		err    error
	)
	Record(t, "observability_log_insights/get", func(c *Client) {
		result, err = c.GetLogInsights(context.TODO(), &GetLogInsightsInput{
			ServiceID:     TestDeliveryServiceID,
			Start:         start,
			End:           end,
			Limit:         ToPointer(limit),
			Visualization: LogInsightsVisualizationTopURLByRequests,
		})
	})
	require.NoError(err)
	require.NotNil(result)
	require.NotNil(result.Meta)
	require.NotNil(result.Meta.Filters)
	require.NotNil(result.Meta.Filters.ServiceID)
	require.NotNil(result.Meta.Filters.Start)
	require.NotNil(result.Meta.Filters.End)
	require.NotNil(result.Meta.Filters.DomainExactMatch)
	require.NotNil(result.Meta.Filters.Limit)

	assert.Equal(TestDeliveryServiceID, *result.Meta.Filters.ServiceID)
	assert.Equal(limit, *result.Meta.Filters.Limit)
	assert.True(*result.Meta.Filters.DomainExactMatch)
	assertRFC3339TimeEqual(t, start, *result.Meta.Filters.Start)
	assertRFC3339TimeEqual(t, end, *result.Meta.Filters.End)

	require.NotNil(result.Data)
	assert.Empty(result.Data)
}
