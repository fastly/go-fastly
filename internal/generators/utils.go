package generators

import (
	"fmt"
	"go/constant"
	"go/types"
	"os"
)

func FailErr(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

func (g *Generator) GetDeclaredString(name string) (string, error) {
	obj := g.Package.Types.Scope().Lookup(name)
	if obj == nil {
		return "", fmt.Errorf("no declaration named '%s' was found in template", name)
	}

	if o, ok := obj.(*types.Const); ok {
		if b, ok := o.Type().(*types.Basic); ok {
			if (b.Info() & types.IsString) != 0 {
				return constant.StringVal(o.Val()), nil
			}
		}
		return "", fmt.Errorf("declaration '%s' must be a string", name)
	}

	return "", fmt.Errorf("declaration '%s' must be a constant string", name)
}

func (g *Generator) FindDefinedTypeStruct(name string) bool {
	obj := g.Package.Types.Scope().Lookup(name)
	if obj == nil {
		return false
	}

	if _, ok := obj.(*types.TypeName); !ok {
		return false
	}

	if _, ok := obj.Type().Underlying().(*types.Struct); ok {
		return true
	}

	return false
}
