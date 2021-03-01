package binocular

import (
	"strings"
	"sync"

	"github.com/kljensen/snowball"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

// Binocular has a thread-safe inverted index
type Binocular struct {
	mut   sync.RWMutex
	index map[string][]interface{}

	stemming        bool
	fuzzyDistance   int
	indexStopWords  bool
	indexShortWords bool
}

// Option alters the behavior of a Binocular
type Option func(b *Binocular)

// New creates a new Binocular with the given Options
func New(options ...Option) *Binocular {
	b := &Binocular{
		index: map[string][]interface{}{},
	}
	for _, opt := range options {
		opt(b)
	}
	return b
}

// WithStemming enables word stemming when calling Binocular.Index
func WithStemming() Option {
	return func(b *Binocular) {
		b.stemming = true
	}
}

// WithFuzzy enables fuzzy searches
func WithFuzzy(distance int) Option {
	return func(b *Binocular) {
		b.fuzzyDistance = distance
	}
}

// WithIndexStopWords enables indexing of stop words
func WithIndexStopWords() Option {
	return func(b *Binocular) {
		b.indexStopWords = true
	}
}

// WithIndexShortWords enables indexing of short words
// This will have no effect if WithStemming is also enabled
func WithIndexShortWords() Option {
	return func(b *Binocular) {
		b.indexShortWords = true
	}
}

// Index splits the given sentence into words and adds them with the reference to the index
func (b *Binocular) Index(sentence string, ref interface{}) {
	for _, word := range strings.Split(sentence, " ") {
		word = stripSpecialChars([]byte(word))
		if b.stemming {
			stemmed, err := snowball.Stem(word, "english", b.indexStopWords)
			if err == nil {
				b.mut.Lock()
				b.index[stemmed] = append(b.index[stemmed], ref)
				b.mut.Unlock()
				continue
			}
		}
		wordLower := strings.ToLower(word)
		if !b.indexShortWords && len(wordLower) <= 2 {
			continue
		}
		if !b.indexStopWords && isStopWord(wordLower) {
			continue
		}
		b.mut.Lock()
		b.index[wordLower] = append(b.index[wordLower], ref)
		b.mut.Unlock()
	}
}

// Search returns a slice of references found for the given word
func (b *Binocular) Search(word string) []interface{} {
	searchWord := strings.ToLower(word)
	if b.stemming {
		stemmed, err := snowball.Stem(searchWord, "english", b.indexStopWords)
		if err == nil {
			searchWord = stemmed
		}
	}
	b.mut.RLock()
	defer b.mut.RUnlock()
	if b.fuzzyDistance > 0 {
		fuzzyRefs := make([]interface{}, 0)
		for k, v := range b.index {
			d := fuzzy.RankMatch(searchWord, k)
			if d > -1 && d <= b.fuzzyDistance {
				fuzzyRefs = append(fuzzyRefs, v...)
			}
		}
		return fuzzyRefs
	}
	return b.index[searchWord]
}

// Remove deletes the reference from the index
func (b *Binocular) Remove(ref interface{}) {
	for word, refs := range b.index {
		for i, refInIndex := range refs {
			if refInIndex == ref {
				// if it's the last ref just delete the entry
				if len(refs) == 1 {
					b.mut.Lock()
					delete(b.index, word)
					b.mut.Unlock()
					continue
				}
				// otherwise create a new slice of refs without the given ref
				b.mut.Lock()
				b.index[word] = removeElementFromSlice(refs, i)
				b.mut.Unlock()
			}
		}
	}
}

// Drop deletes the whole index
func (b *Binocular) Drop() {
	b.mut.Lock()
	defer b.mut.Unlock()
	b.index = map[string][]interface{}{}
}

// copied from snowball package as it's unexported
func isStopWord(word string) bool {
	switch word {
	case "a", "about", "above", "after", "again", "against", "all", "am", "an",
		"and", "any", "are", "as", "at", "be", "because", "been", "before",
		"being", "below", "between", "both", "but", "by", "can", "did", "do",
		"does", "doing", "don", "down", "during", "each", "few", "for", "from",
		"further", "had", "has", "have", "having", "he", "her", "here", "hers",
		"herself", "him", "himself", "his", "how", "i", "if", "in", "into", "is",
		"it", "its", "itself", "just", "me", "more", "most", "my", "myself",
		"no", "nor", "not", "now", "of", "off", "on", "once", "only", "or",
		"other", "our", "ours", "ourselves", "out", "over", "own", "s", "same",
		"she", "should", "so", "some", "such", "t", "than", "that", "the", "their",
		"theirs", "them", "themselves", "then", "there", "these", "they",
		"this", "those", "through", "to", "too", "under", "until", "up",
		"very", "was", "we", "were", "what", "when", "where", "which", "while",
		"who", "whom", "why", "will", "with", "you", "your", "yours", "yourself",
		"yourselves":
		return true
	}
	return false
}

// faster than using regex
// copied from https://stackoverflow.com/questions/54461423/efficient-way-to-remove-all-non-alphanumeric-characters-from-large-text
func stripSpecialChars(s []byte) string {
	j := 0
	for _, b := range s {
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			s[j] = b
			j++
		}
	}
	return string(s[:j])
}

// swap the given index with the last element in the slice and drop it
// copied from https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
func removeElementFromSlice(s []interface{}, i int) []interface{} {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
