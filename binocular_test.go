package binocular

import (
	"fmt"
	"testing"

	"github.com/tjarratt/babble"
)

func TestBinocular(t *testing.T) {
	testdata := []struct {
		name    string
		options []Option
		index   []string
		search  string
		lenRefs int
	}{
		{
			"basic",
			[]Option{},
			[]string{
				"Basic testing",
			},
			"testing",
			1,
		},
		{
			"stop word disabled",
			[]Option{},
			[]string{
				"So testing some stop words",
			},
			"some",
			0,
		},
		{
			"stop word enabled",
			[]Option{WithIndexStopWords()},
			[]string{
				"So testing some stop words",
			},
			"some",
			1,
		},
		{
			"short stop word enabled",
			[]Option{WithIndexStopWords(), WithIndexShortWords()},
			[]string{
				"So testing some stop words",
			},
			"so",
			1,
		},
		{
			"stemming enabled",
			[]Option{WithStemming()},
			[]string{
				"There are too many cats!",
			},
			"cat",
			1,
		},
		{
			"stemming enabled with simplification",
			[]Option{WithStemming()},
			[]string{
				"There are too many cats!",
			},
			"many",
			1,
		},
		{
			"stemming enabled with simplification as input",
			[]Option{WithStemming()},
			[]string{
				"There are too many cats!",
			},
			"mani",
			1,
		},
		{
			"fuzzy search disabled",
			[]Option{},
			[]string{
				"Can we have a dog please?",
			},
			"pls",
			0,
		},
		{
			"fuzzy search enabled",
			[]Option{WithFuzzy(5)},
			[]string{
				"Can we have a dog please?",
			},
			"pls",
			1,
		},
		{
			"fuzzy search and stemming enabled",
			[]Option{WithFuzzy(5), WithStemming()},
			[]string{
				"Please check all the accumulators",
			},
			"accumulator",
			1,
		},
	}
	for _, td := range testdata {
		t.Run(td.name, func(t *testing.T) {
			b := New(td.options...)
			for i, v := range td.index {
				b.Index(v, i)
			}
			result := b.Search(td.search)
			if len(result) != td.lenRefs {
				t.Errorf("expected %d, got %d", td.lenRefs, len(result))
			}
			fmt.Printf("%s: %v\n", td.name, result)
		})
	}
}

func TestRemove(t *testing.T) {
	b := New()
	b.Index("Some testing data", 1)
	b.Remove(1)
	r1 := b.Search("testing")
	if len(r1) != 0 {
		t.Errorf("result should be empty")
	}

	b.Index("Some testing data", 2)
	b.Index("Some testing data", 3)
	b.Remove(2)
	r2 := b.Search("testing")
	if len(r2) != 1 {
		t.Errorf("result should not be empty")
	}
}

func TestDrop(t *testing.T) {
	b := New()
	b.Index("Some testing data", 1)
	b.Drop()
	result := b.Search("testing")
	if len(result) != 0 {
		t.Errorf("index should be empty")
	}
}

func BenchmarkIndex(b *testing.B) {
	testdata := []struct {
		name      string
		options   []Option
		wordCount int
	}{
		{
			"basic",
			[]Option{},
			10,
		},
		{
			"short sentence",
			[]Option{},
			2,
		},
		{
			"stemming",
			[]Option{WithStemming()},
			10,
		},
		{
			"index stop words",
			[]Option{WithIndexStopWords()},
			10,
		},
		{
			"index short words",
			[]Option{WithIndexShortWords()},
			10,
		},
		{
			"all",
			[]Option{WithStemming(), WithIndexShortWords(), WithIndexStopWords()},
			10,
		},
	}
	for _, td := range testdata {
		bin := New(td.options...)
		babbler := babble.NewBabbler()
		babbler.Separator = " "
		babbler.Count = td.wordCount
		b.Run(td.name, func(b *testing.B) {
			sentence := babbler.Babble()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				bin.Index(sentence, i)
			}
		})
	}
}

func BenchmarkSearch(b *testing.B) {
	testdata := []struct {
		name      string
		options   []Option
		indexSize int
		wordCount int
	}{
		{
			"basic",
			[]Option{},
			1e+6,
			10,
		},
		{
			"stemming",
			[]Option{WithStemming()},
			1e+6,
			10,
		},
		{
			"fuzzy",
			[]Option{WithFuzzy(5)},
			1e+6,
			10,
		},
		{
			"all",
			[]Option{WithStemming(), WithFuzzy(5)},
			1e+6,
			10,
		},
	}
	for _, td := range testdata {
		b.Run(td.name, func(b *testing.B) {
			bin := New(td.options...)
			babbler := babble.NewBabbler()
			babbler.Separator = " "
			babbler.Count = td.wordCount
			for i := 0; i < td.indexSize; i++ {
				bin.Index(babbler.Babble(), i)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				bin.Search("hello")
			}
		})
	}
}

func BenchmarkRemove(b *testing.B) {
	testdata := []struct {
		name      string
		indexSize int
		wordCount int
	}{
		{
			"index size 1e+6",
			1e+6,
			10,
		},
		{
			"index size 1e+5",
			1e+5,
			10,
		},
		{
			"index size 1e+4",
			1e+4,
			10,
		},
		{
			"index size 1e+3",
			1e+3,
			10,
		},
	}
	for _, td := range testdata {
		b.Run(td.name, func(b *testing.B) {
			bin := New()
			babbler := babble.NewBabbler()
			babbler.Separator = " "
			babbler.Count = td.wordCount
			for i := 0; i < td.indexSize; i++ {
				bin.Index(babbler.Babble(), i)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				bin.Remove(b.N)
			}
		})
	}
}
