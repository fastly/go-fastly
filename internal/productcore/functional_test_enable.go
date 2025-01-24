package productcore

import (
	"strings"
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/fastly/products"
	"github.com/fastly/go-fastly/v9/internal/test_utils"
)

// EnableTestInput specifies the information needed for the
// NewEnableTest constructor to construct a FunctionalTest object.
//
// Because Enable operations accept input and produce output, this
// struct has two type parameters used to specify the types of the
// input and output. The output type parameter is constrained to match
// the ProductOutput interface so that the test case can validate the
// common portions of the output.
type EnableTestInput[O products.ProductOutput, I any] struct {
	// Phase is used to distinguish between multiple Enable test
	// cases in a sequence of test cases; it is included in the
	// test case's Name and Operation fields
	Phase string
	// OpNoInputFn is the function to be invoked to perform the
	// enable operation, when no input is required (I == NulInput)
	OpNoInputFn func(*fastly.Client, string) (O, error)
	// OpWithInputFn is the function to be invoked to perform the
	// enable operation, when input is required (I != NulInput)
	OpWithInputFn func(*fastly.Client, string, I) (O, error)
	// Input is the input to be provided to OpWithInputFn
	Input I
	// ProductID identifies the product for which information
	// should be obtained on the service; note that it is only
	// used to validate the ProductID in the output from OpFn (if
	// any), it is not provided to OpNoInputFn or OpWithInputFn
	ProductID string
	// ServiceID identifies the service on which the product
	// information should be obtained
	ServiceID string
	// ExpectFailure specifies whether this test case is expected
	// to fail
	ExpectFailure bool
	// IgnoreFailure specifies that errors returned from
	// OpNoInputFn/OpWithInputFn should be ignored
	IgnoreFailure bool
	// CheckOutputFn specifies a function whch will be invoked if
	// OpNoInputFn/OpWithInputFn return normally; it can be used
	// to perform validation of the contents of the output
	CheckOutputFn func(*testing.T, *test_utils.FunctionalTest, O)
}

// NewEnableTest constructs a FunctionalTest object as specified by its
// input.
//
// This function requires the same input type parameter as the
// EnableTestInput struct, and that type is used to determine whether
// OpNoInputFn or OpWithInputFn are called to perform the test.
//
// This function requires the same output type parameter as the
// EnableTestInput struct, and that type is used to construct,
// populate, and validate the output present in the response body.
func NewEnableTest[O products.ProductOutput, I any](i *EnableTestInput[O, I]) *test_utils.FunctionalTest {
	r := test_utils.FunctionalTest{}

	if i.Phase != "" {
		r.Name = "enable " + i.Phase
		r.Operation = "enable-" + strings.ReplaceAll(i.Phase, " ", "-")
	} else {
		r.Name = "enable"
		r.Operation = "enable"
	}

	switch any(i.Input).(type) {
	case products.NullInput:
		r.TestFn = func(t *testing.T, tc *test_utils.FunctionalTest, c *fastly.Client) error {
			result, err := i.OpNoInputFn(c, i.ServiceID)
			if err == nil {
				validateOutput(t, tc, result, i.ProductID, i.ServiceID)
				if i.CheckOutputFn != nil {
					i.CheckOutputFn(t, tc, result)
				}
			}
			return err
		}
	default:
		r.TestFn = func(t *testing.T, tc *test_utils.FunctionalTest, c *fastly.Client) error {
			result, err := i.OpWithInputFn(c, i.ServiceID, i.Input)
			if err == nil {
				validateOutput(t, tc, result, i.ProductID, i.ServiceID)
				if i.CheckOutputFn != nil {
					i.CheckOutputFn(t, tc, result)
				}
			}
			return err
		}
	}

	if i.ExpectFailure {
		// FIXME unclear what an 'expected' failure would be
		// for this operation
		r.WantNoError = true
	} else if !i.IgnoreFailure {
		r.WantNoError = true
	}

	return &r
}
