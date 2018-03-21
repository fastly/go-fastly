package jsonapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestUnmarshall_attrStringSlice(t *testing.T) {
	out := &Book{}
	tags := []string{"fiction", "sale"}
	data := map[string]interface{}{
		"data": map[string]interface{}{
			"type":       "books",
			"id":         "1",
			"attributes": map[string]interface{}{"tags": tags},
		},
	}
	b, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	if err := UnmarshalPayload(bytes.NewReader(b), out); err != nil {
		t.Fatal(err)
	}

	if e, a := len(tags), len(out.Tags); e != a {
		t.Fatalf("Was expecting %d tags, got %d", e, a)
	}

	sort.Strings(tags)
	sort.Strings(out.Tags)

	for i, tag := range tags {
		if e, a := tag, out.Tags[i]; e != a {
			t.Fatalf("At index %d, was expecting %s got %s", i, e, a)
		}
	}
}

func TestUnmarshalToStructWithPointerAttr(t *testing.T) {
	out := new(WithPointer)
	in := map[string]interface{}{
		"name":      "The name",
		"is-active": true,
		"int-val":   8,
		"float-val": 1.1,
	}
	if err := UnmarshalPayload(sampleWithPointerPayload(in), out); err != nil {
		t.Fatal(err)
	}
	if *out.Name != "The name" {
		t.Fatalf("Error unmarshalling to string ptr")
	}
	if !*out.IsActive {
		t.Fatalf("Error unmarshalling to bool ptr")
	}
	if *out.IntVal != 8 {
		t.Fatalf("Error unmarshalling to int ptr")
	}
	if *out.FloatVal != 1.1 {
		t.Fatalf("Error unmarshalling to float ptr")
	}
}

func TestUnmarshalPayload_ptrsAllNil(t *testing.T) {
	out := new(WithPointer)
	if err := UnmarshalPayload(
		strings.NewReader(`{"data": {}}`), out); err != nil {
		t.Fatalf("Error unmarshalling to Foo")
	}

	if out.ID != nil {
		t.Fatalf("Error unmarshalling; expected ID ptr to be nil")
	}
}

func TestUnmarshalPayloadWithPointerID(t *testing.T) {
	out := new(WithPointer)
	attrs := map[string]interface{}{}

	if err := UnmarshalPayload(sampleWithPointerPayload(attrs), out); err != nil {
		t.Fatalf("Error unmarshalling to Foo")
	}

	// these were present in the payload -- expect val to be not nil
	if out.ID == nil {
		t.Fatalf("Error unmarshalling; expected ID ptr to be not nil")
	}
	if e, a := uint64(2), *out.ID; e != a {
		t.Fatalf("Was expecting the ID to have a value of %d, got %d", e, a)
	}
}

func TestUnmarshalPayloadWithPointerAttr_AbsentVal(t *testing.T) {
	out := new(WithPointer)
	in := map[string]interface{}{
		"name":      "The name",
		"is-active": true,
	}

	if err := UnmarshalPayload(sampleWithPointerPayload(in), out); err != nil {
		t.Fatalf("Error unmarshalling to Foo")
	}

	// these were present in the payload -- expect val to be not nil
	if out.Name == nil || out.IsActive == nil {
		t.Fatalf("Error unmarshalling; expected ptr to be not nil")
	}

	// these were absent in the payload -- expect val to be nil
	if out.IntVal != nil || out.FloatVal != nil {
		t.Fatalf("Error unmarshalling; expected ptr to be nil")
	}
}

func TestUnmarshalToStructWithPointerAttr_BadType(t *testing.T) {
	out := new(WithPointer)
	in := map[string]interface{}{
		"name": true, // This is the wrong type.
	}
	expectedErrorMessage := ErrUnsupportedPtrType.Error()

	err := UnmarshalPayload(sampleWithPointerPayload(in), out)

	if err == nil {
		t.Fatalf("Expected error due to invalid type.")
	}
	if err.Error() != expectedErrorMessage {
		t.Fatalf("Unexpected error message: %s", err.Error())
	}
}

