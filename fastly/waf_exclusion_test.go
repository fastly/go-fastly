package fastly

import (
	"math/rand"
	"reflect"
	"strconv"
	"testing"
)

func TestClient_WAF_Exclusion_list(t *testing.T) {
	t.Parallel()

	fixtureBase := "waf_exclusion/list/"

	testService, serviceVersion, responseName := createServiceForWAF(t, fixtureBase)
	waf := createWAFWithRules(t, fixtureBase, testService, serviceVersion, responseName)

	defer deleteTestService(t, fixtureBase+"service/delete", testService.ID)
	defer deleteWAF(t, fixtureBase+"waf/delete", waf.ID, "1")

	excl1In := buildWAFExclusion1(waf.ID, 1)
	excl2In := buildWAFExclusion2(waf.ID, 1)
	createWAFExclusion(t, fixtureBase+"waf_exclusions/create-1", excl1In)
	createWAFExclusion(t, fixtureBase+"waf_exclusions/create-2", excl2In)

	exclResp := listWAFExclusions(t, fixtureBase+"waf_exclusions/list", waf, "waf_rules")
	if len(exclResp.Items) != 2 {
		t.Errorf("expected 2 waf exclusions: got %d", len(exclResp.Items))
	}

	assertWAFExclusionEquals(t, exclResp.Items[0], excl1In.WAFExclusion)
	assertWAFExclusionEquals(t, exclResp.Items[1], excl2In.WAFExclusion)
}

func TestClient_WAF_Exclusion_list_filters(t *testing.T) {
	t.Parallel()

	fixtureBase := "waf_exclusion/list_filters/"

	testService, serviceVersion, responseName := createServiceForWAF(t, fixtureBase)
	waf := createWAFWithRules(t, fixtureBase, testService, serviceVersion, responseName)

	defer deleteTestService(t, fixtureBase+"service/delete", testService.ID)
	defer deleteWAF(t, fixtureBase+"waf/delete", waf.ID, "1")

	excl1In := buildWAFExclusion1(waf.ID, 1)
	excl2In := buildWAFExclusion2(waf.ID, 1)
	excl3In := buildWAFExclusionWAF(waf.ID, 1)
	createWAFExclusion(t, fixtureBase+"waf_exclusions/create-1", excl1In)
	createWAFExclusion(t, fixtureBase+"waf_exclusions/create-2", excl2In)
	createWAFExclusion(t, fixtureBase+"waf_exclusions/create-3", excl3In)

	exclResp := listWAFExclusionsWithFilters(t, fixtureBase+"waf_exclusions/list", &ListAllWAFExclusionsInput{
		WAFID:               waf.ID,
		WAFVersionNumber:    1,
		FilterExclusionType: strToPtr(WAFExclusionTypeRule),
		FilterModSedID:      strToPtr("1010070"),
		Include:             strToPtr("waf_rules"),
	})
	if len(exclResp.Items) != 1 {
		t.Errorf("expected 2 waf exclusions: got %d", len(exclResp.Items))
	}

	assertWAFExclusionEquals(t, exclResp.Items[0], excl2In.WAFExclusion)
}

func TestClient_WAF_Exclusion_create(t *testing.T) {
	t.Parallel()

	fixtureBase := "waf_exclusion/create/"

	testService, serviceVersion, responseName := createServiceForWAF(t, fixtureBase)
	waf := createWAFWithRules(t, fixtureBase, testService, serviceVersion, responseName)

	defer deleteTestService(t, fixtureBase+"service/delete", testService.ID)
	defer deleteWAF(t, fixtureBase+"waf/delete", waf.ID, "1")

	exclIn := buildWAFExclusion1(waf.ID, 1)
	createWAFExclusion(t, fixtureBase+"waf_exclusions/create", exclIn)

	exclResp := listWAFExclusions(t, fixtureBase+"waf_exclusions/list", waf, "waf_rules")
	if len(exclResp.Items) != 1 {
		t.Errorf("expected 2 waf exclusions: got %d", len(exclResp.Items))
	}
	assertWAFExclusionEquals(t, exclResp.Items[0], exclIn.WAFExclusion)
}

