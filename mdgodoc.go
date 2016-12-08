package main

import (
	"bytes"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"text/template"
)

// MdDoc holds the state used for generating documentation.
type MdDoc struct {
	Fset    *token.FileSet
	Pkgs    map[string]*ast.Package
	Dpkg    *doc.Package
	SrcHref *template.Template
}

// NewMdDoc creates new MdDoc instance, srcHref is used to generate link to source.
// * for github: srcHref="https://github.com/go-otserv/mdgodoc/blob/master/{{.Filename}}#L{{.Line}}"
func NewMdDoc(srcHref string) *MdDoc {
	srcHrefTmpl, _ := template.New("hrefTmpl").Parse(srcHref)
	return &MdDoc{token.NewFileSet(), nil, nil, srcHrefTmpl}
}

// ParseDir parses .go files and generate documentation for given package.
func (md *MdDoc) ParseDir(pkgName string) *doc.Package {
	md.Pkgs, _ = parser.ParseDir(md.Fset, ".", nil, parser.ParseComments)
	md.Dpkg = doc.New(md.Pkgs[pkgName], ".", 0)
	return md.Dpkg
}

// GenMdDoc generates markdown documentation from doc.Package instance.
func (md *MdDoc) GenMdDoc(funcs template.FuncMap, templs []string) string {
	funcs["fragment"] = md.genFragmentFunc()
	funcs["srclink"] = md.genSourceFunc()
	var mdBuf bytes.Buffer
	docTpl, _ := template.New("baseTpl").Funcs(funcs).Parse(templs[0])
	for _, templ := range templs[1:] {
		docTpl, _ = template.Must(docTpl.Clone()).Parse(templ)
	}
	docTpl.Execute(&mdBuf, md.Dpkg)
	return normalizeMd(mdBuf.String())
}

func (md *MdDoc) genFragmentFunc() func(token.Pos, token.Pos) string {
	return func(startPos, endPos token.Pos) string {
		tStart := md.Fset.Position(startPos)
		tEnd := md.Fset.Position(endPos)
		fh, _ := os.Open(tStart.Filename)
		defer fh.Close()
		buf := make([]byte, tEnd.Offset-tStart.Offset)
		fh.ReadAt(buf, int64(tStart.Offset))
		return string(buf)
	}
}

func (md *MdDoc) genSourceFunc() func(token.Pos) string {
	return func(pos token.Pos) string {
		var buf bytes.Buffer
		md.SrcHref.Execute(&buf, md.Fset.Position(pos))
		return buf.String()
	}
}

func normalizeMd(doc string) string {
	nlReplace := regexp.MustCompile(`\n(\s)+\n`)
	trimCodes := regexp.MustCompile("\n{2,}```")
	doc = nlReplace.ReplaceAllString(doc, "\n\n")
	doc = trimCodes.ReplaceAllString(doc, "\n```")
	return doc
}
