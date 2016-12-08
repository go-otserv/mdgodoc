package main

import (
	"fmt"

	"github.com/mkideal/cli"
)

type argT struct {
	cli.Helper
	Pkg     string `cli:"*pkg" usage:"name of package which will be documented"`
	SrcHref string `cli:"*src-href" usage:"text template that will be used to \n\t\t   generate links to source code fragments"`
}

func (argv *argT) AutoHelp() bool {
	if argv.Help {
		fmt.Printf("Parse Go source files and produce GoDoc-like markdown documentation.\n\n")
		fmt.Printf("Usage with github src-href template:\n")
		fmt.Printf(" $ mdgodoc --pkg=main --src-href=\"/blob/master/{{.Filename}}#L{{.Line}}\"\n\n")
		fmt.Printf("Object passed to src-href template is of type token.Position.\n\n")
	}
	return argv.Help
}

// go run *.go --pkg=main --src-href="https://github.com/go-otserv/mdgodoc/blob/master/{{.Filename}}#L{{.Line}}"
func main() {
	cli.Run(new(argT), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*argT)

		md := NewMdDoc(argv.SrcHref)
		md.ParseDir(argv.Pkg)
		fmt.Println(md.GenMdDoc(funcs, templs))

		return nil
	})
}
