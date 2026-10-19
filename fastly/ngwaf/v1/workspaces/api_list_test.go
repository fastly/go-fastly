package workspaces

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/fastly/go-fastly/v17/fastly"
)

func TestClient_ListWorkspaces_paginates(t *testing.T) {
	const total = 5

	var pages []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := 1
		if p := r.URL.Query().Get("page"); p != "" {
			page, _ = strconv.Atoi(p)
		}
		pages = append(pages, r.URL.Query().Get("page"))

		// Serve two workspaces per page.
		var data string
		for n := (page-1)*2 + 1; n <= page*2 && n <= total; n++ {
			if data != "" {
				data += ","
			}
			data += fmt.Sprintf(`{"id":"ws-%d"}`, n)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"data":[%s],"meta":{"limit":2,"total":%d}}`, data, total)
	}))
	defer srv.Close()

	c, err := fastly.NewClientForEndpoint("test", srv.URL)
	if err != nil {
		t.Fatal(err)
	}

	limit := 2
	wss, err := List(context.TODO(), c, &ListInput{Limit: &limit})
	if err != nil {
		t.Fatal(err)
	}

	if len(wss) != total {
		t.Fatalf("expected %d workspaces, got %d", total, len(wss))
	}
	for n, ws := range wss {
		if want := fmt.Sprintf("ws-%d", n+1); ws.WorkspaceID != want {
			t.Errorf("workspace %d: expected id %q, got %q", n, want, ws.WorkspaceID)
		}
	}
	// Page 1 omits the param; subsequent pages request 2 and 3.
	if got := fmt.Sprint(pages); got != "[ 2 3]" {
		t.Errorf("unexpected page requests: %s", got)
	}
}
