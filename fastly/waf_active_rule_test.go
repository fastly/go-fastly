package fastly

import (
	"errors"
	"testing"
)

func TestClient_WAF_Active_Rules(t *testing.T) {
	t.Parallel()

	fixtureBase := "waf_active_rules/"

	testService := createTestService(t, fixtureBase+"service/create", "service")
	defer deleteTestService(t, fixtureBase+"/service/delete", *testService.ServiceID)

	tv := CreateTestVersion(t, fixtureBase+"/service/version", *testService.ServiceID)

	createTestLogging(t, fixtureBase+"/logging/create", *testService.ServiceID, *tv.Number)
	defer deleteTestLogging(t, fixtureBase+"/logging/delete", *testService.ServiceID, *tv.Number)

	prefetch := "WAF_Prefetch"
	createTestWAFCondition(t, fixtureBase+"/condition/create", *testService.ServiceID, prefetch, *tv.Number)

	responseName := "WAf_Response"
	createTestWAFResponseObject(t, fixtureBase+"/response_object/create", *testService.ServiceID, responseName, *tv.Number)

	waf := createWAF(t, fixtureBase+"/waf/create", *testService.ServiceID, prefetch, responseName, *tv.Number)
	defer deleteWAF(t, fixtureBase+"/waf/delete", waf.ID)

	var err error

	var rulesResp *WAFActiveRuleResponse
	Record(t, fixtureBase+"/list_empty", func(c *Client) {
		rulesResp, err = c.ListWAFActiveRules(&ListWAFActiveRulesInput{
			WAFID:            waf.ID,
			WAFVersionNumber: 1,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rulesResp.Items) != 0 {
		t.Errorf("expected 0 waf version: got %d", len(rulesResp.Items))
	}

	rulesIn := buildWAFRules("log")
	var rulesOut []*WAFActiveRule
	Record(t, fixtureBase+"/create", func(c *Client) {
		rulesOut, err = c.BatchModificationWAFActiveRules(&BatchModificationWAFActiveRulesInput{
			WAFID:            waf.ID,
			WAFVersionNumber: 1,
			Rules:            rulesIn,
			OP:               UpsertBatchOperation,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rulesOut) != len(rulesIn) {
		t.Errorf("expected 0 waf version: got %d", len(rulesOut))
	}

	Record(t, fixtureBase+"/list_not_empty", func(c *Client) {
		rulesResp, err = c.ListWAFActiveRules(&ListWAFActiveRulesInput{
			WAFID:            waf.ID,
			WAFVersionNumber: 1,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rulesResp.Items) != len(rulesIn) {
		t.Errorf("expected equal slice sizes: got expected %d  actual %d", len(rulesIn), len(rulesResp.Items))
	}
	for i := range rulesIn {
		if rulesIn[i].ModSecID != rulesOut[i].ModSecID {
			t.Errorf("Error matching:\nexpected: %#v\ngot: %#v", rulesIn[i].ModSecID, rulesOut[i].ModSecID)
		}
		if rulesIn[i].Status != rulesOut[i].Status {
			t.Errorf("Error matching:\nexpected: %#v\ngot: %#v", rulesIn[i].Status, rulesOut[i].Status)
		}
	}

	rulesIn = buildWAFRules("block")
	Record(t, fixtureBase+"/update", func(c *Client) {
		rulesOut, err = c.BatchModificationWAFActiveRules(&BatchModificationWAFActiveRulesInput{
			WAFID:            waf.ID,
			WAFVersionNumber: 1,
			Rules:            rulesIn,
			OP:               UpsertBatchOperation,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rulesOut) != len(rulesIn) {
		t.Errorf("expected 0 waf version: got %d", len(rulesOut))
	}

	Record(t, fixtureBase+"/list_not_empty2", func(c *Client) {
		rulesResp, err = c.ListWAFActiveRules(&ListWAFActiveRulesInput{
			WAFID:            waf.ID,
			WAFVersionNumber: 1,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rulesResp.Items) != len(rulesIn) {
		t.Errorf("expected equal slice sizes: got expected %d  actual %d", len(rulesIn), len(rulesResp.Items))
	}
	for i := range rulesIn {
		if rulesIn[i].ModSecID != rulesOut[i].ModSecID {
			t.Errorf("Error matching:\nexpected: %#v\ngot: %#v", rulesIn[i].ModSecID, rulesOut[i].ModSecID)
		}
		if rulesIn[i].Status != rulesOut[i].Status {
			t.Errorf("Error matching:\nexpected: %#v\ngot: %#v", rulesIn[i].Status, rulesOut[i].Status)
		}
	}

	rules := []*WAFActiveRule{{
		ModSecID: 1010070,
	}}
	Record(t, fixtureBase+"/delete_one", func(c *Client) {
		rulesOut, err = c.BatchModificationWAFActiveRules(&BatchModificationWAFActiveRulesInput{
			WAFID:            waf.ID,
			WAFVersionNumber: 1,
			Rules:            rules,
			OP:               DeleteBatchOperation,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	Record(t, fixtureBase+"/list_after_delete", func(c *Client) {
		rulesResp, err = c.ListWAFActiveRules(&ListWAFActiveRulesInput{
			WAFID:            waf.ID,
			WAFVersionNumber: 1,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(rulesResp.Items) != 2 {
		t.Errorf("expected 2 waf rules: got %d", len(rulesResp.Items))
	}
}

func TestClient_ListWAFActiveRules_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListWAFActiveRules(&ListWAFActiveRulesInput{
		WAFID: "",
	})
	if !errors.Is(err, ErrMissingWAFID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListWAFActiveRules(&ListWAFActiveRulesInput{
		WAFID:            "1",
		WAFVersionNumber: 0,
	})
	if !errors.Is(err, ErrMissingWAFVersionNumber) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_ListAllWAFActiveRules_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListAllWAFActiveRules(&ListAllWAFActiveRulesInput{
		WAFID: "",
	})
	if !errors.Is(err, ErrMissingWAFID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListAllWAFActiveRules(&ListAllWAFActiveRulesInput{
		WAFID:            "1",
		WAFVersionNumber: 0,
	})
	if !errors.Is(err, ErrMissingWAFVersionNumber) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateWAFActiveRules_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateWAFActiveRules(&CreateWAFActiveRulesInput{
		WAFID: "",
	})
	if !errors.Is(err, ErrMissingWAFID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateWAFActiveRules(&CreateWAFActiveRulesInput{
		WAFID:            "1",
		WAFVersionNumber: 0,
	})
	if !errors.Is(err, ErrMissingWAFVersionNumber) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateWAFActiveRules(&CreateWAFActiveRulesInput{
		WAFID:            "1",
		WAFVersionNumber: 1,
		Rules:            []*WAFActiveRule{},
	})
	if !errors.Is(err, ErrMissingWAFActiveRule) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_BatchModificationWAFActiveRules_validation(t *testing.T) {
	var err error
	_, err = TestClient.BatchModificationWAFActiveRules(&BatchModificationWAFActiveRulesInput{})
	if err == nil {
		t.Errorf("error expected")
	}

	var rules []*WAFActiveRule
	for i := 0; i <= BatchModifyMaximumOperations; i++ {
		rules = append(rules, &WAFActiveRule{})
	}
	_, err = TestClient.BatchModificationWAFActiveRules(&BatchModificationWAFActiveRulesInput{
		Rules: rules,
	})
	if !errors.Is(err, ErrMaxExceededRules) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteWAFActiveRules_validation(t *testing.T) {
	var err error
	err = TestClient.DeleteWAFActiveRules(&DeleteWAFActiveRulesInput{
		WAFID: "",
	})

	if !errors.Is(err, ErrMissingWAFID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteWAFActiveRules(&DeleteWAFActiveRulesInput{
		WAFID:            "1",
		WAFVersionNumber: 0,
	})

	if !errors.Is(err, ErrMissingWAFVersionNumber) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteWAFActiveRules(&DeleteWAFActiveRulesInput{
		WAFID:            "1",
		WAFVersionNumber: 1,
		Rules:            []*WAFActiveRule{},
	})

	if !errors.Is(err, ErrMissingWAFActiveRule) {
		t.Errorf("bad error: %s", err)
	}
}

func buildWAFRules(status string) []*WAFActiveRule {
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
			ModSecID: 1010070,
			Status:   status,
			Revision: 1,
		},
	}
}
