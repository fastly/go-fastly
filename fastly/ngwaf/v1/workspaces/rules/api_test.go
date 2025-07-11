package rules

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v10/fastly"
)

func TestClient_Rule(t *testing.T) {
	assert := require.New(t)

	var err error

	ruleType := "request"
	description := "test"
	groupOperator := "all"
	enabled := true
	requestLogging := "sampled"

	// Action
	actionType := "block"
	redirectURL := "https://test.com"
	responseCode := 301

	// Single conditions
	conditionType := "single"
	field1 := "ip"
	operator1 := "equals"
	value1 := "127.0.0.1"

	field2 := "path"
	operator2 := "equals"
	value2 := "/login"

	field3 := "agent_name"
	operator3 := "equals"
	value3 := "host-001"

	// Group conditions
	groupConditionType := "group"
	groupOperator1 := "all"
	groupOperator2 := "any"

	field4 := "country"
	operator4 := "equals"
	value4 := "AD"

	field5 := "method"
	operator5 := "equals"
	value5 := "POST"

	field6 := "protocol_version"
	operator6 := "equals"
	value6 := "HTTP/1.0"

	field7 := "method"
	operator7 := "equals"
	value7 := "HEAD"

	field8 := "domain"
	operator8 := "equals"
	value8 := "example.com"

	// List all rules.
	var rs *Rules
	fastly.Record(t, "list_rules", func(c *fastly.Client) {
		rs, err = List(context.TODO(), c, &ListInput{
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(rs)
	assert.NotNil(rs.Data)

	for _, rule := range rs.Data {
		// Ensure we are checking the correct scope
		assert.Equal("workspace", rule.Scope.Type)
		assert.Contains(rule.Scope.AppliesTo, fastly.TestNGWAFWorkspaceID)

		// Assert the rule with description "test" does not exist
		assert.NotEqual(description, rule.Description, "unexpected rule with description 'test' found")
	}

	// Create a test rule.
	var rule *Rule
	fastly.Record(t, "create_rule", func(c *fastly.Client) {
		rule, err = Create(context.TODO(), c, &CreateInput{
			WorkspaceID:    fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
			Type:           &ruleType,
			Description:    &description,
			GroupOperator:  &groupOperator,
			Enabled:        &enabled,
			RequestLogging: &requestLogging,
			Actions: []*CreateAction{
				{
					Type:         &actionType,
					RedirectURL:  &redirectURL,
					ResponseCode: &responseCode,
				},
			},
			Conditions: []*CreateCondition{
				{
					Type:     &conditionType,
					Field:    &field1,
					Operator: &operator1,
					Value:    &value1,
				},
				{
					Type:     &conditionType,
					Field:    &field2,
					Operator: &operator2,
					Value:    &value2,
				},
				{
					Type:     &conditionType,
					Field:    &field3,
					Operator: &operator3,
					Value:    &value3,
				},
			},
			GroupConditions: []*CreateGroupCondition{
				{
					Type:          &groupConditionType,
					GroupOperator: &groupOperator1,
					Conditions: []*CreateCondition{
						{
							Type:     &conditionType,
							Field:    &field4,
							Operator: &operator4,
							Value:    &value4,
						},
						{
							Type:     &conditionType,
							Field:    &field5,
							Operator: &operator5,
							Value:    &value5,
						},
					},
				},
				{
					Type:          &groupConditionType,
					GroupOperator: &groupOperator2,
					Conditions: []*CreateCondition{
						{
							Type:     &conditionType,
							Field:    &field6,
							Operator: &operator6,
							Value:    &value6,
						},
						{
							Type:     &conditionType,
							Field:    &field7,
							Operator: &operator7,
							Value:    &value7,
						},
						{
							Type:     &conditionType,
							Field:    &field8,
							Operator: &operator8,
							Value:    &value8,
						},
					},
				},
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(ruleType, rule.Type)
	assert.Equal(description, rule.Description)
	assert.Equal(groupOperator, rule.GroupOperator)
	assert.Equal(enabled, rule.Enabled)
	assert.Equal(requestLogging, rule.RequestLogging)

	assert.Len(rule.Actions, 1)
	action := rule.Actions[0]
	assert.Equal(actionType, action.Type)
	assert.Equal(redirectURL, action.RedirectURL)
	assert.Equal(responseCode, action.ResponseCode)

	assert.Len(rule.Conditions, 5) // 3 single + 2 group top-level

	var singleConditions []SingleCondition
	for _, cond := range rule.Conditions {
		if cond.Type == conditionType {
			if sc, ok := cond.Fields.(SingleCondition); ok {
				singleConditions = append(singleConditions, sc)
			} else {
				t.Errorf("expected SingleCondition, got %T", cond.Fields)
			}
		}
	}

	assert.Len(singleConditions, 3)
	assert.Contains(singleConditions, SingleCondition{Field: field1, Operator: operator1, Value: value1})
	assert.Contains(singleConditions, SingleCondition{Field: field2, Operator: operator2, Value: value2})
	assert.Contains(singleConditions, SingleCondition{Field: field3, Operator: operator3, Value: value3})

	// Validate group conditions
	var groupConditions []GroupCondition
	for _, cond := range rule.Conditions {
		if cond.Type == groupConditionType {
			if gc, ok := cond.Fields.(GroupCondition); ok {
				groupConditions = append(groupConditions, gc)
			} else {
				t.Errorf("expected SingleCondition, got %T", cond.Fields)
			}
		}
	}

	assert.Len(groupConditions, 2)

	// First group condition
	assert.Equal(groupOperator1, groupConditions[0].GroupOperator)
	assert.Len(groupConditions[0].Conditions, 2)
	assert.Contains(groupConditions[0].Conditions, Condition{Type: conditionType, Field: field4, Operator: operator4, Value: value4})
	assert.Contains(groupConditions[0].Conditions, Condition{Type: conditionType, Field: field5, Operator: operator5, Value: value5})

	// Second group condition
	assert.Equal(groupOperator2, groupConditions[1].GroupOperator)
	assert.Len(groupConditions[1].Conditions, 3)
	assert.Contains(groupConditions[1].Conditions, Condition{Type: conditionType, Field: field6, Operator: operator6, Value: value6})
	assert.Contains(groupConditions[1].Conditions, Condition{Type: conditionType, Field: field7, Operator: operator7, Value: value7})
	assert.Contains(groupConditions[1].Conditions, Condition{Type: conditionType, Field: field8, Operator: operator8, Value: value8})

	// Ensure we delete the test rule at the end.
	defer func() {
		fastly.Record(t, "delete_rule", func(c *fastly.Client) {
			err = Delete(context.TODO(), c, &DeleteInput{
				RuleID:      fastly.ToPointer(rule.RuleID),
				WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
			})
		})
		if err != nil {
			t.Errorf("error during rule cleanup: %v", err)
		}
	}()

	// Get the test rule.
	var testRule *Rule
	fastly.Record(t, "get_rule", func(c *fastly.Client) {
		testRule, err = Get(context.TODO(), c, &GetInput{
			RuleID:      fastly.ToPointer(rule.RuleID),
			WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(ruleType, testRule.Type)
	assert.Equal(description, testRule.Description)
	assert.Equal(groupOperator, testRule.GroupOperator)
	assert.Equal(enabled, testRule.Enabled)
	assert.Equal(requestLogging, testRule.RequestLogging)

	assert.Len(testRule.Actions, 1)
	testRuleAction := rule.Actions[0]
	assert.Equal(actionType, testRuleAction.Type)
	assert.Equal(redirectURL, testRuleAction.RedirectURL)
	assert.Equal(responseCode, testRuleAction.ResponseCode)

	assert.Len(testRule.Conditions, 5) // 3 single + 2 group top-level

	// Validate single conditions
	var testRuleSingleConditions []SingleCondition
	for _, cond := range rule.Conditions {
		if cond.Type == conditionType {
			if sc, ok := cond.Fields.(SingleCondition); ok {
				testRuleSingleConditions = append(testRuleSingleConditions, sc)
			} else {
				t.Errorf("expected SingleCondition, got %T", cond.Fields)
			}
		}
	}

	assert.Len(testRuleSingleConditions, 3)
	assert.Contains(testRuleSingleConditions, SingleCondition{Field: field1, Operator: operator1, Value: value1})
	assert.Contains(testRuleSingleConditions, SingleCondition{Field: field2, Operator: operator2, Value: value2})
	assert.Contains(testRuleSingleConditions, SingleCondition{Field: field3, Operator: operator3, Value: value3})

	// Validate group conditions
	var testRuleGroupConditions []GroupCondition
	for _, cond := range rule.Conditions {
		if cond.Type == groupConditionType {
			if gc, ok := cond.Fields.(GroupCondition); ok {
				testRuleGroupConditions = append(testRuleGroupConditions, gc)
			} else {
				t.Errorf("expected SingleCondition, got %T", cond.Fields)
			}
		}
	}

	assert.Len(testRuleGroupConditions, 2)

	// First group condition
	assert.Equal(groupOperator1, testRuleGroupConditions[0].GroupOperator)
	assert.Len(testRuleGroupConditions[0].Conditions, 2)
	assert.Contains(testRuleGroupConditions[0].Conditions, Condition{Type: conditionType, Field: field4, Operator: operator4, Value: value4})
	assert.Contains(testRuleGroupConditions[0].Conditions, Condition{Type: conditionType, Field: field5, Operator: operator5, Value: value5})

	// Second group condition
	assert.Equal(groupOperator2, testRuleGroupConditions[1].GroupOperator)
	assert.Len(testRuleGroupConditions[1].Conditions, 3)
	assert.Contains(testRuleGroupConditions[1].Conditions, Condition{Type: conditionType, Field: field6, Operator: operator6, Value: value6})
	assert.Contains(testRuleGroupConditions[1].Conditions, Condition{Type: conditionType, Field: field7, Operator: operator7, Value: value7})
	assert.Contains(testRuleGroupConditions[1].Conditions, Condition{Type: conditionType, Field: field8, Operator: operator8, Value: value8})

	// Update rule test variables
	updatedDescription := "updated test"
	updatedGroupOperator := "any"
	updatedEnabled := false
	updatedRequestLogging := "none"

	// Updated Action
	updatedActionType := "allow"
	updatedRedirectURL := ""
	updatedResponseCode := 0

	// Updated Single Conditions
	updatedOperator1 := "does_not_equal"
	updatedValue1 := "10.0.0.1"

	updatedOperator2 := "does_not_equal"
	updatedValue2 := "/admin"

	updatedOperator3 := "matches"
	updatedValue3 := "bot-*"

	// Updated Group Conditions
	updatedGroupOperator1 := "any"
	updatedGroupOperator2 := "all"

	updatedOperator4 := "does_not_equal"
	updatedValue4 := "US"

	updatedOperator5 := "does_not_equal"
	updatedValue5 := "PUT"

	updatedOperator6 := "does_not_equal"
	updatedValue6 := "HTTP/2.0"

	updatedOperator7 := "does_not_equal"
	updatedValue7 := "OPTIONS"

	updatedOperator8 := "does_not_equal"
	updatedValue8 := "internal.example"

	// Update the test rule.
	var updatedRule *Rule
	fastly.Record(t, "update_rule", func(c *fastly.Client) {
		updatedRule, err = Update(context.TODO(), c, &UpdateInput{
			WorkspaceID:    fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
			RuleID:         fastly.ToPointer(rule.RuleID),
			Description:    &updatedDescription,
			GroupOperator:  &updatedGroupOperator,
			Enabled:        &updatedEnabled,
			RequestLogging: &updatedRequestLogging,
			Actions: []*UpdateAction{
				{
					Type:         &updatedActionType,
					RedirectURL:  &updatedRedirectURL,
					ResponseCode: &updatedResponseCode,
				},
			},
			Conditions: []*UpdateCondition{
				{
					Type:     &conditionType,
					Field:    &field1,
					Operator: &updatedOperator1,
					Value:    &updatedValue1,
				},
				{
					Type:     &conditionType,
					Field:    &field2,
					Operator: &updatedOperator2,
					Value:    &updatedValue2,
				},
				{
					Type:     &conditionType,
					Field:    &field3,
					Operator: &updatedOperator3,
					Value:    &updatedValue3,
				},
			},
			GroupConditions: []*UpdateGroupCondition{
				{
					Type:          &groupConditionType,
					GroupOperator: &updatedGroupOperator1,
					Conditions: []*UpdateCondition{
						{
							Type:     &conditionType,
							Field:    &field4,
							Operator: &updatedOperator4,
							Value:    &updatedValue4,
						},
						{
							Type:     &conditionType,
							Field:    &field5,
							Operator: &updatedOperator5,
							Value:    &updatedValue5,
						},
					},
				},
				{
					Type:          &groupConditionType,
					GroupOperator: &updatedGroupOperator2,
					Conditions: []*UpdateCondition{
						{
							Type:     &conditionType,
							Field:    &field6,
							Operator: &updatedOperator6,
							Value:    &updatedValue6,
						},
						{
							Type:     &conditionType,
							Field:    &field7,
							Operator: &updatedOperator7,
							Value:    &updatedValue7,
						},
						{
							Type:     &conditionType,
							Field:    &field8,
							Operator: &updatedOperator8,
							Value:    &updatedValue8,
						},
					},
				},
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Assertions
	assert.Equal(ruleType, updatedRule.Type)
	assert.Equal(updatedDescription, updatedRule.Description)
	assert.Equal(updatedGroupOperator, updatedRule.GroupOperator)
	assert.Equal(updatedEnabled, updatedRule.Enabled)
	assert.Equal(updatedRequestLogging, updatedRule.RequestLogging)

	assert.Len(updatedRule.Actions, 1)
	updatedAction := updatedRule.Actions[0]
	assert.Equal(updatedActionType, updatedAction.Type)
	assert.Equal(updatedRedirectURL, updatedAction.RedirectURL)
	assert.Equal(updatedResponseCode, updatedAction.ResponseCode)

	assert.Len(updatedRule.Conditions, 5) // 3 single + 2 group top-level

	// Validate single conditions
	var updatedSingleConditions []SingleCondition
	for _, cond := range updatedRule.Conditions {
		if cond.Type == conditionType {
			if sc, ok := cond.Fields.(SingleCondition); ok {
				updatedSingleConditions = append(updatedSingleConditions, sc)
			} else {
				t.Errorf("expected SingleCondition, got %T", cond.Fields)
			}
		}
	}

	assert.Len(updatedSingleConditions, 3)
	assert.Contains(updatedSingleConditions, SingleCondition{Field: field1, Operator: updatedOperator1, Value: updatedValue1})
	assert.Contains(updatedSingleConditions, SingleCondition{Field: field2, Operator: updatedOperator2, Value: updatedValue2})
	assert.Contains(updatedSingleConditions, SingleCondition{Field: field3, Operator: updatedOperator3, Value: updatedValue3})

	// Validate group conditions
	var updatedGroupConditions []GroupCondition
	for _, cond := range updatedRule.Conditions {
		if cond.Type == groupConditionType {
			if gc, ok := cond.Fields.(GroupCondition); ok {
				updatedGroupConditions = append(updatedGroupConditions, gc)
			} else {
				t.Errorf("expected SingleCondition, got %T", cond.Fields)
			}
		}
	}

	assert.Len(updatedGroupConditions, 2)

	// First group condition
	assert.Equal(updatedGroupOperator1, updatedGroupConditions[0].GroupOperator)
	assert.Len(updatedGroupConditions[0].Conditions, 2)
	assert.Contains(updatedGroupConditions[0].Conditions, Condition{Type: conditionType, Field: field4, Operator: updatedOperator4, Value: updatedValue4})
	assert.Contains(updatedGroupConditions[0].Conditions, Condition{Type: conditionType, Field: field5, Operator: updatedOperator5, Value: updatedValue5})

	// Second group condition
	assert.Equal(updatedGroupOperator2, updatedGroupConditions[1].GroupOperator)
	assert.Len(updatedGroupConditions[1].Conditions, 3)
	assert.Contains(updatedGroupConditions[1].Conditions, Condition{Type: conditionType, Field: field6, Operator: updatedOperator6, Value: updatedValue6})
	assert.Contains(updatedGroupConditions[1].Conditions, Condition{Type: conditionType, Field: field7, Operator: updatedOperator7, Value: updatedValue7})
	assert.Contains(updatedGroupConditions[1].Conditions, Condition{Type: conditionType, Field: field8, Operator: updatedOperator8, Value: updatedValue8})
}

func TestClient_CreateRule_validation(t *testing.T) {
	var err error
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		Type:        nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingType) {
		t.Errorf("expected ErrMissingType: got %s", err)
	}
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		Description: nil,
		Type:        fastly.ToPointer("request"),
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingDescription) {
		t.Errorf("expected ErrMissingDescription: got %s", err)
	}
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		Conditions:      []*CreateCondition{},
		GroupConditions: []*CreateGroupCondition{},
		Description:     fastly.ToPointer("test"),
		Type:            fastly.ToPointer("request"),
		WorkspaceID:     fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingConditions) {
		t.Errorf("expected ErrMissingConditions: got %s", err)
	}
}

func TestClient_GetRule_validation(t *testing.T) {
	var err error
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		RuleID:      nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingRuleID) {
		t.Errorf("expected ErrMissingRuleID: got %s", err)
	}
}

func TestClient_ListRules_validation(t *testing.T) {
	var err error
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
}

func TestClient_UpdateRule_validation(t *testing.T) {
	var err error
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		RuleID:      nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingRuleID) {
		t.Errorf("expected ErrMissingRuleID: got %s", err)
	}
}

func TestClient_DeleteRule_validation(t *testing.T) {
	var err error
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		RuleID:      nil,
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
	})
	if !errors.Is(err, fastly.ErrMissingRuleID) {
		t.Errorf("expected ErrMissingRuleID: got %s", err)
	}
}