func TestStringPointerField(t *testing.T) {
	// Build Book payload
	description := "Hello World!"
	data := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "books",
			"id":   "5",
			"attributes": map[string]interface{}{
				"author":      "aren55555",
				"description": description,
				"isbn":        "",
			},
		},
	}
	payload, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	// Parse JSON API payload
	book := new(Book)
	if err := UnmarshalPayload(bytes.NewReader(payload), book); err != nil {
		t.Fatal(err)
	}

	if book.Description == nil {
		t.Fatal("Was not expecting a nil pointer for book.Description")
	}
	if expected, actual := description, *book.Description; expected != actual {
		t.Fatalf("Was expecting descript to be `%s`, got `%s`", expected, actual)
	}
}

func TestMalformedTag(t *testing.T) {
	out := new(BadModel)
	err := UnmarshalPayload(samplePayload(), out)
	if err == nil || err != ErrBadJSONAPIStructTag {
		t.Fatalf("Did not error out with wrong number of arguments in tag")
	}
}

func TestUnmarshalInvalidJSON(t *testing.T) {
	in := strings.NewReader("{}")
	out := new(Blog)

	err := UnmarshalPayload(in, out)

	if err == nil {
		t.Fatalf("Did not error out the invalid JSON.")
	}
}

func TestUnmarshalInvalidJSON_BadType(t *testing.T) {
	var badTypeTests = []struct {
		Field    string
		BadValue interface{}
		Error    error
	}{ // The `Field` values here correspond to the `ModelBadTypes` jsonapi fields.
		{Field: "string_field", BadValue: 0, Error: ErrUnknownFieldNumberType},  // Expected string.
		{Field: "float_field", BadValue: "A string.", Error: ErrInvalidType},    // Expected float64.
		{Field: "time_field", BadValue: "A string.", Error: ErrInvalidTime},     // Expected int64.
		{Field: "time_ptr_field", BadValue: "A string.", Error: ErrInvalidTime}, // Expected *time / int64.
	}
	for _, test := range badTypeTests {
		t.Run(fmt.Sprintf("Test_%s", test.Field), func(t *testing.T) {
			out := new(ModelBadTypes)
			in := map[string]interface{}{}
			in[test.Field] = test.BadValue
			expectedErrorMessage := test.Error.Error()

			err := UnmarshalPayload(samplePayloadWithBadTypes(in), out)

			if err == nil {
				t.Fatalf("Expected error due to invalid type.")
			}
			if err.Error() != expectedErrorMessage {
				t.Fatalf("Unexpected error message: %s", err.Error())
			}
		})
	}
}

func TestUnmarshalSetsID(t *testing.T) {
	in := samplePayloadWithID()
	out := new(Blog)

	if err := UnmarshalPayload(in, out); err != nil {
		t.Fatal(err)
	}

	if out.ID != 2 {
		t.Fatalf("Did not set ID on dst interface")
	}
}

func TestUnmarshal_nonNumericID(t *testing.T) {
	data := samplePayloadWithoutIncluded()
	data["data"].(map[string]interface{})["id"] = "non-numeric-id"
	payload, _ := payload(data)
	in := bytes.NewReader(payload)
	out := new(Post)

	if err := UnmarshalPayload(in, out); err != ErrBadJSONAPIID {
		t.Fatalf(
			"Was expecting a `%s` error, got `%s`",
			ErrBadJSONAPIID,
			err,
		)
	}
}

func TestUnmarshalSetsAttrs(t *testing.T) {
	out, err := unmarshalSamplePayload()
	if err != nil {
		t.Fatal(err)
	}

	if out.CreatedAt.IsZero() {
		t.Fatalf("Did not parse time")
	}

	if out.ViewCount != 1000 {
		t.Fatalf("View count not properly serialized")
	}
}

func TestUnmarshalParsesISO8601(t *testing.T) {
	payload := &OnePayload{
		Data: &Node{
			Type: "timestamps",
			Attributes: map[string]interface{}{
				"timestamp": "2016-08-17T08:27:12Z",
			},
		},
	}

	in := bytes.NewBuffer(nil)
	json.NewEncoder(in).Encode(payload)

	out := new(Timestamp)

	if err := UnmarshalPayload(in, out); err != nil {
		t.Fatal(err)
	}

	expected := time.Date(2016, 8, 17, 8, 27, 12, 0, time.UTC)

	if !out.Time.Equal(expected) {
		t.Fatal("Parsing the ISO8601 timestamp failed")
	}
}

