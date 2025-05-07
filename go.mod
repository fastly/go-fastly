module github.com/fastly/go-fastly/v10

go 1.24.0

toolchain go1.24.2

require (
	github.com/dnaeon/go-vcr v1.2.0
	github.com/google/go-cmp v0.7.0
	github.com/google/go-querystring v1.1.0
	github.com/google/jsonapi v1.0.0
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/mitchellh/mapstructure v1.5.0
	github.com/peterhellberg/link v1.2.0
	github.com/stretchr/testify v1.10.0
	golang.org/x/crypto v0.38.0
)

require (
	github.com/BurntSushi/toml v1.5.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/exp/typeparams v0.0.0-20250408133849-7e4ce0ab07d0 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/tools v0.33.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	honnef.co/go/tools v0.6.1 // indirect
)

tool (
	golang.org/x/tools/cmd/goimports
	honnef.co/go/tools/cmd/staticcheck
)
