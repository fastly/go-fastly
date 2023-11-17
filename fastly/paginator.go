package fastly

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/peterhellberg/link"
)

// TODO: In go 1.18 (Feb 2022) use generics to reduce the duplicated code.

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

func NewPaginator[T ACLEntry | DictionaryItem | Service](client *Client, input *ListInput, path string) Paginator[T] {
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
type Paginator[T ACLEntry | DictionaryItem | Service] interface {
	HasNext() bool
	Remaining() int
	GetNext() ([]*T, error)
}

// ListPaginator implements the generic Paginator[N] interface.
type ListPaginator[T ACLEntry | DictionaryItem | Service] struct {
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
		// indicates the Link response header contained the next page instruction
		if l.Rel == "next" {
			u, _ := url.Parse(l.URI)
			query := u.Query()
			p.NextPage, _ = strconv.Atoi(query["page"][0])
		}
		// indicates the Link response header contained the last page instruction
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

	// slices.SortStableFunc(s, func(a, b *T) int {
	// 	fmt.Printf("%#v\n", a.Name)
	// 	return 0
	// })

	return s, nil
}

// The following code was copied verbatim from:
// https://cs.opensource.google/go/go/+/refs/tags/go1.21.4:src/cmp/cmp.go;l=40
//
// The only modification was to make the types private so we can replace them at
// a later date without requiring a breaking interface change/release.
//
// This functionality is exposed publicly in the go1.21 release.
// But at the time of writing we're constrained to using go1.19.

// Ordered is a constraint that permits any ordered type: any type
// that supports the operators < <= >= >.
// If future releases of Go add new ordered types,
// this constraint will be modified to include them.
//
// Note that floating-point types may contain NaN ("not-a-number") values.
// An operator such as == or < will always report false when
// comparing a NaN value with any other value, NaN or not.
// See the [Compare] function for a consistent way to compare NaN values.
type ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~string
}

// Compare returns
//
//	-1 if x is less than y,
//	 0 if x equals y,
//	+1 if x is greater than y.
//
// For floating-point types, a NaN is considered less than any non-NaN,
// a NaN is considered equal to a NaN, and -0.0 is equal to 0.0.
func compare[T ordered](x, y T) int {
	xNaN := isNaN(x)
	yNaN := isNaN(y)
	if xNaN && yNaN {
		return 0
	}
	if xNaN || x < y {
		return -1
	}
	if yNaN || x > y {
		return +1
	}
	return 0
}

// isNaN reports whether x is a NaN without requiring the math package.
// This will always return false if T is not floating-point.
func isNaN[T ordered](x T) bool {
	return x != x
}