func TestUnmarshalParsesISO8601TimePointer(t *testing.T) {
	payload := &OnePayload{
		Data: &Node{
			Type: "timestamps",
			Attributes: map[string]interface{}{
				"next": "2016-08-17T08:27:12Z",
			},
		},
	}

	in := bytes.NewBuffer(nil)
	json.NewEncoder(in).Encode(payload)

	out := new(Timestamp)

	if err := UnmarshalPayload(in, out); err != nil {
		t.Fatal(err)
	}

	expected := time.Date(2016, 8, 17, 8, 27, 12, 0, time.UTC)

	if !out.Next.Equal(expected) {
		t.Fatal("Parsing the ISO8601 timestamp failed")
	}
}

func TestUnmarshalInvalidISO8601(t *testing.T) {
	payload := &OnePayload{
		Data: &Node{
			Type: "timestamps",
			Attributes: map[string]interface{}{
				"timestamp": "17 Aug 16 08:027 MST",
			},
		},
	}

	in := bytes.NewBuffer(nil)
	json.NewEncoder(in).Encode(payload)

	out := new(Timestamp)

	if err := UnmarshalPayload(in, out); err != ErrInvalidISO8601 {
		t.Fatalf("Expected ErrInvalidISO8601, got %v", err)
	}
}

func TestUnmarshalRelationshipsWithoutIncluded(t *testing.T) {
	data, _ := payload(samplePayloadWithoutIncluded())
	in := bytes.NewReader(data)
	out := new(Post)

	if err := UnmarshalPayload(in, out); err != nil {
		t.Fatal(err)
	}

	// Verify each comment has at least an ID
	for _, comment := range out.Comments {
		if comment.ID == 0 {
			t.Fatalf("The comment did not have an ID")
		}
	}
}

func TestUnmarshalRelationships(t *testing.T) {
	out, err := unmarshalSamplePayload()
	if err != nil {
		t.Fatal(err)
	}

	if out.CurrentPost == nil {
		t.Fatalf("Current post was not materialized")
	}

	if out.CurrentPost.Title != "Bas" || out.CurrentPost.Body != "Fuubar" {
		t.Fatalf("Attributes were not set")
	}

	if len(out.Posts) != 2 {
		t.Fatalf("There should have been 2 posts")
	}
}

func TestUnmarshalNullRelationship(t *testing.T) {
	sample := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "posts",
			"id":   "1",
			"attributes": map[string]interface{}{
				"body":  "Hello",
				"title": "World",
			},
			"relationships": map[string]interface{}{
				"latest_comment": map[string]interface{}{
					"data": nil, // empty to-one relationship
				},
			},
		},
	}
	data, err := json.Marshal(sample)
	if err != nil {
		t.Fatal(err)
	}

	in := bytes.NewReader(data)
	out := new(Post)

	if err := UnmarshalPayload(in, out); err != nil {
		t.Fatal(err)
	}

	if out.LatestComment != nil {
		t.Fatalf("Latest Comment was not set to nil")
	}
}

func TestUnmarshalNullRelationshipInSlice(t *testing.T) {
	sample := map[string]interface{}{
		"data": map[string]interface{}{
			"type": "posts",
			"id":   "1",
			"attributes": map[string]interface{}{
				"body":  "Hello",
				"title": "World",
			},
			"relationships": map[string]interface{}{
				"comments": map[string]interface{}{
					"data": []interface{}{}, // empty to-many relationships
				},
			},
		},
	}
	data, err := json.Marshal(sample)
	if err != nil {
		t.Fatal(err)
	}

	in := bytes.NewReader(data)
	out := new(Post)

	if err := UnmarshalPayload(in, out); err != nil {
		t.Fatal(err)
	}

	if len(out.Comments) != 0 {
		t.Fatalf("Wrong number of comments; Comments should be empty")
	}
}

