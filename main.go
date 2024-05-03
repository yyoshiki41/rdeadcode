package main

import (
	"bufio"
	"flag"
	"format/gofmt"
	"go/ast"
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

	scanner := bufio.NewScanner(file)
	var lines []string
	currentLine := 1
	for scanner.Scan() {
		if currentLine < startLine || currentLine > endLine {
			lines = append(lines, scanner.Text())
		}
		currentLine++
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	output, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Failed to open file for writing: %v", err)
	}
	defer output.Close()

	// TODO: Restore original file if we encounter an error
	writer := bufio.NewWriter(output)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			log.Fatalf("Failed to write to file: %v", err)
		}
	}
	if err := writer.Flush(); err != nil {
		log.Fatalf("Failed to flush writer: %v", err)
	}

	// run gofmt on the file
	if err := gofmt.Format(fileName); err != nil {
		log.Fatalf("Failed to run gofmt: %v", err)
	}
}
