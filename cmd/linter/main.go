package main

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

var UnhandledExitCheckAnalyzer = &analysis.Analyzer{
	Name: "unhandledexit",
	Doc:  "Check for unhandled exit calls, e.g. panic(). log.Fatal and os.Exit allowed only in the main func.",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		isMainPackage := file.Name.Name == "main"

		logIdentName := getImportedNameForPackage("log", file.Imports)
		osIdentName := getImportedNameForPackage("os", file.Imports)

		var stack []ast.Node
		preorderStack(file, stack, func(n ast.Node, stack []ast.Node) bool {
			switch n.(type) {
			case *ast.CallExpr:
				c := n.(*ast.CallExpr)

				switch c.Fun.(type) {
				case *ast.SelectorExpr:
					s := c.Fun.(*ast.SelectorExpr)

					xIdent, ok := s.X.(*ast.Ident)
					if !ok {
						return true
					}

					isMainFunc := false
					if isMainPackage {
						for i := len(stack) - 1; i >= 0; i-- {
							f, ok := stack[i].(*ast.FuncDecl)
							if ok {
								isMainFunc = f.Name.Name == "main"
								break
							}

							_, ok = stack[i].(*ast.FuncLit)
							if ok {
								isMainFunc = false
								break
							}
						}
					}

					if xIdent.Name == logIdentName && s.Sel.Name == "Fatal" && !isMainFunc {
						pass.Reportf(n.Pos(), "unexpected log.Fatal outside of main package")
					}

					if xIdent.Name == osIdentName && s.Sel.Name == "Exit" && !isMainFunc {
						pass.Reportf(n.Pos(), "unexpected os.Exit outside of main package")
					}
				case *ast.Ident:
					i := c.Fun.(*ast.Ident)

					if i.Name == "panic" {
						pass.Reportf(n.Pos(), "unexpected panic outside of main package")
					}
				}

				return true
			}
			return true
		})
	}

	return nil, nil
}

func getImportedNameForPackage(name string, imports []*ast.ImportSpec) string {
	for _, i := range imports {
		if i.Name == nil {
			continue
		}

		if i.Path.Value == fmt.Sprintf("\"%s\"", name) {
			return i.Name.Name
		}
	}

	return name
}

// preorderStack это копия функции ast.PreorderStack из go1.25
// https://cs.opensource.google/go/x/tools/+/refs/tags/v0.39.0:internal/astutil/util.go;l=34
func preorderStack(root ast.Node, stack []ast.Node, f func(n ast.Node, stack []ast.Node) bool) {
	before := len(stack)
	ast.Inspect(root, func(n ast.Node) bool {
		if n != nil {
			if !f(n, stack) {
				// Do not push, as there will be no corresponding pop.
				return false
			}
			stack = append(stack, n) // push
		} else {
			stack = stack[:len(stack)-1] // pop
		}
		return true
	})
	if len(stack) != before {
		panic("push/pop mismatch")
	}
}

func main() {
	singlechecker.Main(UnhandledExitCheckAnalyzer)
}