func TestUnmarshalNestedRelationships(t *testing.T) {
	out, err := unmarshalSamplePayload()
	if err != nil {
		t.Fatal(err)
	}

	if out.CurrentPost == nil {
		t.Fatalf("Current post was not materialized")
	}

	if out.CurrentPost.Comments == nil {
		t.Fatalf("Did not materialize nested records, comments")
	}

	if len(out.CurrentPost.Comments) != 2 {
		t.Fatalf("Wrong number of comments")
	}
}

func TestUnmarshalRelationshipsSerializedEmbedded(t *testing.T) {
	out := sampleSerializedEmbeddedTestModel()

	if out.CurrentPost == nil {
		t.Fatalf("Current post was not materialized")
	}

	if out.CurrentPost.Title != "Foo" || out.CurrentPost.Body != "Bar" {
		t.Fatalf("Attributes were not set")
	}

	if len(out.Posts) != 2 {
		t.Fatalf("There should have been 2 posts")
	}

	if out.Posts[0].LatestComment.Body != "foo" {
		t.Fatalf("The comment body was not set")
	}
}

func TestUnmarshalNestedRelationshipsEmbedded(t *testing.T) {
	out := bytes.NewBuffer(nil)
	if err := MarshalOnePayloadEmbedded(out, testModel()); err != nil {
		t.Fatal(err)
	}

	model := new(Blog)

	if err := UnmarshalPayload(out, model); err != nil {
		t.Fatal(err)
	}

	if model.CurrentPost == nil {
		t.Fatalf("Current post was not materialized")
	}

	if model.CurrentPost.Comments == nil {
		t.Fatalf("Did not materialize nested records, comments")
	}

	if len(model.CurrentPost.Comments) != 2 {
		t.Fatalf("Wrong number of comments")
	}

	if model.CurrentPost.Comments[0].Body != "foo" {
		t.Fatalf("Comment body not set")
	}
}

func TestUnmarshalRelationshipsSideloaded(t *testing.T) {
	payload := samplePayloadWithSideloaded()
	out := new(Blog)

	if err := UnmarshalPayload(payload, out); err != nil {
		t.Fatal(err)
	}

	if out.CurrentPost == nil {
		t.Fatalf("Current post was not materialized")
	}

	if out.CurrentPost.Title != "Foo" || out.CurrentPost.Body != "Bar" {
		t.Fatalf("Attributes were not set")
	}

	if len(out.Posts) != 2 {
		t.Fatalf("There should have been 2 posts")
	}
}

func TestUnmarshalNestedRelationshipsSideloaded(t *testing.T) {
	payload := samplePayloadWithSideloaded()
	out := new(Blog)

	if err := UnmarshalPayload(payload, out); err != nil {
		t.Fatal(err)
	}

	if out.CurrentPost == nil {
		t.Fatalf("Current post was not materialized")
	}

	if out.CurrentPost.Comments == nil {
		t.Fatalf("Did not materialize nested records, comments")
	}

	if len(out.CurrentPost.Comments) != 2 {
		t.Fatalf("Wrong number of comments")
	}

	if out.CurrentPost.Comments[0].Body != "foo" {
		t.Fatalf("Comment body not set")
	}
}

func TestUnmarshalNestedRelationshipsEmbedded_withClientIDs(t *testing.T) {
	model := new(Blog)

	if err := UnmarshalPayload(samplePayload(), model); err != nil {
		t.Fatal(err)
	}

	if model.Posts[0].ClientID == "" {
		t.Fatalf("ClientID not set from request on related record")
	}
}

func unmarshalSamplePayload() (*Blog, error) {
	in := samplePayload()
	out := new(Blog)

	if err := UnmarshalPayload(in, out); err != nil {
		return nil, err
	}

	return out, nil
}

