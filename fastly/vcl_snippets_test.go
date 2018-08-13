package fastly

import "testing"

func TestClient_Snippets(t *testing.T) {
	t.Parallel()

	var err error
	const (
		tv              = 688
		testSnippetID   = "XwYm8yjU5QLLDvgXXXXXX"
		testSnippetName = "testsnip5"
	)

	// var tv *Version
	// record(t, "vcls/version", func(c *Client) {
	// 	tv = testVersion(t, c)
	// })

	content := `
	 # testing EdgeACL6 and EdgeDictionary6
	  declare local var.number6 STRING;
	  set var.number6 = table.lookup(demoDICTtest, client.as.number);

	  if (var.number6 == "true") {
	    set req.http.securityruleid = "num6-block";
	 error 403 "Access Denied";
	  }
	`

	returnedContent := " # testing EdgeACL6 and EdgeDictionary6\n  declare local var.asn_number6 STRING;\n  set var.asn_number6 = table.lookup(demoDICTtest, client.as.number);\n\n  if (var.asn_number6 == \"true\") {\n    set req.http.securityruleid = \"num6-block\";\n error 403 \"Access Denied\";\n  }"

	// Create
	var snippet *Snippet
	record(t, "vcl_snippets/create", func(c *Client) {
		snippet, err = c.CreateSnippet(&CreateSnippetInput{
			ServiceID: testServiceID,
			Version:   tv,
			Name:      "test-vcl-snippet",
			Type:      "recv",
			Priority:  100,
			Dynamic:   1,
			Content:   content,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if snippet.Name != "test-vcl-snippet" {
		t.Errorf("bad name: %q", snippet.Name)
	}
	if snippet.Type != "recv" {
		t.Errorf("bad type: %q", snippet.Type)
	}

	// Update
	var updateSnippet *UpdateSnippet
	record(t, "vcl_snippets/update", func(c *Client) {
		updateSnippet, err = c.UpdateSnippet(&UpdateSnippetInput{
			ServiceID: testServiceID,
			SnippetID: testSnippetID,
			Content:   content,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if updateSnippet.Content != returnedContent {
		t.Errorf("bad content: %q", updateSnippet.Content)
	}

	// Delete
	record(t, "vcl_snippets/delete", func(c *Client) {
		err = c.DeleteSnippet(&DeleteSnippetInput{
			ServiceID:   testServiceID,
			SnippetName: testSnippetName,
			Version:     tv,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}
