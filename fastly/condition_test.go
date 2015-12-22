package fastly

import "testing"

func TestClient_Conditions(t *testing.T) {
	t.Parallel()

	tv := testVersion(t)

	// Create
	c, err := testClient.CreateCondition(&CreateConditionInput{
		Service:   testServiceID,
		Version:   tv.Number,
		Name:      "test-condition",
		Statement: "req.url~+\"index.html\"",
		Type:      "REQUEST",
		Priority:  1,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		testClient.DeleteCondition(&DeleteConditionInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-condition",
		})
	}()

	if c.Name != "test-condition" {
		t.Errorf("bad name: %q", c.Name)
	}
	if c.Statement != "req.url~+\"index.html\"" {
		t.Errorf("bad statement: %q", c.Statement)
	}
	if c.Type != "REQUEST" {
		t.Errorf("bad type: %d", c.Type)
	}
	if c.Priority != 1 {
		t.Errorf("bad priority: %d", c.Priority)
	}

	// List
	cs, err := testClient.ListConditions(&ListConditionsInput{
		Service: testServiceID,
		Version: tv.Number,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(cs) < 1 {
		t.Errorf("bad conditions: %v", cs)
	}

	// Get
	nc, err := testClient.GetCondition(&GetConditionInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "test-condition",
	})
	if err != nil {
		t.Fatal(err)
	}
	if c.Name != nc.Name {
		t.Errorf("bad name: %q (%q)", c.Name, nc.Name)
	}
	if c.Statement != "req.url~+\"index.html\"" {
		t.Errorf("bad statement: %q", c.Statement)
	}
	if c.Type != "REQUEST" {
		t.Errorf("bad type: %d", c.Type)
	}
	if c.Priority != 1 {
		t.Errorf("bad priority: %d", c.Priority)
	}

	// Update
	uc, err := testClient.UpdateCondition(&UpdateConditionInput{
		Service:   testServiceID,
		Version:   tv.Number,
		Name:      "test-condition",
		Statement: "req.url~+\"updated.html\"",
	})
	if err != nil {
		t.Fatal(err)
	}
	if uc.Statement != "req.url~+\"updated.html\"" {
		t.Errorf("bad statement: %q", uc.Statement)
	}

	// Delete
	if err := testClient.DeleteCondition(&DeleteConditionInput{
		Service: testServiceID,
		Version: tv.Number,
		Name:    "test-condition",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListConditions_validation(t *testing.T) {
	var err error
	_, err = testClient.ListConditions(&ListConditionsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListConditions(&ListConditionsInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateCondition_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateCondition(&CreateConditionInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateCondition(&CreateConditionInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetCondition_validation(t *testing.T) {
	var err error
	_, err = testClient.GetCondition(&GetConditionInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetCondition(&GetConditionInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetCondition(&GetConditionInput{
		Service: "foo",
		Version: "1",
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateCondition_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateCondition(&UpdateConditionInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateCondition(&UpdateConditionInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateCondition(&UpdateConditionInput{
		Service: "foo",
		Version: "1",
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteCondition_validation(t *testing.T) {
	var err error
	err = testClient.DeleteCondition(&DeleteConditionInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteCondition(&DeleteConditionInput{
		Service: "foo",
		Version: "",
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteCondition(&DeleteConditionInput{
		Service: "foo",
		Version: "1",
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
