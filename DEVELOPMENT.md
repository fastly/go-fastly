# Development

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

## Compute

Support for the [Fastly Compute](https://www.fastly.com/products/edge-compute) platform is still in development.

There are known issues with the use of Go's `reflect` package and for TinyGo support to mature.

> **NOTE:** The go-fastly API client uses [github.com/mitchellh/mapstructure](https://github.com/mitchellh/mapstructure)

If using standard Go (not TinyGo) then a usable client can be achieved with:

```go
client, err := fastly.NewClient("FASTLY_API_KEY")
client.HTTPClient.Transport = fsthttp.NewTransport("fastly")
```

This presumes you have a backend named `fastly` pointing to `https://api.fastly.com`
