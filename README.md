# rdeadcode

Parse the output of [deadcode](https://go.dev/blog/deadcode) and remove(rewrite) the dead code from the source files.

## Installation

```shell
go get -u github.com/yyoshiki41/rdeadcode
```

## Usage

```shell
deadcode -json ./path/to/your/project | rdeadcode
// equivalent to
rdeadcode -json deadcode.json
```

or passed the args to `rdeadcode` directly.

```shell
// remove the dead function from the file
rdeadcode -file path/to/your/file.go -function deadFunction
```

## Known issues

> [!WARNING] > _deadcode_ detects methods that implement an interface as dead code if it is not used in the project.

```go
// Verify interface compliance at compile time
var _ fmt.Stringer = myString{}

type myString struct {
	Value string
}

func (s myString) String() string {
	return s.Value
}
```

Please verify the interface compliance at compile time and restore the method when rdeadcode removes it.
