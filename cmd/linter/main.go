package main

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

var Analyzer = &analysis.Analyzer{
	Name: "panicfatal",
	Doc:  "Сообщает о вызовах panic, а также log.Fatal и os.Exit вне функции main пакета main",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	var funcStack []string

	isInMainMain := func() bool {
		return pass.Pkg != nil &&
			pass.Pkg.Name() == "main" &&
			len(funcStack) > 0 &&
			funcStack[len(funcStack)-1] == "main"
	}

	checkCall := func(call *ast.CallExpr, inMainMain bool) {
		switch fun := call.Fun.(type) {
		case *ast.Ident:
			if fun.Name == "panic" {
				if obj, ok := pass.TypesInfo.Uses[fun]; ok {
					if _, ok := obj.(*types.Builtin); ok {
						pass.Reportf(fun.Pos(), "обнаружен вызов встроенной функции panic")
					}
				}
			}
		case *ast.SelectorExpr:
			sel := fun.Sel
			pkgIdent, ok := fun.X.(*ast.Ident)
			if !ok {
				return
			}
			pkgObj, ok := pass.TypesInfo.Uses[pkgIdent].(*types.PkgName)
			if !ok {
				return
			}
			switch pkgObj.Imported().Path() {
			case "log":
				if sel.Name == "Fatal" && !inMainMain {
					pass.Reportf(sel.Pos(), "вызов log.Fatal вне функции main пакета main")
				}
			case "os":
				if sel.Name == "Exit" && !inMainMain {
					pass.Reportf(sel.Pos(), "вызов os.Exit вне функции main пакета main")
				}
			}
		}
	}

	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.FuncDecl:
				funcStack = append(funcStack, node.Name.Name)
				if node.Body != nil {
					ast.Inspect(node.Body, func(nn ast.Node) bool {
						if call, ok := nn.(*ast.CallExpr); ok {
							checkCall(call, isInMainMain())
						}
						return true
					})
				}
				funcStack = funcStack[:len(funcStack)-1]
				return false
			case *ast.CallExpr:
				checkCall(node, isInMainMain())
			}
			return true
		})
	}

	return nil, nil
}

func main() {
	singlechecker.Main(Analyzer)
}
