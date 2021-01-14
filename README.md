# Go Fastly

[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]

[godocs]: https://pkg.go.dev/github.com/fastly/go-fastly/fastly?tab=doc
[v2]: https://pkg.go.dev/github.com/fastly/go-fastly/v2
[v1]: https://pkg.go.dev/github.com/fastly/go-fastly

Go Fastly is a Golang API client for interacting with most facets of the
[Fastly API](https://docs.fastly.com/api).

## Installation

This is a client library, so there is nothing to install. But, it uses Go modules,
so you must be running Go 1.11 or higher.

## Usage

```go
import "github.com/fastly/go-fastly/v2/fastly"
```

## Migrating from v1 to v2

The move from major version [1][v1] to [2][v2] has resulted in a couple of fundamental changes to the library:

- Consistent field name format for IDs and Versions (e.g. `DictionaryID`, `PoolID`, `ServiceID`, `ServiceVersion` etc).
- Input struct fields (for write/update operations) that are optional (i.e. `omitempty`) and use basic types, are now defined as pointers.

The move to more consistent field names in some cases will have resulted in the corresponding sentinel error name to be updated also. For example, `ServiceID` has resulted in a change from `ErrMissingService` to `ErrMissingServiceID`.

The change in type for [basic types](https://tour.golang.org/basics/11) that are optional on input structs related to write/update operations is designed to avoid unexpected behaviours when dealing with their zero value (see [this reference](https://willnorris.com/2014/05/go-rest-apis-and-pointers/) for more details). As part of this change we now provide [helper functions](./fastly/basictypes_helper.go) to assist with generating the new pointer types required.

## Examples

Fastly's API is designed to work in the following manner:

1. Create (or clone) a new configuration version for the service
2. Make any changes to the version
3. Validate the version
4. Activate the version

This flow using the Golang client looks like this:

```go
// Create a client object. The client has no state, so it can be persisted
// and re-used. It is also safe to use concurrently due to its lack of state.
// There is also a DefaultClient() method that reads an environment variable.
// Please see the documentation for more information and details.
client, err := fastly.NewClient("YOUR_FASTLY_API_KEY")
if err != nil {
  log.Fatal(err)
}

// You can find the service ID in the Fastly web console.
var serviceID = "SERVICE_ID"

// Get the latest active version
latest, err := client.LatestVersion(&fastly.LatestVersionInput{
  ServiceID: serviceID,
})
if err != nil {
  log.Fatal(err)
}

// Clone the latest version so we can make changes without affecting the
// active configuration.
version, err := client.CloneVersion(&fastly.CloneVersionInput{
  ServiceID:      serviceID,
  ServiceVersion: latest.Number,
})
if err != nil {
  log.Fatal(err)
}

// Now you can make any changes to the new version. In this example, we will add
// a new domain.
domain, err := client.CreateDomain(&fastly.CreateDomainInput{
  ServiceID:      serviceID,
  ServiceVersion: version.Number,
  Name:           "example.com",
})
if err != nil {
  log.Fatal(err)
}

// Output: "example.com"
fmt.Println(domain.Name)

// And we will also add a new backend.
backend, err := client.CreateBackend(&fastly.CreateBackendInput{
  ServiceID:      serviceID,
  ServiceVersion: version.Number,
  Name:           "example-backend",
  Address:        "127.0.0.1",
  Port:           80,
})
if err != nil {
  log.Fatal(err)
}

// Output: "example-backend"
fmt.Println(backend.Name)

// Now we can validate that our version is valid.
valid, _, err := client.ValidateVersion(&fastly.ValidateVersionInput{
  ServiceID:      serviceID,
  ServiceVersion: version.Number,
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
  ServiceVersion: version.Number,
})
if err != nil {
  log.Fatal(err)
}

// Output: true
fmt.Printf("%t\n", activeVersion.Locked)
```

More information can be found in the
[Fastly Godoc][godocs].

## Developing

1. Clone the project to your preferred directory, using your preferred method.
2. Download the module and accompanying developer tooling.

  ```bash
  $ go mod download
  ```

3. Make changes.
4. Verify those changes.

  ```bash
  $ make all
  ```

## Testing

Go Fastly uses [go-vcr](https://github.com/dnaeon/go-vcr) to "record" and "replay" API request fixtures to improve the speed and portability of integration tests. The test suite uses a single test service ID for all test fixtures.

Contributors without access to the test service can still update the fixtures but with some additional steps required. Below is an example workflow for updating a set of fixture files (where `...` should be replaced with an appropriate value):

```sh
# Remove all yaml fixture files from the specified directory.
#
rm -r fastly/fixtures/.../*

# Run a subsection of the tests.
# This will cause the deleted fixtures to be recreated.
# 
# FASTLY_TEST_SERVICE_ID: should correspond to a real service you control.
# FASTLY_API_KEY: should be a real token associated with the Service you control.
# TESTARGS: allows you to use the -run flag of the 'go test' command.
# 
make test FASTLY_TEST_SERVICE_ID="..." FASTLY_API_KEY="..." TESTARGS="-run=..."
```

> **NOTE**: to run the tests with go-vcr disabled, set `VCR_DISABLE=1` (`make test-full` does this).

When adding or updating client code and integration tests, contributors should record a new set of fixtures. Before submitting a pull request with new or updated fixtures, we ask that contributors update them to use the default service ID by running `make fix-fixtures` with `FASTLY_TEST_SERVICE_ID` set to the same value used to run your tests.

```sh
make fix-fixtures FASTLY_TEST_SERVICE_ID="..."
```

## Contributing

Refer to [CONTRIBUTING.md](./CONTRIBUTING.md)

## License

```
Copyright 2015 Seth Vargo

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
