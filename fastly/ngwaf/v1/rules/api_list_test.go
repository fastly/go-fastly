package rules

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/fastly/go-fastly/v18/fastly"
	"github.com/fastly/go-fastly/v18/fastly/ngwaf/v1/scope"
)

func TestClient_ListRules_paginates(t *testing.T) {
	const total = 5

	var pages []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := 1
		if p := r.URL.Query().Get("page"); p != "" {
			page, _ = strconv.Atoi(p)
		}
		pages = append(pages, r.URL.Query().Get("page"))

		// Serve two rules per page.
		var data string
		for n := (page-1)*2 + 1; n <= page*2 && n <= total; n++ {
			if data != "" {
				data += ","
			}
			data += fmt.Sprintf(`{"id":"rule-%d"}`, n)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"data":[%s],"meta":{"limit":2,"total":%d}}`, data, total)
	}))
	defer srv.Close()

	c, err := fastly.NewClientForEndpoint("test", srv.URL)
	if err != nil {
		t.Fatal(err)
	}

	rs, err := List(context.TODO(), c, &ListInput{
		Scope: &scope.Scope{Type: scope.ScopeTypeWorkspace, AppliesTo: []string{"ws"}},
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(rs) != total {
		t.Fatalf("expected %d rules, got %d", total, len(rs))
	}
	for n, r := range rs {
		if want := fmt.Sprintf("rule-%d", n+1); r.RuleID != want {
			t.Errorf("rule %d: expected id %q, got %q", n, want, r.RuleID)
		}
	}
	// Page 1 omits the param; subsequent pages request 2 and 3.
	if got := fmt.Sprint(pages); got != "[ 2 3]" {
		t.Errorf("unexpected page requests: %s", got)
	}
}

func TestClient_ListRules_limit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Always serve two rules per page, regardless of the requested limit.
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"data":[{"id":"a"},{"id":"b"}],"meta":{"limit":2,"total":10}}`)
	}))
	defer srv.Close()

	c, err := fastly.NewClientForEndpoint("test", srv.URL)
	if err != nil {
		t.Fatal(err)
	}

	limit := 3
	rs, err := List(context.TODO(), c, &ListInput{
		Limit: &limit,
		Scope: &scope.Scope{Type: scope.ScopeTypeWorkspace, AppliesTo: []string{"ws"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rs) != limit {
		t.Fatalf("expected %d rules, got %d", limit, len(rs))
	}
}
