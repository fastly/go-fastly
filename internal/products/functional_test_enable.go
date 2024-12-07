package products

import (
	"strings"
	"github.com/fastly/go-fastly/v9/fastly"
)

type TestEnableInput[I, O any] struct {
	Phase string
	ExecutorNoInput func(*fastly.Client, string) (*O, error)
	ExecutorWithInput func(*fastly.Client, string, *I) (*O, error)
	Input *I
	ServiceID string
	ExpectFailure bool
}

func TestEnable[I, O any] (i *TestEnableInput[I, O]) *fastly.FunctionalTestInput {
	r := fastly.FunctionalTestInput{}

	if i.Phase != "" {
		r.Name = "enable " + i.Phase
		r.Operation = "enable-" + strings.ReplaceAll(i.Phase, " ", "-")
	} else {
		r.Name = "enable"
		r.Operation = "enable"
	}

	switch any(i.Input).(type) {
	case *NullInput:
		r.Execute = func(c *fastly.Client) error {
			_, err := i.ExecutorNoInput(c, i.ServiceID)
			return err
		}
	default:
		r.Execute = func(c *fastly.Client) error {
			_, err := i.ExecutorWithInput(c, i.ServiceID, i.Input)
			return err
		}
	}

	if i.ExpectFailure {
		// TODO
		r.WantNoError = true
	} else {
		r.WantNoError = true
	}

	return &r
}
