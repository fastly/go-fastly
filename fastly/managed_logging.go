package fastly

import (
	"fmt"
	"net/http"
)

type (
	// ManagedLogging represents a managed logging endpoint for a service.
	ManagedLogging struct {
		ServiceID string `mapstructure:"service_id"`
	}

	// ManagedLoggingKind type represents multiple kinds of log streams the
	// managed logging feature will support.
	ManagedLoggingKind uint

	// CreateManagedLoggingInput is used as input to the CreateManagedLogging function.
	CreateManagedLoggingInput struct {
		// Kind is the kind of managed logging we are creating (required).
		Kind ManagedLoggingKind
		// ServiceID is the ID of the service (required).
		ServiceID string
	}
)

const (
	// ManagedLoggingUnset is a log stream variant.
	ManagedLoggingUnset ManagedLoggingKind = iota
	// ManagedLoggingInstanceOutput is a log stream variant.
	ManagedLoggingInstanceOutput
)

// CreateManagedLogging creates a new resource.
func (c *Client) CreateManagedLogging(i *CreateManagedLoggingInput) (*ManagedLogging, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	var path string

	switch i.Kind {
	case ManagedLoggingUnset:
		return nil, ErrMissingKind
	case ManagedLoggingInstanceOutput:
		path = fmt.Sprintf("/service/%s/log_stream/managed/instance_output", i.ServiceID)
	default:
		return nil, ErrNotImplemented
	}

	// nosemgrep: trailofbits.go.invalid-usage-of-modified-variable.invalid-usage-of-modified-variable
	resp, err := c.Post(path, nil)
	// If the service already has managed logging enabled, it will respond
	// with a 409. Handle this case specially so users can decide if this is
	// truly an error.
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusConflict {
			return nil, ErrManagedLoggingEnabled
		}
		return nil, err
	}
	defer resp.Body.Close()

	var m *ManagedLogging
	if err := decodeBodyMap(resp.Body, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// DeleteManagedLoggingInput is used as input to the DeleteManagedLogging function.
type DeleteManagedLoggingInput struct {
	// Kind is the kind of managed logging we are removing (required).
	Kind ManagedLoggingKind
	// ServiceID is the ID of the service (required).
	ServiceID string
}

// DeleteManagedLogging deletes the specified resource.
func (c *Client) DeleteManagedLogging(i *DeleteManagedLoggingInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	var path string

	switch i.Kind {
	case ManagedLoggingUnset:
		return ErrMissingKind
	case ManagedLoggingInstanceOutput:
		path = fmt.Sprintf("/service/%s/log_stream/managed/instance_output", i.ServiceID)
	default:
		return ErrNotImplemented
	}

	_, err := c.Delete(path, nil)
	return err
}
