package fastly

import (
	"testing"
)

func Test_Snippets(t *testing.T) {
	t.Parallel()

	const (
		svName            = "snipver"
		sdName            = "snipdyn"
		svNameUpdated     = "snipverUpdated"
		defaultPriority   = 100
		defaultDynamic    = 0
		vclContent        = "#vcl"
		vclContentUpdated = "#vclUpdated"
	)

	var tv *Version
	record(t, "vcl_snippets/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	var err error
	var cs *Snippet

	record(t, "vcl_snippets/create_with_required_fields_only", func(c *Client) {
		cs, err = c.CreateSnippet(&CreateSnippetInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           svName,
			Content:        vclContent,
			Type:           SnippetTypeFetch,
		})
	})

	if err != nil {
		t.Fatal(err)
	}
	if testServiceID != cs.ServiceID {
		t.Errorf("incorrect ServiceID: want %v, have %q", testServiceID, cs.ServiceID)
	}
	if svName != cs.Name {
		t.Errorf("incorrect Name: want %v, have %q", svName, cs.Name)
	}
	if defaultPriority != cs.Priority {
		t.Errorf("incorrect Priority: want %v, have %q", defaultPriority, cs.Priority)
	}
	if defaultDynamic != cs.Dynamic {
		t.Errorf("incorrect Dynamic: want %v, have %q", defaultDynamic, cs.Dynamic)
	}
	if vclContent != cs.Content {
		t.Errorf("incorrect Content: want %v, have %q", vclContent, cs.Content)
	}
	if SnippetTypeFetch != cs.Type {
		t.Errorf("incorrect Name: want %v, have %q", SnippetTypeFetch, cs.Type)
	}

	dynamic := 1
	priority := 123

	record(t, "vcl_snippets/create_with_all_fields", func(c *Client) {
		cs, err = c.CreateSnippet(&CreateSnippetInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           sdName,
			Content:        vclContent,
			Type:           SnippetTypeFetch,
			Dynamic:        dynamic,
			Priority:       priority,
		})
	})

	if err != nil {
		t.Fatal(err)
	}
	if testServiceID != cs.ServiceID {
		t.Errorf("incorrect ServiceID: want %v, have %q", testServiceID, cs.ServiceID)
	}
	if sdName != cs.Name {
		t.Errorf("incorrect Name: want %v, have %q", sdName, cs.Name)
	}
	if priority != cs.Priority {
		t.Errorf("incorrect Priority: want %v, have %q", priority, cs.Priority)
	}
	if dynamic != cs.Dynamic {
		t.Errorf("incorrect Dynamic: want %v, have %q", dynamic, cs.Dynamic)
	}
	if "" != cs.Content {
		t.Errorf("incorrect Content: want %v, have %q", "", cs.Content) // dynamic snippets don't return content
	}
	if SnippetTypeFetch != cs.Type {
		t.Errorf("incorrect Name: want %v, have %q", SnippetTypeFetch, cs.Type)
	}

	var ls []*Snippet

	record(t, "vcl_snippets/list", func(c *Client) {
		ls, err = c.ListSnippets(&ListSnippetsInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
		})
	})

	if err != nil {
		t.Fatal(err)
	}

	for _, s := range ls {
		if testServiceID != s.ServiceID {
			t.Errorf("incorrect ServiceID: want %v, have %q", testServiceID, s.ServiceID)
		}
		if SnippetTypeFetch != s.Type {
			t.Errorf("incorrect Name: want %v, have %q", SnippetTypeFetch, s.Type)
		}
		if s.Dynamic == defaultDynamic {
			if svName != s.Name {
				t.Errorf("incorrect Name: want %v, have %q", svName, s.Name)
			}
			if defaultPriority != s.Priority {
				t.Errorf("incorrect Priority: want %v, have %q", defaultPriority, s.Priority)
			}
			if vclContent != s.Content {
				t.Errorf("incorrect Content: want %v, have %q", vclContent, s.Content)
			}
		} else {
			if sdName != s.Name {
				t.Errorf("incorrect Name: want %v, have %q", sdName, s.Name)
			}
			if priority != s.Priority {
				t.Errorf("incorrect Priority: want %v, have %q", priority, s.Priority)
			}
			if "" != s.Content {
				t.Errorf("incorrect Content: want %v, have %q", "", s.Content) // dynamic snippets don't return content
			}
		}
	}

	var vs *Snippet

	record(t, "vcl_snippets/get_versioned", func(c *Client) {
		vs, err = c.GetSnippet(&GetSnippetInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           svName,
		})
	})

	if err != nil {
		t.Fatal(err)
	}
	if testServiceID != vs.ServiceID {
		t.Errorf("incorrect ServiceID: want %v, have %q", testServiceID, vs.ServiceID)
	}
	if svName != vs.Name {
		t.Errorf("incorrect Name: want %v, have %q", svName, vs.Name)
	}
	if defaultPriority != vs.Priority {
		t.Errorf("incorrect Priority: want %v, have %q", defaultPriority, vs.Priority)
	}
	if defaultDynamic != vs.Dynamic {
		t.Errorf("incorrect Dynamic: want %v, have %q", defaultDynamic, vs.Dynamic)
	}
	if vclContent != vs.Content {
		t.Errorf("incorrect Content: want %v, have %q", vclContent, vs.Content)
	}
	if SnippetTypeFetch != vs.Type {
		t.Errorf("incorrect Name: want %v, have %q", SnippetTypeFetch, vs.Type)
	}

	var ds *DynamicSnippet

	record(t, "vcl_snippets/get_dynamic", func(c *Client) {
		ds, err = c.GetDynamicSnippet(&GetDynamicSnippetInput{
			ServiceID: testServiceID,
			ID:        cs.ID,
		})
	})

	if err != nil {
		t.Fatal(err)
	}
	if testServiceID != ds.ServiceID {
		t.Errorf("incorrect ServiceID: want %v, have %q", testServiceID, ds.ServiceID)
	}
	if cs.ID != ds.ID {
		t.Errorf("incorrect ID: want %v, have %q", cs.ID, ds.ID)
	}
	if vclContent != ds.Content {
		t.Errorf("incorrect Content: want %v, have %q", vclContent, ds.Content)
	}

	priority = 456
	hit := SnippetTypeHit

	record(t, "vcl_snippets/update_versioned", func(c *Client) {
		vs, err = c.UpdateSnippet(&UpdateSnippetInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           svName,
			NewName:        String(svNameUpdated),
			Priority:       Int(priority),
			Content:        String(vclContentUpdated),
			Type:           &hit,
		})
	})

	if err != nil {
		t.Fatal(err)
	}
	if testServiceID != vs.ServiceID {
		t.Errorf("incorrect ServiceID: want %v, have %q", testServiceID, vs.ServiceID)
	}
	if svNameUpdated != vs.Name {
		t.Errorf("incorrect Name: want %v, have %q", svNameUpdated, vs.Name)
	}
	if priority != vs.Priority {
		t.Errorf("incorrect Priority: want %v, have %q", priority, vs.Priority)
	}
	if defaultDynamic != vs.Dynamic {
		t.Errorf("incorrect Dynamic: want %v, have %q", defaultDynamic, vs.Dynamic)
	}
	if vclContentUpdated != vs.Content {
		t.Errorf("incorrect Content: want %v, have %q", vclContentUpdated, vs.Content)
	}
	if hit != vs.Type {
		t.Errorf("incorrect Name: want %v, have %q", hit, vs.Type)
	}

	record(t, "vcl_snippets/update_dynamic", func(c *Client) {
		ds, err = c.UpdateDynamicSnippet(&UpdateDynamicSnippetInput{
			ServiceID: testServiceID,
			ID:        cs.ID,
			Content:   String(vclContentUpdated),
		})
	})

	if err != nil {
		t.Fatal(err)
	}
	if testServiceID != ds.ServiceID {
		t.Errorf("incorrect ServiceID: want %v, have %q", testServiceID, ds.ServiceID)
	}
	if cs.ID != ds.ID {
		t.Errorf("incorrect ID: want %v, have %q", cs.ID, ds.ID)
	}
	if vclContentUpdated != ds.Content {
		t.Errorf("incorrect Content: want %v, have %q", vclContentUpdated, ds.Content)
	}

	record(t, "vcl_snippets/delete_versioned", func(c *Client) {
		err = c.DeleteSnippet(&DeleteSnippetInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           svNameUpdated,
		})
	})

	if err != nil {
		t.Fatal(err)
	}

	record(t, "vcl_snippets/delete_dynamic", func(c *Client) {
		err = c.DeleteSnippet(&DeleteSnippetInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           sdName,
		})
	})

	if err != nil {
		t.Fatal(err)
	}
}