func TestClient_WAF_Exclusion_update(t *testing.T) {
	t.Parallel()

	fixtureBase := "waf_exclusion/update/"

	testService, serviceVersion, responseName := createServiceForWAF(t, fixtureBase)
	waf := createWAFWithRules(t, fixtureBase, testService, serviceVersion, responseName)

	defer deleteTestService(t, fixtureBase+"service/delete", testService.ID)
	defer deleteWAF(t, fixtureBase+"waf/delete", waf.ID, "1")

	exclIn := buildWAFExclusion1(waf.ID, 1)
	exclOut := createWAFExclusion(t, fixtureBase+"waf_exclusions/create", exclIn)

	exclUpdateIn := buildWAFExclusion1ForUpdate(waf.ID, 1, *exclOut.Number)
	updateWAFExclusion(t, fixtureBase+"waf_exclusions/update", exclUpdateIn)

	exclResp := listWAFExclusions(t, fixtureBase+"waf_exclusions/list", waf, "waf_rules")
	if len(exclResp.Items) != 1 {
		t.Errorf("expected 1 waf version: got %d", len(exclResp.Items))
	}
	assertWAFExclusionEquals(t, exclResp.Items[0], exclIn.WAFExclusion)
}

func TestClient_WAF_Exclusion_delete(t *testing.T) {
	t.Parallel()

	fixtureBase := "waf_exclusion/delete/"

	testService, serviceVersion, responseName := createServiceForWAF(t, fixtureBase)
	waf := createWAFWithRules(t, fixtureBase, testService, serviceVersion, responseName)

	defer deleteTestService(t, fixtureBase+"service/delete", testService.ID)
	defer deleteWAF(t, fixtureBase+"waf/delete", waf.ID, "1")

	exclIn := buildWAFExclusion1(waf.ID, 1)
	exclOut := createWAFExclusion(t, fixtureBase+"waf_exclusions/create", exclIn)

	exclUpdateIn := buildWAFExclusion1ForDeletion(waf.ID, 1, *exclOut.Number)
	deleteWAFExclusion(t, fixtureBase+"waf_exclusions/update", exclUpdateIn)

	exclResp := listWAFExclusions(t, fixtureBase+"waf_exclusions/list", waf, "waf_rules")
	if len(exclResp.Items) != 0 {
		t.Errorf("expected 0 waf version: got %d", len(exclResp.Items))
	}
}

func TestClient_ListWAFExclusion_validation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input         *ListWAFExclusionsInput
		expectedError error
	}{
		{
			input: &ListWAFExclusionsInput{
				WAFID: "",
			},
			expectedError: ErrMissingWAFID,
		},
		{
			input: &ListWAFExclusionsInput{
				WAFID:            "1",
				WAFVersionNumber: 0,
			},
			expectedError: ErrMissingWAFVersionNumber,
		},
	}
	for _, c := range cases {
		if _, err := testClient.ListWAFExclusions(c.input); err != c.expectedError {
			t.Errorf("bad error: %s", err)
		}
	}
}

func TestClient_ListAllWAFExclusion_validation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input         *ListAllWAFExclusionsInput
		expectedError error
	}{
		{
			input: &ListAllWAFExclusionsInput{
				WAFID: "",
			},
			expectedError: ErrMissingWAFID,
		},
		{
			input: &ListAllWAFExclusionsInput{
				WAFID:            "1",
				WAFVersionNumber: 0,
			},
			expectedError: ErrMissingWAFVersionNumber,
		},
	}
	for _, c := range cases {
		if _, err := testClient.ListAllWAFExclusions(c.input); err != c.expectedError {
			t.Errorf("bad error: %s", err)
		}
	}
}

