package fastly

import (
	"math/rand"
	"reflect"
	"strconv"
	"testing"
)

func TestClient_WAF_Rule_Exclusion_list(t *testing.T) {
	t.Parallel()

	fixtureBase := "waf_rule_exclusion/list/"

	testService, serviceVersion, responseName := createServiceForWAF(t, fixtureBase)
	waf := createWAFWithRulesForExclusion(t, fixtureBase, testService, serviceVersion, responseName)

	defer deleteTestService(t, fixtureBase+"service/delete", testService.ID)
	defer deleteWAF(t, fixtureBase+"waf/delete", waf.ID, 1)

	excl1In := buildWAFRuleExclusion1(waf.ID, 1)
	excl2In := buildWAFRuleExclusion2(waf.ID, 1)
	createWAFRuleExclusion(t, fixtureBase+"waf_rule_exclusions/create-1", excl1In)
	createWAFRuleExclusion(t, fixtureBase+"waf_rule_exclusions/create-2", excl2In)

	exclResp := listWAFRuleExclusions(t, fixtureBase+"waf_rule_exclusions/list", waf, "waf_rules")
	if len(exclResp.Items) != 2 {
		t.Errorf("expected 2 waf rule exclusions: got %d", len(exclResp.Items))
	}

	assertWAFRuleExclusionEquals(t, exclResp.Items[0], excl1In.WAFRuleExclusion)
	assertWAFRuleExclusionEquals(t, exclResp.Items[1], excl2In.WAFRuleExclusion)
}

func TestClient_WAF_Rule_Exclusion_list_filters(t *testing.T) {
	t.Parallel()

	fixtureBase := "waf_rule_exclusion/list_filters/"

	testService, serviceVersion, responseName := createServiceForWAF(t, fixtureBase)
	waf := createWAFWithRulesForExclusion(t, fixtureBase, testService, serviceVersion, responseName)

	defer deleteTestService(t, fixtureBase+"service/delete", testService.ID)
	defer deleteWAF(t, fixtureBase+"waf/delete", waf.ID, 1)

	excl1In := buildWAFRuleExclusion1(waf.ID, 1)
	excl2In := buildWAFRuleExclusion2(waf.ID, 1)
	excl3In := buildWAFExclusionWAF(waf.ID, 1)
	createWAFRuleExclusion(t, fixtureBase+"waf_rule_exclusions/create-1", excl1In)
	createWAFRuleExclusion(t, fixtureBase+"waf_rule_exclusions/create-2", excl2In)
	createWAFRuleExclusion(t, fixtureBase+"waf_rule_exclusions/create-3", excl3In)

	exclResp := listWAFExclusionsWithFilters(t, fixtureBase+"waf_rule_exclusions/list", &ListAllWAFRuleExclusionsInput{
		WAFID:               waf.ID,
		WAFVersionNumber:    1,
		FilterExclusionType: strToPtr(WAFRuleExclusionTypeRule),
		FilterModSedID:      strToPtr("21032607"),
		Include:             []string{"waf_rules"},
	})
	if len(exclResp.Items) != 1 {
		t.Errorf("expected 2 waf rule exclusions: got %d", len(exclResp.Items))
	}

	assertWAFRuleExclusionEquals(t, exclResp.Items[0], excl2In.WAFRuleExclusion)
}

func TestClient_WAF_Rule_Exclusion_create(t *testing.T) {
	t.Parallel()

	fixtureBase := "waf_rule_exclusion/create/"

	testService, serviceVersion, responseName := createServiceForWAF(t, fixtureBase)
	waf := createWAFWithRulesForExclusion(t, fixtureBase, testService, serviceVersion, responseName)

	defer deleteTestService(t, fixtureBase+"service/delete", testService.ID)
	defer deleteWAF(t, fixtureBase+"waf/delete", waf.ID, 1)

	exclIn := buildWAFRuleExclusion1(waf.ID, 1)
	createWAFRuleExclusion(t, fixtureBase+"waf_rule_exclusions/create", exclIn)

	exclResp := listWAFRuleExclusions(t, fixtureBase+"waf_rule_exclusions/list", waf, "waf_rules")
	if len(exclResp.Items) != 1 {
		t.Errorf("expected 2 waf rule exclusions: got %d", len(exclResp.Items))
	}
	assertWAFRuleExclusionEquals(t, exclResp.Items[0], exclIn.WAFRuleExclusion)
}

