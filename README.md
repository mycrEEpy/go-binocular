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
cpu: Intel(R) Core(TM) i7-6820HQ CPU @ 2.70GHz
BenchmarkIndex
BenchmarkIndex/basic
BenchmarkIndex/basic-8         	                  591122	      2158 ns/op
BenchmarkIndex/short_sentence
BenchmarkIndex/short_sentence-8         	 2605183	     477.5 ns/op
BenchmarkIndex/stemming
BenchmarkIndex/stemming-8               	   32940	     34007 ns/op
BenchmarkIndex/index_stop_words
BenchmarkIndex/index_stop_words-8       	  639454	      1962 ns/op
BenchmarkIndex/index_short_words
BenchmarkIndex/index_short_words-8      	  617394	      1881 ns/op
BenchmarkIndex/all
BenchmarkIndex/all-8                    	   34707	     35180 ns/op
BenchmarkSearch
BenchmarkSearch/basic
BenchmarkSearch/basic-8                 	31879255	     37.43 ns/op
BenchmarkSearch/stemming
BenchmarkSearch/stemming-8              	  387381	      3409 ns/op
BenchmarkSearch/fuzzy
BenchmarkSearch/fuzzy-8                 	      52	  30109980 ns/op
BenchmarkSearch/all
BenchmarkSearch/all-8                   	     100	  18037148 ns/op
BenchmarkRemove
BenchmarkRemove/index_size_1e+6
BenchmarkRemove/index_size_1e+6-8       	       3	 360625101 ns/op
BenchmarkRemove/index_size_1e+5
BenchmarkRemove/index_size_1e+5-8       	      39	  33768126 ns/op
BenchmarkRemove/index_size_1e+4
BenchmarkRemove/index_size_1e+4-8       	     304	   4012673 ns/op
BenchmarkRemove/index_size_1e+3
BenchmarkRemove/index_size_1e+3-8       	    4255	    261944 ns/op
PASS
ok  	github.com/mycreepy/go-binocular	345.043s
```
