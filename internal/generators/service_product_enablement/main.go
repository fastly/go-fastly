package main

import (
	"github.com/fastly/go-fastly/v9/internal/generators"
)

func main() {
	var (
		err error
		g   *generators.Generator
	)

	if g, err = generators.Setup("service_product_enablement"); err != nil {
		generators.FailErr(err)
	}

	if g.ProductName, err = g.GetDeclaredString("ProductName"); err != nil {
		generators.FailErr(err)
	}

	if err = generateAPI(g); err != nil {
		generators.FailErr(err)
	}
}
