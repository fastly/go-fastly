package main

import (
	"github.com/fastly/go-fastly/v9/internal/generators"
)

type Generator struct {
	base        *generators.Generator
	productID   string
	productName string
}

func main() {
	var (
		err error
		g   Generator
	)

	if g.base, err = generators.Setup("service_linked_product"); err != nil {
		generators.FailErr(err)
	}

	if g.productID, err = g.base.GetDeclaredString(g.base.APIPackage, "ProductID"); err != nil {
		generators.FailErr(err)
	}

	if g.productName, err = g.base.GetDeclaredString(g.base.APIPackage, "ProductName"); err != nil {
		generators.FailErr(err)
	}

	if err = generate_api(&g); err != nil {
		generators.FailErr(err)
	}

	if err = generate_api_tests(&g); err != nil {
		generators.FailErr(err)
	}
}
