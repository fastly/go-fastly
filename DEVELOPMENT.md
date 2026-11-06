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

The client compiles with TinyGo and its unit tests pass. When running on
Compute you should supply your own HTTP transport that routes requests through
a Compute backend, for example with `fsthttp` from the
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
make test-tinygo
```

This runs `tinygo test -target=wasip1 -tags=safe` against all packages. The
`TESTARGS` variable can be used to select tests, as with `make test`.

> **NOTE:** The tests require a version of TinyGo that includes the `reflect`
> additions from [tinygo-org/tinygo#6081](https://github.com/tinygo-org/tinygo/pull/6081),
> which was merged after the 0.44.0 release. Earlier versions fail with
> `unimplemented: reflect.SliceOf()`.

If you want to run pieces by hand, the flags are:

- `-tags=safe` stops `go-spew` from accessing `reflect` internals, which differ
  between Go and TinyGo.
- `-target=wasip1` is not required, but it guards against code that assumes
  ints are 64 bits wide.

Tests that rely on features TinyGo does not support (such as `net/http/httptest`)
skip themselves when built with TinyGo.
