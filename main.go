package main

import (
	"bufio"
	"bytes"
	"flag"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
)

var (
	fileName, functionName string
)

func main() {
	flag.StringVar(&fileName, "f", "", "File to remove function from (shorthand)")
	flag.StringVar(&fileName, "file", "", "File to remove function from")
	flag.StringVar(&functionName, "fn", "", "File to remove function from (shorthand)")
	flag.StringVar(&functionName, "function", "", "Function to remove")
	flag.Parse()

	if fileName == "" || functionName == "" {
		flag.Usage()
		os.Exit(1)
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("Failed to parse file: %v", err)
	}

	var startLine, endLine int
	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if ok && fn.Name.Name == functionName {
			startLine = fset.Position(fn.Pos()).Line
			endLine = fset.Position(fn.End()).Line
			return false
		}
		return true
	})

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	scanner := bufio.NewScanner(file)
	for l := 1; scanner.Scan(); {
		if l < startLine || l > endLine {
			if _, err := buf.Write(append(scanner.Bytes(), '\n')); err != nil {
				log.Fatalf("Failed to write to buffer: %v", err)
			}
		}
		l++
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// run gofmt
	b, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatalf("Failed to format source: %v", err)
	}
	// write to file
	if err := os.WriteFile(fileName, b, 0644); err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}
}
