package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
)

// flags
var (
	jsonFlag       = flag.String("json", "", "JSON file generated by deadcode")
	fileFlag       = flag.String("file", "", "File to remove function from")
	functionFlag   = flag.String("function", "", "Function to remove")
	ignoreFileFlag = flag.String("ignore", "", "Ignore files matching this regex")
)

func main() {
	flag.Parse()

	// run specific fixer
	if f, fn := *fileFlag, *functionFlag; f != "" && fn != "" {
		if err := fix(f, fn); err != nil {
			log.Fatal(err)
		}
		return
	}
	// run fixer by reading JSON
	var r io.Reader
	switch v := *jsonFlag; v {
	case "", "-":
		r = os.Stdin
	default:
		f, err := os.Open(v)
		if err != nil {
			log.Fatalf("failed to open file %s: %v", v, err)
		}
		defer f.Close()
		r = f
	}
	if err := run(r); err != nil {
		log.Fatal(err)
	}
}

func run(r io.Reader) error {
	packages := []Package{}
	if err := json.NewDecoder(r).Decode(&packages); err != nil {
		return err
	}
	for _, p := range packages {
		for _, f := range p.Funcs {
			if err := fix(f.Position.File, f.Name); err != nil {
				return err
			}
		}
	}
	return nil
}
