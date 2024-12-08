package products

import (
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/stretchr/testify/require"
)

// validateOutput provides common validation for all responses to
// product-specific API operations.
//
// All of these operations return a JSON object containing the
// ProductID and ServiceID on which the operation was invoked; this
// function confirms that they were returned and that they have the
// expected values.
func validateOutput[O ProductOutput](t *testing.T, tc *fastly.FunctionalTest, o O, productID, serviceID string) {
	require.NotNilf(t, o, "test '%s'", tc.Name)
	require.NotNilf(t, o.ProductID(), "test '%s'", tc.Name)
	require.NotNilf(t, o.ServiceID(), "test '%s'", tc.Name)
	require.Equalf(t, productID, *o.ProductID(), "test '%s'", tc.Name)
	require.Equalf(t, serviceID, *o.ServiceID(), "test '%s'", tc.Name)
}
