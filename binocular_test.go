package binocular

import (
	"fmt"
	"testing"
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
