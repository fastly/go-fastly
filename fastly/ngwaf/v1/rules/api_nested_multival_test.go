package rules

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v12/fastly"
	"github.com/fastly/go-fastly/v12/fastly/ngwaf/v1/scope"
)

func TestClient_Rule_NestedMultivalConditions_WorkspaceScope(t *testing.T) {
	runNestedMultivalTest(t, scope.ScopeTypeWorkspace, fastly.TestNGWAFWorkspaceID)
}

func runNestedMultivalTest(t *testing.T, scopeType scope.Type, appliesToID string) {
	assert := require.New(t)

	var err error

	ruleType := "request"
	description := "nested_multival_test"
	groupOperator := "any"
	enabled := true
	requestLogging := "sampled"

	// Action
	actionType := "block"

	// Group condition with nested multival
	groupOperator1 := "all"

	// Single condition within group
	field1 := "ip"
	operator1 := "in_list"
	value1 := "site.blacklist"

	// Multival condition within group
	multivalField := "request_header"
	multivalOperator := "exists"
	multivalGroupOperator := "all"

	field2 := "name"
	operator2 := "equals"
	value2 := "x-login"

	field3 := "value_string"
	operator3 := "equals"
	value3 := "x-updated"

	// Create a test rule with nested multival conditions in a group.
	var rule *Rule
	fastly.Record(t, fmt.Sprintf("%s_create_nested_multival_rule", scopeType), func(c *fastly.Client) {
		rule, err = Create(context.TODO(), c, &CreateInput{
			Scope: &scope.Scope{
				Type:      scopeType,
				AppliesTo: []string{appliesToID},
			},
			Type:           &ruleType,
			Description:    &description,
			GroupOperator:  &groupOperator,
			Enabled:        &enabled,
			RequestLogging: &requestLogging,
			Actions: []*CreateAction{
				{
					Type: &actionType,
				},
			},
			GroupConditions: []*CreateGroupCondition{
				{
					GroupOperator: &groupOperator1,
					Conditions: []*CreateCondition{
						{
							Field:    &field1,
							Operator: &operator1,
							Value:    &value1,
						},
					},
					MultivalConditions: []*CreateMultivalCondition{
						{
							Field:         &multivalField,
							Operator:      &multivalOperator,
							GroupOperator: &multivalGroupOperator,
							Conditions: []*CreateConditionMult{
								{
									Field:    &field2,
									Operator: &operator2,
									Value:    &value2,
								},
								{
									Field:    &field3,
									Operator: &operator3,
									Value:    &value3,
								},
							},
						},
					},
				},
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(string(scopeType), rule.Scope.Type)
	assert.Contains(rule.Scope.AppliesTo, appliesToID)
	assert.Equal(ruleType, rule.Type)
	assert.Equal(description, rule.Description)
	assert.Equal(groupOperator, rule.GroupOperator)
	assert.Equal(enabled, rule.Enabled)
	assert.Equal(requestLogging, rule.RequestLogging)

	assert.Len(rule.Actions, 1)
	action := rule.Actions[0]
	assert.Equal(actionType, action.Type)

	assert.Len(rule.Conditions, 1) // 1 group condition

	// Validate group condition
	var groupConditions []GroupCondition
	for _, cond := range rule.Conditions {
		if cond.Type == "group" {
			if gc, ok := cond.Fields.(GroupCondition); ok {
				groupConditions = append(groupConditions, gc)
			} else {
				t.Errorf("expected GroupCondition, got %T", cond.Fields)
			}
		}
	}

	assert.Len(groupConditions, 1)
	assert.Equal(groupOperator1, groupConditions[0].GroupOperator)
	assert.Len(groupConditions[0].Conditions, 2) // 1 single + 1 multival

	// Validate nested single condition
	var singleConditions []Condition
	for _, gci := range groupConditions[0].Conditions {
		if gci.Type == "single" {
			if c, ok := gci.Fields.(Condition); ok {
				singleConditions = append(singleConditions, c)
			}
		}
	}
	assert.Len(singleConditions, 1)
	assert.Equal(field1, singleConditions[0].Field)
	assert.Equal(operator1, singleConditions[0].Operator)
	assert.Equal(value1, singleConditions[0].Value)

	// Validate nested multival condition
	var multivalConditions []MultivalCondition
	for _, gci := range groupConditions[0].Conditions {
		if gci.Type == "multival" {
			if mc, ok := gci.Fields.(MultivalCondition); ok {
				multivalConditions = append(multivalConditions, mc)
			}
		}
	}
	assert.Len(multivalConditions, 1)
	assert.Equal(multivalField, multivalConditions[0].Field)
	assert.Equal(multivalOperator, multivalConditions[0].Operator)
	assert.Equal(multivalGroupOperator, multivalConditions[0].GroupOperator)
	assert.Len(multivalConditions[0].Conditions, 2)

	assert.Equal(field2, multivalConditions[0].Conditions[0].Field)
	assert.Equal(operator2, multivalConditions[0].Conditions[0].Operator)
	assert.Equal(value2, multivalConditions[0].Conditions[0].Value)

	assert.Equal(field3, multivalConditions[0].Conditions[1].Field)
	assert.Equal(operator3, multivalConditions[0].Conditions[1].Operator)
	assert.Equal(value3, multivalConditions[0].Conditions[1].Value)

	// Ensure we delete the test rule at the end.
	defer func() {
		fastly.Record(t, fmt.Sprintf("%s_delete_nested_multival_rule", scopeType), func(c *fastly.Client) {
			err = Delete(context.TODO(), c, &DeleteInput{
				RuleID: fastly.ToPointer(rule.RuleID),
				Scope: &scope.Scope{
					Type:      scopeType,
					AppliesTo: []string{appliesToID},
				},
			})
		})
		if err != nil {
			t.Errorf("error during rule cleanup: %v", err)
		}
	}()

	// Update rule test
	updatedDescription := "nested_multival_updated"
	updatedField1 := "ip"
	updatedOperator1 := "in_list"
	updatedValue1 := "site.blacklist-updated"

	updatedMultivalField := "request_header"
	updatedMultivalOperator := "does_not_exist"
	updatedMultivalGroupOperator := "any"

	updatedField2 := "name"
	updatedOperator2 := "equals"
	updatedValue2 := "x-updated"

	updatedField3 := "value_string"
	updatedOperator3 := "equals"
	updatedValue3 := "def-456"

	// Update the test rule.
	var updatedRule *Rule
	fastly.Record(t, fmt.Sprintf("%s_update_nested_multival_rule", scopeType), func(c *fastly.Client) {
		updatedRule, err = Update(context.TODO(), c, &UpdateInput{
			Scope: &scope.Scope{
				Type:      scopeType,
				AppliesTo: []string{appliesToID},
			},
			RuleID:         fastly.ToPointer(rule.RuleID),
			Description:    &updatedDescription,
			GroupOperator:  &groupOperator,
			Enabled:        &enabled,
			RequestLogging: &requestLogging,
			Actions: []*UpdateAction{
				{
					Type: &actionType,
				},
			},
			GroupConditions: []*UpdateGroupCondition{
				{
					GroupOperator: &groupOperator1,
					Conditions: []*UpdateCondition{
						{
							Field:    &updatedField1,
							Operator: &updatedOperator1,
							Value:    &updatedValue1,
						},
					},
					MultivalConditions: []*UpdateMultivalCondition{
						{
							Field:         &updatedMultivalField,
							Operator:      &updatedMultivalOperator,
							GroupOperator: &updatedMultivalGroupOperator,
							Conditions: []*UpdateConditionMult{
								{
									Field:    &updatedField2,
									Operator: &updatedOperator2,
									Value:    &updatedValue2,
								},
								{
									Field:    &updatedField3,
									Operator: &updatedOperator3,
									Value:    &updatedValue3,
								},
							},
						},
					},
				},
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Assertions for updated rule
	assert.Equal(string(scopeType), updatedRule.Scope.Type)
	assert.Contains(updatedRule.Scope.AppliesTo, appliesToID)
	assert.Equal(ruleType, updatedRule.Type)
	assert.Equal(updatedDescription, updatedRule.Description)

	assert.Len(updatedRule.Conditions, 1) // 1 group condition

	// Validate updated group condition
	var updatedGroupConditions []GroupCondition
	for _, cond := range updatedRule.Conditions {
		if cond.Type == "group" {
			if gc, ok := cond.Fields.(GroupCondition); ok {
				updatedGroupConditions = append(updatedGroupConditions, gc)
			}
		}
	}

	assert.Len(updatedGroupConditions, 1)
	assert.Len(updatedGroupConditions[0].Conditions, 2) // 1 single + 1 multival

	// Validate updated nested single condition
	var updatedSingleConditions []Condition
	for _, gci := range updatedGroupConditions[0].Conditions {
		if gci.Type == "single" {
			if c, ok := gci.Fields.(Condition); ok {
				updatedSingleConditions = append(updatedSingleConditions, c)
			}
		}
	}
	assert.Len(updatedSingleConditions, 1)
	assert.Equal(updatedField1, updatedSingleConditions[0].Field)
	assert.Equal(updatedOperator1, updatedSingleConditions[0].Operator)
	assert.Equal(updatedValue1, updatedSingleConditions[0].Value)

	// Validate updated nested multival condition
	var updatedMultivalConditions []MultivalCondition
	for _, gci := range updatedGroupConditions[0].Conditions {
		if gci.Type == "multival" {
			if mc, ok := gci.Fields.(MultivalCondition); ok {
				updatedMultivalConditions = append(updatedMultivalConditions, mc)
			}
		}
	}
	assert.Len(updatedMultivalConditions, 1)
	assert.Equal(updatedMultivalField, updatedMultivalConditions[0].Field)
	assert.Equal(updatedMultivalOperator, updatedMultivalConditions[0].Operator)
	assert.Equal(updatedMultivalGroupOperator, updatedMultivalConditions[0].GroupOperator)
	assert.Len(updatedMultivalConditions[0].Conditions, 2)

	assert.Equal(updatedField2, updatedMultivalConditions[0].Conditions[0].Field)
	assert.Equal(updatedOperator2, updatedMultivalConditions[0].Conditions[0].Operator)
	assert.Equal(updatedValue2, updatedMultivalConditions[0].Conditions[0].Value)

	assert.Equal(updatedField3, updatedMultivalConditions[0].Conditions[1].Field)
	assert.Equal(updatedOperator3, updatedMultivalConditions[0].Conditions[1].Operator)
	assert.Equal(updatedValue3, updatedMultivalConditions[0].Conditions[1].Value)
}
