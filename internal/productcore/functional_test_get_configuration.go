package productcore

import (
	"strings"
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products"
	"github.com/fastly/go-fastly/v9/internal/test_utils"
)

// GetConfigurationTestInput specifies the information needed for the
// NewGetConfigurationTest constructor to construct a FunctionalTest
// object.
//
// Because GetConfiguration operations produce output, this struct has
// a type parameter used to specify the type of the output. The type
// parameter is constrained to match the ProductOutput interface so
// that the test case can validate the common portions of the output.
type GetConfigurationTestInput[O products.ProductOutput] struct {
	// Phase is used to distinguish between multiple
	// GetConfiguration test cases in a sequence of test cases; it
	// is included in the test case's Name and Operation fields
	Phase string
	// OpFn is the function to be invoked to perform the
	// get operation
	OpFn func(*fastly.Client, string) (O, error)
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

// NewGetConfigurationTest constructs a FunctionalTest object as
// specified by its input.
func NewGetConfigurationTest[O products.ProductOutput](i *GetConfigurationTestInput[O]) *test_utils.FunctionalTest {
	r := test_utils.FunctionalTest{}

	if i.Phase != "" {
		r.Name = "get configuration " + i.Phase
		r.Operation = "get-configuration-" + strings.ReplaceAll(i.Phase, " ", "-")
	} else {
		r.Name = "get configuration"
		r.Operation = "get-configuration"
	}

	r.TestFn = func(t *testing.T, tc *test_utils.FunctionalTest, c *fastly.Client) error {
		result, err := i.OpFn(c, i.ServiceID)
		if err == nil {
			validateOutput(t, tc, result, i.ProductID, i.ServiceID)
			if i.CheckOutputFn != nil {
				i.CheckOutputFn(t, tc, result)
			}
		}
		return err
	}

	if i.ExpectFailure {
		// FIXME need to determine the expected error here,
		// probably 404
		r.WantNoError = true
	} else if !i.IgnoreFailure {
		r.WantNoError = true
	}

	return &r
}
