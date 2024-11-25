package main

import (
	"path/filepath"

	. "github.com/dave/jennifer/jen"
	"github.com/fastly/go-fastly/v9/internal/generators"
)

var enableInputStructName string

func generateAPI(g *generators.Generator) error {
	var err error

	targetFilename := filepath.Join("..", "..", "..", "fastly", g.PackageName+"_enablement.go")

	f := NewFilePath("fastly")

	g.Header(f)

	f.Line()

	generateGetFunction(g, f)

	if g.FindDefinedTypeStruct("EnableInput") {
		enableInputStructName = "EnableProduct" + g.ProductName + "Input"
		generateEnableStruct(g, f)
	}
	generateEnableFunction(g, f)

	generateDisableFunction(g, f)

	if err = f.Save(targetFilename); err != nil {
		return err
	}

	return nil
}

type generateFunctionInput struct {
	f              *File
	method         string
	commentAction  string
	parameters     []Code
	inputValidator []Code
	returns        []Code
	apiCall        Code
	bodyHandler    []Code
}

func generateFunction(g *generators.Generator, i *generateFunctionInput) {
	parameters := []Code{Id("serviceID").String()}
	parameters = append(parameters, i.parameters...)

	functionName := i.method + "Product" + g.ProductName

	errReturn := func(exp Code) []Code {
		if len(i.bodyHandler) > 0 {
			return []Code{Nil(), exp}
		} else {
			return []Code{exp}
		}
	}

	var body []Code

	body = append(body, If(Id("serviceID").Op("==").Lit("")).Block(Return(errReturn(Qual("fastly", "ErrMissingServiceID"))...)))
	body = append(body, i.inputValidator...)
	body = append(body, Line())
	body = append(body, Id("path").Op(":=").Qual("fastly", "ToSafeURL").Params(Lit("enabled-products"), Qual(g.PackagePath, "ProductID"), Lit("services"), Id("serviceID")))
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

	i.f.Commentf("%s %s the %s product on the service.", functionName, i.commentAction, g.ProductName)
	i.f.Func().Add(generators.ClientReceiver).Id(functionName).Add(Params(parameters...)).Add(Params(i.returns...)).Block(body...)
	i.f.Line()
}

func generateGetFunction(g *generators.Generator, f *File) {
	generateFunction(g, &generateFunctionInput{
		f:             f,
		method:        "Get",
		commentAction: "gets the status of",
		returns:       []Code{Op("*").Qual("fastly", "ProductEnablement"), Error()},
		apiCall:       List(Id("resp"), Err()).Op(":=").Id("c").Dot("Get").Params(Id("path"), Nil()),
		bodyHandler: []Code{
			Var().Id("h").Op("*").Qual("fastly", "ProductEnablement"),
			If(Err().Op(":=").Qual("fastly", "decodeBodyMap").Call(Id("resp").Dot("Body"), Op("&").Id("h")), Err().Op("!=").Nil()).Block(Return(Nil(), Err())),
			Return(Id("h"), Nil()),
		},
	})
}

func generateEnableStruct(g *generators.Generator, f *File) {
	f.Type().Id(enableInputStructName).Op("=").Qual(g.PackagePath, "EnableInput")
	f.Line()
}

func generateEnableFunction(g *generators.Generator, f *File) {
	i := generateFunctionInput{
		f:             f,
		method:        "Enable",
		commentAction: "enables",
		returns:       []Code{Op("*").Qual("fastly", "ProductEnablement"), Error()},
		apiCall:       List(Id("resp"), Err()).Op(":=").Id("c").Dot("Put").Params(Id("path"), Nil()),
		bodyHandler: []Code{
			Var().Id("h").Op("*").Qual("fastly", "ProductEnablement"),
			If(Err().Op(":=").Qual("fastly", "decodeBodyMap").Call(Id("resp").Dot("Body"), Op("&").Id("h")), Err().Op("!=").Nil()).Block(Return(Nil(), Err())),
			Return(Id("h"), Nil()),
		},
	}

	if enableInputStructName != "" {
		i.parameters = []Code{Id("i").Op("*").Id(enableInputStructName)}
		i.inputValidator = []Code{
			If(Err().Op(":=").Id("i").Dot("Validate").Call(), Err().Op("!=").Nil()).Block(Return(Nil(), Err())),
		}
		i.apiCall = List(Id("resp"), Err()).Op(":=").Id("c").Dot("PutJSON").Params(Id("path"), Id("i"), Nil())
	}

	generateFunction(g, &i)
}

func generateDisableFunction(g *generators.Generator, f *File) {
	generateFunction(g, &generateFunctionInput{
		f:             f,
		method:        "Disable",
		commentAction: "disables",
		returns:       []Code{Error()},
		apiCall:       List(Id("resp"), Err()).Op(":=").Id("c").Dot("Delete").Params(Id("path"), Nil()),
		bodyHandler:   []Code{},
	})
}