func TestUnmarshalManyPayload(t *testing.T) {
	sample := map[string]interface{}{
		"data": []interface{}{
			map[string]interface{}{
				"type": "posts",
				"id":   "1",
				"attributes": map[string]interface{}{
					"body":  "First",
					"title": "Post",
				},
			},
			map[string]interface{}{
				"type": "posts",
				"id":   "2",
				"attributes": map[string]interface{}{
					"body":  "Second",
					"title": "Post",
				},
			},
		},
	}

	data, err := json.Marshal(sample)
	if err != nil {
		t.Fatal(err)
	}
	in := bytes.NewReader(data)

	posts, err := UnmarshalManyPayload(in, reflect.TypeOf(new(Post)))
	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != 2 {
		t.Fatal("Wrong number of posts")
	}

	for _, p := range posts {
		_, ok := p.(*Post)
		if !ok {
			t.Fatal("Was expecting a Post")
		}
	}
}

func TestManyPayload_withLinks(t *testing.T) {
	firstPageURL := "http://somesite.com/movies?page[limit]=50&page[offset]=50"
	prevPageURL := "http://somesite.com/movies?page[limit]=50&page[offset]=0"
	nextPageURL := "http://somesite.com/movies?page[limit]=50&page[offset]=100"
	lastPageURL := "http://somesite.com/movies?page[limit]=50&page[offset]=500"

	sample := map[string]interface{}{
		"data": []interface{}{
			map[string]interface{}{
				"type": "posts",
				"id":   "1",
				"attributes": map[string]interface{}{
					"body":  "First",
					"title": "Post",
				},
			},
			map[string]interface{}{
				"type": "posts",
				"id":   "2",
				"attributes": map[string]interface{}{
					"body":  "Second",
					"title": "Post",
				},
			},
		},
		"links": map[string]interface{}{
			KeyFirstPage:    firstPageURL,
			KeyPreviousPage: prevPageURL,
			KeyNextPage:     nextPageURL,
			KeyLastPage:     lastPageURL,
		},
	}

	data, err := json.Marshal(sample)
	if err != nil {
		t.Fatal(err)
	}
	in := bytes.NewReader(data)

	payload := new(ManyPayload)
	if err = json.NewDecoder(in).Decode(payload); err != nil {
		t.Fatal(err)
	}

	if payload.Links == nil {
		t.Fatal("Was expecting a non nil ptr Link field")
	}

	links := *payload.Links

	first, ok := links[KeyFirstPage]
	if !ok {
		t.Fatal("Was expecting a non nil ptr Link field")
	}
	if e, a := firstPageURL, first; e != a {
		t.Fatalf("Was expecting links.%s to have a value of %s, got %s", KeyFirstPage, e, a)
	}

	prev, ok := links[KeyPreviousPage]
	if !ok {
		t.Fatal("Was expecting a non nil ptr Link field")
	}
	if e, a := prevPageURL, prev; e != a {
		t.Fatalf("Was expecting links.%s to have a value of %s, got %s", KeyPreviousPage, e, a)
	}

	next, ok := links[KeyNextPage]
	if !ok {
		t.Fatal("Was expecting a non nil ptr Link field")
	}
	if e, a := nextPageURL, next; e != a {
		t.Fatalf("Was expecting links.%s to have a value of %s, got %s", KeyNextPage, e, a)
	}

	last, ok := links[KeyLastPage]
	if !ok {
		t.Fatal("Was expecting a non nil ptr Link field")
	}
	if e, a := lastPageURL, last; e != a {
		t.Fatalf("Was expecting links.%s to have a value of %s, got %s", KeyLastPage, e, a)
	}
}

func samplePayloadWithoutIncluded() map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]interface{}{
			"type": "posts",
			"id":   "1",
			"attributes": map[string]interface{}{
				"body":  "Hello",
				"title": "World",
			},
			"relationships": map[string]interface{}{
				"comments": map[string]interface{}{
					"data": []interface{}{
						map[string]interface{}{
							"type": "comments",
							"id":   "123",
						},
						map[string]interface{}{
							"type": "comments",
							"id":   "456",
						},
					},
				},
				"latest_comment": map[string]interface{}{
					"data": map[string]interface{}{
						"type": "comments",
						"id":   "55555",
					},
				},
			},
		},
	}
}

func payload(data map[string]interface{}) (result []byte, err error) {
	result, err = json.Marshal(data)
	return
}

