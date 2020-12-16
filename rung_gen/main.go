package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	. "github.com/dave/jennifer/jen"
)

var flagRegex = regexp.MustCompile(`-(?P<flag>[a-zA-Z0-9]+)(?:/-(?P<flag_alt>[a-zA-Z0-9]+))?:(?P<type>[a-z0-9]*)`)

func main() {
	parsedFlags, err := parseFlags(os.Args)
	if err != nil {
		log.Fatal("rung: error parsing flags:", err)
	}

	flagStruct := generateFlagStruct(parsedFlags)

	f := NewFile("main")
	f.ImportName("github.com/jamespwilliams/rung", "")
	f.ImportName("flag", "")
	f.Add(flagStruct)
	f.Add(generateMain(parsedFlags))
	f.Save("rung_gen.go")
}

type parsedFlag struct {
	flag         string
	name         string
	flagType     string
	defaultValue string
	usage        string
}

func parseFlags(args []string) ([]parsedFlag, error) {
	var parsedFlags []parsedFlag

	for i := 1; i < len(args); i += 4 {
		if i+3 >= len(args) {
			return nil, errors.New("not enough arguments (format must be " +
				"`-flag1/-1:type defaultValue 'usage'`, and at least one of these components is missing)")
		}

		flagSpec := args[i]
		name := args[i+1]
		defaultValue := args[i+2]
		usage := args[i+3]

		components := findNamedMatches(flagRegex, flagSpec)
		flag, ok := components["flag"]
		if !ok {
			return nil, fmt.Errorf("flag specifier %v is missing the flag component", flagSpec)
		}

		flagType, ok := components["type"]
		if !ok {
			return nil, fmt.Errorf("flag specifier %v is missing the type component", flagSpec)
		}

		parsedFlags = append(parsedFlags, parsedFlag{
			flag:         flag,
			name:         name,
			flagType:     flagType,
			defaultValue: defaultValue,
			usage:        usage,
		})
	}

	return parsedFlags, nil
}

func generateFlagStruct(flags []parsedFlag) *Statement {
	var fields []Code
	for _, flag := range flags {
		fields = append(fields, Id(flag.name).Qual("github.com/jamespwilliams/rung", strings.Title(flag.flagType)+"Flag"))
	}

	c := Type().Id("flagSet").Struct(fields...)
	return c
}

func generateMain(flags []parsedFlag) *Statement {
	var flagPtrDefinitions []Code
	for _, flag := range flags {
		flagPtrDefinitions = append(flagPtrDefinitions, Id(flag.name+"Ptr").Op(":=").Id("flag."+strings.Title(flag.flagType)).Call(
			Lit(flag.flag),
			Id(flag.defaultValue),
			Lit(flag.usage),
		))
	}

	var flagDefinitions []Code
	for _, flag := range flags {
		flagDefinitions = append(flagDefinitions, Id(flag.name).Op(":=").Qual("github.com/jamespwilliams/rung", strings.Title(flag.flagType)+"Flag").Values(Dict{
			Id("Value"): Op("*").Id(flag.name + "Ptr"),
		}))
	}

	flagStructLiteralFields := make(Dict)
	for _, flag := range flags {
		flagStructLiteralFields[Id(flag.name)] = Id(flag.name)
	}

	flagStructLiteral := Id("flags").Op(":=").Id("flagSet").Values(flagStructLiteralFields)

	blockCode := []Code{}
	blockCode = append(blockCode, flagPtrDefinitions...)
	blockCode = append(blockCode, Line())
	blockCode = append(blockCode, Qual("flag", "Parse").Call())
	blockCode = append(blockCode, Line())
	blockCode = append(blockCode, flagDefinitions...)
	blockCode = append(blockCode, Line())
	blockCode = append(blockCode, generateVisitor(flags))
	blockCode = append(blockCode, Line())
	blockCode = append(blockCode, flagStructLiteral)
	blockCode = append(blockCode, Line())
	blockCode = append(blockCode, Id("run").Call(Qual("os", "Stdin"), Qual("os", "Stdout"), Id("flags")))

	return Func().Id("main").Params().Block(
		blockCode...,
	)
}

func generateVisitor(flags []parsedFlag) *Statement {
	return Qual("flag", "Visit").Call(
		Func().Params(Id("flg").Op("*").Qual("flag", "Flag")).BlockFunc(func(g *Group) {
			g.Qual("fmt", "Println").Call(Id("flg"))
			g.Switch(Id("flg").Dot("Name")).BlockFunc(func(g *Group) {
				for _, f := range flags {
					g.Case(Lit(f.flag)).Block(
						Id(f.name).Dot("WasSet").Op("=").Lit(true),
					)
				}
			})
		}),
	)
}