func TestClient_CreateWAFExclusion_validation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input         *CreateWAFExclusionInput
		expectedError error
	}{
		{
			input: &CreateWAFExclusionInput{
				WAFID: "",
			},
			expectedError: ErrMissingWAFID,
		},
		{
			input: &CreateWAFExclusionInput{
				WAFID:            "1",
				WAFVersionNumber: 0,
			},
			expectedError: ErrMissingWAFVersionNumber,
		},
		{
			input: &CreateWAFExclusionInput{
				WAFID:            "1",
				WAFVersionNumber: 1,
			},
			expectedError: ErrMissingWAFExclusion,
		},
	}
	for _, c := range cases {
		if _, err := testClient.CreateWAFExclusion(c.input); err != c.expectedError {
			t.Errorf("bad error: %s", err)
		}
	}
}

func TestClient_UpdateWAFExclusion_validation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input         *UpdateWAFExclusionInput
		expectedError error
	}{
		{
			input: &UpdateWAFExclusionInput{
				WAFID: "",
			},
			expectedError: ErrMissingWAFID,
		},
		{
			input: &UpdateWAFExclusionInput{
				WAFID:            "1",
				WAFVersionNumber: 0,
			},
			expectedError: ErrMissingWAFVersionNumber,
		},
		{
			input: &UpdateWAFExclusionInput{
				WAFID:            "1",
				WAFVersionNumber: 1,
			},
			expectedError: ErrMissingWAFExclusionNumber,
		},
		{
			input: &UpdateWAFExclusionInput{
				WAFID:            "1",
				WAFVersionNumber: 1,
				Number:           1,
			},
			expectedError: ErrMissingWAFExclusion,
		},
	}
	for _, c := range cases {
		if _, err := testClient.UpdateWAFExclusion(c.input); err != c.expectedError {
			t.Errorf("bad error: %s", err)
		}
	}
}

func TestClient_DeleteWAFExclusion_validation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input         *DeleteWAFExclusionInput
		expectedError error
	}{
		{
			input: &DeleteWAFExclusionInput{
				WAFID: "",
			},
			expectedError: ErrMissingWAFID,
		},
		{
			input: &DeleteWAFExclusionInput{
				WAFID:            "1",
				WAFVersionNumber: 0,
			},
			expectedError: ErrMissingWAFVersionNumber,
		},
		{
			input: &DeleteWAFExclusionInput{
				WAFID:            "1",
				WAFVersionNumber: 1,
			},
			expectedError: ErrMissingWAFExclusionNumber,
		},
	}
	for _, c := range cases {
		if err := testClient.DeleteWAFExclusion(c.input); err != c.expectedError {
			t.Errorf("bad error: %s", err)
		}
	}
}

