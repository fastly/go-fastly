# Development

1. Clone the project to your preferred directory, using your preferred method.
   The module requires the Go version declared in `go.mod`.
2. Download the module dependencies:

   ```bash
   go mod download
   ```

3. Make changes.
4. Verify those changes:

   ```bash
   make all
   ```

   `make all` downloads the module, runs `go mod tidy` (which may modify `go.mod`
   and `go.sum`), formats the code, lints it with a pinned version of
   `golangci-lint` (installed into `./bin`), runs the test suite and runs
   `semgrep`. On macOS, the `semgrep` target runs `brew install semgrep`; on
   other platforms semgrep only runs if it is already installed.

See [TESTING.md](./TESTING.md) for the environment variables and fixtures used
by the test suite.

## Compute

Support for the [Fastly Compute](https://www.fastly.com/products/edge-compute) platform is experimental.

The client compiles with TinyGo, but TinyGo's `reflect` package is incomplete
and parts of the client that decode API responses (via
[github.com/mitchellh/mapstructure](https://github.com/mitchellh/mapstructure))
can panic at runtime with errors such as `unimplemented: reflect.SliceOf()`.

Under TinyGo the default HTTP transport is minimal, as TinyGo's `net/http` does
not support the dialer and HTTP/2 options used with standard Go. You should
supply your own transport that routes requests through a Compute backend, for
example with `fsthttp` from the
[Compute Go SDK](https://github.com/fastly/compute-sdk-go):

```go
client, err := fastly.NewClient(apiKey)
if err != nil {
	return err
}
client.HTTPClient.Transport = fsthttp.NewTransport("fastly")
```

This presumes you have a backend named `fastly` that points to
`https://api.fastly.com`.

### Running the tests with TinyGo

```bash
cd fastly
tinygo test -target=wasip1 -tags=safe -skip=TestClient
```

- `-tags=safe` stops `go-spew` from accessing `reflect` internals, which differ
  between Go and TinyGo.
- `-skip=TestClient` skips tests that use `net/http/httptest`, as TinyGo's
  networking support is incomplete.

The full suite does not currently pass under TinyGo because of the `reflect`
limitations described above (for example `TestStatsClient_*` and
`TestDatacenters`). Use `-run` to select tests that do.
