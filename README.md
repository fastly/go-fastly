# Go Fastly

[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][latest]

A Go client module for interacting with most facets of the [Fastly
API](https://docs.fastly.com/api).

> **NOTE:** This API client may not function correctly when used on
> the [Fastly Compute](https://www.fastly.com/products/edge-compute)
> platform. Support for Compute is on the roadmap but has not yet been
> prioritised ([details](./DEVELOPMENT.md#compute)).

## Versioning and Release Schedules

The maintainers of this module strive to maintain [semantic versioning
(SemVer)](https://semver.org/). This means that breaking changes
(removal of functionality, or incompatible changes to existing
functionality) will be released in a version with the first version
component (`major`) incremented. Feature additions will increment the
second version component (`minor`), and bug fixes which do not affect
compatibility will increment the third version component (`patch`).

On the first Wednesday of each month, a release will be published
including all breaking, feature, and bug-fix changes that are ready
for release. If that Wednesday should happen to be a US holiday, the
release will be delayed until the next available working day.

If critical or urgent bug fixes are ready for release in between those
primary releases, patch releases will be made as needed to make those
fixes available.

## Usage

```go
import "github.com/fastly/go-fastly/v9/fastly"
```

## Reference

- [CONTRIBUTING.md](./CONTRIBUTING.md)
- [DEVELOPMENT.md](./DEVELOPMENT.md)
- [EXAMPLES.md](./EXAMPLES.md)
- [TESTING.md](./TESTING.md)

[latest]: https://pkg.go.dev/github.com/fastly/go-fastly/v9/fastly
