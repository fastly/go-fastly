package test_utils

import (
	"fmt"
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/stretchr/testify/require"
)

type FunctionalTest struct {
	// Name is the name of the test case, used in error messages
	Name string
	// Operation is the operation being performed, which will be
	// used to construct the recorded fixture file's name
	Operation string
	// TestFn is the function to be called to run the test
	TestFn func(*testing.T, *FunctionalTest, *fastly.Client) error
	// WantNoError indicates that the test case should fail if any
	// error is returned by TestFn
	WantNoError bool
	// WantError indicates the error (value, not type) that must
	// be returned by TestFn for this test case
	WantError error
	// CheckErrorFn is a function which can be used to perform
	// additional validation of the error returned by TestFn; the
	// function will be passed a formatted string containing the
	// test case's name, along with the returned error
	CheckErrorFn func(*testing.T, string, error)
}

func ExecuteFunctionalTests(t *testing.T, tests []*FunctionalTest) {
	t.Parallel()

	var err error

	for _, tc := range tests {
		fastly.Record(t, tc.Operation, func(c *fastly.Client) {
			err = tc.TestFn(t, tc, c)
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
