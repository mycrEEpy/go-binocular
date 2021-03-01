package binocular

import (
	"fmt"
	"testing"

	"github.com/tjarratt/babble"
)

func TestIndex_PutAndSearch(t *testing.T) {
	testdata := []struct {
		name      string
		options   []IndexOption
		indexData []string
		search    string
		distance  int
		lenRefs   int
	}{
		{
			"basic",
			[]IndexOption{},
			[]string{
				"Basic testing",
			},
			"testing",
			0,
			1,
		},
		{
			"stop word disabled",
			[]IndexOption{},
			[]string{
				"So testing some stop words",
			},
			"some",
			0,
			0,
		},
		{
			"stop word enabled",
			[]IndexOption{WithStopWords()},
			[]string{
				"So testing some stop words",
			},
			"some",
			0,
			1,
		},
		{
			"short stop word enabled",
			[]IndexOption{WithStopWords(), WithShortWords()},
			[]string{
				"So testing some stop words",
			},
			"so",
			0,
			1,
		},
		{
			"stemming enabled",
			[]IndexOption{WithStemming()},
			[]string{
				"There are too many cats!",
			},
			"cat",
			0,
			1,
		},
		{
			"stemming enabled with simplification",
			[]IndexOption{WithStemming()},
			[]string{
				"There are too many cats!",
			},
			"many",
			0,
			1,
		},
		{
			"stemming enabled with simplification as input",
			[]IndexOption{WithStemming()},
			[]string{
				"There are too many cats!",
			},
			"mani",
			0,
			1,
		},
		{
			"fuzzy search disabled",
			[]IndexOption{},
			[]string{
				"Can we have a dog please?",
			},
			"pls",
			0,
			0,
		},
		{
			"fuzzy search enabled",
			[]IndexOption{},
			[]string{
				"Can we have a dog please?",
			},
			"pls",
			5,
			1,
		},
		{
			"fuzzy search and stemming enabled",
			[]IndexOption{WithStemming()},
			[]string{
				"Please check all the accumulators",
			},
			"accumulator",
			5,
			1,
		},
	}
	for _, td := range testdata {
		t.Run(td.name, func(t *testing.T) {
			index := NewIndex(td.options...)
			for i, v := range td.indexData {
				index.Put(v, i)
			}
			var result []interface{}
			if td.distance > 0 {
				result = index.FuzzySearch(td.search, td.distance)
			} else {
				result = index.Search(td.search)
			}
			if len(result) != td.lenRefs {
				t.Errorf("expected %d, got %d", td.lenRefs, len(result))
			}
			fmt.Printf("%s: %v\n", td.name, result)
		})
	}
}

func TestIndex_Remove(t *testing.T) {
	index := NewIndex()
	index.Put("Some testing data", 1)
	index.Remove(1)
	r1 := index.Search("testing")
	if len(r1) != 0 {
		t.Errorf("result should be empty")
	}

	index.Put("Some testing data", 2)
	index.Put("Some testing data", 3)
	index.Remove(2)
	r2 := index.Search("testing")
	if len(r2) != 1 {
		t.Errorf("result should not be empty")
	}
}

func TestIndex_Drop(t *testing.T) {
	index := NewIndex()
	index.Put("Some testing data", 1)
	index.Drop()
	result := index.Search("testing")
	if len(result) != 0 {
		t.Errorf("data should be empty")
	}
}

func BenchmarkIndex_Put(b *testing.B) {
	testdata := []struct {
		name      string
		options   []IndexOption
		wordCount int
	}{
		{
			"basic",
			[]IndexOption{},
			10,
		},
		{
			"short sentence",
			[]IndexOption{},
			2,
		},
		{
			"stemming",
			[]IndexOption{WithStemming()},
			10,
		},
		{
			"data stop words",
			[]IndexOption{WithStopWords()},
			10,
		},
		{
			"data short words",
			[]IndexOption{WithShortWords()},
			10,
		},
		{
			"all",
			[]IndexOption{WithStemming(), WithShortWords(), WithStopWords()},
			10,
		},
	}
	for _, td := range testdata {
		index := NewIndex(td.options...)
		babbler := babble.NewBabbler()
		babbler.Separator = " "
		babbler.Count = td.wordCount
		b.Run(td.name, func(b *testing.B) {
			sentence := babbler.Babble()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				index.Put(sentence, i)
			}
		})
	}
}

func BenchmarkIndex_Search(b *testing.B) {
	testdata := []struct {
		name      string
		options   []IndexOption
		indexSize int
		wordCount int
	}{
		{
			"basic",
			[]IndexOption{},
			1e+6,
			10,
		},
		{
			"stemming",
			[]IndexOption{WithStemming()},
			1e+6,
			10,
		},
	}
	for _, td := range testdata {
		b.Run(td.name, func(b *testing.B) {
			index := NewIndex(td.options...)
			babbler := babble.NewBabbler()
			babbler.Separator = " "
			babbler.Count = td.wordCount
			for i := 0; i < td.indexSize; i++ {
				index.Put(babbler.Babble(), i)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				index.Search("hello")
			}
		})
	}
}

func BenchmarkIndex_FuzzySearch(b *testing.B) {
	testdata := []struct {
		name      string
		options   []IndexOption
		distance  int
		indexSize int
		wordCount int
	}{
		{
			"basic",
			[]IndexOption{},
			5,
			1e+6,
			10,
		},
		{
			"stemming",
			[]IndexOption{WithStemming()},
			5,
			1e+6,
			10,
		},
	}
	for _, td := range testdata {
		b.Run(td.name, func(b *testing.B) {
			index := NewIndex(td.options...)
			babbler := babble.NewBabbler()
			babbler.Separator = " "
			babbler.Count = td.wordCount
			for i := 0; i < td.indexSize; i++ {
				index.Put(babbler.Babble(), i)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				index.FuzzySearch("hello", td.distance)
			}
		})
	}
}

func BenchmarkIndex_Remove(b *testing.B) {
	testdata := []struct {
		name      string
		indexSize int
		wordCount int
	}{
		{
			"data size 1e+6",
			1e+6,
			10,
		},
		{
			"data size 1e+5",
			1e+5,
			10,
		},
		{
			"data size 1e+4",
			1e+4,
			10,
		},
		{
			"data size 1e+3",
			1e+3,
			10,
		},
	}
	for _, td := range testdata {
		b.Run(td.name, func(b *testing.B) {
			index := NewIndex()
			babbler := babble.NewBabbler()
			babbler.Separator = " "
			babbler.Count = td.wordCount
			for i := 0; i < td.indexSize; i++ {
				index.Put(babbler.Babble(), i)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				index.Remove(b.N)
			}
		})
	}
}
