package fastly

import (
	"net/url"
	"strconv"

	"github.com/peterhellberg/link"
)

// PaginatorACLEntries represents a paginator.
type PaginatorACLEntries interface {
	HasNext() bool
	Remaining() int
	GetNext() ([]*ACLEntry, error)
}

// PaginatorDictionaryItems represents a paginator.
type PaginatorDictionaryItems interface {
	HasNext() bool
	Remaining() int
	GetNext() ([]*DictionaryItem, error)
}

// PaginatorServices represents a paginator.
type PaginatorServices interface {
	HasNext() bool
	Remaining() int
	GetNext() ([]*Service, error)
}

// PaginatorKVStoreEntries represents a paginator.
type PaginatorKVStoreEntries interface {
	Next() bool
	Keys() []string
	Err() error
}

// NOTE: We can't identify the underlying type of the type parameter T.
// This is because we don't assign it to any of the defined function parameters.
// If we did, then we could do this: https://go.dev/play/p/dfTMGjaSSAX.
//
// This means we have to have the user pass the API path.
// To make things easier we expose each supported path as a constant.
// e.g. fastly.ServicePath, fastly.ACLEntriesPath etc.
func NewPaginator[T any](client *Client, input *ListInput, path string) Paginator[T] {
	return &ListPaginator[T]{
		client: client,
		input:  input,
		path:   path,
	}
}

type ListInput struct {
	// Direction is the direction in which to sort results.
	Direction string
	// Page is the current page.
	Page int
	// PerPage is the number of records per page.
	PerPage int
	// Sort is the field on which to sort.
	Sort string
}

// Paginator represents a generic paginator.
type Paginator[T any] interface {
	HasNext() bool
	Remaining() int
	GetNext() ([]*T, error)
}

// ListPaginator implements the generic Paginator[N] interface.
type ListPaginator[T any] struct {
	CurrentPage int
	LastPage    int
	NextPage    int

	// Private
	client   *Client
	consumed bool
	input    *ListInput
	path     string
}

// HasNext returns a boolean indicating whether more pages are available.
func (p *ListPaginator[T]) HasNext() bool {
	return !p.consumed || p.Remaining() != 0
}

// Remaining returns the remaining page count.
func (p *ListPaginator[T]) Remaining() int {
	if p.LastPage == 0 {
		return 0
	}
	return p.LastPage - p.CurrentPage
}

// GetNext retrieves data in the next page.
func (p *ListPaginator[T]) GetNext() ([]*T, error) {
	var perPage int
	const maxPerPage = 100
	if p.input.PerPage <= 0 {
		perPage = maxPerPage
	} else {
		perPage = p.input.PerPage
	}

	// page is not specified, fetch from the beginning
	if p.input.Page <= 0 && p.CurrentPage == 0 {
		p.CurrentPage = 1
	} else {
		// page is specified, fetch from a given page
		if !p.consumed {
			p.CurrentPage = p.input.Page
		} else {
			p.CurrentPage++
		}
	}

	requestOptions := &RequestOptions{
		Params: map[string]string{
			"per_page": strconv.Itoa(perPage),
			"page":     strconv.Itoa(p.CurrentPage),
		},
	}

	if p.input.Direction != "" {
		requestOptions.Params["direction"] = p.input.Direction
	}
	if p.input.Sort != "" {
		requestOptions.Params["sort"] = p.input.Sort
	}

	resp, err := p.client.Get(p.path, requestOptions)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	for _, l := range link.ParseResponse(resp) {
		// Indicates the Link response header contained the next page instruction
		if l.Rel == "next" {
			u, _ := url.Parse(l.URI)
			query := u.Query()
			p.NextPage, _ = strconv.Atoi(query["page"][0])
		}
		// Indicates the Link response header contained the last page instruction
		if l.Rel == "last" {
			u, _ := url.Parse(l.URI)
			query := u.Query()
			p.LastPage, _ = strconv.Atoi(query["page"][0])
		}
	}

	p.consumed = true

	var s []*T
	if err := decodeBodyMap(resp.Body, &s); err != nil {
		return nil, err
	}

	return s, nil
}
