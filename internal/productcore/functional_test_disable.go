package productcore

import (
	"strings"
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/fastly/go-fastly/v9/internal/test_utils"
)

// DisableTestInput specifies the information needed for the
// NewDisableTest constructor to construct a FunctionalTest object.
type DisableTestInput struct {
	// Phase is used to distinguish between multiple Disable test
	// cases in a sequence of test cases; it is included in the
	// test case's Name and Operation fields
	Phase string
	// OpFn is the function to be invoked to perform the
	// disablement
	OpFn func(*fastly.Client, string) error
	// ServiceID identifies the service on which the product
	// should be disabled
	ServiceID string
	// ExpectFailure specifies whether this test case is expected
	// to fail
	ExpectFailure bool
	// IgnoreFailure specifies that errors returned from OpFn
	// should be ignored
	IgnoreFailure bool
}

// NewDisableTest constructs a FunctionalTest object as specified by
// its input.
func NewDisableTest(i *DisableTestInput) *test_utils.FunctionalTest {
	r := test_utils.FunctionalTest{}

	if i.Phase != "" {
		r.Name = "disable " + i.Phase
		r.Operation = "disable-" + strings.ReplaceAll(i.Phase, " ", "-")
	} else {
		r.Name = "disable"
		r.Operation = "disable"
	}

	r.TestFn = func(_ *testing.T, _ *test_utils.FunctionalTest, c *fastly.Client) error {
		err := i.OpFn(c, i.ServiceID)
		return err
	}

	if i.ExpectFailure {
		// FIXME: not clear what to do here, as the returned
		// error may not be consistent but is probably 404
		r.WantNoError = true
	} else if !i.IgnoreFailure {
		r.WantNoError = true
	}

	return &r
}
