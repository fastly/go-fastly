package main

import (
	. "github.com/dave/jennifer/jen"
	"github.com/fastly/go-fastly/v9/internal/generators"
)

func generate_api(g *Generator) error {
	var err error

	f := NewFile(g.APIPackage.Name)

	g.Header(f)

	f.Line()

	generateGetFunction(g, f)

	generateEnableFunction(g, f, g.requiresEnableInput)

	generateDisableFunction(g, f)

	if g.supportsConfiguration {
		generateGetConfigurationFunction(g, f)
		generateUpdateConfigurationFunction(g, f)
	}

	if err = f.Save("api.go"); err != nil {
		return err
	}

	return nil
}

type generateFunctionInput struct {
	f              *File
	operation      string
	urlComponents  Statement
	commentAction  string
	parameters     Statement
	inputValidator Code
	returns        Statement
	apiCall        Code
	bodyHandler    Statement
}

func generateFunction(g *Generator, i *generateFunctionInput) {
	parameters := Statement{generators.FastlyClientParameter, Id("serviceID").String()}
	parameters = append(parameters, i.parameters...)

	errReturn := func(exp Code) Statement {
		if len(i.bodyHandler) > 0 {
			return Statement{Nil(), exp}
		} else {
			return Statement{exp}
		}
	}

	var body Statement

	body = append(body, Var().Err().Id("error"))

	if i.inputValidator != nil {
		body = append(body, Var().Id("pendingErrors").Index().Id("error"))
		body = append(body, If(Id("serviceID").Op("==").Lit("")).Block(Id("pendingErrors").Op("=").Id("append").Call(Id("pendingErrors"), generators.FastlyPackageId("ErrMissingServiceID"))))

		body = append(body, If(i.inputValidator, Err().Op("!=").Nil()).Block(Id("pendingErrors").Op("=").Id("append").Call(Id("pendingErrors"), Err())))

		var cases Statement

		cases = append(cases, Case(Lit(0)).Block())
		cases = append(cases, Case(Lit(1)).Block(Return(errReturn(Id("pendingErrors").Index(Lit(0)))...)))
		cases = append(cases, Default().Block(Return(errReturn(Qual("errors", "Join").Call(Id("pendingErrors").Op("...")))...)))

		body = append(body, Switch(Id("len").Call(Id("pendingErrors"))).Block(cases...))
	} else {
		body = append(body, If(Id("serviceID").Op("==").Lit("")).Block(Return(errReturn(generators.FastlyPackageId("ErrMissingServiceID"))...)))
	}

	body = append(body, Line())

	urlComponents := Statement{Lit("enabled-products"), Lit("v1"), Lit(g.productID), Lit("services"), Id("serviceID")}
	urlComponents = append(urlComponents, i.urlComponents...)
	body = append(body, Id("path").Op(":=").Add(generators.FastlyPackageId("ToSafeURL")).Params(urlComponents...))
	body = append(body, Line())

	body = append(body, i.apiCall)
	body = append(body, If(Err().Op("!=").Nil()).Block(Return(errReturn(Err())...)))
	body = append(body, Defer().Id("resp").Dot("Body").Dot("Close").Call())
	body = append(body, Line())

	if len(i.bodyHandler) > 0 {
		body = append(body, i.bodyHandler...)
	} else {
		body = append(body, Return(Nil()))
	}

	i.f.Commentf("%s %s the %s product on the service.", i.operation, i.commentAction, g.productName)
	i.f.Func().Id(i.operation).Add(Params(parameters...)).Add(Params(i.returns...)).Block(body...)
	i.f.Line()
}

func generateGetFunction(g *Generator, f *File) {
	generateFunction(g, &generateFunctionInput{
		f:             f,
		operation:     "Get",
		commentAction: "gets the status of",
		returns:       Statement{Op("*").Add(generators.FastlyPackageId("ProductEnablement")), Error()},
		apiCall:       List(Id("resp"), Err()).Op(":=").Id("c").Dot("Get").Params(Id("path"), Nil()),
		bodyHandler: Statement{
			Var().Id("h").Op("*").Add(generators.FastlyPackageId("ProductEnablement")),
			If(Err().Op("=").Add(generators.FastlyPackageId("DecodeBodyMap")).Call(Id("resp").Dot("Body"), Op("&").Id("h")), Err().Op("!=").Nil()).Block(Return(Nil(), Err())),
			Return(Id("h"), Nil()),
		},
	})
}

