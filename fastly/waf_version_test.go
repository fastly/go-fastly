package fastly

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestClient_WAF_Versions(t *testing.T) {
	t.Parallel()

	fixtureBase := "waf_versions/"

	testService := createTestService(t, fixtureBase+"service/create", "service3")
	defer deleteTestService(t, fixtureBase+"/service/delete", testService.ID)

	tv := createTestVersion(t, fixtureBase+"/service/version", testService.ID)

	createTestLogging(t, fixtureBase+"/logging/create", testService.ID, tv.Number)
	defer deleteTestLogging(t, fixtureBase+"/logging/delete", testService.ID, tv.Number)

	prefetch := "WAF_Prefetch"
	condition := createTestWAFCondition(t, fixtureBase+"/condition/create", testService.ID, prefetch, tv.Number)
	defer deleteTestCondition(t, fixtureBase+"/condition/delete", testService.ID, prefetch, tv.Number)

	responseName := "WAf_Response"
	ro := createTestWAFResponseObject(t, fixtureBase+"/response_object/create", testService.ID, responseName, tv.Number)
	defer deleteTestResponseObject(t, fixtureBase+"/response_object/delete", testService.ID, responseName, tv.Number)

	var err error
	var waf *WAF
	record(t, fixtureBase+"/waf/create", func(c *Client) {
		waf, err = c.CreateWAF(&CreateWAFInput{
			Service:           testService.ID,
			Version:           strconv.Itoa(tv.Number),
			PrefetchCondition: condition.Name,
			Response:          ro.Name,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		record(t, fixtureBase+"/waf/delete", func(c *Client) {
			if err := c.DeleteWAF(&DeleteWAFInput{
				Version: strconv.Itoa(tv.Number),
				ID:      waf.ID,
			}); err != nil {
				t.Fatal(err)
			}
		})
	}()

	var wafVerResp *WAFVersionResponse
	record(t, fixtureBase+"/list", func(c *Client) {
		wafVerResp, err = c.ListWAFVersions(&ListWAFVersionsInput{
			WAFID: waf.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(wafVerResp.Items) != 1 {
		t.Errorf("expected 1 waf: got %d", len(wafVerResp.Items))
	}
	if wafVerResp.Items[0].LastDeploymentStatus != "" {
		t.Errorf("unexpected waf deployment status: \"%s\"", wafVerResp.Items[0].LastDeploymentStatus)
	}

	record(t, fixtureBase+"/deploy", func(c *Client) {
		err = c.DeployWAFVersion(&DeployWAFVersionInput{
			WAFID:            waf.ID,
			WAFVersionNumber: 1,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	var wafVerPD *WAFVersion
	for i := 0; i < 120 && (wafVerPD == nil || wafVerPD.LastDeploymentStatus != WAFVersionDeploymentStatusCompleted); i++ {
		record(t, fmt.Sprintf("%s/getPostDeploy_%d", fixtureBase, i), func(c *Client) {
			wafVerPD, err = c.GetWAFVersion(&GetWAFVersionInput{
				WAFID:            waf.ID,
				WAFVersionNumber: 1,
			})
		})
		if err != nil {
			t.Fatal(err)
			break
		}
		if wafVerPD == nil {
			t.Error("expected waf, got nil")
			break
		}
		if wafVerPD.LastDeploymentStatus != WAFVersionDeploymentStatusPending &&
			wafVerPD.LastDeploymentStatus != WAFVersionDeploymentStatusInProgress &&
			wafVerPD.LastDeploymentStatus != WAFVersionDeploymentStatusFailed &&
			wafVerPD.LastDeploymentStatus != WAFVersionDeploymentStatusCompleted {
			t.Errorf("unexpected waf deployment status: \"%s\"", wafVerPD.LastDeploymentStatus)
			break
		}

		if wafVerPD.LastDeploymentStatus != WAFVersionDeploymentStatusCompleted {
			time.Sleep(500 * time.Millisecond)
		}
	}
	if wafVerPD.LastDeploymentStatus != WAFVersionDeploymentStatusCompleted {
		t.Error("waf deployment did not reach completed status")
	}

	var wafVer *WAFVersion
	record(t, fixtureBase+"/clone", func(c *Client) {
		wafVer, err = c.CloneWAFVersion(&CloneWAFVersionInput{
			WAFID:            waf.ID,
			WAFVersionNumber: 1,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if wafVer == nil {
		t.Errorf("expected 1 waf: got %d", len(wafVerResp.Items))
	}

	record(t, fixtureBase+"/get", func(c *Client) {
		wafVer, err = c.GetWAFVersion(&GetWAFVersionInput{
			WAFID:            waf.ID,
			WAFVersionNumber: 2,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if wafVer == nil {
		t.Error("expected waf, got nil")
	}

	input := buildUpdateInput()
	input.WAFID = &waf.ID
	input.WAFVersionNumber = intToPtr(2)
	input.WAFVersionID = &wafVer.ID
	record(t, fixtureBase+"/update", func(c *Client) {
		wafVer, err = c.UpdateWAFVersion(input)
	})
	if err != nil {
		t.Fatal(err)
	}
	if wafVer == nil {
		t.Error("expected waf, got nil")
	}
	verifyWAFVersionUpdate(t, input, wafVer)

	record(t, fixtureBase+"/lock", func(c *Client) {
		wafVer, err = c.LockWAFVersion(&LockWAFVersionInput{
			WAFID:            waf.ID,
			WAFVersionNumber: 2,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if wafVer == nil {
		t.Error("expected waf, got nil")
	}
	if !wafVer.Locked {
		t.Errorf("expected locked = true waf: got locked == %v", wafVer.Locked)
	}

	record(t, fixtureBase+"/list_all", func(c *Client) {
		wafVerResp, err = c.ListAllWAFVersions(&ListAllWAFVersionsInput{
			WAFID: waf.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(wafVerResp.Items) != 2 {
		t.Errorf("expected 2 waf: got %d", len(wafVerResp.Items))
	}

	record(t, fixtureBase+"/create_empty", func(c *Client) {
		wafVer, err = c.CreateEmptyWAFVersion(&CreateEmptyWAFVersionInput{
			WAFID: waf.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if wafVer == nil {
		t.Error("expected waf, got nil")
	}
	verifyEmptyWAFVersion(t, wafVer)
}

func verifyWAFVersionUpdate(t *testing.T, i *UpdateWAFVersionInput, o *WAFVersion) {

	if *i.WAFVersionID != o.ID {
		t.Errorf("expected %s waf: got %s", *i.WAFVersionID, o.ID)
	}
	if *i.AllowedHTTPVersions != o.AllowedHTTPVersions {
		t.Errorf("expected %s waf: got %s", *i.AllowedHTTPVersions, o.AllowedHTTPVersions)
	}
	if *i.AllowedMethods != o.AllowedMethods {
		t.Errorf("expected %s waf: got %s", *i.AllowedMethods, o.AllowedMethods)
	}
	if *i.AllowedRequestContentType != o.AllowedRequestContentType {
		t.Errorf("expected %s waf: got %s", *i.AllowedRequestContentType, o.AllowedRequestContentType)
	}
	if *i.AllowedRequestContentTypeCharset != o.AllowedRequestContentTypeCharset {
		t.Errorf("expected %s waf: got %s", *i.AllowedRequestContentTypeCharset, o.AllowedRequestContentTypeCharset)
	}
	if *i.ArgLength != o.ArgLength {
		t.Errorf("expected %d waf: got %d", *i.ArgLength, o.ArgLength)
	}
	if *i.ArgNameLength != o.ArgNameLength {
		t.Errorf("expected %d waf: got %d", *i.ArgNameLength, o.ArgNameLength)
	}
	if *i.CombinedFileSizes != o.CombinedFileSizes {
		t.Errorf("expected %d waf: got %d", *i.CombinedFileSizes, o.CombinedFileSizes)
	}
	if *i.CriticalAnomalyScore != o.CriticalAnomalyScore {
		t.Errorf("expected %d waf: got %d", *i.CriticalAnomalyScore, o.CriticalAnomalyScore)
	}
	if *i.CRSValidateUTF8Encoding != o.CRSValidateUTF8Encoding {
		t.Errorf("expected %v waf: got %v", *i.CRSValidateUTF8Encoding, o.CRSValidateUTF8Encoding)
	}
	if *i.ErrorAnomalyScore != o.ErrorAnomalyScore {
		t.Errorf("expected %d waf: got %d", *i.ErrorAnomalyScore, o.ErrorAnomalyScore)
	}
	if *i.HighRiskCountryCodes != o.HighRiskCountryCodes {
		t.Errorf("expected %s waf: got %s", *i.HighRiskCountryCodes, o.HighRiskCountryCodes)
	}
	if *i.HTTPViolationScoreThreshold != o.HTTPViolationScoreThreshold {
		t.Errorf("expected %d waf: got %d", *i.HTTPViolationScoreThreshold, o.HTTPViolationScoreThreshold)
	}
	if *i.InboundAnomalyScoreThreshold != o.InboundAnomalyScoreThreshold {
		t.Errorf("expected %d waf: got %d", *i.InboundAnomalyScoreThreshold, o.InboundAnomalyScoreThreshold)
	}
	if *i.LFIScoreThreshold != o.LFIScoreThreshold {
		t.Errorf("expected %d waf: got %d", *i.LFIScoreThreshold, o.LFIScoreThreshold)
	}
	if *i.MaxFileSize != o.MaxFileSize {
		t.Errorf("expected %d waf: got %d", *i.MaxFileSize, o.MaxFileSize)
	}
	if *i.MaxNumArgs != o.MaxNumArgs {
		t.Errorf("expected %d waf: got %d", *i.MaxNumArgs, o.MaxNumArgs)
	}
	if *i.NoticeAnomalyScore != o.NoticeAnomalyScore {
		t.Errorf("expected %d waf: got %d", *i.NoticeAnomalyScore, o.NoticeAnomalyScore)
	}
	if *i.ParanoiaLevel != o.ParanoiaLevel {
		t.Errorf("expected %d waf: got %d", *i.ParanoiaLevel, o.ParanoiaLevel)
	}
	if *i.PHPInjectionScoreThreshold != o.PHPInjectionScoreThreshold {
		t.Errorf("expected %d waf: got %d", *i.PHPInjectionScoreThreshold, o.PHPInjectionScoreThreshold)
	}
	if *i.RCEScoreThreshold != o.RCEScoreThreshold {
		t.Errorf("expected %d waf: got %d", *i.RCEScoreThreshold, o.RCEScoreThreshold)
	}
	if *i.RestrictedExtensions != o.RestrictedExtensions {
		t.Errorf("expected %s waf: got %s", *i.RestrictedExtensions, o.RestrictedExtensions)
	}
	if *i.RestrictedHeaders != o.RestrictedHeaders {
		t.Errorf("expected %s waf: got %s", *i.RestrictedHeaders, o.RestrictedHeaders)
	}
	if *i.RFIScoreThreshold != o.RFIScoreThreshold {
		t.Errorf("expected %d waf: got %d", *i.RFIScoreThreshold, o.RFIScoreThreshold)
	}
	if *i.SessionFixationScoreThreshold != o.SessionFixationScoreThreshold {
		t.Errorf("expected %d waf: got %d", *i.SessionFixationScoreThreshold, o.SessionFixationScoreThreshold)
	}
	if *i.SQLInjectionScoreThreshold != o.SQLInjectionScoreThreshold {
		t.Errorf("expected %d waf: got %d", *i.SQLInjectionScoreThreshold, o.SQLInjectionScoreThreshold)
	}
	if *i.TotalArgLength != o.TotalArgLength {
		t.Errorf("expected %d waf: got %d", *i.TotalArgLength, o.TotalArgLength)
	}
	if *i.WarningAnomalyScore != o.WarningAnomalyScore {
		t.Errorf("expected %d waf: got %d", *i.WarningAnomalyScore, o.WarningAnomalyScore)
	}
	if *i.XSSScoreThreshold != o.XSSScoreThreshold {
		t.Errorf("expected %d waf: got %d", *i.XSSScoreThreshold, o.XSSScoreThreshold)
	}
}

func verifyEmptyWAFVersion(t *testing.T, o *WAFVersion) {

	threshold := 999
	if threshold != o.HTTPViolationScoreThreshold {
		t.Errorf("expected  %d HTTPViolationScoreThreshold: got %d", threshold, o.HTTPViolationScoreThreshold)
	}
	if threshold != o.InboundAnomalyScoreThreshold {
		t.Errorf("expected %d InboundAnomalyScoreThreshold: got %d", threshold, o.InboundAnomalyScoreThreshold)
	}
	if threshold != o.InboundAnomalyScoreThreshold {
		t.Errorf("expected %d InboundAnomalyScoreThreshold: got %d", threshold, o.InboundAnomalyScoreThreshold)
	}
	if threshold != o.PHPInjectionScoreThreshold {
		t.Errorf("expected %d PHPInjectionScoreThreshold: got %d", threshold, o.PHPInjectionScoreThreshold)
	}
	if threshold != o.RCEScoreThreshold {
		t.Errorf("expected %d RCEScoreThreshold: got %d", threshold, o.RCEScoreThreshold)
	}
	if threshold != o.RFIScoreThreshold {
		t.Errorf("expected %d waf: RFIScoreThreshold %d", threshold, o.RFIScoreThreshold)
	}
	if threshold != o.SessionFixationScoreThreshold {
		t.Errorf("expected %d SessionFixationScoreThreshold: got %d", threshold, o.SessionFixationScoreThreshold)
	}
	if threshold != o.SQLInjectionScoreThreshold {
		t.Errorf("expected %d SQLInjectionScoreThreshold: got %d", threshold, o.SQLInjectionScoreThreshold)
	}
	if threshold != o.XSSScoreThreshold {
		t.Errorf("expected %d XSSScoreThreshold: got %d", threshold, o.XSSScoreThreshold)
	}

	totalRules := o.ActiveRulesFastlyBlockCount + o.ActiveRulesFastlyLogCount + o.ActiveRulesOWASPBlockCount +
		o.ActiveRulesOWASPLogCount + o.ActiveRulesOWASPScoreCount + o.ActiveRulesTrustwaveBlockCount + o.ActiveRulesTrustwaveLogCount

	if totalRules != 0 {
		t.Errorf("expected no active rules rules: got %d", totalRules)
	}
}

func buildUpdateInput() *UpdateWAFVersionInput {
	return &UpdateWAFVersionInput{
		Comment:                          strToPtr("my comment"),
		AllowedHTTPVersions:              strToPtr("HTTP/1.0 HTTP/1.1"),
		AllowedMethods:                   strToPtr("GET HEAD POST"),
		AllowedRequestContentType:        strToPtr("application/x-www-form-urlencoded|multipart/form-data|text/xml|application/xml"),
		AllowedRequestContentTypeCharset: strToPtr("utf-8|iso-8859-1"),
		ArgLength:                        intToPtr(800),
		ArgNameLength:                    intToPtr(200),
		CombinedFileSizes:                intToPtr(20000000),
		CriticalAnomalyScore:             intToPtr(12),
		CRSValidateUTF8Encoding:          boolToPtr(true),
		ErrorAnomalyScore:                intToPtr(10),
		HighRiskCountryCodes:             strToPtr("gb"),
		HTTPViolationScoreThreshold:      intToPtr(20),
		InboundAnomalyScoreThreshold:     intToPtr(20),
		LFIScoreThreshold:                intToPtr(20),
		MaxFileSize:                      intToPtr(20000000),
		MaxNumArgs:                       intToPtr(510),
		NoticeAnomalyScore:               intToPtr(8),
		ParanoiaLevel:                    intToPtr(2),
		PHPInjectionScoreThreshold:       intToPtr(20),
		RCEScoreThreshold:                intToPtr(20),
		RestrictedExtensions:             strToPtr(".asa/ .asax/ .ascx/ .axd/ .backup/ .bak/ .bat/ .cdx/ .cer/ .cfg/ .cmd/ .com/"),
		RestrictedHeaders:                strToPtr("/proxy/ /lock-token/"),
		RFIScoreThreshold:                intToPtr(20),
		SessionFixationScoreThreshold:    intToPtr(20),
		SQLInjectionScoreThreshold:       intToPtr(20),
		TotalArgLength:                   intToPtr(12800),
		WarningAnomalyScore:              intToPtr(20),
		XSSScoreThreshold:                intToPtr(20),
	}
}

func TestClient_listWAFVersions_formatFilters(t *testing.T) {
	cases := []struct {
		remote *ListWAFVersionsInput
		local  map[string]string
	}{
		{
			remote: &ListWAFVersionsInput{
				PageSize:   2,
				PageNumber: 2,
				Include:    "included",
			},
			local: map[string]string{
				"page[size]":   "2",
				"page[number]": "2",
				"include":      "included",
			},
		},
	}
	for _, c := range cases {
		out := c.remote.formatFilters()
		if !reflect.DeepEqual(out, c.local) {
			t.Fatalf("Error matching:\nexpected: %#v\n     got: %#v", c.local, out)
		}
	}
}

func TestClient_ListWAFVersions_validation(t *testing.T) {
	var err error
	_, err = testClient.ListWAFVersions(&ListWAFVersionsInput{
		WAFID: "",
	})
	if err != ErrMissingWAFID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_ListAllWAFVersions_validation(t *testing.T) {
	var err error
	_, err = testClient.ListAllWAFVersions(&ListAllWAFVersionsInput{
		WAFID: "",
	})
	if err != ErrMissingWAFID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetWAFVersion_validation(t *testing.T) {
	var err error
	_, err = testClient.GetWAFVersion(&GetWAFVersionInput{
		WAFID: "",
	})
	if err != ErrMissingWAFID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetWAFVersion(&GetWAFVersionInput{
		WAFID:            "1",
		WAFVersionNumber: 0,
	})
	if err != ErrMissingWAFVersionNumber {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateWAFVersion_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateWAFVersion(&UpdateWAFVersionInput{
		WAFID: strToPtr(""),
	})
	if err != ErrMissingWAFID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateWAFVersion(&UpdateWAFVersionInput{
		WAFID:            strToPtr("1"),
		WAFVersionNumber: intToPtr(0),
	})
	if err != ErrMissingWAFVersionNumber {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateWAFVersion(&UpdateWAFVersionInput{
		WAFID:            strToPtr("1"),
		WAFVersionNumber: intToPtr(1),
		WAFVersionID:     strToPtr(""),
	})
	if err != ErrMissingWAFVersionID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_LockWAFVersion_validation(t *testing.T) {
	var err error
	_, err = testClient.LockWAFVersion(&LockWAFVersionInput{
		WAFID: "",
	})
	if err != ErrMissingWAFID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.LockWAFVersion(&LockWAFVersionInput{
		WAFID:            "1",
		WAFVersionNumber: 0,
	})
	if err != ErrMissingWAFVersionNumber {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CloneWAFVersion_validation(t *testing.T) {
	var err error
	_, err = testClient.CloneWAFVersion(&CloneWAFVersionInput{
		WAFID: "",
	})
	if err != ErrMissingWAFID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CloneWAFVersion(&CloneWAFVersionInput{
		WAFID:            "1",
		WAFVersionNumber: 0,
	})
	if err != ErrMissingWAFVersionNumber {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeployWAFVersion_validation(t *testing.T) {
	var err error
	if err = testClient.DeployWAFVersion(&DeployWAFVersionInput{
		WAFID: "",
	}); err != ErrMissingWAFID {
		t.Errorf("bad error: %s", err)
	}

	if err = testClient.DeployWAFVersion(&DeployWAFVersionInput{
		WAFID:            "1",
		WAFVersionNumber: 0,
	}); err != ErrMissingWAFVersionNumber {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateWAFVersionInput_HasChanges(t *testing.T) {

	cases := []struct {
		in  UpdateWAFVersionInput
		out bool
	}{
		{
			in: UpdateWAFVersionInput{
				WAFID:            strToPtr("ID"),
				WAFVersionNumber: intToPtr(1),
				WAFVersionID:     strToPtr("versionID"),
			},
			out: false,
		},
		{
			in: UpdateWAFVersionInput{
				WAFID:            strToPtr("ID"),
				WAFVersionNumber: intToPtr(1),
				WAFVersionID:     strToPtr("versionID"),
				AllowedMethods:   strToPtr("any"),
			},
			out: true,
		},
	}
	for _, c := range cases {
		empty := c.in.HasChanges()
		if empty != c.out {
			t.Fatalf("Error matching:\nexpected: %#v\n     got: %#v", c.out, empty)
		}
	}
}

func TestClient_CreateEmptyWAFVersion_validation(t *testing.T) {
	var err error
	if _, err = testClient.CreateEmptyWAFVersion(&CreateEmptyWAFVersionInput{
		WAFID: "",
	}); err != ErrMissingWAFID {
		t.Errorf("bad error: %s", err)
	}
}

func strToPtr(s string) *string {
	return &s
}

func intToPtr(i int) *int {
	return &i
}

func boolToPtr(i bool) *bool {
	return &i
}
