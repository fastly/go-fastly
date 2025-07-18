package productcore

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v10/fastly"
	"github.com/fastly/go-fastly/v10/fastly/products"
	"github.com/fastly/go-fastly/v10/internal/test_utils"
)

// GetTestInput specifies the information needed for the NewGetTest
// constructor to construct a FunctionalTest object.
//
// Because Get operations produce output, this struct has a type
// parameter used to specify the type of the output. The type
// parameter is constrained to match the ProductOutput interface so
// that the test case can validate the common portions of the output.
type GetTestInput[O products.ProductOutput] struct {
	// Phase is used to distinguish between multiple Get test
	// cases in a sequence of test cases; it is included in the
	// test case's Name and Operation fields
	Phase string
	// OpFn is the function to be invoked to perform the
	// get operation
	OpFn func(context.Context, *fastly.Client, string) (O, error)
	// ProductID identifies the product for which information
	// should be obtained on the service; note that it is only
	// used to validate the ProductID in the output from OpFn (if
	// any), it is not provided to OpFn
	ProductID string
	// ServiceID identifies the service on which the product
	// information should be obtained
	ServiceID string
	// ExpectFailure specifies whether this test case is expected
	// to fail
	ExpectFailure bool
	// IgnoreFailure specifies that errors returned from OpFn
	// should be ignored
	IgnoreFailure bool
	// CheckOutputFn specifies a function whch will be invoked if
	// the OpFn returns normally; it can be used to perform
	// validation of the contents of the output
	CheckOutputFn func(*testing.T, *test_utils.FunctionalTest, O)
}

// NewGetTest constructs a FunctionalTest object as specified by its
// input.
//
// This function requires the same type parameter as the GetTestInput
// struct, and that type is used to construct, populate, and validate
// the output present in the response body.
//
// If the input indicates that failure is expected, the test case will
// ensure that the error received from the API matches the documented
// error returned when a product is not enabled on a service. If any
// other error is returned by the API, the test case will report
// failure.
func NewGetTest[O products.ProductOutput](i *GetTestInput[O]) *test_utils.FunctionalTest {
	r := test_utils.FunctionalTest{}

	if i.Phase != "" {
		r.Name = "get status " + i.Phase
		r.Operation = "get-" + strings.ReplaceAll(i.Phase, " ", "-")
	} else {
		r.Name = "get status"
		r.Operation = "get"
	}

	r.TestFn = func(t *testing.T, tc *test_utils.FunctionalTest, c *fastly.Client) error {
		result, err := i.OpFn(context.TODO(), c, i.ServiceID)
		if err == nil {
			validateOutput(t, tc, result, i.ProductID, i.ServiceID)
			if i.CheckOutputFn != nil {
				i.CheckOutputFn(t, tc, result)
			}
		}
		return err
	}

	if i.ExpectFailure {
		r.CheckErrorFn = func(t *testing.T, testName string, err error) {
			// the API returns "Bad Request" instead of
			// "Not Found" when a product has not been
			// enabled on a service
			var herr *fastly.HTTPError
			require.ErrorAs(t, err, &herr, testName)
			require.Truef(t, herr.IsBadRequest(), "%s expected HTTP 'Bad Request'", testName)
		}
	} else if !i.IgnoreFailure {
		r.WantNoError = true
	}

	return &r
}
