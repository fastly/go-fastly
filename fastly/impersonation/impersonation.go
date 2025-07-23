package impersonation

import (
	"context"
)

const QueryParam = "customer_id"

type impersonationKey struct{}

var key impersonationKey

// NewContextForCustomerID returns a [context.Context] which contains
// the customer ID specified in the id parameter.
//
// This [context.Context] should be passed to the various request
// methods of [fastly.Client] so that those request methods can use
// the included customer ID to impersonate the customer when the
// request is made.
func NewContextForCustomerID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, key, id)
}

// CustomerIDFromContext returns the customer ID contained within the
// provided [context.Context]. If the provided [context.Context] does
// not contain a customer ID, then `false` is returned along with an
// empty string.
func CustomerIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(key).(string)
	return id, ok
}
