# Development

1. Clone the project to your preferred directory, using your preferred method.

2. This package requires Go 1.22, specifically.
   1. If you have only Go 1.22 installed, then no further steps are needed.
   2. If you have Go 1.22 and other versions installed, and are using
      any form of Go version manager to select the 'default' version
      to be used, then select Go 1.22 and proceed.
   3. If you have Go 1.22 and other versions installed, and are using
      Go's built-in support for multiple versions, and Go 1.22 is not
      the highest-numbered version installed, then you need to set the
      `GO_FOR_BUILD` environment variable to ensure that Go 1.22 is
      used to build and test the package.

	  ```bash
	  $ export GO_FOR_BUILD=go1.22.9
	  ```

3. Download the modules that this package depends on, and necessary developer tooling.

  ```bash
  $ make mod-download dev-dependencies
  ```

4. Make changes.

5. Verify those changes.

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
