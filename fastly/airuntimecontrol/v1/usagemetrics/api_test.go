package usagemetrics

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v18/fastly"
)

func TestClient_UsageMetrics(t *testing.T) {
	ctx := context.TODO()

	var err error

	var metrics []UsageMetric
	fastly.Record(t, "list", func(c *fastly.Client) {
		metrics, err = List(ctx, c, &ListInput{
			Limit: new(1),
		})
	})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(metrics), 2, "expected at least two pages worth of usage metrics")

	var csv []byte
	fastly.Record(t, "export", func(c *fastly.Client) {
		csv, err = Export(ctx, c, &ListInput{
			Sort: new("-date"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, csv)
}
