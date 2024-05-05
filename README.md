# rdeadcode

Parse the output of [deadcode](https://go.dev/blog/deadcode) and remove(rewrite) the dead code from the source files.

## Installation

```shell
go install github.com/yyoshiki41/rdeadcode@latest
```

## Usage

```shell
deadcode -json -test ./path/to/your/project | rdeadcode
```

or

```shell
rdeadcode -json deadcode.json
```

### Independent features of `rdeadcode`

passed the argument `-file` and `-function` to remove the dead code from the file.

```shell
rdeadcode -file path/to/your/file.go -function deadFunction
```

## Known issues

> [!WARNING]
> _deadcode_ detects methods that implement an interface as dead code if it is not used in the project.

Please verify the interface compliance at compile time and restore the method when rdeadcode removes it.

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
