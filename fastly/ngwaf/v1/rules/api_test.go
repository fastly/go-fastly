package rules

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v11/fastly"
	"github.com/fastly/go-fastly/v11/fastly/ngwaf/v1/common"
	"github.com/fastly/go-fastly/v11/fastly/ngwaf/v1/signals"
)

func TestClient_Rule_WorkspaceScope(t *testing.T) {
	runRuleTest(t, common.ScopeTypeWorkspace, fastly.TestNGWAFWorkspaceID)
}

func TestClient_Rule_AccountScope(t *testing.T) {
	runRuleTest(t, common.ScopeTypeAccount, "*") // assuming TestNGWAFAccountID exists
}

func runRuleTest(t *testing.T, scopeType common.ScopeType, appliesToID string) {
	assert := require.New(t)

	var err error

	ruleType := "request"
	description := "test"
	groupOperator := "all"
	enabled := true
	requestLogging := "sampled"

	// Action
	actionType := "block"

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

	// Multival conditions
	multivalConditionType := "multival"
	multivalGroupOperator1 := "any"
	multivalField := "request_cookie"
	multivalOperator := "exists"

	field9 := "name"
	operator9 := "equals"
	value9 := "fooCookie"

	field10 := "value"
	operator10 := "equals"
	value10 := "barCookie"

	// List all rules.
	var rs *Rules
	fastly.Record(t, fmt.Sprintf("%s_list_rules", scopeType), func(c *fastly.Client) {
		rs, err = List(context.TODO(), c, &ListInput{
			Scope: &common.Scope{
				Type:      scopeType,
				AppliesTo: []string{appliesToID},
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(rs)
	assert.NotNil(rs.Data)

	for _, rule := range rs.Data {
		// Ensure we are checking the correct scope
		assert.Equal(string(scopeType), rule.Scope.Type)
		assert.Contains(rule.Scope.AppliesTo, appliesToID)

		// Assert the rule with description "test" does not exist
		assert.NotEqual(description, rule.Description, "unexpected rule with description 'test' found")
	}

	// Create a test rule.
	var rule *Rule
	fastly.Record(t, fmt.Sprintf("%s_create_rule", scopeType), func(c *fastly.Client) {
		rule, err = Create(context.TODO(), c, &CreateInput{
			Scope: &common.Scope{
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
			MultivalConditions: []*CreateMultivalCondition{
				{
					Field:         &multivalField,
					GroupOperator: &multivalGroupOperator1,
					Operator:      &multivalOperator,
					Conditions: []*CreateConditionMult{
						{
							Field:    &field9,
							Operator: &operator9,
							Value:    &value9,
						},
						{
							Field:    &field10,
							Operator: &operator10,
							Value:    &value10,
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

	assert.Len(rule.Conditions, 6) // 3 single + 2 group top-level + 1 multival condition

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

	// Validate multival conditions
	var multivalConditions []MultivalCondition
	for _, cond := range rule.Conditions {
		if cond.Type == multivalConditionType {
			if mc, ok := cond.Fields.(MultivalCondition); ok {
				multivalConditions = append(multivalConditions, mc)
			} else {
				t.Errorf("expected MultivalCondition, got %T", cond.Fields)
			}
		}
	}

	assert.Len(multivalConditions, 1)

	// First multival condition
	assert.Equal(multivalGroupOperator1, multivalConditions[0].GroupOperator)
	assert.Len(multivalConditions[0].Conditions, 2)
	assert.Contains(multivalConditions[0].Conditions, ConditionMul{Type: conditionType, Field: field9, Operator: operator9, Value: value9})
	assert.Contains(multivalConditions[0].Conditions, ConditionMul{Type: conditionType, Field: field10, Operator: operator10, Value: value10})

	// Ensure we delete the test rule at the end.
	defer func() {
		fastly.Record(t, fmt.Sprintf("%s_delete_rule", scopeType), func(c *fastly.Client) {
			err = Delete(context.TODO(), c, &DeleteInput{
				RuleID: fastly.ToPointer(rule.RuleID),
				Scope: &common.Scope{
					Type:      scopeType,
					AppliesTo: []string{appliesToID},
				},
			})
		})
		if err != nil {
			t.Errorf("error during rule cleanup: %v", err)
		}
	}()

	// Get the test rule.
	var testRule *Rule
	fastly.Record(t, fmt.Sprintf("%s_get_rule", scopeType), func(c *fastly.Client) {
		testRule, err = Get(context.TODO(), c, &GetInput{
			RuleID: fastly.ToPointer(rule.RuleID),
			Scope: &common.Scope{
				Type:      scopeType,
				AppliesTo: []string{appliesToID},
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(string(scopeType), testRule.Scope.Type)
	assert.Contains(testRule.Scope.AppliesTo, appliesToID)
	assert.Equal(ruleType, testRule.Type)
	assert.Equal(description, testRule.Description)
	assert.Equal(groupOperator, testRule.GroupOperator)
	assert.Equal(enabled, testRule.Enabled)
	assert.Equal(requestLogging, testRule.RequestLogging)

	assert.Len(testRule.Actions, 1)
	testRuleAction := rule.Actions[0]
	assert.Equal(actionType, testRuleAction.Type)

	assert.Len(testRule.Conditions, 6) // 3 single + 2 group top-level + 1 multival condition

	// Validate single conditions
	var testRuleSingleConditions []SingleCondition
	for _, cond := range testRule.Conditions {
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
	for _, cond := range testRule.Conditions {
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

	// Validate multival conditions
	var testMultivalConditions []MultivalCondition
	for _, cond := range testRule.Conditions {
		if cond.Type == multivalConditionType {
			if mc, ok := cond.Fields.(MultivalCondition); ok {
				testMultivalConditions = append(testMultivalConditions, mc)
			} else {
				t.Errorf("expected MultivalCondition, got %T", cond.Fields)
			}
		}
	}

	assert.Len(testMultivalConditions, 1)

	// First multival condition
	assert.Equal(multivalGroupOperator1, testMultivalConditions[0].GroupOperator)
	assert.Len(testMultivalConditions[0].Conditions, 2)
	assert.Contains(testMultivalConditions[0].Conditions, ConditionMul{Type: conditionType, Field: field9, Operator: operator9, Value: value9})
	assert.Contains(testMultivalConditions[0].Conditions, ConditionMul{Type: conditionType, Field: field10, Operator: operator10, Value: value10})

	// Update rule test variables
	updatedDescription := "updated test"
	updatedGroupOperator := "any"
	updatedEnabled := false
	updatedRequestLogging := "none"

	// Updated Action
	updatedActionType := "allow"

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

	// Updated multival conditions
	updatedMultivalGroupOperator1 := "all"
	updatedMultivalOperator := "does_not_exist"

	updatedField9 := "name"
	updatedOperator9 := "does_not_equal"
	updatedValue9 := "fooCookieUpdated"

	updatedField10 := "value"
	updatedOperator10 := "does_not_equal"
	updatedValue10 := "barCookieUpdated"

	// Update the test rule.
	var updatedRule *Rule
	fastly.Record(t, fmt.Sprintf("%s_update_rule", scopeType), func(c *fastly.Client) {
		updatedRule, err = Update(context.TODO(), c, &UpdateInput{
			Scope: &common.Scope{
				Type:      scopeType,
				AppliesTo: []string{appliesToID},
			},
			RuleID:         fastly.ToPointer(rule.RuleID),
			Description:    &updatedDescription,
			GroupOperator:  &updatedGroupOperator,
			Enabled:        &updatedEnabled,
			RequestLogging: &updatedRequestLogging,
			Actions: []*UpdateAction{
				{
					Type: &updatedActionType,
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
			MultivalConditions: []*UpdateMultivalCondition{
				{
					Field:         &multivalField,
					GroupOperator: &updatedMultivalGroupOperator1,
					Operator:      &updatedMultivalOperator,
					Conditions: []*UpdateConditionMult{
						{
							Field:    &updatedField9,
							Operator: &updatedOperator9,
							Value:    &updatedValue9,
						},
						{
							Field:    &updatedField10,
							Operator: &updatedOperator10,
							Value:    &updatedValue10,
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
	assert.Equal(string(scopeType), updatedRule.Scope.Type)
	assert.Contains(updatedRule.Scope.AppliesTo, appliesToID)
	assert.Equal(ruleType, updatedRule.Type)
	assert.Equal(updatedDescription, updatedRule.Description)
	assert.Equal(updatedGroupOperator, updatedRule.GroupOperator)
	assert.Equal(updatedEnabled, updatedRule.Enabled)
	assert.Equal(updatedRequestLogging, updatedRule.RequestLogging)

	assert.Len(updatedRule.Actions, 1)
	updatedAction := updatedRule.Actions[0]
	assert.Equal(updatedActionType, updatedAction.Type)

	assert.Len(updatedRule.Conditions, 6) // 3 single + 2 group top-level + 1 multival condition

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

	// Validate multival conditions
	var updatedMultivalConditions []MultivalCondition
	for _, cond := range updatedRule.Conditions {
		if cond.Type == multivalConditionType {
			if mc, ok := cond.Fields.(MultivalCondition); ok {
				updatedMultivalConditions = append(updatedMultivalConditions, mc)
			} else {
				t.Errorf("expected MultivalCondition, got %T", cond.Fields)
			}
		}
	}

	assert.Len(updatedMultivalConditions, 1)

	// First multival condition
	assert.Equal(updatedMultivalGroupOperator1, updatedMultivalConditions[0].GroupOperator)
	assert.Len(updatedMultivalConditions[0].Conditions, 2)
	assert.Contains(updatedMultivalConditions[0].Conditions, ConditionMul{Type: conditionType, Field: updatedField9, Operator: updatedOperator9, Value: updatedValue9})
	assert.Contains(updatedMultivalConditions[0].Conditions, ConditionMul{Type: conditionType, Field: updatedField10, Operator: updatedOperator10, Value: updatedValue10})
}

func TestClient_Rate_Limit_Rule_WorkspaceScope(t *testing.T) {
	runRateLimitRuleTest(t, common.ScopeTypeWorkspace, fastly.TestNGWAFWorkspaceID)
}

func runRateLimitRuleTest(t *testing.T, scopeType common.ScopeType, appliesToID string) {
	assert := require.New(t)

	var err error

	ruleType := "rate_limit"
	description := "rate_limit_test"
	groupOperator := "all"
	enabled := true

	// Action
	actionType := "block_signal"

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

	// Multival conditions
	multivalConditionType := "multival"
	multivalGroupOperator1 := "any"
	multivalField := "request_cookie"
	multivalOperator := "exists"

	field9 := "name"
	operator9 := "equals"
	value9 := "fooCookie"

	field10 := "value"
	operator10 := "equals"
	value10 := "barCookie"

	// List all rules.
	var rs *Rules
	fastly.Record(t, fmt.Sprintf("%s_rate_limit_list_rules", scopeType), func(c *fastly.Client) {
		rs, err = List(context.TODO(), c, &ListInput{
			Scope: &common.Scope{
				Type:      scopeType,
				AppliesTo: []string{appliesToID},
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(rs)
	assert.NotNil(rs.Data)

	for _, rule := range rs.Data {
		// Ensure we are checking the correct scope
		assert.Equal(string(scopeType), rule.Scope.Type)
		assert.Contains(rule.Scope.AppliesTo, appliesToID)

		// Assert the rule with description "test" does not exist
		assert.NotEqual(description, rule.Description, "unexpected rule with description 'rate_limit_test' found")
	}

	testSignalName := "A Real Name " + string(scopeType)
	testDescription := "This is a description"

	// Create a test signal.
	var signal *signals.Signal

	fastly.Record(t, fmt.Sprintf("%s_rate_limit_create_signal", scopeType), func(c *fastly.Client) {
		signal, err = signals.Create(context.TODO(), c, &signals.CreateInput{
			Description: fastly.ToPointer(testDescription),
			Name:        fastly.ToPointer(testSignalName),
			Scope: &common.Scope{
				Type:      scopeType,
				AppliesTo: []string{appliesToID},
			},
		})
	})

	// Rate Limit Parameters
	threshold := 1
	interval := 60
	duration := 300
	clientIdentifier := "ip"
	clientIdentifierName := "name"

	// Create a test rule.
	var rule *Rule
	fastly.Record(t, fmt.Sprintf("%s_rate_limit_create_rule", scopeType), func(c *fastly.Client) {
		rule, err = Create(context.TODO(), c, &CreateInput{
			Scope: &common.Scope{
				Type:      scopeType,
				AppliesTo: []string{appliesToID},
			},
			Type:          &ruleType,
			Description:   &description,
			GroupOperator: &groupOperator,
			Enabled:       &enabled,
			RateLimit: &CreateRateLimit{
				Signal:    &signal.ReferenceID,
				Threshold: &threshold,
				Interval:  &interval,
				Duration:  &duration,
				ClientIdentifiers: []*CreateClientIdentifier{
					{
						Type: &clientIdentifier,
						Name: &clientIdentifierName,
					},
				},
			},
			Actions: []*CreateAction{
				{
					Type:   &actionType,
					Signal: &signal.ReferenceID,
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
			MultivalConditions: []*CreateMultivalCondition{
				{
					Field:         &multivalField,
					GroupOperator: &multivalGroupOperator1,
					Operator:      &multivalOperator,
					Conditions: []*CreateConditionMult{
						{
							Field:    &field9,
							Operator: &operator9,
							Value:    &value9,
						},
						{
							Field:    &field10,
							Operator: &operator10,
							Value:    &value10,
						},
					},
				},
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure we delete the test rule and signal at the end.
	defer func() {
		fastly.Record(t, fmt.Sprintf("%s_rate_limit_delete_rule", scopeType), func(c *fastly.Client) {
			err = Delete(context.TODO(), c, &DeleteInput{
				RuleID: fastly.ToPointer(rule.RuleID),
				Scope: &common.Scope{
					Type:      scopeType,
					AppliesTo: []string{appliesToID},
				},
			})
		})
		if err != nil {
			t.Errorf("error during rule cleanup: %v", err)
		}
		fastly.Record(t, fmt.Sprintf("%s_rate_limit_delete_signal", scopeType), func(c *fastly.Client) {
			err = signals.Delete(context.TODO(), c, &signals.DeleteInput{
				Scope: &common.Scope{
					Type:      scopeType,
					AppliesTo: []string{appliesToID},
				},
				SignalID: fastly.ToPointer(signal.SignalID),
			})
		})
		if err != nil {
			t.Errorf("error during signal cleanup: %v", err)
		}
	}()

	assert.Equal(string(scopeType), rule.Scope.Type)
	assert.Contains(rule.Scope.AppliesTo, appliesToID)
	assert.Equal(ruleType, rule.Type)
	assert.Equal(description, rule.Description)
	assert.Equal(groupOperator, rule.GroupOperator)
	assert.Equal(enabled, rule.Enabled)

	assert.Len(rule.Actions, 1)
	action := rule.Actions[0]
	assert.Equal(actionType, action.Type)
	assert.Equal(signal.ReferenceID, action.Signal)

	assert.Equal(duration, rule.RateLimit.Duration)
	assert.Equal(clientIdentifier, rule.RateLimit.ClientIdentifiers[0].Type)
	assert.Equal(clientIdentifierName, rule.RateLimit.ClientIdentifiers[0].Name)
	assert.Equal(interval, rule.RateLimit.Interval)
	assert.Equal(signal.ReferenceID, rule.RateLimit.Signal)
	assert.Equal(threshold, rule.RateLimit.Threshold)

	assert.Len(rule.Conditions, 6) // 3 single + 2 group top-level + 1 multival condition

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

	// Validate multival conditions
	var multivalConditions []MultivalCondition
	for _, cond := range rule.Conditions {
		if cond.Type == multivalConditionType {
			if mc, ok := cond.Fields.(MultivalCondition); ok {
				multivalConditions = append(multivalConditions, mc)
			} else {
				t.Errorf("expected MultivalCondition, got %T", cond.Fields)
			}
		}
	}

	assert.Len(multivalConditions, 1)

	// First multival condition
	assert.Equal(multivalGroupOperator1, multivalConditions[0].GroupOperator)
	assert.Len(multivalConditions[0].Conditions, 2)
	assert.Contains(multivalConditions[0].Conditions, ConditionMul{Type: conditionType, Field: field9, Operator: operator9, Value: value9})
	assert.Contains(multivalConditions[0].Conditions, ConditionMul{Type: conditionType, Field: field10, Operator: operator10, Value: value10})

	// Get the test rule.
	var testRule *Rule
	fastly.Record(t, fmt.Sprintf("%s_rate_limit_get_rule", scopeType), func(c *fastly.Client) {
		testRule, err = Get(context.TODO(), c, &GetInput{
			RuleID: fastly.ToPointer(rule.RuleID),
			Scope: &common.Scope{
				Type:      scopeType,
				AppliesTo: []string{appliesToID},
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(string(scopeType), testRule.Scope.Type)
	assert.Contains(testRule.Scope.AppliesTo, appliesToID)
	assert.Equal(ruleType, testRule.Type)
	assert.Equal(description, testRule.Description)
	assert.Equal(groupOperator, testRule.GroupOperator)
	assert.Equal(enabled, testRule.Enabled)

	assert.Len(testRule.Actions, 1)
	testRuleAction := rule.Actions[0]
	assert.Equal(actionType, testRuleAction.Type)
	assert.Equal(signal.ReferenceID, testRuleAction.Signal)

	assert.Equal(duration, rule.RateLimit.Duration)
	assert.Equal(clientIdentifier, rule.RateLimit.ClientIdentifiers[0].Type)
	assert.Equal(clientIdentifierName, rule.RateLimit.ClientIdentifiers[0].Name)
	assert.Equal(interval, rule.RateLimit.Interval)
	assert.Equal(signal.ReferenceID, rule.RateLimit.Signal)
	assert.Equal(threshold, rule.RateLimit.Threshold)

	assert.Len(testRule.Conditions, 6) // 3 single + 2 group top-level + 1 multival condition

	// Validate single conditions
	var testRuleSingleConditions []SingleCondition
	for _, cond := range testRule.Conditions {
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
	for _, cond := range testRule.Conditions {
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

	// Validate multival conditions
	var testMultivalConditions []MultivalCondition
	for _, cond := range testRule.Conditions {
		if cond.Type == multivalConditionType {
			if mc, ok := cond.Fields.(MultivalCondition); ok {
				testMultivalConditions = append(testMultivalConditions, mc)
			} else {
				t.Errorf("expected MultivalCondition, got %T", cond.Fields)
			}
		}
	}

	assert.Len(testMultivalConditions, 1)

	// First multival condition
	assert.Equal(multivalGroupOperator1, testMultivalConditions[0].GroupOperator)
	assert.Len(testMultivalConditions[0].Conditions, 2)
	assert.Contains(testMultivalConditions[0].Conditions, ConditionMul{Type: conditionType, Field: field9, Operator: operator9, Value: value9})
	assert.Contains(testMultivalConditions[0].Conditions, ConditionMul{Type: conditionType, Field: field10, Operator: operator10, Value: value10})

	// Update rule test variables
	updatedDescription := "updated test"
	updatedGroupOperator := "any"
	updatedEnabled := false

	// Updated Rate Limit Parameters
	updatedThreshold := 2
	updatedInterval := 600
	updatedDuration := 600

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

	// Updated multival conditions
	updatedMultivalGroupOperator1 := "all"
	updatedMultivalOperator := "does_not_exist"

	updatedField9 := "name"
	updatedOperator9 := "does_not_equal"
	updatedValue9 := "fooCookieUpdated"

	updatedField10 := "value"
	updatedOperator10 := "does_not_equal"
	updatedValue10 := "barCookieUpdated"

	// Update the test rule.
	var updatedRule *Rule
	fastly.Record(t, fmt.Sprintf("%s_rate_limit_update_rule", scopeType), func(c *fastly.Client) {
		updatedRule, err = Update(context.TODO(), c, &UpdateInput{
			Scope: &common.Scope{
				Type:      scopeType,
				AppliesTo: []string{appliesToID},
			},
			RuleID:        fastly.ToPointer(rule.RuleID),
			Description:   &updatedDescription,
			GroupOperator: &updatedGroupOperator,
			Enabled:       &updatedEnabled,
			RateLimit: &UpdateRateLimit{
				Signal:    &signal.ReferenceID,
				Threshold: &updatedThreshold,
				Interval:  &updatedInterval,
				Duration:  &updatedDuration,
				ClientIdentifiers: []*UpdateClientIdentifier{
					{
						Name: &clientIdentifierName,
						Type: &clientIdentifier,
					},
				},
			},
			Actions: []*UpdateAction{
				{
					Type: &actionType,
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
			MultivalConditions: []*UpdateMultivalCondition{
				{
					Field:         &multivalField,
					GroupOperator: &updatedMultivalGroupOperator1,
					Operator:      &updatedMultivalOperator,
					Conditions: []*UpdateConditionMult{
						{
							Field:    &updatedField9,
							Operator: &updatedOperator9,
							Value:    &updatedValue9,
						},
						{
							Field:    &updatedField10,
							Operator: &updatedOperator10,
							Value:    &updatedValue10,
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
	assert.Equal(string(scopeType), updatedRule.Scope.Type)
	assert.Contains(updatedRule.Scope.AppliesTo, appliesToID)
	assert.Equal(ruleType, updatedRule.Type)
	assert.Equal(updatedDescription, updatedRule.Description)
	assert.Equal(updatedGroupOperator, updatedRule.GroupOperator)
	assert.Equal(updatedEnabled, updatedRule.Enabled)
	assert.Equal("", updatedRule.RequestLogging)

	assert.Len(updatedRule.Actions, 1)
	updatedAction := updatedRule.Actions[0]
	assert.Equal(actionType, updatedAction.Type)

	// Validate Rate Limit
	assert.Equal(updatedDuration, updatedRule.RateLimit.Duration)
	assert.Equal(clientIdentifier, updatedRule.RateLimit.ClientIdentifiers[0].Type)
	assert.Equal(updatedInterval, updatedRule.RateLimit.Interval)
	assert.Equal(signal.ReferenceID, updatedRule.RateLimit.Signal)
	assert.Equal(updatedThreshold, updatedRule.RateLimit.Threshold)

	assert.Len(updatedRule.Conditions, 6) // 3 single + 2 group top-level + 1 multival condition

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

	// Validate multival conditions
	var updatedMultivalConditions []MultivalCondition
	for _, cond := range updatedRule.Conditions {
		if cond.Type == multivalConditionType {
			if mc, ok := cond.Fields.(MultivalCondition); ok {
				updatedMultivalConditions = append(updatedMultivalConditions, mc)
			} else {
				t.Errorf("expected MultivalCondition, got %T", cond.Fields)
			}
		}
	}

	assert.Len(updatedMultivalConditions, 1)

	// First multival condition
	assert.Equal(updatedMultivalGroupOperator1, updatedMultivalConditions[0].GroupOperator)
	assert.Len(updatedMultivalConditions[0].Conditions, 2)
	assert.Contains(updatedMultivalConditions[0].Conditions, ConditionMul{Type: conditionType, Field: updatedField9, Operator: updatedOperator9, Value: updatedValue9})
	assert.Contains(updatedMultivalConditions[0].Conditions, ConditionMul{Type: conditionType, Field: updatedField10, Operator: updatedOperator10, Value: updatedValue10})
}

func TestClient_Deception_Rule_WorkspaceScope(t *testing.T) {
	runDeceptionRuleTest(t, common.ScopeTypeWorkspace, fastly.TestNGWAFWorkspaceID)
}

func runDeceptionRuleTest(t *testing.T, scopeType common.ScopeType, appliesToID string) {
	assert := require.New(t)

	var err error

	ruleType := "request"
	description := "deception_test"
	groupOperator := "all"
	enabled := true

	// Action
	actionType := "deception"
	deceptionType := "invalid_login_response"

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

	// Multival conditions
	multivalConditionType := "multival"
	multivalGroupOperator1 := "any"
	multivalField := "request_cookie"
	multivalOperator := "exists"

	field9 := "name"
	operator9 := "equals"
	value9 := "fooCookie"

	field10 := "value"
	operator10 := "equals"
	value10 := "barCookie"

	// List all rules.
	var rs *Rules
	fastly.Record(t, fmt.Sprintf("%s_deception_list_rules", scopeType), func(c *fastly.Client) {
		rs, err = List(context.TODO(), c, &ListInput{
			Scope: &common.Scope{
				Type:      scopeType,
				AppliesTo: []string{appliesToID},
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(rs)
	assert.NotNil(rs.Data)

	for _, rule := range rs.Data {
		// Ensure we are checking the correct scope
		assert.Equal(string(scopeType), rule.Scope.Type)
		assert.Contains(rule.Scope.AppliesTo, appliesToID)

		// Assert the rule with description "test" does not exist
		assert.NotEqual(description, rule.Description, "unexpected rule with description 'deception_test' found")
	}

	// Create a test rule.
	var rule *Rule
	fastly.Record(t, fmt.Sprintf("%s_deception_create_rule", scopeType), func(c *fastly.Client) {
		rule, err = Create(context.TODO(), c, &CreateInput{
			Scope: &common.Scope{
				Type:      scopeType,
				AppliesTo: []string{appliesToID},
			},
			Type:          &ruleType,
			Description:   &description,
			GroupOperator: &groupOperator,
			Enabled:       &enabled,
			Actions: []*CreateAction{
				{
					DeceptionType: &deceptionType,
					Type:          &actionType,
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
			MultivalConditions: []*CreateMultivalCondition{
				{
					Field:         &multivalField,
					GroupOperator: &multivalGroupOperator1,
					Operator:      &multivalOperator,
					Conditions: []*CreateConditionMult{
						{
							Field:    &field9,
							Operator: &operator9,
							Value:    &value9,
						},
						{
							Field:    &field10,
							Operator: &operator10,
							Value:    &value10,
						},
					},
				},
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure we delete the test rule and signal at the end.
	defer func() {
		fastly.Record(t, fmt.Sprintf("%s_deception_delete_rule", scopeType), func(c *fastly.Client) {
			err = Delete(context.TODO(), c, &DeleteInput{
				RuleID: fastly.ToPointer(rule.RuleID),
				Scope: &common.Scope{
					Type:      scopeType,
					AppliesTo: []string{appliesToID},
				},
			})
		})
		if err != nil {
			t.Errorf("error during rule cleanup: %v", err)
		}
	}()

	assert.Equal(string(scopeType), rule.Scope.Type)
	assert.Contains(rule.Scope.AppliesTo, appliesToID)
	assert.Equal(ruleType, rule.Type)
	assert.Equal(description, rule.Description)
	assert.Equal(groupOperator, rule.GroupOperator)
	assert.Equal(enabled, rule.Enabled)

	assert.Len(rule.Actions, 1)
	action := rule.Actions[0]
	assert.Equal(actionType, action.Type)
	assert.Equal(deceptionType, action.DeceptionType)

	assert.Len(rule.Conditions, 6) // 3 single + 2 group top-level + 1 multival condition

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

	// Validate multival conditions
	var multivalConditions []MultivalCondition
	for _, cond := range rule.Conditions {
		if cond.Type == multivalConditionType {
			if mc, ok := cond.Fields.(MultivalCondition); ok {
				multivalConditions = append(multivalConditions, mc)
			} else {
				t.Errorf("expected MultivalCondition, got %T", cond.Fields)
			}
		}
	}

	assert.Len(multivalConditions, 1)

	// First multival condition
	assert.Equal(multivalGroupOperator1, multivalConditions[0].GroupOperator)
	assert.Len(multivalConditions[0].Conditions, 2)
	assert.Contains(multivalConditions[0].Conditions, ConditionMul{Type: conditionType, Field: field9, Operator: operator9, Value: value9})
	assert.Contains(multivalConditions[0].Conditions, ConditionMul{Type: conditionType, Field: field10, Operator: operator10, Value: value10})

	// Get the test rule.
	var testRule *Rule
	fastly.Record(t, fmt.Sprintf("%s_deception_get_rule", scopeType), func(c *fastly.Client) {
		testRule, err = Get(context.TODO(), c, &GetInput{
			RuleID: fastly.ToPointer(rule.RuleID),
			Scope: &common.Scope{
				Type:      scopeType,
				AppliesTo: []string{appliesToID},
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(string(scopeType), testRule.Scope.Type)
	assert.Contains(testRule.Scope.AppliesTo, appliesToID)
	assert.Equal(ruleType, testRule.Type)
	assert.Equal(description, testRule.Description)
	assert.Equal(groupOperator, testRule.GroupOperator)
	assert.Equal(enabled, testRule.Enabled)

	assert.Len(testRule.Actions, 1)
	testRuleAction := rule.Actions[0]
	assert.Equal(actionType, testRuleAction.Type)
	assert.Equal(deceptionType, action.DeceptionType)

	assert.Len(testRule.Conditions, 6) // 3 single + 2 group top-level + 1 multival condition

	// Validate single conditions
	var testRuleSingleConditions []SingleCondition
	for _, cond := range testRule.Conditions {
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
	for _, cond := range testRule.Conditions {
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

	// Validate multival conditions
	var testMultivalConditions []MultivalCondition
	for _, cond := range testRule.Conditions {
		if cond.Type == multivalConditionType {
			if mc, ok := cond.Fields.(MultivalCondition); ok {
				testMultivalConditions = append(testMultivalConditions, mc)
			} else {
				t.Errorf("expected MultivalCondition, got %T", cond.Fields)
			}
		}
	}

	assert.Len(testMultivalConditions, 1)

	// First multival condition
	assert.Equal(multivalGroupOperator1, testMultivalConditions[0].GroupOperator)
	assert.Len(testMultivalConditions[0].Conditions, 2)
	assert.Contains(testMultivalConditions[0].Conditions, ConditionMul{Type: conditionType, Field: field9, Operator: operator9, Value: value9})
	assert.Contains(testMultivalConditions[0].Conditions, ConditionMul{Type: conditionType, Field: field10, Operator: operator10, Value: value10})

	// Update rule test variables
	updatedDescription := "updated test"
	updatedGroupOperator := "any"
	updatedEnabled := false

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

	// Updated multival conditions
	updatedMultivalGroupOperator1 := "all"
	updatedMultivalOperator := "does_not_exist"

	updatedField9 := "name"
	updatedOperator9 := "does_not_equal"
	updatedValue9 := "fooCookieUpdated"

	updatedField10 := "value"
	updatedOperator10 := "does_not_equal"
	updatedValue10 := "barCookieUpdated"

	// Update the test rule.
	var updatedRule *Rule
	fastly.Record(t, fmt.Sprintf("%s_deception_update_rule", scopeType), func(c *fastly.Client) {
		updatedRule, err = Update(context.TODO(), c, &UpdateInput{
			Scope: &common.Scope{
				Type:      scopeType,
				AppliesTo: []string{appliesToID},
			},
			RuleID:        fastly.ToPointer(rule.RuleID),
			Description:   &updatedDescription,
			GroupOperator: &updatedGroupOperator,
			Enabled:       &updatedEnabled,
			Actions: []*UpdateAction{
				{
					Type:          &actionType,
					DeceptionType: &deceptionType,
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
			MultivalConditions: []*UpdateMultivalCondition{
				{
					Field:         &multivalField,
					GroupOperator: &updatedMultivalGroupOperator1,
					Operator:      &updatedMultivalOperator,
					Conditions: []*UpdateConditionMult{
						{
							Field:    &updatedField9,
							Operator: &updatedOperator9,
							Value:    &updatedValue9,
						},
						{
							Field:    &updatedField10,
							Operator: &updatedOperator10,
							Value:    &updatedValue10,
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
	assert.Equal(string(scopeType), updatedRule.Scope.Type)
	assert.Contains(updatedRule.Scope.AppliesTo, appliesToID)
	assert.Equal(ruleType, updatedRule.Type)
	assert.Equal(updatedDescription, updatedRule.Description)
	assert.Equal(updatedGroupOperator, updatedRule.GroupOperator)
	assert.Equal(updatedEnabled, updatedRule.Enabled)
	assert.Equal("sampled", updatedRule.RequestLogging)

	assert.Len(updatedRule.Actions, 1)
	updatedAction := updatedRule.Actions[0]
	assert.Equal(actionType, updatedAction.Type)

	assert.Len(updatedRule.Conditions, 6) // 3 single + 2 group top-level + 1 multival condition

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

	// Validate multival conditions
	var updatedMultivalConditions []MultivalCondition
	for _, cond := range updatedRule.Conditions {
		if cond.Type == multivalConditionType {
			if mc, ok := cond.Fields.(MultivalCondition); ok {
				updatedMultivalConditions = append(updatedMultivalConditions, mc)
			} else {
				t.Errorf("expected MultivalCondition, got %T", cond.Fields)
			}
		}
	}

	assert.Len(updatedMultivalConditions, 1)

	// First multival condition
	assert.Equal(updatedMultivalGroupOperator1, updatedMultivalConditions[0].GroupOperator)
	assert.Len(updatedMultivalConditions[0].Conditions, 2)
	assert.Contains(updatedMultivalConditions[0].Conditions, ConditionMul{Type: conditionType, Field: updatedField9, Operator: updatedOperator9, Value: updatedValue9})
	assert.Contains(updatedMultivalConditions[0].Conditions, ConditionMul{Type: conditionType, Field: updatedField10, Operator: updatedOperator10, Value: updatedValue10})
}

func TestClient_CreateRule_validation(t *testing.T) {
	var err error
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		Type: nil,
	})
	if !errors.Is(err, fastly.ErrMissingType) {
		t.Errorf("expected ErrMissingType: got %s", err)
	}
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		Type:        fastly.ToPointer("request"),
		Description: nil,
	})
	if !errors.Is(err, fastly.ErrMissingDescription) {
		t.Errorf("expected ErrMissingDescription: got %s", err)
	}
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		Type:        fastly.ToPointer("request"),
		Description: fastly.ToPointer("test"),
		Scope:       nil,
	})
	if !errors.Is(err, fastly.ErrMissingScope) {
		t.Errorf("expected ErrMissingScope: got %s", err)
	}
	_, err = Create(context.TODO(), fastly.TestClient, &CreateInput{
		Type:        fastly.ToPointer("request"),
		Description: fastly.ToPointer("test"),
		Scope: &common.Scope{
			Type:      common.ScopeTypeWorkspace,
			AppliesTo: []string{"123"},
		},
		Conditions:      []*CreateCondition{},
		GroupConditions: []*CreateGroupCondition{},
	})
	if !errors.Is(err, fastly.ErrMissingConditions) {
		t.Errorf("expected ErrMissingConditions: got %s", err)
	}
}

func TestClient_GetRule_validation(t *testing.T) {
	var err error
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		RuleID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingRuleID) {
		t.Errorf("expected ErrMissingRuleID: got %s", err)
	}
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		RuleID: fastly.ToPointer("123"),
		Scope:  nil,
	})
	if !errors.Is(err, fastly.ErrMissingScope) {
		t.Errorf("expected ErrMissingScope: got %s", err)
	}
}

func TestClient_ListRules_validation(t *testing.T) {
	var err error
	_, err = List(context.TODO(), fastly.TestClient, &ListInput{
		Scope: nil,
	})
	if !errors.Is(err, fastly.ErrMissingScope) {
		t.Errorf("expected ErrMissingScope: got %s", err)
	}
}

func TestClient_UpdateRule_validation(t *testing.T) {
	var err error
	_, err = Update(context.TODO(), fastly.TestClient, &UpdateInput{
		RuleID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingRuleID) {
		t.Errorf("expected ErrMissingRuleID: got %s", err)
	}
	_, err = Update(context.TODO(), fastly.TestClient, &UpdateInput{
		RuleID: fastly.ToPointer("123"),
		Scope:  nil,
	})
	if !errors.Is(err, fastly.ErrMissingScope) {
		t.Errorf("expected ErrMissingScope: got %s", err)
	}
}

func TestClient_DeleteRule_validation(t *testing.T) {
	var err error
	err = Delete(context.TODO(), fastly.TestClient, &DeleteInput{
		RuleID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingRuleID) {
		t.Errorf("expected ErrMissingRuleID: got %s", err)
	}
	err = Delete(context.TODO(), fastly.TestClient, &DeleteInput{
		RuleID: fastly.ToPointer("123"),
		Scope:  nil,
	})
	if !errors.Is(err, fastly.ErrMissingScope) {
		t.Errorf("expected ErrMissingScope: got %s", err)
	}
}
