package testnamestyle

import (
	"go/ast"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

var TestFuncs = &analysis.Analyzer{
	Name:       "testfuncs",
	Doc:        "extract test functions",
	Run:        extractTestFuncs,
	ResultType: reflect.TypeOf(make(TestFuncList, 0)),
}

type TestFuncList []*ast.FuncDecl

func extractTestFuncs(pass *analysis.Pass) (interface{}, error) {
	var testFuncs TestFuncList
	for _, file := range pass.Files {
		filename := pass.Fset.File(file.Pos()).Name()
		isTestFile := strings.HasSuffix(filename, "_test.go")
		if !isTestFile {
			continue
		}

		for _, decl := range file.Decls {
			f, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if !strings.HasPrefix(f.Name.Name, "Test") {
				continue
			}
			r, sz := utf8.DecodeRune([]byte(strings.TrimPrefix(f.Name.Name, "Test")))
			if sz == 0 {
				testFuncs = append(testFuncs, f)
				continue
			}
			if r == utf8.RuneError {
				continue
			}
			if unicode.IsLower(r) {
				continue
			}
			testFuncs = append(testFuncs, f)
		}
	}
	return testFuncs, nil
}
