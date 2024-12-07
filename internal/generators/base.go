package generators

import (
	"flag"
	"fmt"
	"os"

	"golang.org/x/tools/go/packages"
)

type Generator struct {
	ThisGenerator    string
	GenerateAPI      bool
	GenerateAPITests bool
	APIPackage       *packages.Package
}

func (g *Generator) Setup() error {
	flag.BoolVar(&g.GenerateAPI, "api", false, "Generate API functions")
	flag.BoolVar(&g.GenerateAPITests, "tests", false, "Generate API test functions")
	flag.Parse()

	packageName := os.Getenv("GOPACKAGE")

	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedName}
	pkgs, err := packages.Load(cfg, ".")
	if err != nil {
		return fmt.Errorf("loading template for inspection: %w", err)
	}

	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	for i, pkg := range pkgs {
		if pkg.Name == packageName {
			g.APIPackage = pkgs[i]
			break
		}
	}

	if g.APIPackage == nil {
		return fmt.Errorf("failure loading package '%s' from the template", packageName)
	}

	return nil
}
