package main

import (
	"fmt"
	"regexp"
	"strings"
	"text/template"
	"time"
)

var (
	nlToSpaces = regexp.MustCompile(`\n`)

	funcs = template.FuncMap{
		"inline": func(txt string) string {
			return fmt.Sprintf("`%s`", txt)
		},
		"code": func(lang, code string) string {
			return fmt.Sprintf("```%s\n%s\n```", lang, code)
		},
		"genDate": func() string {
			return time.Now().Format("2 Jan 2006")
		},
		"docstring": func(doc string) string {
			return nlToSpaces.ReplaceAllString(doc, "  \n")
		},
		"join": strings.Join,
	}

	baseTpl = `# Package {{inline .Name}}
## Overview
{{.Doc}}
## Index
{{block "index" .}}_no index_{{end}}
{{block "constants" .}}_no constants_{{end}}
{{block "variables" .}}_no variables_{{end}}
{{block "functions" .}}_no functions_{{end}}
{{block "types" .}}_no types_{{end}}
***
_Last updated {{genDate}}_`

	indexTpl = `{{define "index"}}

{{if gt (len .Consts) 0}}
* Constants
{{end}}

{{if gt (len .Vars) 0}}
* Variables
{{end}}

{{if gt (len .Funcs) 0}}
* Functions{{range $idx, $fn := .Funcs}}
  * [{{fragment $fn.Decl.Pos $fn.Decl.End}}](#{{$fn.Name}}){{end}}
{{end}}

{{if gt (len .Types) 0}}
* Types{{range $idx, $ty := .Types}}
  * [{{$ty.Name}}](#{{$ty.Name}}){{range $idx, $fn := $ty.Funcs}}
	 * [{{fragment $fn.Decl.Pos $fn.Decl.End}}](#{{$fn.Name}}){{end}}{{range $idx, $mt := $ty.Methods}}
	 * [{{fragment $mt.Decl.Pos $mt.Decl.End}}](#{{$ty.Name}}-{{$mt.Name}}){{end}}
  {{end}}
{{end}}

{{end}}`

	constantsTpl = `{{define "constants"}}
{{if gt (len .Consts) 0}}

## Constants
  {{range $idx, $val := .Consts}}
{{code "go" (fragment $val.Decl.Pos $val.Decl.End)}}
{{$val.Doc}}
  {{end}}
{{end}}

{{end}}`

	variablesTpl = `{{define "variables"}}
{{if gt (len .Vars) 0}}

## Variables
  {{range $idx, $val := .Vars}}
{{code "go" (fragment $val.Decl.Pos $val.Decl.End)}}
{{$val.Doc}}
  {{end}}
{{end}}

{{end}}`

	functionsTpl = `{{define "functions"}}
{{if gt (len .Funcs) 0}}

## Functions
  {{range $idx, $fn := .Funcs}}
### func <a href="{{srclink $fn.Decl.Pos}}" name="{{$fn.Name}}">{{$fn.Name}}</a> [¶](#{{$fn.Name}})
{{code "go" (fragment $fn.Decl.Pos $fn.Decl.End)}}
{{$fn.Doc}}

  {{end}}
{{end}}

{{end}}`

	typesTpl = `{{define "types"}}
{{if gt (len .Types) 0}}

## Types
  {{range $idx, $ty := .Types}}
### type <a href="{{srclink $ty.Decl.Pos}}" name="{{$ty.Name}}">{{$ty.Name}}</a> [¶](#{{$ty.Name}})
{{code "go" (fragment $ty.Decl.Pos $ty.Decl.End)}}
{{docstring $ty.Doc}}

{{range $idx, $fn := $ty.Funcs}}
#### func <a href="{{srclink $fn.Decl.Pos}}" name="{{$fn.Name}}">{{$fn.Name}}</a> [¶](#{{$fn.Name}})
{{code "go" (fragment $fn.Decl.Pos $fn.Decl.End)}}
{{$fn.Doc}}

{{end}}

{{range $idx, $mt := $ty.Methods}}
#### func <a href="{{srclink $mt.Decl.Pos}}" name="{{$ty.Name}}-{{$mt.Name}}">{{$mt.Name}}</a> [¶](#{{$ty.Name}}-{{$mt.Name}})
{{code "go" (fragment $mt.Decl.Pos $mt.Decl.End)}}
{{$mt.Doc}}

{{end}}

  {{end}}
{{end}}

{{end}}`
	templs = []string{baseTpl, indexTpl, constantsTpl, variablesTpl, functionsTpl, typesTpl}
)
