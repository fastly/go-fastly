package usagemetrics

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v17/fastly"
)

func TestClient_UsageMetrics(t *testing.T) {
	ctx := context.TODO()

	var err error

	var metrics *UsageMetrics
	fastly.Record(t, "list", func(c *fastly.Client) {
		metrics, err = List(ctx, c, &ListInput{})
	})
	require.NoError(t, err)
	require.NotNil(t, metrics)

	var csv []byte
	fastly.Record(t, "export", func(c *fastly.Client) {
		csv, err = Export(ctx, c, &ListInput{})
	})
	require.NoError(t, err)
	require.NotNil(t, csv)
}
