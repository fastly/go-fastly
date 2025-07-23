// Package impersonation provides a method for using a
// [context.Context] object to provide a Fastly CID that will be used
// for API requests made using the context. This method is used by
// Fastlyans to issue API requests 'on behalf of' (impersonating) a
// Fastly customer.
//
// Note: This impersonation mechanism is only usable by Fastly employees.
package impersonation
