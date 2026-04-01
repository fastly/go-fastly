package operations

import (
	"context"
	"strconv"

	"github.com/fastly/go-fastly/v14/fastly"
)

const defaultPageLimit = 100

func applyPaginationParams(opts *fastly.RequestOptions, page, limit *int) {
	if limit != nil {
		opts.Params["limit"] = strconv.Itoa(*limit)
	}
	if page != nil {
		opts.Params["page"] = strconv.Itoa(*page)
	}
}

func normalizePageLimit(page, limit *int) (int, int) {
	p := 0
	if page != nil && *page >= 0 {
		p = *page
	}

	l := defaultPageLimit
	if limit != nil && *limit > 0 {
		l = *limit
	}

	return p, l
}

type pageFetcher[T any] func(ctx context.Context, c *fastly.Client, page, limit int) ([]T, int, error)

// Paginator paginates a page+limit API where the response includes meta.total.
// It does NOT rely on Link headers.
type Paginator[T any] struct {
	ctx   context.Context
	c     *fastly.Client
	fetch pageFetcher[T]

	nextPage int
	limit    int

	total   int
	fetched int
	started bool
	done    bool
}

func newPaginator[T any](ctx context.Context, c *fastly.Client, startPage, limit int, fetch pageFetcher[T]) *Paginator[T] {
	return &Paginator[T]{
		ctx:      ctx,
		c:        c,
		fetch:    fetch,
		nextPage: startPage,
		limit:    limit,
	}
}

// SetClient swaps the underlying client. Useful in tests when using fastly.Record.
func (p *Paginator[T]) SetClient(c *fastly.Client) {
	p.c = c
}

// HasNext reports whether another page fetch should be attempted.
func (p *Paginator[T]) HasNext() bool {
	if p.done {
		return false
	}
	if !p.started {
		return true
	}
	if p.total > 0 && p.fetched >= p.total {
		return false
	}
	return true
}

// GetNext fetches the next page of results.
func (p *Paginator[T]) GetNext() ([]T, error) {
	if !p.HasNext() {
		return nil, nil
	}

	p.started = true

	page := p.nextPage
	limit := p.limit

	items, total, err := p.fetch(p.ctx, p.c, page, limit)
	if err != nil {
		p.done = true
		return nil, err
	}

	p.total = total
	p.fetched += len(items)
	p.nextPage++

	// If we get an empty page, consider pagination exhausted.
	if len(items) == 0 {
		p.done = true
		return []T{}, nil
	}

	// If total is known and we've reached it, stop.
	if p.total > 0 && p.fetched >= p.total {
		p.done = true
	}

	return items, nil
}

// ---- Typed wrappers for API Security: Operations (/operations and /discovered-operations) ----

// OperationPaginator paginates GET /operations using page+limit.
type OperationPaginator = Paginator[Operation]

// NewOperationPaginator returns a paginator that iterates over operations pages.
// It respects any filters set on the input (TagID/Method/Domain/Path).
func NewOperationPaginator(ctx context.Context, c *fastly.Client, i *ListOperationsInput) *OperationPaginator {
	page, limit := normalizePageLimit(i.Page, i.Limit)

	// Copy input so callers can reuse their struct without it being mutated.
	cp := *i
	cp.Page = nil
	cp.Limit = nil

	fetch := func(ctx context.Context, c *fastly.Client, page, limit int) ([]Operation, int, error) {
		req := cp
		req.Page = &page
		req.Limit = &limit

		resp, err := ListOperations(ctx, c, &req)
		if err != nil {
			return nil, 0, err
		}
		return resp.Data, resp.Meta.Total, nil
	}

	return newPaginator[Operation](ctx, c, page, limit, fetch)
}

// ListOperationsAll retrieves all operations across pages.
func ListOperationsAll(ctx context.Context, c *fastly.Client, i *ListOperationsInput) ([]Operation, error) {
	p := NewOperationPaginator(ctx, c, i)
	var out []Operation
	for p.HasNext() {
		pageData, err := p.GetNext()
		if err != nil {
			return nil, err
		}
		out = append(out, pageData...)
	}
	return out, nil
}

// DiscoveredOperationPaginator paginates GET /discovered-operations using page+limit.
type DiscoveredOperationPaginator = Paginator[DiscoveredOperation]

// NewDiscoveredOperationPaginator returns a paginator that iterates over discovered operation pages.
// It respects any filters set on the input (Status/Method/Domain/Path).
func NewDiscoveredOperationPaginator(ctx context.Context, c *fastly.Client, i *ListDiscoveredInput) *DiscoveredOperationPaginator {
	page, limit := normalizePageLimit(i.Page, i.Limit)

	cp := *i
	cp.Page = nil
	cp.Limit = nil

	fetch := func(ctx context.Context, c *fastly.Client, page, limit int) ([]DiscoveredOperation, int, error) {
		req := cp
		req.Page = &page
		req.Limit = &limit

		resp, err := ListDiscovered(ctx, c, &req)
		if err != nil {
			return nil, 0, err
		}
		return resp.Data, resp.Meta.Total, nil
	}

	return newPaginator[DiscoveredOperation](ctx, c, page, limit, fetch)
}

// ListDiscoveredAll retrieves all discovered operations across pages.
func ListDiscoveredAll(ctx context.Context, c *fastly.Client, i *ListDiscoveredInput) ([]DiscoveredOperation, error) {
	p := NewDiscoveredOperationPaginator(ctx, c, i)
	var out []DiscoveredOperation
	for p.HasNext() {
		pageData, err := p.GetNext()
		if err != nil {
			return nil, err
		}
		out = append(out, pageData...)
	}
	return out, nil
}

// ---- Typed wrappers for API Security: Tags (/tags) ----

// TagPaginator paginates GET /tags using page+limit.
type TagPaginator = Paginator[OperationTag]

// NewTagPaginator returns a paginator that iterates over tag pages.
func NewTagPaginator(ctx context.Context, c *fastly.Client, i *ListTagsInput) *TagPaginator {
	page, limit := normalizePageLimit(i.Page, i.Limit)

	cp := *i
	cp.Page = nil
	cp.Limit = nil

	fetch := func(ctx context.Context, c *fastly.Client, page, limit int) ([]OperationTag, int, error) {
		req := cp
		req.Page = &page
		req.Limit = &limit

		resp, err := ListTags(ctx, c, &req)
		if err != nil {
			return nil, 0, err
		}
		return resp.Data, resp.Meta.Total, nil
	}

	return newPaginator[OperationTag](ctx, c, page, limit, fetch)
}

// ListTagsAll retrieves all tags across pages.
func ListTagsAll(ctx context.Context, c *fastly.Client, i *ListTagsInput) ([]OperationTag, error) {
	p := NewTagPaginator(ctx, c, i)
	var out []OperationTag
	for p.HasNext() {
		pageData, err := p.GetNext()
		if err != nil {
			return nil, err
		}
		out = append(out, pageData...)
	}
	return out, nil
}
