package binocular

import "testing"

func TestWithDefaultIndex(t *testing.T) {
	idxName := "new_default_idx"
	b := New(WithDefaultIndex(idxName))
	if b.indices[idxName] == nil {
		t.Errorf("index %s should not be nil", idxName)
	}
	if b.indices[DefaultIndex] != nil {
		t.Error("default index should be nil")
	}
}

func TestWithIndex(t *testing.T) {
	idxName := "new_idx"
	b := New(WithIndex(idxName))
	if b.indices[idxName] == nil {
		t.Errorf("index %s should not be nil", idxName)
	}
	if b.indices[DefaultIndex] == nil {
		t.Error("default index should not be nil")
	}
}

func TestBinocular_Add_String(t *testing.T) {
	b := New()
	testdata := "testdata"
	id := b.Add(testdata)
	if b.docs[id].Data != testdata {
		t.Errorf("wrong data")
	}
	if b.indices[DefaultIndex].data[testdata][0] != id {
		t.Errorf("wrong id")
	}
}

func TestBinocular_AddWithID_String(t *testing.T) {
	b := New()
	testdata := "testdata"
	id := "123"
	b.AddWithID(id, testdata)
	if b.docs[id].Data != testdata {
		t.Error("wrong data")
	}
	if b.indices[DefaultIndex].data[testdata][0] != id {
		t.Error("wrong id")
	}
}

func TestBinocular_Add_Struct(t *testing.T) {
	b := New()
	testdata := struct {
		s1 string `binocular:"default"`
		s2 string `binocular:"idx1"`
		i1 int
		x1 struct {
			s3 string `binocular:"idx2"`
			f1 float64
			b1 bool
		}
	}{
		"test",
		"idx1data",
		123,
		struct {
			s3 string `binocular:"idx2"`
			f1 float64
			b1 bool
		}{
			"idx2data",
			1.0,
			true,
		},
	}
	id := b.Add(testdata)
	if b.docs[id].Data != testdata {
		t.Error("wrong data")
	}
	if b.indices[DefaultIndex].data["test"][0] != id {
		t.Error("wrong id")
	}
	if b.indices["idx1"].data["idx1data"][0] != id {
		t.Error("wrong id")
	}
}

func TestBinocular_Get(t *testing.T) {
	b := New()
	testdata := "testdata"
	id := b.Add(testdata)
	data, err := b.Get(id)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if data != testdata {
		t.Error("wrong data")
	}
	data, err = b.Get("unknown_id")
	if err == nil {
		t.Error("expected error but got nil")
	}
	if err != ErrRefNotFound {
		t.Errorf("wrong error: %s", err)
	}
	if data != nil {
		t.Error("data should have been nil")
	}
}

func TestBinocular_Remove(t *testing.T) {
	b := New()
	testdata := "testdata"
	id := b.Add(testdata)
	err := b.Remove(id)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if b.docs[id] != nil {
		t.Error("document should not exist")
	}
	err = b.Remove("unknown_id")
	if err == nil {
		t.Error("expected error but got nil")
	}
	if err != ErrRefNotFound {
		t.Errorf("wrong error: %s", err)
	}
}

func TestBinocular_Search(t *testing.T) {
	b := New()
	testdata := "Lorem ipsum dolor sit amet"
	b.Add(testdata)
	b.Add("consetetur sadipscing elitr")
	result, err := b.Search("ipsum", DefaultIndex)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if len(result.Refs()) != 1 {
		t.Error("wrong result len")
	}
	data, err := result.Collect()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if data[0] != testdata {
		t.Error("wrong data")
	}
}

func TestBinocular_FuzzySearch(t *testing.T) {
	b := New()
	testdata := "Lorem ipsum dolor sit amet"
	b.Add(testdata)
	b.Add("consetetur sadipscing elitr")
	result, err := b.FuzzySearch("ipm", DefaultIndex, 3)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if len(result.Refs()) != 1 {
		t.Error("wrong result len")
	}
	data, err := result.Collect()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if data[0] != testdata {
		t.Error("wrong data")
	}
}

func TestBinocular_Search_IndexNotFound(t *testing.T) {
	b := New()
	result, err := b.Search("test", "unknown_idx")
	if err == nil {
		t.Error("expected error but got nil")
	}
	if err != ErrIndexNotFound {
		t.Errorf("wrong error: %s", err)
	}
	if result != nil {
		t.Error("result should be nil")
	}
}

func TestBinocular_FuzzySearch_IndexNotFound(t *testing.T) {
	b := New()
	result, err := b.FuzzySearch("test", "unknown_idx", 3)
	if err == nil {
		t.Error("expected error but got nil")
	}
	if err != ErrIndexNotFound {
		t.Errorf("wrong error: %s", err)
	}
	if result != nil {
		t.Error("result should be nil")
	}
}

func TestSearchResult_Collect_ErrRefNotFound(t *testing.T) {
	b := New()
	id := b.Add("testdata")
	result, err := b.Search("testdata", DefaultIndex)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	err = b.Remove(id)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	data, err := result.Collect()
	if err == nil {
		t.Error("expected error but got nil")
	}
	if err != ErrRefNotFound {
		t.Errorf("wrong error: %s", err)
	}
	if data != nil {
		t.Error("data should be nil")
	}
}
