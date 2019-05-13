package fastly

import "testing"

func TestClient_Conditions(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "conditions/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var condition *Condition
	record(t, "conditions/create", func(c *Client) {
		condition, err = c.CreateCondition(&CreateConditionInput{
			Service:   testServiceID,
			Version:   tv.Number,
			Name:      "test/condition",
			Statement: "req.url~+\"index.html\"",
			Type:      "REQUEST",
			Priority:  1,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// // Ensure deleted
	defer func() {
		record(t, "conditions/cleanup", func(c *Client) {
			c.DeleteCondition(&DeleteConditionInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test/condition",
			})
		})
	}()

	if condition.Name != "test/condition" {
		t.Errorf("bad name: %q", condition.Name)
	}
	if condition.Statement != "req.url~+\"index.html\"" {
		t.Errorf("bad statement: %q", condition.Statement)
	}
	if condition.Type != "REQUEST" {
		t.Errorf("bad type: %s", condition.Type)
	}
	if condition.Priority != 1 {
		t.Errorf("bad priority: %d", condition.Priority)
	}

	// List
	var conditions []*Condition
	record(t, "conditions/list", func(c *Client) {
		conditions, err = c.ListConditions(&ListConditionsInput{
			Service: testServiceID,
			Version: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(conditions) < 1 {
		t.Errorf("bad conditions: %v", conditions)
	}

	// Get
	var newCondition *Condition
	record(t, "conditions/get", func(c *Client) {
		newCondition, err = c.GetCondition(&GetConditionInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test/condition",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if condition.Name != newCondition.Name {
		t.Errorf("bad name: %q (%q)", condition.Name, newCondition.Name)
	}
	if condition.Statement != "req.url~+\"index.html\"" {
		t.Errorf("bad statement: %q", condition.Statement)
	}
	if condition.Type != "REQUEST" {
		t.Errorf("bad type: %s", condition.Type)
	}
	if condition.Priority != 1 {
		t.Errorf("bad priority: %d", condition.Priority)
	}

	// Update
	var updatedCondition *Condition
	record(t, "conditions/update", func(c *Client) {
		updatedCondition, err = c.UpdateCondition(&UpdateConditionInput{
			Service:   testServiceID,
			Version:   tv.Number,
			Name:      "test/condition",
			Statement: "req.url~+\"updated.html\"",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if updatedCondition.Statement != "req.url~+\"updated.html\"" {
		t.Errorf("bad statement: %q", updatedCondition.Statement)
	}

	// Delete
	record(t, "conditions/delete", func(c *Client) {
		err = c.DeleteCondition(&DeleteConditionInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test/condition",
		})
	})
	if err != nil {
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
		Version: 0,
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
		Version: 0,
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
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetCondition(&GetConditionInput{
		Service: "foo",
		Version: 1,
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
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateCondition(&UpdateConditionInput{
		Service: "foo",
		Version: 1,
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
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteCondition(&DeleteConditionInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
