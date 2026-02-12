package test_utils //nolint: revive

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v13/fastly"
)

type FunctionalTest struct {
	// CheckErrorFn is a function which can be used to perform
	// additional validation of the error returned by TestFn; the
	// function will be passed a formatted string containing the
	// test case's name, along with the returned error
	CheckErrorFn func(*testing.T, string, error)
	// Name is the name of the test case, used in error messages
	Name string
	// Operation is the operation being performed, which will be
	// used to construct the recorded fixture file's name
	Operation string
	// TestFn is the function to be called to run the test
	TestFn func(*testing.T, *FunctionalTest, *fastly.Client, string) error
	// WantNoError indicates that the test case should fail if any
	// error is returned by TestFn
	WantNoError bool
	// WantError indicates the error (value, not type) that must
	// be returned by TestFn for this test case
	WantError error
}

func ExecuteFunctionalTests(t *testing.T, tests []*FunctionalTest, serviceID string) {
	t.Parallel()

	var err error

	var serviceType string

	if serviceID == fastly.TestDeliveryServiceID {
		serviceType = "Delivery"
	}

	if serviceID == fastly.TestComputeServiceID {
		serviceType = "Compute"
	}

	for _, tc := range tests {
		fastly.Record(t, tc.Operation+"_"+serviceType, func(c *fastly.Client) {
			err = tc.TestFn(t, tc, c, serviceID)
		})
		if tc.WantNoError {
			require.NoErrorf(t, err, "test '%s'", tc.Name)
		}
		if tc.WantError != nil {
			require.ErrorIsf(t, err, tc.WantError, "test '%s'", tc.Name)
		}
		if tc.CheckErrorFn != nil {
			tc.CheckErrorFn(t, fmt.Sprintf("test '%s", tc.Name), err)
		}
	}
}
