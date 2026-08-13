package fastly

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_GetLogRecords_validation(t *testing.T) {
	assert := assert.New(t)

	_, err := TestClient.GetLogRecords(context.TODO(), &GetLogRecordsInput{
		Start: "2026-08-12T15:00:00Z",
		End:   "2026-08-13T15:00:00Z",
	})
	assert.ErrorIs(err, ErrMissingServiceID)

	_, err = TestClient.GetLogRecords(context.TODO(), &GetLogRecordsInput{
		ServiceID: TestDeliveryServiceID,
		End:       "2026-08-13T15:00:00Z",
	})
	assert.ErrorIs(err, ErrMissingStart)

	_, err = TestClient.GetLogRecords(context.TODO(), &GetLogRecordsInput{
		ServiceID: TestDeliveryServiceID,
		Start:     "2026-08-12T15:00:00Z",
	})
	assert.ErrorIs(err, ErrMissingEnd)

	_, err = TestClient.GetLogRecords(context.TODO(), &GetLogRecordsInput{
		ServiceID: TestDeliveryServiceID,
		Start:     "2026-08-12T15:00:00Z",
		End:       "2026-08-13T15:00:00Z",
		Filters: []LogExplorerFilter{
			{
				Operator: LogExplorerFilterOperatorGTE,
				Value:    "0",
			},
		},
	})
	assert.ErrorIs(err, ErrMissingField)

	_, err = TestClient.GetLogRecords(context.TODO(), &GetLogRecordsInput{
		ServiceID: TestDeliveryServiceID,
		Start:     "2026-08-12T15:00:00Z",
		End:       "2026-08-13T15:00:00Z",
		Filters: []LogExplorerFilter{
			{
				Field: LogExplorerFilterFieldResponseTime,
				Value: "0",
			},
		},
	})
	assert.ErrorIs(err, ErrMissingOperator)

	_, err = TestClient.GetLogRecords(context.TODO(), &GetLogRecordsInput{
		ServiceID: TestDeliveryServiceID,
		Start:     "2026-08-12T15:00:00Z",
		End:       "2026-08-13T15:00:00Z",
		Filters: []LogExplorerFilter{
			{
				Field:    LogExplorerFilterFieldResponseTime,
				Operator: LogExplorerFilterOperatorGTE,
				Value:    "0",
			},
			{
				Field:    LogExplorerFilterFieldResponseTime,
				Operator: LogExplorerFilterOperatorGTE,
				Value:    "1",
			},
		},
	})
	assert.ErrorIs(err, ErrDuplicateLogExplorerFilter)
}

func TestClient_GetLogRecords(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	t.Parallel()

	// NOTE: Update these to a recent time range when regenerating the test
	// fixture, otherwise the data may be outside Log Explorer retention.
	const (
		start = "2026-08-12T15:00:00Z"
		end   = "2026-08-13T15:00:00Z"
	)
	const limit = 5

	var (
		result *LogRecordsResponse
		err    error
	)
	Record(t, "observability_log_explorer/get", func(c *Client) {
		result, err = c.GetLogRecords(context.TODO(), &GetLogRecordsInput{
			ServiceID: TestDeliveryServiceID,
			Start:     start,
			End:       end,
			Limit:     ToPointer(limit),
			Filters: []LogExplorerFilter{
				{
					Field:    LogExplorerFilterFieldResponseTime,
					Operator: LogExplorerFilterOperatorGTE,
					Value:    "0",
				},
			},
		})
	})
	require.NoError(err)
	require.NotNil(result)
	require.NotNil(result.Meta)
	require.NotNil(result.Meta.Filters)

	filters := result.Meta.Filters
	require.NotNil(filters.ServiceID)
	require.NotNil(filters.Start)
	require.NotNil(filters.End)
	require.NotNil(filters.Limit)

	assert.Equal(TestDeliveryServiceID, *filters.ServiceID)
	assert.Equal(limit, *filters.Limit)
	assertRFC3339TimeEqual(t, start, *filters.Start)
	assertRFC3339TimeEqual(t, end, *filters.End)

	require.Len(filters.FieldFilters, 1)
	require.NotNil(filters.FieldFilters[0])
	require.NotNil(filters.FieldFilters[0].Field)
	require.NotNil(filters.FieldFilters[0].Operator)
	assert.Equal(LogExplorerFilterFieldResponseTime, *filters.FieldFilters[0].Field)
	assert.Equal(LogExplorerFilterOperatorGTE, *filters.FieldFilters[0].Operator)
	assert.Equal(float64(0), filters.FieldFilters[0].Value)

	assert.Nil(result.Meta.NextCursor)
	assert.Nil(result.Data)
}

func assertRFC3339TimeEqual(t *testing.T, want, got string) {
	t.Helper()

	require := require.New(t)
	assert := assert.New(t)

	wantTime, err := time.Parse(time.RFC3339Nano, want)
	require.NoError(err)

	gotTime, err := time.Parse(time.RFC3339Nano, got)
	require.NoError(err)

	assert.True(gotTime.Equal(wantTime), "got time %q, want %q", got, want)
}