func TestClient_WAF_Rule_Exclusion_update(t *testing.T) {
	t.Parallel()

	fixtureBase := "waf_rule_exclusion/update/"

	testService, serviceVersion, responseName := createServiceForWAF(t, fixtureBase)
	waf := createWAFWithRulesForExclusion(t, fixtureBase, testService, serviceVersion, responseName)

	defer deleteTestService(t, fixtureBase+"service/delete", testService.ID)
	defer deleteWAF(t, fixtureBase+"waf/delete", waf.ID, 1)

	exclIn := buildWAFRuleExclusion1(waf.ID, 1)
	exclOut := createWAFRuleExclusion(t, fixtureBase+"waf_rule_exclusions/create", exclIn)

	exclUpdateIn := buildWAFExclusion1ForUpdate(waf.ID, 1, *exclOut.Number)
	updateWAFExclusion(t, fixtureBase+"waf_rule_exclusions/update", exclUpdateIn)

	exclResp := listWAFRuleExclusions(t, fixtureBase+"waf_rule_exclusions/list", waf, "waf_rules")
	if len(exclResp.Items) != 1 {
		t.Errorf("expected 1 waf version: got %d", len(exclResp.Items))
	}
	assertWAFRuleExclusionEquals(t, exclResp.Items[0], exclIn.WAFRuleExclusion)
}

func TestClient_WAF_Rule_Exclusion_delete(t *testing.T) {
	t.Parallel()

	fixtureBase := "waf_rule_exclusion/delete/"

	testService, serviceVersion, responseName := createServiceForWAF(t, fixtureBase)
	waf := createWAFWithRulesForExclusion(t, fixtureBase, testService, serviceVersion, responseName)

	defer deleteTestService(t, fixtureBase+"service/delete", testService.ID)
	defer deleteWAF(t, fixtureBase+"waf/delete", waf.ID, 1)

	exclIn := buildWAFRuleExclusion1(waf.ID, 1)
	exclOut := createWAFRuleExclusion(t, fixtureBase+"waf_rule_exclusions/create", exclIn)

	exclUpdateIn := buildWAFExclusion1ForDeletion(waf.ID, 1, *exclOut.Number)
	deleteWAFExclusion(t, fixtureBase+"waf_rule_exclusions/update", exclUpdateIn)

	exclResp := listWAFRuleExclusions(t, fixtureBase+"waf_rule_exclusions/list", waf, "waf_rules")
	if len(exclResp.Items) != 0 {
		t.Errorf("expected 0 waf version: got %d", len(exclResp.Items))
	}
}

func TestClient_WAF_Rule_Exclusion_list_validation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input         *ListWAFRuleExclusionsInput
		expectedError string
	}{
		{
			input: &ListWAFRuleExclusionsInput{
				WAFID: "",
			},
			expectedError: ErrMissingWAFID.Error(),
		},
		{
			input: &ListWAFRuleExclusionsInput{
				WAFID:            "1",
				WAFVersionNumber: 0,
			},
			expectedError: ErrMissingWAFVersionNumber.Error(),
		},
	}
	for _, c := range cases {
		if _, err := testClient.ListWAFRuleExclusions(c.input); err.Error() != c.expectedError {
			t.Errorf("bad error: %s", err)
		}
	}
}

func TestClient_WAF_Rule_Exclusion_list_all_validation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input         *ListAllWAFRuleExclusionsInput
		expectedError string
	}{
		{
			input: &ListAllWAFRuleExclusionsInput{
				WAFID: "",
			},
			expectedError: ErrMissingWAFID.Error(),
		},
		{
			input: &ListAllWAFRuleExclusionsInput{
				WAFID:            "1",
				WAFVersionNumber: 0,
			},
			expectedError: ErrMissingWAFVersionNumber.Error(),
		},
	}
	for _, c := range cases {
		if _, err := testClient.ListAllWAFRuleExclusions(c.input); err.Error() != c.expectedError {
			t.Errorf("bad error: %s", err)
		}
	}
}

func TestClient_WAF_Rule_Exclusion_create_validation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input         *CreateWAFRuleExclusionInput
		expectedError string
	}{
		{
			input: &CreateWAFRuleExclusionInput{
				WAFID: "",
			},
			expectedError: ErrMissingWAFID.Error(),
		},
		{
			input: &CreateWAFRuleExclusionInput{
				WAFID:            "1",
				WAFVersionNumber: 0,
			},
			expectedError: ErrMissingWAFVersionNumber.Error(),
		},
		{
			input: &CreateWAFRuleExclusionInput{
				WAFID:            "1",
				WAFVersionNumber: 1,
			},
			expectedError: ErrMissingWAFRuleExclusion.Error(),
		},
	}
	for _, c := range cases {
		if _, err := testClient.CreateWAFRuleExclusion(c.input); err.Error() != c.expectedError {
			t.Errorf("bad error: %s", err)
		}
	}
}

