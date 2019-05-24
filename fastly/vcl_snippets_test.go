package fastly

import "testing"

func Test_Snippets(t *testing.T) {
	t.Parallel()

	const (
		testDynSnippetName = "testsnip5"
		testSnippetName    = "testsnip0"
	)

	var err error
	var tv *Version
	record(t, "vcl_snippets/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	content := `
	# testing EdgeACL2 and EdgeDictionary2
	 declare local var.number2 STRING;
	 set var.number2 = table.lookup(demoDICTtest, client.as.number);

	 if (var.number2 == "true") {
	   set req.http.securityruleid = "num2-block";
	 error 403 "Access Denied";
	  }
    `

	dynContent := `
	 # testing EdgeACL6 and EdgeDictionary6
	  declare local var.number6 STRING;
	  set var.number6 = table.lookup(demoDICTtest, client.as.number);

	  if (var.number6 == "true") {
	    set req.http.securityruleid = "num6-block";
	 error 403 "Access Denied";
	  }
	`
	updatedDynContent := `
	 # testing EdgeACL5 and EdgeDictionary5
	 declare local var.number5 STRING;
	 set var.number5 = table.lookup(demoDICTtest, client.as.number);

	 if (var.number5 == "true") {
	 set req.http.securityruleid = "num5-block";
	 error 403 "Access Denied";
	  }
    `

	// Create
	var cs *Snippet
	record(t, "vcl_snippets/create", func(c *Client) {
		cs, err = c.CreateSnippet(&CreateSnippetInput{
			Service:  testServiceID,
			Version:  tv.Number,
			Name:     testSnippetName,
			Type:     SnippetTypeRecv,
			Priority: 100,
			Dynamic:  0,
			Content:  content,
		})
	})

	if err != nil {
		t.Fatal(err)
	}
	if cs.ServiceID != testServiceID {
		t.Errorf("bad sID: %q", cs.ServiceID)
	}
	if cs.Name != testSnippetName {
		t.Errorf("bad name: %q", cs.Name)
	}
	if cs.Type != SnippetTypeRecv {
		t.Errorf("bad type: %q", cs.Type)
	}
	if cs.Content != content {
		t.Errorf("bad content: %q", cs.Content)
	}

	// Create Dynamic
	var cds *Snippet
	record(t, "vcl_snippets/create_dyn", func(c *Client) {
		cds, err = c.CreateSnippet(&CreateSnippetInput{
			Service:  testServiceID,
			Version:  tv.Number,
			Name:     testDynSnippetName,
			Type:     SnippetTypeRecv,
			Priority: 100,
			Dynamic:  1,
			Content:  dynContent,
		})
	})

	if err != nil {
		t.Fatal(err)
	}
	if cds.ServiceID != testServiceID {
		t.Errorf("bad sID: %q", cds.ServiceID)
	}
	if cds.Name != testDynSnippetName {
		t.Errorf("bad name: %q", cds.Name)
	}
	if cds.Type != SnippetTypeRecv {
		t.Errorf("bad type: %q", cds.Type)
	}

	// Update Dynamic
	var uds *DynamicSnippet
	record(t, "vcl_snippets/update_dyn", func(c *Client) {
		uds, err = c.UpdateDynamicSnippet(&UpdateDynamicSnippetInput{
			Service: testServiceID,
			ID:      cds.ID,
			Content: updatedDynContent,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if uds.Content != updatedDynContent {
		t.Errorf("bad content: %q", uds.Content)
	}

	// Delete
	record(t, "vcl_snippets/delete", func(c *Client) {
		err = c.DeleteSnippet(&DeleteSnippetInput{
			Service: testServiceID,
			Name:    testDynSnippetName,
			Version: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// GETDynamicSnippet
	var ds *DynamicSnippet
	record(t, "vcl_snippets/get_dynamic", func(c *Client) {
		ds, err = c.GetDynamicSnippet(&GetDynamicSnippetInput{
			Service: testServiceID,
			ID:      cds.ID,
		})

	})
	if err != nil {
		t.Fatal(err)
	}
	if ds.ServiceID != testServiceID {
		t.Errorf("bad sID: %q", ds.ServiceID)
	}
	if ds.ID != cds.ID {
		t.Errorf("bad snipID: %q", ds.ID)
	}

	// GETSnippet
	var gs *Snippet
	record(t, "vcl_snippets/get", func(c *Client) {
		gs, err = c.GetSnippet(&GetSnippetInput{
			Service: testServiceID,
			Name:    testSnippetName,
			Version: tv.Number,
		})

	})
	if err != nil {
		t.Fatal(err)
	}
	if gs.Name != testSnippetName {
		t.Errorf("bad name: %q", gs.Name)
	}
	if gs.ServiceID != testServiceID {
		t.Errorf("bad service: %q", gs.ServiceID)
	}
	if gs.Content != content {
		t.Errorf("bad content: %q", gs.Content)
	}

	// Update
	var us *Snippet
	record(t, "vcl_snippets/update", func(c *Client) {
		us, err = c.UpdateSnippet(&UpdateSnippetInput{
			Service:  testServiceID,
			Name:     testSnippetName,
			NewName:  "newTestSnippetName",
			Content:  updatedDynContent,
			Dynamic:  0,
			Type:     "none",
			Priority: 50,
			Version:  tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if us.Name != "newTestSnippetName" {
		t.Errorf("bad updated name")
	}
	if us.Priority != 50 {
		t.Errorf("bad priority: %d", us.Priority)
	}

	if us.Content != updatedDynContent {
		t.Errorf("bad content: %q", us.Content)
	}

	if us.Type != "none" {
		t.Errorf("bad type: %s", us.Type)
	}

	// ListSnippets
	var sl []*Snippet
	record(t, "vcl_snippets/list", func(c *Client) {
		sl, err = c.ListSnippets(&ListSnippetsInput{
			Service: testServiceID,
			Version: tv.Number,
		})

	})
	if err != nil {
		t.Fatal(err)
	}
	for _, x := range sl {
		if x.ServiceID != testServiceID {
			t.Errorf("bad service: %q", x.ServiceID)
		}
		if x.Version != tv.Number {
			t.Errorf("bad Version: %q", x.Version)
		}
	}

	_, err = testClient.GetDynamicSnippet(&GetDynamicSnippetInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}
	_, err = testClient.GetDynamicSnippet(&GetDynamicSnippetInput{
		Service: testServiceID,
		ID:      "",
	})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateSnippet(&CreateSnippetInput{
		Service:  testServiceID,
		Version:  tv.Number,
		Name:     testSnippetName,
		Type:     SnippetTypeRecv,
		Priority: 100,
		Dynamic:  0,
		Content:  "",
	})

	if err != ErrMissingContent {
		t.Errorf("bad error: %s", err)
	}
}
