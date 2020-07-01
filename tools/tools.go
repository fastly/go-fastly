// +build tools

// This follows the agreed-upon current best approach for adding developer tooling to your Module.
//
// For more details, refer to this section in the Go Wiki:
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

package tools

import (
	_ "golang.org/x/tools/cmd/goimports"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