func assertWAFExclusionEquals(t *testing.T, actual *WAFExclusion, expected *WAFExclusion) {
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

func createWAFWithRules(t *testing.T, fixtureBase string, testService *Service, version *Version, responseName string) *WAF {
	waf := createWAF(t, fixtureBase+"waf/create", testService.ID, strconv.Itoa(version.Number), "", responseName)

	var err error
	rulesIn := buildWAFRules("log")
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

func createServiceForWAF(t *testing.T, fixtureBase string) (*Service, *Version, string) {
	service := createTestService(t, fixtureBase+"service/create", "service-"+strconv.Itoa(rand.Int()))
	version := createTestVersion(t, fixtureBase+"service/version", service.ID)
	responseName := "WAf_Response"
	createTestWAFResponseObject(t, fixtureBase+"response_object/create", service.ID, responseName, version.Number)
	return service, version, responseName
}

func createWAFExclusion(t *testing.T, fixture string, excl1In *CreateWAFExclusionInput) *WAFExclusion {
	var excl1Out *WAFExclusion
	var err error
	record(t, fixture, func(c *Client) {
		excl1Out, err = c.CreateWAFExclusion(excl1In)
	})
	if err != nil {
		t.Fatal(err)
	}
	return excl1Out
}

func updateWAFExclusion(t *testing.T, fixture string, updateExcl *UpdateWAFExclusionInput) *WAFExclusion {
	var err error
	var out *WAFExclusion
	record(t, fixture, func(c *Client) {
		out, err = c.UpdateWAFExclusion(updateExcl)
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func deleteWAFExclusion(t *testing.T, fixture string, updateExcl *DeleteWAFExclusionInput) {
	var err error
	record(t, fixture, func(c *Client) {
		err = c.DeleteWAFExclusion(updateExcl)
	})
	if err != nil {
		t.Fatal(err)
	}
}

func listWAFExclusions(t *testing.T, fixture string, waf *WAF, include string) *WAFExclusionResponse {
	var err error
	var exclResp *WAFExclusionResponse
	record(t, fixture, func(c *Client) {
		exclResp, err = c.ListAllWAFExclusions(&ListAllWAFExclusionsInput{
			WAFID:            waf.ID,
			WAFVersionNumber: 1,
			Include:          &include,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	return exclResp
}

func listWAFExclusionsWithFilters(t *testing.T, fixture string, request *ListAllWAFExclusionsInput) *WAFExclusionResponse {
	var err error
	var exclResp *WAFExclusionResponse
	record(t, fixture, func(c *Client) {
		exclResp, err = c.ListAllWAFExclusions(request)
	})
	if err != nil {
		t.Fatal(err)
	}
	return exclResp
}

func buildWAFExclusion1(wafID string, version int) *CreateWAFExclusionInput {
	return &CreateWAFExclusionInput{
		WAFID:            wafID,
		WAFVersionNumber: version,
		WAFExclusion: &WAFExclusion{
			Name:          strToPtr("index page"),
			ExclusionType: strToPtr(WAFExclusionTypeRule),
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

func buildWAFExclusion2(wafID string, version int) *CreateWAFExclusionInput {
	return &CreateWAFExclusionInput{
		WAFID:            wafID,
		WAFVersionNumber: version,
		WAFExclusion: &WAFExclusion{
			Name:          strToPtr("index php"),
			ExclusionType: strToPtr(WAFExclusionTypeRule),
			Condition:     strToPtr("req.url.basename == \"index.php\""),
			Rules: []*WAFRule{
				{
					ID:       "1010070",
					ModSecID: 1010070,
				},
			},
		},
	}
}

func buildWAFExclusionWAF(wafID string, version int) *CreateWAFExclusionInput {
	return &CreateWAFExclusionInput{
		WAFID:            wafID,
		WAFVersionNumber: version,
		WAFExclusion: &WAFExclusion{
			Name:          strToPtr("index asp"),
			ExclusionType: strToPtr(WAFExclusionTypeWAF),
			Condition:     strToPtr("req.url.basename == \"index.asp\""),
		},
	}
}

func buildWAFExclusion1ForUpdate(wafID string, version, number int) *UpdateWAFExclusionInput {
	return &UpdateWAFExclusionInput{
		WAFID:            wafID,
		WAFVersionNumber: version,
		Number:           number,
		WAFExclusion: &WAFExclusion{
			Name:          strToPtr("index page"),
			ExclusionType: strToPtr(WAFExclusionTypeRule),
			Condition:     strToPtr("req.url.basename == \"index.html\""),
			Rules: []*WAFRule{
				{
					ID:       "2029718",
					ModSecID: 2029718,
				},
				{
					ID:       "1010070",
					ModSecID: 1010070,
				},
			},
		},
	}
}

func buildWAFExclusion1ForDeletion(wafID string, version, number int) *DeleteWAFExclusionInput {
	return &DeleteWAFExclusionInput{
		WAFID:            wafID,
		WAFVersionNumber: version,
		Number:           number,
	}
}
