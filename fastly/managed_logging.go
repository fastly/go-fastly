package fastly

import (
	"context"
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
func (c *Client) CreateManagedLogging(ctx context.Context, i *CreateManagedLoggingInput) (*ManagedLogging, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	var path string

	switch i.Kind {
	case ManagedLoggingUnset:
		return nil, ErrMissingKind
	case ManagedLoggingInstanceOutput:
		path = ToSafeURL("service", i.ServiceID, "log_stream", "managed", "instance_output")
	default:
		return nil, ErrNotImplemented
	}

	// nosemgrep: trailofbits.go.invalid-usage-of-modified-variable.invalid-usage-of-modified-variable
	resp, err := c.Post(ctx, path, CreateRequestOptions())
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
	if err := DecodeBodyMap(resp.Body, &m); err != nil {
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
func (c *Client) DeleteManagedLogging(ctx context.Context, i *DeleteManagedLoggingInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	var path string

	switch i.Kind {
	case ManagedLoggingUnset:
		return ErrMissingKind
	case ManagedLoggingInstanceOutput:
		path = ToSafeURL("service", i.ServiceID, "log_stream", "managed", "instance_output")
	default:
		return ErrNotImplemented
	}

	ignored, err := c.Delete(ctx, path, CreateRequestOptions())
	if err != nil {
		return err
	}
	defer ignored.Body.Close()
	return nil
}