func samplePayload() io.Reader {
	payload := &OnePayload{
		Data: &Node{
			Type: "blogs",
			Attributes: map[string]interface{}{
				"title":      "New blog",
				"created_at": 1436216820,
				"view_count": 1000,
			},
			Relationships: map[string]interface{}{
				"posts": &RelationshipManyNode{
					Data: []*Node{
						&Node{
							Type: "posts",
							Attributes: map[string]interface{}{
								"title": "Foo",
								"body":  "Bar",
							},
							ClientID: "1",
						},
						&Node{
							Type: "posts",
							Attributes: map[string]interface{}{
								"title": "X",
								"body":  "Y",
							},
							ClientID: "2",
						},
					},
				},
				"current_post": &RelationshipOneNode{
					Data: &Node{
						Type: "posts",
						Attributes: map[string]interface{}{
							"title": "Bas",
							"body":  "Fuubar",
						},
						ClientID: "3",
						Relationships: map[string]interface{}{
							"comments": &RelationshipManyNode{
								Data: []*Node{
									&Node{
										Type: "comments",
										Attributes: map[string]interface{}{
											"body": "Great post!",
										},
										ClientID: "4",
									},
									&Node{
										Type: "comments",
										Attributes: map[string]interface{}{
											"body": "Needs some work!",
										},
										ClientID: "5",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	out := bytes.NewBuffer(nil)
	json.NewEncoder(out).Encode(payload)

	return out
}

func samplePayloadWithID() io.Reader {
	payload := &OnePayload{
		Data: &Node{
			ID:   "2",
			Type: "blogs",
			Attributes: map[string]interface{}{
				"title":      "New blog",
				"view_count": 1000,
			},
		},
	}

	out := bytes.NewBuffer(nil)
	json.NewEncoder(out).Encode(payload)

	return out
}

func samplePayloadWithBadTypes(m map[string]interface{}) io.Reader {
	payload := &OnePayload{
		Data: &Node{
			ID:         "2",
			Type:       "badtypes",
			Attributes: m,
		},
	}

	out := bytes.NewBuffer(nil)
	json.NewEncoder(out).Encode(payload)

	return out
}

func sampleWithPointerPayload(m map[string]interface{}) io.Reader {
	payload := &OnePayload{
		Data: &Node{
			ID:         "2",
			Type:       "with-pointers",
			Attributes: m,
		},
	}

	out := bytes.NewBuffer(nil)
	json.NewEncoder(out).Encode(payload)

	return out
}

func testModel() *Blog {
	return &Blog{
		ID:        5,
		ClientID:  "1",
		Title:     "Title 1",
		CreatedAt: time.Now(),
		Posts: []*Post{
			&Post{
				ID:    1,
				Title: "Foo",
				Body:  "Bar",
				Comments: []*Comment{
					&Comment{
						ID:   1,
						Body: "foo",
					},
					&Comment{
						ID:   2,
						Body: "bar",
					},
				},
				LatestComment: &Comment{
					ID:   1,
					Body: "foo",
				},
			},
			&Post{
				ID:    2,
				Title: "Fuubar",
				Body:  "Bas",
				Comments: []*Comment{
					&Comment{
						ID:   1,
						Body: "foo",
					},
					&Comment{
						ID:   3,
						Body: "bas",
					},
				},
				LatestComment: &Comment{
					ID:   1,
					Body: "foo",
				},
			},
		},
		CurrentPost: &Post{
			ID:    1,
			Title: "Foo",
			Body:  "Bar",
			Comments: []*Comment{
				&Comment{
					ID:   1,
					Body: "foo",
				},
				&Comment{
					ID:   2,
					Body: "bar",
				},
			},
			LatestComment: &Comment{
				ID:   1,
				Body: "foo",
			},
		},
	}
}

func samplePayloadWithSideloaded() io.Reader {
	testModel := testModel()

	out := bytes.NewBuffer(nil)
	MarshalPayload(out, testModel)

	return out
}

func sampleSerializedEmbeddedTestModel() *Blog {
	out := bytes.NewBuffer(nil)
	MarshalOnePayloadEmbedded(out, testModel())

	blog := new(Blog)
	UnmarshalPayload(out, blog)

	return blog
}
