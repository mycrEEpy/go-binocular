package binocular

import (
	"strings"
	"sync"

	"github.com/kljensen/snowball"
)

// Binocular has a thread-safe inverted index
type Binocular struct {
	mut   sync.RWMutex
	index map[string][]interface{}

	stemming        bool
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
	searchWord := word
	if b.stemming {
		stemmed, err := snowball.Stem(word, "english", b.indexStopWords)
		if err == nil {
			searchWord = stemmed
		}
	}
	b.mut.RLock()
	defer b.mut.RUnlock()
	return b.index[searchWord]
}

// Remove deletes the word from the index
func (b *Binocular) Remove(word string) {
	b.mut.Lock()
	defer b.mut.Unlock()
	delete(b.index, word)
}

// Reset recreates the index
// This will delete the current index
func (b *Binocular) Reset() {
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
