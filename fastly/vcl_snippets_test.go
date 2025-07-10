package fastly

import (
	"testing"
)

func TestClient_Snippets(t *testing.T) {
	t.Parallel()

	const (
		svName            = "snipver"
		sdName            = "snipdyn"
		svNameUpdated     = "snipverUpdated"
		defaultPriority   = "100"
		defaultDynamic    = 0
		vclContent        = "#vcl"
		vclContentUpdated = "#vclUpdated"
	)

	var tv *Version
	Record(t, "vcl_snippets/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	var err error
	var cs *Snippet

	Record(t, "vcl_snippets/create_with_required_fields_only", func(c *Client) {
		cs, err = c.CreateSnippet(&CreateSnippetInput{
			Content:        ToPointer(vclContent),
			Dynamic:        ToPointer(0),
			Name:           ToPointer(svName),
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Type:           ToPointer(SnippetTypeFetch),
		})
	})

	if err != nil {
		t.Fatal(err)
	}
	if *cs.ServiceID != TestDeliveryServiceID {
		t.Errorf("incorrect ServiceID: want %v, have %q", TestDeliveryServiceID, *cs.ServiceID)
	}
	if *cs.Name != svName {
		t.Errorf("incorrect Name: want %v, have %q", svName, *cs.Name)
	}
	if *cs.Priority != defaultPriority {
		t.Errorf("incorrect Priority: want %v, have %q", defaultPriority, *cs.Priority)
	}
	if *cs.Dynamic != defaultDynamic {
		t.Errorf("incorrect Dynamic: want %v, have %q", defaultDynamic, *cs.Dynamic)
	}
	if *cs.Content != vclContent {
		t.Errorf("incorrect Content: want %v, have %q", vclContent, *cs.Content)
	}
	if *cs.Type != SnippetTypeFetch {
		t.Errorf("incorrect Name: want %v, have %q", SnippetTypeFetch, *cs.Type)
	}

	dynamic := 1
	priority := "123"

	Record(t, "vcl_snippets/create_with_all_fields", func(c *Client) {
		cs, err = c.CreateSnippet(&CreateSnippetInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer(sdName),
			Content:        ToPointer(vclContent),
			Type:           ToPointer(SnippetTypeFetch),
			Dynamic:        ToPointer(dynamic),
			Priority:       ToPointer(priority),
		})
	})

	if err != nil {
		t.Fatal(err)
	}
	if *cs.ServiceID != TestDeliveryServiceID {
		t.Errorf("incorrect ServiceID: want %v, have %q", TestDeliveryServiceID, *cs.ServiceID)
	}
	if *cs.Name != sdName {
		t.Errorf("incorrect Name: want %v, have %q", sdName, *cs.Name)
	}
	if *cs.Priority != priority {
		t.Errorf("incorrect Priority: want %v, have %q", priority, *cs.Priority)
	}
	if *cs.Dynamic != dynamic {
		t.Errorf("incorrect Dynamic: want %v, have %q", dynamic, *cs.Dynamic)
	}
	if cs.Content != nil {
		t.Errorf("incorrect Content: want %v, have %q", "", *cs.Content) // dynamic snippets don't return content
	}
	if *cs.Type != SnippetTypeFetch {
		t.Errorf("incorrect Name: want %v, have %q", SnippetTypeFetch, *cs.Type)
	}

	var ls []*Snippet

	Record(t, "vcl_snippets/list", func(c *Client) {
		ls, err = c.ListSnippets(&ListSnippetsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
		})
	})

	if err != nil {
		t.Fatal(err)
	}

	for _, s := range ls {
		if *s.ServiceID != TestDeliveryServiceID {
			t.Errorf("incorrect ServiceID: want %v, have %q", TestDeliveryServiceID, *s.ServiceID)
		}
		if *s.Type != SnippetTypeFetch {
			t.Errorf("incorrect Name: want %v, have %q", SnippetTypeFetch, *s.Type)
		}
		if defaultDynamic == *s.Dynamic {
			if svName != *s.Name {
				t.Errorf("incorrect Name: want %v, have %q", svName, *s.Name)
			}
			if defaultPriority != *s.Priority {
				t.Errorf("incorrect Priority: want %v, have %q", defaultPriority, *s.Priority)
			}
			if vclContent != *s.Content {
				t.Errorf("incorrect Content: want %v, have %q", vclContent, *s.Content)
			}
		} else {
			if *s.Name != sdName {
				t.Errorf("incorrect Name: want %v, have %q", sdName, *s.Name)
			}
			if *s.Priority != priority {
				t.Errorf("incorrect Priority: want %v, have %q", priority, *s.Priority)
			}
			if s.Content != nil {
				t.Errorf("incorrect Content: want %v, have %q", "", *s.Content) // dynamic snippets don't return content
			}
		}
	}

	var vs *Snippet

	Record(t, "vcl_snippets/get_versioned", func(c *Client) {
		vs, err = c.GetSnippet(&GetSnippetInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           svName,
		})
	})

	if err != nil {
		t.Fatal(err)
	}
	if *vs.ServiceID != TestDeliveryServiceID {
		t.Errorf("incorrect ServiceID: want %v, have %q", TestDeliveryServiceID, *vs.ServiceID)
	}
	if *vs.Name != svName {
		t.Errorf("incorrect Name: want %v, have %q", svName, *vs.Name)
	}
	if *vs.Priority != defaultPriority {
		t.Errorf("incorrect Priority: want %v, have %q", defaultPriority, *vs.Priority)
	}
	if *vs.Dynamic != defaultDynamic {
		t.Errorf("incorrect Dynamic: want %v, have %q", defaultDynamic, *vs.Dynamic)
	}
	if *vs.Content != vclContent {
		t.Errorf("incorrect Content: want %v, have %q", vclContent, *vs.Content)
	}
	if *vs.Type != SnippetTypeFetch {
		t.Errorf("incorrect Name: want %v, have %q", SnippetTypeFetch, *vs.Type)
	}

	var ds *DynamicSnippet

	Record(t, "vcl_snippets/get_dynamic", func(c *Client) {
		ds, err = c.GetDynamicSnippet(&GetDynamicSnippetInput{
			ServiceID: TestDeliveryServiceID,
			SnippetID: *cs.SnippetID,
		})
	})

	if err != nil {
		t.Fatal(err)
	}
	if *ds.ServiceID != TestDeliveryServiceID {
		t.Errorf("incorrect ServiceID: want %v, have %q", TestDeliveryServiceID, *ds.ServiceID)
	}
	if *ds.SnippetID != *cs.SnippetID {
		t.Errorf("incorrect ID: want %v, have %q", *cs.SnippetID, *ds.SnippetID)
	}
	if *ds.Content != vclContent {
		t.Errorf("incorrect Content: want %v, have %q", vclContent, *ds.Content)
	}

	priority = "456"
	hit := SnippetTypeHit

	Record(t, "vcl_snippets/update_versioned", func(c *Client) {
		vs, err = c.UpdateSnippet(&UpdateSnippetInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           svName,
			NewName:        ToPointer(svNameUpdated),
			Priority:       ToPointer(priority),
			Content:        ToPointer(vclContentUpdated),
			Type:           &hit,
		})
	})

	if err != nil {
		t.Fatal(err)
	}
	if *vs.ServiceID != TestDeliveryServiceID {
		t.Errorf("incorrect ServiceID: want %v, have %q", TestDeliveryServiceID, *vs.ServiceID)
	}
	if *vs.Name != svNameUpdated {
		t.Errorf("incorrect Name: want %v, have %q", svNameUpdated, *vs.Name)
	}
	if *vs.Priority != priority {
		t.Errorf("incorrect Priority: want %v, have %q", priority, *vs.Priority)
	}
	if *vs.Dynamic != defaultDynamic {
		t.Errorf("incorrect Dynamic: want %v, have %q", defaultDynamic, *vs.Dynamic)
	}
	if *vs.Content != vclContentUpdated {
		t.Errorf("incorrect Content: want %v, have %q", vclContentUpdated, *vs.Content)
	}
	if *vs.Type != hit {
		t.Errorf("incorrect Name: want %v, have %q", hit, *vs.Type)
	}

	Record(t, "vcl_snippets/update_dynamic", func(c *Client) {
		ds, err = c.UpdateDynamicSnippet(&UpdateDynamicSnippetInput{
			ServiceID: TestDeliveryServiceID,
			SnippetID: *cs.SnippetID,
			Content:   ToPointer(vclContentUpdated),
		})
	})

	if err != nil {
		t.Fatal(err)
	}
	if *ds.ServiceID != TestDeliveryServiceID {
		t.Errorf("incorrect ServiceID: want %v, have %q", TestDeliveryServiceID, *ds.ServiceID)
	}
	if *ds.SnippetID != *cs.SnippetID {
		t.Errorf("incorrect ID: want %v, have %q", cs.SnippetID, *ds.SnippetID)
	}
	if *ds.Content != vclContentUpdated {
		t.Errorf("incorrect Content: want %v, have %q", vclContentUpdated, *ds.Content)
	}

	Record(t, "vcl_snippets/delete_versioned", func(c *Client) {
		err = c.DeleteSnippet(&DeleteSnippetInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           svNameUpdated,
		})
	})

	if err != nil {
		t.Fatal(err)
	}

	Record(t, "vcl_snippets/delete_dynamic", func(c *Client) {
		err = c.DeleteSnippet(&DeleteSnippetInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           sdName,
		})
	})

	if err != nil {
		t.Fatal(err)
	}
}