func generateEnableFunction(g *Generator, f *File, needsEnableInput bool) {
	i := generateFunctionInput{
		f:             f,
		operation:     "Enable",
		commentAction: "enables",
		returns:       Statement{Op("*").Add(generators.FastlyPackageId("ProductEnablement")), Error()},
		apiCall:       List(Id("resp"), Err()).Op(":=").Id("c").Dot("Put").Params(Id("path"), Nil()),
		bodyHandler: Statement{
			Var().Id("h").Op("*").Add(generators.FastlyPackageId("ProductEnablement")),
			If(Err().Op("=").Add(generators.FastlyPackageId("DecodeBodyMap")).Call(Id("resp").Dot("Body"), Op("&").Id("h")), Err().Op("!=").Nil()).Block(Return(Nil(), Err())),
			Return(Id("h"), Nil()),
		},
	}

	if needsEnableInput {
		i.parameters = Statement{Id("i").Op("*").Id("EnableInput")}
		i.inputValidator = Err().Op("=").Id("i").Dot("Validate").Call()
		i.apiCall = List(Id("resp"), Err()).Op(":=").Id("c").Dot("PutJSON").Params(Id("path"), Id("i"), Nil())
	}

	generateFunction(g, &i)
}

func generateDisableFunction(g *Generator, f *File) {
	generateFunction(g, &generateFunctionInput{
		f:             f,
		operation:     "Disable",
		commentAction: "disables",
		returns:       Statement{Error()},
		apiCall:       List(Id("resp"), Err()).Op(":=").Id("c").Dot("Delete").Params(Id("path"), Nil()),
		bodyHandler:   Statement{},
	})
}

func generateGetConfigurationFunction(g *Generator, f *File) {
	generateFunction(g, &generateFunctionInput{
		f:             f,
		operation:     "GetConfiguration",
		urlComponents: Statement{Lit("configuration")},
		commentAction: "gets the configuration of",
		returns:       Statement{Op("*").Id("ConfigureOutput"), Error()},
		apiCall:       List(Id("resp"), Err()).Op(":=").Id("c").Dot("Get").Params(Id("path"), Nil()),
		bodyHandler: Statement{
			Var().Id("h").Op("*").Id("ConfigureOutput"),
			If(Err().Op("=").Add(generators.FastlyPackageId("DecodeBodyMap")).Call(Id("resp").Dot("Body"), Op("&").Id("h")), Err().Op("!=").Nil()).Block(Return(Nil(), Err())),
			Return(Id("h"), Nil()),
		},
	})
}

func generateUpdateConfigurationFunction(g *Generator, f *File) {
	i := generateFunctionInput{
		f:              f,
		operation:      "UpdateConfiguration",
		urlComponents:  Statement{Lit("configuration")},
		commentAction:  "updates the configuration of",
		parameters:     Statement{Id("i").Op("*").Id("ConfigureInput")},
		inputValidator: Err().Op("=").Id("i").Dot("Validate").Call(),
		returns:        Statement{Op("*").Id("ConfigureOutput"), Error()},
		apiCall:        List(Id("resp"), Err()).Op(":=").Id("c").Dot("PatchJSON").Params(Id("path"), Id("i"), Nil()),
		bodyHandler: Statement{
			Var().Id("h").Op("*").Id("ConfigureOutput"),
			If(Err().Op("=").Add(generators.FastlyPackageId("DecodeBodyMap")).Call(Id("resp").Dot("Body"), Op("&").Id("h")), Err().Op("!=").Nil()).Block(Return(Nil(), Err())),
			Return(Id("h"), Nil()),
		},
	}

	generateFunction(g, &i)
}
