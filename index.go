package binocular

import (
	"strings"
	"sync"

	"github.com/kljensen/snowball"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

// Index is a thread-safe inverted index.
type Index struct {
	mut  sync.RWMutex
	data map[string][]string

	stemming       bool
	keepStopWords  bool
	keepShortWords bool
}

// IndexOption alters the indexing behavior of an Index.
type IndexOption func(index *Index)

// NewIndex creates a new Index with the given Options.
func NewIndex(options ...IndexOption) *Index {
	index := &Index{
		data: make(map[string][]string),
	}
	for _, opt := range options {
		opt(index)
	}
	return index
}

// WithStemming enables word stemming when reading and writing data to the Index.
func WithStemming() IndexOption {
	return func(index *Index) {
		index.stemming = true
	}
}

// WithStopWords enables indexing of stop words.
func WithStopWords() IndexOption {
	return func(index *Index) {
		index.keepStopWords = true
	}
}

// WithShortWords enables indexing of short words.
// This will have no effect if WithStemming is also enabled.
func WithShortWords() IndexOption {
	return func(index *Index) {
		index.keepShortWords = true
	}
}

// Add splits the given sentence into words and adds them with the reference to the data map.
func (index *Index) Add(sentence string, ref string) {
	for _, word := range strings.Split(sentence, " ") {
		word = stripSpecialChars([]byte(word))
		if index.stemming {
			stemmed, err := snowball.Stem(word, "english", index.keepStopWords)
			if err == nil {
				index.mut.Lock()
				index.data[stemmed] = append(index.data[stemmed], ref)
				index.mut.Unlock()
				continue
			}
		}
		wordLower := strings.ToLower(word)
		if !index.keepShortWords && len(wordLower) <= 2 {
			continue
		}
		if !index.keepStopWords && isStopWord(wordLower) {
			continue
		}
		index.mut.Lock()
		index.data[wordLower] = append(index.data[wordLower], ref)
		index.mut.Unlock()
	}
}

// search returns a slice of references found for the given word.
func (index *Index) search(word string) []string {
	searchWord := strings.ToLower(word)
	if index.stemming {
		stemmed, err := snowball.Stem(searchWord, "english", index.keepStopWords)
		if err == nil {
			searchWord = stemmed
		}
	}
	index.mut.RLock()
	defer index.mut.RUnlock()
	return index.data[searchWord]
}

// Search returns a slice of references found for the given word.
// Distance is the Levenshtein distance.
func (index *Index) Search(word string, distance int) []string {
	if distance <= 0 {
		return index.search(word)
	}

	searchWord := strings.ToLower(word)
	if index.stemming {
		stemmed, err := snowball.Stem(searchWord, "english", index.keepStopWords)
		if err == nil {
			searchWord = stemmed
		}
	}
	index.mut.RLock()
	defer index.mut.RUnlock()
	refs := make([]string, 0)
	for k, v := range index.data {
		d := fuzzy.RankMatch(searchWord, k)
		if d > -1 && d <= distance {
			refs = append(refs, v...)
		}
	}
	return refs
}

// Remove deletes the reference from the Index.
func (index *Index) Remove(ref string) {
	for word, refs := range index.data {
		for i, refInIndex := range refs {
			if refInIndex == ref {
				// if it's the last ref just delete the entry
				if len(refs) == 1 {
					index.mut.Lock()
					delete(index.data, word)
					index.mut.Unlock()
					continue
				}
				// otherwise create a new slice of refs without the given ref
				index.mut.Lock()
				index.data[word] = removeElementFromSlice(refs, i)
				index.mut.Unlock()
			}
		}
	}
}

// Drop deletes the indexed data.
func (index *Index) Drop() {
	index.mut.Lock()
	defer index.mut.Unlock()
	index.data = make(map[string][]string)
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

// swap the given data with the last element in the slice and drop it
// copied from https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
func removeElementFromSlice(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
