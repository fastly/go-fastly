package products

import (
	"strings"
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
	"github.com/stretchr/testify/require"
)

type TestGetInput[O any] struct {
	Phase         string
	Executor      func(*fastly.Client, string) (*O, error)
	ServiceID     string
	ExpectFailure bool
}

func TestGet[O any](i *TestGetInput[O]) *fastly.FunctionalTestInput {
	r := fastly.FunctionalTestInput{}

	if i.Phase != "" {
		r.Name = "get status " + i.Phase
		r.Operation = "get-" + strings.ReplaceAll(i.Phase, " ", "-")
	} else {
		r.Name = "get status"
		r.Operation = "get"
	}

	r.Execute = func(c *fastly.Client) error {
		_, err := i.Executor(c, i.ServiceID)
		return err
	}

	if i.ExpectFailure {
		r.CheckError = func(testName string, t *testing.T, err error) {
			// the API returns "Bad Request" instead of
			// "Not Found" when a product has not been
			// enabled on a service
			var herr *fastly.HTTPError
			require.ErrorAs(t, err, &herr, testName)
			require.Truef(t, herr.IsBadRequest(), "%s expected HTTP 'Bad Request'", testName)
		}
	} else {
		r.WantNoError = true
	}

	return &r
}
