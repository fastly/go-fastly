# Examples

Fastly's API is designed to work in the following manner:

1. Create (or clone) a new configuration version for the service
1. Make any changes to the version
1. Validate the version
1. Activate the version

This flow using the Golang client looks like this:

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fastly/go-fastly/v10/fastly"
)

func main() {
	// Create a client object. The client has no state, so it can be persisted
	// and re-used. It is also safe to use concurrently due to its lack of state.
	// There is also a DefaultClient() method that reads an environment variable.
	// Please see the documentation for more information and details.
	client, err := fastly.NewClient(os.Getenv("FASTLY_API_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	// You can find the service ID in the Fastly web console.
	var serviceID = "SERVICE_ID"

	// We'll get the latest 'active' version by inspecting the service metadata and
	// then finding which available version is the 'active' version.
	service, err := client.GetService(&fastly.GetServiceInput{
		ServiceID: serviceID,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Let's acquire a service version to clone from. We'll start by searching for
	// the latest 'active' version available, and if there are no active versions,
	// then we'll clone from whatever is the latest version.
	latest := service.Versions[len(service.Versions)-1]
	for _, version := range service.Versions {
		if *version.Active {
			latest = version
			break
		}
	}

	// Clone the latest version so we can make changes without affecting the
	// active configuration.
	version, err := client.CloneVersion(&fastly.CloneVersionInput{
		ServiceID:      serviceID,
		ServiceVersion: *latest.Number,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Now you can make any changes to the new version. In this example, we will add
	// a new domain.
	domain, err := client.CreateDomain(&fastly.CreateDomainInput{
		ServiceID:      serviceID,
		ServiceVersion: *version.Number,
		Name:           fastly.ToPointer("example.com"),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Output: "example.com"
	fmt.Println("domain.Name:", domain.Name)

	// And we will also add a new backend.
	backend, err := client.CreateBackend(&fastly.CreateBackendInput{
		ServiceID:      serviceID,
		ServiceVersion: *version.Number,
		Name:           fastly.ToPointer("example-backend"),
		Address:        fastly.ToPointer("127.0.0.1"),
		Port:           fastly.ToPointer(80),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Output: "example-backend"
	fmt.Println("backend.Name:", backend.Name)

	// Now we can validate that our version is valid.
	valid, _, err := client.ValidateVersion(&fastly.ValidateVersionInput{
		ServiceID:      serviceID,
		ServiceVersion: *version.Number,
	})
	if err != nil {
		log.Fatal(err)
	}
	if !valid {
		log.Fatal("not valid version")
	}

	// Finally, activate this new version.
	activeVersion, err := client.ActivateVersion(&fastly.ActivateVersionInput{
		ServiceID:      serviceID,
		ServiceVersion: *version.Number,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Output: true
	fmt.Println("activeVersion.Locked:", activeVersion.Locked)
}
```
