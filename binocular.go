package binocular

import (
	"errors"
	"reflect"

	"github.com/fatih/structtag"
	"github.com/google/uuid"
)

// DefaultIndex is the default Index when adding data to a Binocular.
const DefaultIndex = "default"

// ErrIndexNotFound indicates that the given index does not exist.
var ErrIndexNotFound = errors.New("index not found")

// ErrRefNotFound indicates that the given reference does not exist.
var ErrRefNotFound = errors.New("ref not found")

// Binocular holds you data and can use multiple Indices for searching it.
// DefaultIndex is the currently configured default Index for the given Binocular instance.
type Binocular struct {
	docs         map[string]*document
	indices      map[string]*Index
	DefaultIndex string
}

type document struct {
	Data          interface{}
	recordLocator map[string]struct{}
}

// Option can alter the behavior if a Binocular instance.
type Option func(binocular *Binocular)

// New will create a new Binocular instance with the given Options.
func New(options ...Option) *Binocular {
	binocular := &Binocular{
		docs:         map[string]*document{},
		indices:      map[string]*Index{},
		DefaultIndex: DefaultIndex,
	}
	for _, opt := range options {
		opt(binocular)
	}
	if _, ok := binocular.indices[binocular.DefaultIndex]; !ok {
		binocular.indices[DefaultIndex] = NewIndex()
	}
	return binocular
}

// WithDefaultIndex creates a new default Index with the given name and IndexOptions.
func WithDefaultIndex(name string, options ...IndexOption) Option {
	return func(binocular *Binocular) {
		binocular.DefaultIndex = name
		binocular.indices[name] = NewIndex(options...)
	}
}

// WithIndex creates a new Index with the given name and IndexOptions.
func WithIndex(name string, options ...IndexOption) Option {
	return func(binocular *Binocular) {
		binocular.indices[name] = NewIndex(options...)
	}
}

// Add will create a new id for your data and adds it to the Binocular instance.
func (binocular *Binocular) Add(data interface{}) string {
	id := uuid.New().String()
	binocular.AddWithID(id, data)
	return id
}

// AddWithID adds the data with the given id to the Binocular instance.
func (binocular *Binocular) AddWithID(id string, data interface{}) {
	doc := document{
		Data:          data,
		recordLocator: make(map[string]struct{}),
	}
	binocular.docs[id] = &doc

	switch v := doc.Data.(type) {
	case string:
		binocular.indices[binocular.DefaultIndex].Add(v, id)
		doc.recordLocator[binocular.DefaultIndex] = struct{}{}
	default:
		t := reflect.TypeOf(doc.Data)
		if t.Kind() == reflect.Struct {
			binocular.parseStruct(id, &doc, t)
		}
	}
}

// Get will retrieve the data at the given id.
// ErrRefNotFound is returned if the data does not exist.
func (binocular *Binocular) Get(id string) (interface{}, error) {
	doc, ok := binocular.docs[id]
	if !ok {
		return nil, ErrRefNotFound
	}
	return doc.Data, nil
}

// Search will search the given index with the given word and returns a SearchResult.
// ErrIndexNotFound is returned if the given index does not exist.
func (binocular *Binocular) Search(word string, index string) (*SearchResult, error) {
	i, ok := binocular.indices[index]
	if !ok {
		return nil, ErrIndexNotFound
	}
	result := binocular.newSearchResult()
	result.refs = i.Search(word, 0)
	return result, nil
}

// FuzzySearch will use the distance to search the given index with the given word and returns a SearchResult.
// ErrIndexNotFound is returned if the given index does not exist.
func (binocular *Binocular) FuzzySearch(word string, index string, distance int) (*SearchResult, error) {
	i, ok := binocular.indices[index]
	if !ok {
		return nil, ErrIndexNotFound
	}
	result := binocular.newSearchResult()
	result.refs = i.Search(word, distance)
	return result, nil
}

// Remove deletes the given id from all indices and the internal data map.
// ErrRefNotFound is returned if the given id does not exist.
func (binocular *Binocular) Remove(id string) error {
	doc, ok := binocular.docs[id]
	if !ok {
		return ErrRefNotFound
	}
	for i := range doc.recordLocator {
		binocular.indices[i].Remove(id)
	}
	delete(binocular.docs, id)
	return nil
}

func (binocular *Binocular) newSearchResult() *SearchResult {
	return &SearchResult{
		binocular: binocular,
	}
}

// SearchResult holds the resulting references of your search.
type SearchResult struct {
	binocular *Binocular
	refs      []string
}

// Refs returns the list of references found for your search.
func (searchResult *SearchResult) Refs() []string {
	return searchResult.refs
}

// Collect will use the found references and returns the data associated with it.
// ErrRefNotFound is returned if a reference does not exist.
func (searchResult *SearchResult) Collect() ([]interface{}, error) {
	data := make([]interface{}, len(searchResult.refs))
	for i, ref := range searchResult.refs {
		doc, err := searchResult.binocular.Get(ref)
		if err != nil {
			return nil, ErrRefNotFound
		}
		data[i] = doc
	}
	return data, nil
}

// parses the given document for `binocular` tags and adds them to their respective Index.
func (binocular *Binocular) parseStruct(id string, doc *document, t reflect.Type) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		switch f.Type.Kind() {
		case reflect.String:
			bt := parseBinocularFieldTag(f.Tag)
			if bt == nil {
				continue
			}
			val := reflect.ValueOf(doc.Data).Field(i).String()
			binocular.indexString(id, bt.Name, val, doc)
		case reflect.Struct:
			binocular.parseStruct(id, doc, f.Type)
		case reflect.Slice:
			switch f.Type.Elem().Kind() {
			case reflect.String:
				bt := parseBinocularFieldTag(f.Tag)
				if bt == nil {
					continue
				}
				s := reflect.ValueOf(f)
				for j := 0; j < s.Len(); j++ {
					val := reflect.ValueOf(s.Index(j)).Field(i).String()
					binocular.indexString(id, bt.Name, val, doc)
				}
			case reflect.Struct:
				s := reflect.ValueOf(f)
				for j := 0; j < s.Len(); j++ {
					binocular.parseStruct(id, doc, s.Type())
				}
			}
		}
	}
}

func parseBinocularFieldTag(tag reflect.StructTag) *structtag.Tag {
	tags, err := structtag.Parse(string(tag))
	if err != nil {
		return nil
	}
	bt, err := tags.Get("binocular")
	if err != nil {
		return nil
	}
	return bt
}

func (binocular *Binocular) indexString(id, indexName, data string, doc *document) {
	if _, ok := binocular.indices[indexName]; !ok {
		binocular.indices[indexName] = NewIndex()
	}
	binocular.indices[indexName].Add(data, id)
	doc.recordLocator[indexName] = struct{}{}
}