func TestClient_WAF_Rule_Exclusion_update_validation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input         *UpdateWAFRuleExclusionInput
		expectedError string
	}{
		{
			input: &UpdateWAFRuleExclusionInput{
				WAFID: "",
			},
			expectedError: ErrMissingWAFID.Error(),
		},
		{
			input: &UpdateWAFRuleExclusionInput{
				WAFID:            "1",
				WAFVersionNumber: 0,
			},
			expectedError: ErrMissingWAFVersionNumber.Error(),
		},
		{
			input: &UpdateWAFRuleExclusionInput{
				WAFID:            "1",
				WAFVersionNumber: 1,
			},
			expectedError: ErrMissingWAFRuleExclusionNumber.Error(),
		},
		{
			input: &UpdateWAFRuleExclusionInput{
				WAFID:            "1",
				WAFVersionNumber: 1,
				Number:           1,
			},
			expectedError: ErrMissingWAFRuleExclusion.Error(),
		},
	}
	for _, c := range cases {
		if _, err := testClient.UpdateWAFRuleExclusion(c.input); err.Error() != c.expectedError {
			t.Errorf("bad error: %s", err)
		}
	}
}

func TestClient_WAF_Rule_Exclusion_delete_validation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input         *DeleteWAFRuleExclusionInput
		expectedError string
	}{
		{
			input: &DeleteWAFRuleExclusionInput{
				WAFID: "",
			},
			expectedError: ErrMissingWAFID.Error(),
		},
		{
			input: &DeleteWAFRuleExclusionInput{
				WAFID:            "1",
				WAFVersionNumber: 0,
			},
			expectedError: ErrMissingWAFVersionNumber.Error(),
		},
		{
			input: &DeleteWAFRuleExclusionInput{
				WAFID:            "1",
				WAFVersionNumber: 1,
			},
			expectedError: ErrMissingNumber.Error(),
		},
	}
	for _, c := range cases {
		if err := testClient.DeleteWAFRuleExclusion(c.input); err.Error() != c.expectedError {
			t.Errorf("bad error: %s", err)
		}
	}
}

func assertWAFRuleExclusionEquals(t *testing.T, actual *WAFRuleExclusion, expected *WAFRuleExclusion) {
	if *actual.ExclusionType != *expected.ExclusionType {
		t.Errorf("expected ExclusionType to be %s got %s", *expected.ExclusionType, *actual.ExclusionType)
	}
	if *actual.Condition != *expected.Condition {
		t.Errorf("expected Condition to be %s got %s", *expected.Condition, *actual.Condition)
	}
	if *actual.Name != *expected.Name {
		t.Errorf("expected Name to be %s got %s", *expected.Name, *actual.Name)
	}
	if len(actual.Rules) != len(expected.Rules) {
		t.Errorf("expected Rules to be of size %d got %d", len(actual.Rules), len(expected.Rules))
	}

	reflect.DeepEqual(extractModSecId(actual.Rules), extractModSecId(expected.Rules))
}

func extractModSecId(wafRule []*WAFRule) []int {
	var modSecIDs []int
	for _, rule := range wafRule {
		modSecIDs = append(modSecIDs, rule.ModSecID)
	}
	return modSecIDs
}

