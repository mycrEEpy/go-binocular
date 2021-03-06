[![Go Reference](https://pkg.go.dev/badge/github.com/mycreepy/go-binocular.svg)](https://pkg.go.dev/github.com/mycreepy/go-binocular)
[![Go Report Card](https://goreportcard.com/badge/github.com/mycreepy/go-binocular?style=flat-square)](https://goreportcard.com/report/github.com/mycreepy/go-binocular)

# go-binocular

Some sort of in-memory, record-level inverted index 🤷

## Example

Using a `Binocular` instance:

```go
package main

import (
	"fmt"
	"github.com/mycreepy/go-binocular"
)

func main() {
	b := binocular.New()
	b.AddWithID("Always look on the bright side of life", "123")
	b.AddWithID("Houston we have a problem", "456")
	result, err := b.Search("life", binocular.DefaultIndex)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.Refs()) // ["123"]
	data, err := result.Collect()
	if err != nil {
		panic(err)
	}
	fmt.Println(data) // ["Always look on the bright side of life"]
}
```

Just using a standalone `Index`:

```go
package main

import (
	"fmt"
	"github.com/mycreepy/go-binocular"
)

func main() {
	index := binocular.NewIndex()
	index.Add("Always look on the bright side of life", "123")
	index.Add("Houston we have a problem", "456")
	result := index.Search("life", 0)
	fmt.Println(result) // ["123"]
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
BenchmarkIndex/basic-8   	                  526098	      2530 ns/op
BenchmarkIndex/short_sentence
BenchmarkIndex/short_sentence-8         	 2106877	     578.6 ns/op
BenchmarkIndex/stemming
BenchmarkIndex/stemming-8               	   25466	     45394 ns/op
BenchmarkIndex/index_stop_words
BenchmarkIndex/index_stop_words-8       	  494758	      2308 ns/op
BenchmarkIndex/index_short_words
BenchmarkIndex/index_short_words-8      	  468702	      2550 ns/op
BenchmarkIndex/all
BenchmarkIndex/all-8                    	   28339	     44917 ns/op
BenchmarkSearch
BenchmarkSearch/basic
BenchmarkSearch/basic-8                 	22643508	     48.56 ns/op
BenchmarkSearch/stemming
BenchmarkSearch/stemming-8              	  315206	      4089 ns/op
BenchmarkFuzzySearch
BenchmarkFuzzySearch/basic
BenchmarkFuzzySearch/basic-8            	      42	  39916250 ns/op
BenchmarkFuzzySearch/stemming
BenchmarkFuzzySearch/stemming-8         	     100	  17308908 ns/op
BenchmarkRemove
BenchmarkRemove/index_size_1e+6
BenchmarkRemove/index_size_1e+6-8       	       3	 421132080 ns/op
BenchmarkRemove/index_size_1e+5
BenchmarkRemove/index_size_1e+5-8       	      28	  45965739 ns/op
BenchmarkRemove/index_size_1e+4
BenchmarkRemove/index_size_1e+4-8       	     204	   5386136 ns/op
BenchmarkRemove/index_size_1e+3
BenchmarkRemove/index_size_1e+3-8       	    3774	    314521 ns/op
PASS
ok  	github.com/mycreepy/go-binocular	423.974s
```
