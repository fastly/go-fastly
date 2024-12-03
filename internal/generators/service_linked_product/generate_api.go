package main

import (
	. "github.com/dave/jennifer/jen"
	"github.com/fastly/go-fastly/v9/internal/generators"
)

func generate_api(g *Generator) error {
	var err error

	f := NewFile(g.base.APIPackage.Name)

	g.base.Header(f)

	f.Line()

	generateGetFunction(g, f)

	generateEnableFunction(g, f, g.base.FindDefinedTypeStruct(g.base.APIPackage, "EnableInput"))

	generateDisableFunction(g, f)

	if g.base.FindDefinedTypeStruct(g.base.APIPackage, "ConfigureInput") {
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
	method         string
	urlComponents  []Code
	commentAction  string
	parameters     []Code
	inputValidator []Code
	returns        []Code
	apiCall        Code
	bodyHandler    []Code
}

func generateFunction(g *Generator, i *generateFunctionInput) {
	parameters := []Code{generators.ClientParameter, Id("serviceID").String()}
	parameters = append(parameters, i.parameters...)

	errReturn := func(exp Code) []Code {
		if len(i.bodyHandler) > 0 {
			return []Code{Nil(), exp}
		} else {
			return []Code{exp}
		}
	}

	var body []Code

	body = append(body, If(Id("serviceID").Op("==").Lit("")).Block(Return(errReturn(Qual(generators.FastlyPackagePath, "ErrMissingServiceID"))...)))
	body = append(body, i.inputValidator...)
	body = append(body, Line())
	urlComponents := []Code{Lit("enabled-products"), Lit("v1"), Lit(g.productID), Lit("services"), Id("serviceID")}
	urlComponents = append(urlComponents, i.urlComponents...)
	body = append(body, Id("path").Op(":=").Qual(generators.FastlyPackagePath, "ToSafeURL").Params(urlComponents...))
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

	i.f.Commentf("%s %s the %s product on the service.", i.method, i.commentAction, g.productName)
	i.f.Func().Id(i.method).Add(Params(parameters...)).Add(Params(i.returns...)).Block(body...)
	i.f.Line()
}

func generateGetFunction(g *Generator, f *File) {
	generateFunction(g, &generateFunctionInput{
		f:             f,
		method:        "Get",
		commentAction: "gets the status of",
		returns:       []Code{Op("*").Qual(generators.FastlyPackagePath, "ProductEnablement"), Error()},
		apiCall:       List(Id("resp"), Err()).Op(":=").Id("c").Dot("Get").Params(Id("path"), Nil()),
		bodyHandler: []Code{
			Var().Id("h").Op("*").Qual(generators.FastlyPackagePath, "ProductEnablement"),
			If(Err().Op(":=").Qual(generators.FastlyPackagePath, "DecodeBodyMap").Call(Id("resp").Dot("Body"), Op("&").Id("h")), Err().Op("!=").Nil()).Block(Return(Nil(), Err())),
			Return(Id("h"), Nil()),
		},
	})
}

func generateEnableFunction(g *Generator, f *File, needsEnableInput bool) {
	i := generateFunctionInput{
		f:             f,
		method:        "Enable",
		commentAction: "enables",
		returns:       []Code{Op("*").Qual(generators.FastlyPackagePath, "ProductEnablement"), Error()},
		apiCall:       List(Id("resp"), Err()).Op(":=").Id("c").Dot("Put").Params(Id("path"), Nil()),
		bodyHandler: []Code{
			Var().Id("h").Op("*").Qual(generators.FastlyPackagePath, "ProductEnablement"),
			If(Err().Op(":=").Qual(generators.FastlyPackagePath, "DecodeBodyMap").Call(Id("resp").Dot("Body"), Op("&").Id("h")), Err().Op("!=").Nil()).Block(Return(Nil(), Err())),
			Return(Id("h"), Nil()),
		},
	}

	if needsEnableInput {
		i.parameters = []Code{Id("i").Op("*").Id("EnableInput")}
		i.inputValidator = []Code{
			If(Err().Op(":=").Id("i").Dot("Validate").Call(), Err().Op("!=").Nil()).Block(Return(Nil(), Err())),
		}
		i.apiCall = List(Id("resp"), Err()).Op(":=").Id("c").Dot("PutJSON").Params(Id("path"), Id("i"), Nil())
	}

	generateFunction(g, &i)
}

func generateDisableFunction(g *Generator, f *File) {
	generateFunction(g, &generateFunctionInput{
		f:             f,
		method:        "Disable",
		commentAction: "disables",
		returns:       []Code{Error()},
		apiCall:       List(Id("resp"), Err()).Op(":=").Id("c").Dot("Delete").Params(Id("path"), Nil()),
		bodyHandler:   []Code{},
	})
}

func generateGetConfigurationFunction(g *Generator, f *File) {
	generateFunction(g, &generateFunctionInput{
		f:             f,
		method:        "GetConfiguration",
		urlComponents: []Code{Lit("configuration")},
		commentAction: "gets the configuration of",
		returns:       []Code{Op("*").Id("ConfigureOutput"), Error()},
		apiCall:       List(Id("resp"), Err()).Op(":=").Id("c").Dot("Get").Params(Id("path"), Nil()),
		bodyHandler: []Code{
			Var().Id("h").Op("*").Id("ConfigureOutput"),
			If(Err().Op(":=").Qual(generators.FastlyPackagePath, "DecodeBodyMap").Call(Id("resp").Dot("Body"), Op("&").Id("h")), Err().Op("!=").Nil()).Block(Return(Nil(), Err())),
			Return(Id("h"), Nil()),
		},
	})
}

func generateUpdateConfigurationFunction(g *Generator, f *File) {
	i := generateFunctionInput{
		f:             f,
		method:        "UpdateConfiguration",
		urlComponents: []Code{Lit("configuration")},
		commentAction: "updates the configuration of",
		parameters:    []Code{Id("i").Op("*").Id("ConfigureInput")},
		inputValidator: []Code{
			If(Err().Op(":=").Id("i").Dot("Validate").Call(), Err().Op("!=").Nil()).Block(Return(Nil(), Err())),
		},
		returns: []Code{Op("*").Id("ConfigureOutput"), Error()},
		apiCall: List(Id("resp"), Err()).Op(":=").Id("c").Dot("PutJSON").Params(Id("path"), Id("i"), Nil()),
		bodyHandler: []Code{
			Var().Id("h").Op("*").Id("ConfigureOutput"),
			If(Err().Op(":=").Qual(generators.FastlyPackagePath, "DecodeBodyMap").Call(Id("resp").Dot("Body"), Op("&").Id("h")), Err().Op("!=").Nil()).Block(Return(Nil(), Err())),
			Return(Id("h"), Nil()),
		},
	}

	generateFunction(g, &i)
}
