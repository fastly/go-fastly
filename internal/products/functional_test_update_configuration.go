package products

import (
	"strings"

	"github.com/fastly/go-fastly/v9/fastly"
)

type TestUpdateConfigurationInput[I, O any] struct {
	Phase         string
	Executor      func(*fastly.Client, string, *I) (*O, error)
	Input         *I
	ServiceID     string
	ExpectFailure bool
}

func TestUpdateConfiguration[I, O any](i *TestUpdateConfigurationInput[I, O]) *fastly.FunctionalTestInput {
	r := fastly.FunctionalTestInput{}

	if i.Phase != "" {
		r.Name = "update configuration " + i.Phase
		r.Operation = "update-configuration-" + strings.ReplaceAll(i.Phase, " ", "-")
	} else {
		r.Name = "update configuration"
		r.Operation = "update-configuration"
	}

	r.Execute = func(c *fastly.Client) error {
		_, err := i.Executor(c, i.ServiceID, i.Input)
		return err
	}

	if i.ExpectFailure {
		// TODO
		r.WantNoError = true
	} else {
		r.WantNoError = true
	}

	return &r
}
