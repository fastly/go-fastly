package fastly

import (
	"errors"
	"testing"
)

func TestClient_Conditions(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "conditions/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var condition *Condition
	Record(t, "conditions/create", func(c *Client) {
		condition, err = c.CreateCondition(&CreateConditionInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test/condition"),
			Statement:      ToPointer("req.url~+\"index.html\""),
			Type:           ToPointer("REQUEST"),
			Priority:       ToPointer(1),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// // Ensure deleted
	defer func() {
		Record(t, "conditions/cleanup", func(c *Client) {
			_ = c.DeleteCondition(&DeleteConditionInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test/condition",
			})
		})
	}()

	if *condition.Name != "test/condition" {
		t.Errorf("bad name: %q", *condition.Name)
	}
	if *condition.Statement != "req.url~+\"index.html\"" {
		t.Errorf("bad statement: %q", *condition.Statement)
	}
	if *condition.Type != "REQUEST" {
		t.Errorf("bad type: %s", *condition.Type)
	}
	if *condition.Priority != 1 {
		t.Errorf("bad priority: %d", *condition.Priority)
	}

	// List
	var conditions []*Condition
	Record(t, "conditions/list", func(c *Client) {
		conditions, err = c.ListConditions(&ListConditionsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
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
	Record(t, "conditions/get", func(c *Client) {
		newCondition, err = c.GetCondition(&GetConditionInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test/condition",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *condition.Name != *newCondition.Name {
		t.Errorf("bad name: %q (%q)", *condition.Name, *newCondition.Name)
	}
	if *condition.Statement != "req.url~+\"index.html\"" {
		t.Errorf("bad statement: %q", *condition.Statement)
	}
	if *condition.Type != "REQUEST" {
		t.Errorf("bad type: %s", *condition.Type)
	}
	if *condition.Priority != 1 {
		t.Errorf("bad priority: %d", *condition.Priority)
	}

	// Update
	var updatedCondition *Condition
	Record(t, "conditions/update", func(c *Client) {
		updatedCondition, err = c.UpdateCondition(&UpdateConditionInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test/condition",
			Statement:      ToPointer("req.url~+\"updated.html\""),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *updatedCondition.Statement != "req.url~+\"updated.html\"" {
		t.Errorf("bad statement: %q", *updatedCondition.Statement)
	}

	// Delete
	Record(t, "conditions/delete", func(c *Client) {
		err = c.DeleteCondition(&DeleteConditionInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test/condition",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListConditions_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListConditions(&ListConditionsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListConditions(&ListConditionsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateCondition_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateCondition(&CreateConditionInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateCondition(&CreateConditionInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetCondition_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetCondition(&GetConditionInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetCondition(&GetConditionInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetCondition(&GetConditionInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateCondition_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateCondition(&UpdateConditionInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateCondition(&UpdateConditionInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateCondition(&UpdateConditionInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteCondition_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteCondition(&DeleteConditionInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteCondition(&DeleteConditionInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteCondition(&DeleteConditionInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
