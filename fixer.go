package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

var (
	ErrFuncNotFound = fmt.Errorf("not found")
)

func fix(filename, functionName string) error {
	start, end, err := lookup(filename, functionName)
	if err != nil {
		return err
	}

	if err := replaceLines(filename, start, end); err != nil {
		return err
	}
	return nil
}

func lookup(filename, functionName string) (int, int, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return 0, 0, err
	}
	funcParts := strings.Split(functionName, ".")

	var start, end int
	ast.Inspect(node, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.FuncDecl:
			if len(funcParts) == 2 && t.Recv != nil {
				// handle method
				for _, field := range t.Recv.List {
					var ident *ast.Ident
					switch field.Type.(type) {
					case *ast.Ident:
						ident = field.Type.(*ast.Ident)
					case *ast.StarExpr:
						ident = field.Type.(*ast.StarExpr).X.(*ast.Ident)
					}
					if ident != nil && ident.Name == funcParts[0] &&
						t.Name.Name == funcParts[1] {
						start = fset.Position(t.Pos()).Line
						end = fset.Position(t.End()).Line
						return false
					}
				}
			} else {
				// handle function
				if t.Name.Name == functionName && t.Recv == nil {
					start = fset.Position(t.Pos()).Line
					end = fset.Position(t.End()).Line
					return false
				}
			}
		}
		return true
	})
	if start == 0 || end == 0 {
		return 0, 0, fmt.Errorf("%w, func %s in %s", ErrFuncNotFound, functionName, filename)
	}
	return start, end, nil
}

func replaceLines(filename string, start, end int) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	scanner := bufio.NewScanner(file)
	for l := 1; scanner.Scan(); {
		if l < start || l > end {
			if _, err := buf.Write(append(scanner.Bytes(), '\n')); err != nil {
				return err
			}
		}
		l++
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	// run gofmt
	b, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	// write to file
	return writeFile(file, b)
}

func writeFile(file *os.File, b []byte) error {
	// write to file
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if err := os.WriteFile(file.Name(), b, info.Mode()); err != nil {
		return err
	}
	return nil
}
