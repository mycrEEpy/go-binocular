[![Go Reference](https://pkg.go.dev/badge/github.com/mycreepy/go-binocular.svg)](https://pkg.go.dev/github.com/mycreepy/go-binocular)
[![Go Report Card](https://goreportcard.com/badge/github.com/mycreepy/go-binocular?style=flat-square)](https://goreportcard.com/report/github.com/mycreepy/go-binocular)

# go-binocular

Some sort of in-memory, record-level inverted index 🤷

## Example

```go
package main

import (
	"fmt"
	"github.com/mycreepy/go-binocular"
)

func main() {
	b := binocular.New()
	b.Index("Always look on the bright side of life", 123)
	b.Index("Houston we have a problem", 456)
	result := b.Search("life")
	fmt.Println(result) // [123]
}
```

## Benchmarks

```text
go test -v -bench=. -run=^$
goos: linux
goarch: amd64
pkg: github.com/mycreepy/go-binocular
BenchmarkIndex
BenchmarkIndex/basic
BenchmarkIndex/basic-4                            590067              2169 ns/op
BenchmarkIndex/short_sentence
BenchmarkIndex/short_sentence-4                  2249628               510 ns/op
BenchmarkIndex/stemming
BenchmarkIndex/stemming-4                          31075             35668 ns/op
BenchmarkIndex/index_stop_words
BenchmarkIndex/index_stop_words-4                 667306              2051 ns/op
BenchmarkIndex/index_short_words
BenchmarkIndex/index_short_words-4                633146              2200 ns/op
BenchmarkIndex/all
BenchmarkIndex/all-4                               30289             39230 ns/op
BenchmarkSearch
BenchmarkSearch/basic
BenchmarkSearch/basic-4                         36972064              35.1 ns/op
BenchmarkSearch/stemming
BenchmarkSearch/stemming-4                        356125              4571 ns/op
BenchmarkSearch/fuzzy
BenchmarkSearch/fuzzy-4                               21          47637475 ns/op
BenchmarkSearch/all
BenchmarkSearch/all-4                                 42          23815802 ns/op
BenchmarkRemove
BenchmarkRemove/index_size_1e+6
BenchmarkRemove/index_size_1e+6-4                      3         479793637 ns/op
BenchmarkRemove/index_size_1e+5
BenchmarkRemove/index_size_1e+5-4                     26          45604739 ns/op
BenchmarkRemove/index_size_1e+4
BenchmarkRemove/index_size_1e+4-4                    252           5232758 ns/op
BenchmarkRemove/index_size_1e+3
BenchmarkRemove/index_size_1e+3-4                   4766            248763 ns/op
PASS
ok      github.com/mycreepy/go-binocular        443.910s
```