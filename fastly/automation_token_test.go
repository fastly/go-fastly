package fastly

import (
	"context"
	"testing"
)

func TestClient_ListAutomationTokens(t *testing.T) {
	t.Parallel()

	var tokens []*AutomationToken
	var err error
	Record(t, "automation_tokens/list", func(c *Client) {
		tokens, err = c.ListAutomationTokens(context.TODO())
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) < 1 {
		t.Errorf("bad tokens: %v", tokens)
	}
}

func TestClient_GetAutomationToken(t *testing.T) {
	t.Parallel()

	input := &GetAutomationTokenInput{
		TokenID: "XXXXXXXXXXXXXXXXXXXXXX",
	}

	var token *AutomationToken
	var err error
	Record(t, "automation_tokens/get", func(c *Client) {
		token, err = c.GetAutomationToken(context.TODO(), input)
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", token)
}

func TestClient_CreateAutomationToken(t *testing.T) {
	t.Parallel()

	input := &CreateAutomationTokenInput{
		Name:     "my-test-token",
		Role:     EngineerRole,
		Scope:    ToPointer(GlobalScope),
		Username: ToPointer("XXXXXXXXXXXXXXXXXXXXXX"),
		Password: ToPointer("XXXXXXXXXXXXXXXXXXXXXX"),
	}

	var token *AutomationToken
	var err error
	Record(t, "automation_tokens/create", func(c *Client) {
		token, err = c.CreateAutomationToken(context.TODO(), input)
	})
	if err != nil {
		t.Fatal(err)
	}

	if *token.Name != input.Name {
		t.Errorf("returned invalid name, got %s, want %s", *token.Name, input.Name)
	}
	if *token.Scope != *input.Scope {
		t.Errorf("returned invalid scope, got %s, want %s", *token.Scope, *input.Scope)
	}
}

func TestClient_DeleteAutomationToken(t *testing.T) {
	t.Parallel()

	input := &DeleteAutomationTokenInput{
		TokenID: "XXXXXXXXXXXXXXXXXXXXXX",
	}

	var err error
	Record(t, "automation_tokens/delete", func(c *Client) {
		err = c.DeleteAutomationToken(context.TODO(), input)
	})
	if err != nil {
		t.Fatal(err)
	}
}
