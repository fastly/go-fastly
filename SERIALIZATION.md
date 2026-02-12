# Serialization of mutating requests

Some resources managed via the Fastly Control Plane API (notably
versionless resources on Delivery services, and Next-Gen WAF
workspaces) cannot be safely modified using concurrent requests. If
your usage of go-fastly doesn't modify resources, or does not modify
resources using concurrent goroutines, this issue does not affect you
and you don't need to read the rest of this document!

## Automatic serialization

In the default case, go-fastly will serialize all mutating requests
(Create, Update, and Delete), regardless of the resources they are
modifying. While this is safe, it can also degrade performance, as
concurrent mutating requests against distinct resources will not be
executed in parallel.

## Resource-level serialization

As an alternative, you can provide a 'resource ID' to go-fastly which
it can use to identify requests which would mutate the same resource;
these requests will be serialized against each other, but not against
any other requests.

To do this, use the `fastly.NewContextForResourceID` function to
create a `context.Context` object which contains the resource's unique
ID (e.g. a Service ID/SID or Workspace ID), and then set the `Context`
field of the relevant `Input` structure for the operation you are
planning to invoke.

For example, creating a backend named `test` for `example.com` on a
service with SID `1234abcd` might look like this:

```go
  import (
	"context"
	"github.com/fastly/go-fastly/v13/fastly"
  )

  client := fastly.DefaultClient()
  serviceID := '1234abcd'
  requestContext := fastly.NewContextForResourceID(context.TODO(), serviceID)
  client.CreateBackend(&fastly.CreateBackendInput {
    ServiceID: serviceID,
    ServiceVersion: 5,
    Name: fastly.ToPointer('test'),
    Address: fastly.ToPointer('example.com'),
    Context: requestContext
  })
```
