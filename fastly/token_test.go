package fastly

import "testing"

func TestClient_ListTokens(t *testing.T) {
	t.Parallel()

	var tokens []*Token
	var err error
	record(t, "tokens/list", func(c *Client) {
		tokens, err = c.ListTokens()
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) < 1 {
		t.Errorf("bad tokens: %v", tokens)
	}
}

func TestClient_ListCustomerTokens(t *testing.T) {
	t.Parallel()

	var tokens []*Token
	var err error
	record(t, "tokens/list_customer", func(c *Client) {
		tokens, err = c.ListCustomerTokens(&ListCustomerTokensInput{
			ID: "xxxxxxxxxxxxxxxxxxxx",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) < 1 {
		t.Errorf("bad tokens: %v", tokens)
	}
}

func TestClient_GetTokenSelf(t *testing.T) {
	t.Parallel()

	var token *Token
	var err error
	record(t, "tokens/get_self", func(c *Client) {
		token, err = c.GetTokenSelf()
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", token)
}

func TestClient_CreateToken(t *testing.T) {
	t.Parallel()

	input := &CreateTokenInput{
		Name:     "my-test-token",
		Scope:    "global",
		Username: "xxxxxxxxxxxxxxxxxxxx",
		Password: "xxxxxxxxxxxxxxxxxxxx",
	}

	var token *Token
	var err error
	record(t, "tokens/create", func(c *Client) {
		token, err = c.CreateToken(input)
	})
	if err != nil {
		t.Fatal(err)
	}

	if token.Name != input.Name {
		t.Errorf("returned invalid name, got %s, want %s", token.Name, input.Name)
	}

	if token.Scope != input.Scope {
		t.Errorf("returned invalid scope, got %s, want %s", token.Scope, input.Scope)
	}
}

func TestClient_DeleteToken(t *testing.T) {
	t.Parallel()

	input := &DeleteTokenInput{
		ID: "xxxxxxxxxxxxxxxxxxxxx",
	}

	var err error
	record(t, "tokens/delete", func(c *Client) {
		err = c.DeleteToken(input)
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_DeleteTokenSelf(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "tokens/delete_self", func(c *Client) {
		err = c.DeleteTokenSelf()
	})
	if err != nil {
		t.Fatal(err)
	}
}
