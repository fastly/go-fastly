package products

import (
	"strings"

	"github.com/fastly/go-fastly/v9/fastly"
)

type TestDisableInput struct {
	Phase         string
	Executor      func(*fastly.Client, string) error
	ServiceID     string
	ExpectFailure bool
}

func TestDisable(i *TestDisableInput) *fastly.FunctionalTestInput {
	r := fastly.FunctionalTestInput{}

	if i.Phase != "" {
		r.Name = "disable " + i.Phase
		r.Operation = "disable-" + strings.ReplaceAll(i.Phase, " ", "-")
	} else {
		r.Name = "disable"
		r.Operation = "disable"
	}

	r.Execute = func(c *fastly.Client) error {
		err := i.Executor(c, i.ServiceID)
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