func createWAFWithRulesForExclusion(t *testing.T, fixtureBase string, testService *Service, version *Version, responseName string) *WAF {
	waf := createWAF(t, fixtureBase+"waf/create", testService.ID, "", responseName, version.Number)

	var err error
	rulesIn := buildWAFRulesForExclusion("log")
	record(t, fixtureBase+"active_rules/create", func(c *Client) {
		_, err = c.BatchModificationWAFActiveRules(&BatchModificationWAFActiveRulesInput{
			WAFID:            waf.ID,
			WAFVersionNumber: 1,
			Rules:            rulesIn,
			OP:               UpsertBatchOperation,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	return waf
}

func buildWAFRulesForExclusion(status string) []*WAFActiveRule {

	return []*WAFActiveRule{
		{
			ModSecID: 2029718,
			Status:   status,
			Revision: 1,
		},
		{
			ModSecID: 2037405,
			Status:   status,
			Revision: 1,
		},
		{
			ModSecID: 21032607,
			Status:   status,
			Revision: 1,
		},
	}
}

func createServiceForWAF(t *testing.T, fixtureBase string) (*Service, *Version, string) {
	service := createTestService(t, fixtureBase+"service/create", "service-"+strconv.Itoa(rand.Int()))
	version := createTestVersion(t, fixtureBase+"service/version", service.ID)
	responseName := "WAf_Response"
	createTestWAFResponseObject(t, fixtureBase+"response_object/create", service.ID, responseName, version.Number)
	return service, version, responseName
}

func createWAFRuleExclusion(t *testing.T, fixture string, excl1In *CreateWAFRuleExclusionInput) *WAFRuleExclusion {
	var excl1Out *WAFRuleExclusion
	var err error
	record(t, fixture, func(c *Client) {
		excl1Out, err = c.CreateWAFRuleExclusion(excl1In)
	})
	if err != nil {
		t.Fatal(err)
	}
	return excl1Out
}

func updateWAFExclusion(t *testing.T, fixture string, updateExcl *UpdateWAFRuleExclusionInput) *WAFRuleExclusion {
	var err error
	var out *WAFRuleExclusion
	record(t, fixture, func(c *Client) {
		out, err = c.UpdateWAFRuleExclusion(updateExcl)
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func deleteWAFExclusion(t *testing.T, fixture string, updateExcl *DeleteWAFRuleExclusionInput) {
	var err error
	record(t, fixture, func(c *Client) {
		err = c.DeleteWAFRuleExclusion(updateExcl)
	})
	if err != nil {
		t.Fatal(err)
	}
}

func listWAFRuleExclusions(t *testing.T, fixture string, waf *WAF, include string) *WAFRuleExclusionResponse {
	var err error
	var exclResp *WAFRuleExclusionResponse
	record(t, fixture, func(c *Client) {
		exclResp, err = c.ListAllWAFRuleExclusions(&ListAllWAFRuleExclusionsInput{
			WAFID:            waf.ID,
			WAFVersionNumber: 1,
			Include:          []string{include},
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	return exclResp
}

func listWAFExclusionsWithFilters(t *testing.T, fixture string, request *ListAllWAFRuleExclusionsInput) *WAFRuleExclusionResponse {
	var err error
	var exclResp *WAFRuleExclusionResponse
	record(t, fixture, func(c *Client) {
		exclResp, err = c.ListAllWAFRuleExclusions(request)
	})
	if err != nil {
		t.Fatal(err)
	}
	return exclResp
}

func buildWAFRuleExclusion1(wafID string, version int) *CreateWAFRuleExclusionInput {
	return &CreateWAFRuleExclusionInput{
		WAFID:            wafID,
		WAFVersionNumber: version,
		WAFRuleExclusion: &WAFRuleExclusion{
			Name:          strToPtr("index page"),
			ExclusionType: strToPtr(WAFRuleExclusionTypeRule),
			Condition:     strToPtr("req.url.basename == \"index.html\""),
			Rules: []*WAFRule{
				{
					ID:       "2029718",
					ModSecID: 2029718,
				},
				{
					ID:       "2037405",
					ModSecID: 2037405,
				},
			},
		},
	}
}

func buildWAFRuleExclusion2(wafID string, version int) *CreateWAFRuleExclusionInput {
	return &CreateWAFRuleExclusionInput{
		WAFID:            wafID,
		WAFVersionNumber: version,
		WAFRuleExclusion: &WAFRuleExclusion{
			Name:          strToPtr("index php"),
			ExclusionType: strToPtr(WAFRuleExclusionTypeRule),
			Condition:     strToPtr("req.url.basename == \"index.php\""),
			Rules: []*WAFRule{
				{
					ID:       "21032607",
					ModSecID: 21032607,
				},
			},
		},
	}
}

func buildWAFExclusionWAF(wafID string, version int) *CreateWAFRuleExclusionInput {
	return &CreateWAFRuleExclusionInput{
		WAFID:            wafID,
		WAFVersionNumber: version,
		WAFRuleExclusion: &WAFRuleExclusion{
			Name:          strToPtr("index asp"),
			ExclusionType: strToPtr(WAFRuleExclusionTypeWAF),
			Condition:     strToPtr("req.url.basename == \"index.asp\""),
		},
	}
}

func buildWAFExclusion1ForUpdate(wafID string, version, number int) *UpdateWAFRuleExclusionInput {
	return &UpdateWAFRuleExclusionInput{
		WAFID:            wafID,
		WAFVersionNumber: version,
		Number:           number,
		WAFRuleExclusion: &WAFRuleExclusion{
			Name:          strToPtr("index page"),
			ExclusionType: strToPtr(WAFRuleExclusionTypeRule),
			Condition:     strToPtr("req.url.basename == \"index.html\""),
			Rules: []*WAFRule{
				{
					ID:       "2029718",
					ModSecID: 2029718,
				},
				{
					ID:       "21032607",
					ModSecID: 21032607,
				},
			},
		},
	}
}

func buildWAFExclusion1ForDeletion(wafID string, version, number int) *DeleteWAFRuleExclusionInput {
	return &DeleteWAFRuleExclusionInput{
		WAFID:            wafID,
		WAFVersionNumber: version,
		Number:           number,
	}
}
