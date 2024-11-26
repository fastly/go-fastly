package main

import (
	. "github.com/dave/jennifer/jen"
	"github.com/fastly/go-fastly/v9/internal/generators"
)

func generate_api_tests(g *Generator) error {
	var err error

	f := NewFile(g.base.Package.Name + "_test")

	g.base.Header(f)

	f.Line()

	validateGetFunction(g, f)

	validateEnableFunction(g, f, g.base.FindDefinedTypeStruct("EnableInput"))

	validateDisableFunction(g, f)

	if err = f.Save("api_test.go"); err != nil {
		return err
	}

	return nil
}

type validateFunctionInput struct {
	f           *File
	method      string
	returnsBody bool
	needsInput  bool
}

func validateFunction(g *Generator, i *validateFunctionInput) {
	parameter := Id("t").Op("*").Qual("testing", "T")

	var returnVals []Code

	if i.returnsBody {
		returnVals = []Code{Id("_"), Err()}
	} else {
		returnVals = []Code{Err()}
	}

	var body []Code

	testResult := Err().Op("!=").Qual(generators.FastlyPackagePath, "ErrMissingServiceID")
	generateError := Id("t").Dot("Fatalf").Call(Lit("expected '%s', got: '%s'"), Qual(generators.FastlyPackagePath, "ErrMissingServiceID"), Err())

	if !i.needsInput {
		getResult := List(returnVals...).Op(":=").Qual(g.base.Package.PkgPath, i.method).Call(Qual(generators.FastlyPackagePath, "TestClient"), Lit(""))

		body = append(body, If(getResult, testResult).Block(generateError))
	} else {
		getResult := List(returnVals...).Op(":=").Qual(g.base.Package.PkgPath, i.method).Call(Qual(generators.FastlyPackagePath, "TestClient"), Lit(""), Op("&").Id("tc").Dot("Input"))
		// testResult := []Code{Id("tc").Dot("wantError").Op("!=").Nil(), Err().Op("!=").Id("tc").Dot("wantError")}
		// generateError := Id("t").Dot("Fatalf").Call(Lit("expected '%s', got: '%s'"), Id("tc").Dot("wantError"), Err())

		body = append(body, For(List(Id("_"), Id("tc")).Op(":=").Range().Qual(g.base.Package.PkgPath, "EnableInputTestCases").Index(Lit("valid"))).Block(If(getResult, testResult).Block(generateError)))
	}

	//for _, case := range EnableInputTestCases["valid"] {
	//		_, err := Get(&TestClient, "", &case.input)
	//	if case.wantError != nil; err != case.wantError {
	//		t.Failf("expected ... got: %s", err)
	//	}
	//}

	i.f.Func().Id("Test_" + i.method + "_validation").Add(Params(parameter)).Block(body...)
	i.f.Line()
}

func validateGetFunction(g *Generator, f *File) {
	validateFunction(g, &validateFunctionInput{
		f:           f,
		method:      "Get",
		returnsBody: true,
	})
}

func validateEnableFunction(g *Generator, f *File, needsInput bool) {
	i := validateFunctionInput{
		f:           f,
		method:      "Enable",
		returnsBody: true,
		needsInput:  needsInput,
	}

	validateFunction(g, &i)
}

func validateDisableFunction(g *Generator, f *File) {
	validateFunction(g, &validateFunctionInput{
		f:      f,
		method: "Disable",
	})
}
