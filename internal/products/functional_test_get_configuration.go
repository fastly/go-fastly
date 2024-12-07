package products

import (
	"strings"

	"github.com/fastly/go-fastly/v9/fastly"
)

type TestGetConfigurationInput[O any] struct {
	Phase         string
	Executor      func(*fastly.Client, string) (*O, error)
	ServiceID     string
	ExpectFailure bool
}

func TestGetConfiguration[O any](i *TestGetConfigurationInput[O]) *fastly.FunctionalTestInput {
	r := fastly.FunctionalTestInput{}

	if i.Phase != "" {
		r.Name = "get configuration " + i.Phase
		r.Operation = "get-configuration-" + strings.ReplaceAll(i.Phase, " ", "-")
	} else {
		r.Name = "get configuration"
		r.Operation = "get-configuration"
	}

	r.Execute = func(c *fastly.Client) error {
		_, err := i.Executor(c, i.ServiceID)
		return err
	}

	if i.ExpectFailure {
		// TBD
		r.WantNoError = true
	} else {
		r.WantNoError = true
	}

	return &r
}
