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
)

func fix(filename, functionName string) error {
	start, end, err := lookup(filename, functionName)
	if err != nil {
		return err
	}

	if err := rewrite(filename, start, end); err != nil {
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

	var start, end int
	ast.Inspect(node, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if ok && funcDecl.Name.Name == functionName {
			start = fset.Position(funcDecl.Pos()).Line
			end = fset.Position(funcDecl.End()).Line
			return false
		}
		return true
	})
	if start == 0 || end == 0 {
		return 0, 0, fmt.Errorf("func %s not found in %s", functionName, filename)
	}
	return start, end, nil
}

func rewrite(filename string, start, end int) error {
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
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if err := os.WriteFile(filename, b, info.Mode()); err != nil {
		return err
	}
	return nil
}
