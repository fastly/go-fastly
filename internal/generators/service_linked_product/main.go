package main

import (
	"github.com/fastly/go-fastly/v9/internal/generators"
)

type Generator struct {
	generators.Generator
	productID             string
	productName           string
	requiresEnableInput   bool
	supportsConfiguration bool
}

func main() {
	var (
		err error
		g   Generator
	)

	g.ThisGenerator = "service_linked_product"

	if err = g.Setup(); err != nil {
		g.FailErr(err)
	}

	if g.productID, err = g.GetDeclaredString("ProductID"); err != nil {
		g.FailErr(err)
	}

	if g.productName, err = g.GetDeclaredString("ProductName"); err != nil {
		g.FailErr(err)
	}

	g.requiresEnableInput = g.FindDefinedTypeStruct("EnableInput")
	g.supportsConfiguration = g.FindDefinedTypeStruct("ConfigureInput")

	if g.GenerateAPI {
		if err = generate_api(&g); err != nil {
			g.FailErr(err)
		}
	}

	if g.GenerateAPITests {
		if err = generate_api_tests(&g); err != nil {
			g.FailErr(err)
		}
	}
}
