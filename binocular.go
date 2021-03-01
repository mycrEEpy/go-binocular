package binocular

import (
	"reflect"

	"github.com/fatih/structtag"
	"github.com/google/uuid"
)

const (
	DefaultIndex = "_default"
)

type Binocular struct {
	docs         map[string]*Document
	indices      map[string]*Index
	DefaultIndex string
}

type Document struct {
	Data          interface{}
	recordLocator map[string]struct{}
}

type Option func(binocular *Binocular)

func New(options ...Option) *Binocular {
	binocular := &Binocular{
		docs:         map[string]*Document{},
		indices:      map[string]*Index{},
		DefaultIndex: DefaultIndex,
	}
	for _, opt := range options {
		opt(binocular)
	}
	// TODO IndexOptions?
	binocular.indices[binocular.DefaultIndex] = NewIndex()
	return binocular
}

func WithDefaultIndex(name string) Option {
	return func(binocular *Binocular) {
		binocular.DefaultIndex = name
	}
}

func (binocular *Binocular) Add(doc Document) string {
	id := uuid.New().String()
	binocular.AddWithID(id, doc)
	return id
}

func (binocular *Binocular) AddWithID(id string, doc Document) {
	doc.recordLocator = make(map[string]struct{})
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

func (binocular *Binocular) parseStruct(id string, doc *Document, t reflect.Type) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		switch f.Type.Kind() {
		case reflect.String:
			tags, err := structtag.Parse(string(f.Tag))
			if err != nil {
				break
			}
			bt, err := tags.Get("binocular")
			if err != nil {
				break
			}
			if _, ok := binocular.indices[bt.Name]; !ok {
				binocular.indices[bt.Name] = NewIndex()
			}
			binocular.indices[bt.Name].Add(reflect.ValueOf(doc.Data).Field(i).String(), id)
			doc.recordLocator[bt.Name] = struct{}{}
		case reflect.Struct:
			binocular.parseStruct(id, doc, f.Type)
		}
	}
}

func (binocular *Binocular) Get(id string) interface{} {
	doc, ok := binocular.docs[id]
	if !ok {
		return nil
	}
	return doc.Data
}

func (binocular *Binocular) Lookup(id string) (interface{}, bool) {
	doc, ok := binocular.docs[id]
	if !ok {
		return nil, false
	}
	if doc.Data == nil {
		// TODO garbage collection?
		return nil, false
	}
	return doc.Data, true
}

func (binocular *Binocular) Search(word string, index string) []interface{} {
	i, ok := binocular.indices[index]
	if !ok {
		return nil
	}
	return i.Search(word, 0)
}

func (binocular *Binocular) FuzzySearch(word string, index string, distance int) []interface{} {
	i, ok := binocular.indices[index]
	if !ok {
		return nil
	}
	return i.Search(word, distance)
}

func (binocular *Binocular) Remove(id string) {
	doc, ok := binocular.docs[id]
	if !ok {
		return
	}
	for i := range doc.recordLocator {
		binocular.indices[i].Remove(id)
	}
	delete(binocular.docs, id)
}