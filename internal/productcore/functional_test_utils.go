package productcore

import (
	"testing"

	"github.com/fastly/go-fastly/v9/fastly/products"
	"github.com/fastly/go-fastly/v9/internal/test_utils"
	"github.com/stretchr/testify/require"
)

// validateOutput provides common validation for all responses to
// product-specific API operations.
//
// All of the operations return a JSON object containing the ProductID
// and ServiceID on which the operation was invoked; this function
// confirms that they were returned and that they have the expected
// values.
func validateOutput(t *testing.T, tc *test_utils.FunctionalTest, o products.ProductOutput, productID, serviceID string) {
	require.NotNilf(t, o, "test '%s'", tc.Name)
	require.Equalf(t, productID, o.ProductID(), "test '%s'", tc.Name)
	require.Equalf(t, serviceID, o.ServiceID(), "test '%s'", tc.Name)
}
