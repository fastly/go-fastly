package fastly

import (
	"reflect"
	"strconv"
	"testing"
)

func TestClient_WAF_Versions(t *testing.T) {
	t.Parallel()

	fixtureBase := "waf_versions/"

	testService := createTestService(t, fixtureBase+"service/create", "service")
	defer deleteTestService(t, fixtureBase+"/service/delete", testService.ID)

	tv := createTestVersion(t, fixtureBase+"/service/version", testService.ID)

	createTestLogging(t, fixtureBase+"/logging/create", testService.ID, tv.Number)
	defer deleteTestLogging(t, fixtureBase+"/logging/delete", testService.ID, tv.Number)

	prefetch := "WAF_Prefetch"
	condition := createTestWAFCondition(t, fixtureBase+"/condition/create", testService.ID, prefetch, tv.Number)
	defer deleteTestWAFCondition(t, fixtureBase+"/condition/delete", testService.ID, prefetch, tv.Number)

	responseName := "WAf_Response"
	ro := createTestResponseObject(t, fixtureBase+"/response_object/create", testService.ID, responseName, tv.Number)
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

	record(t, fixtureBase+"/deploy", func(c *Client) {
		err = c.DeployWAFVersion(&DeployWAFVersionInput{
			WAFID:            waf.ID,
			WAFVersionNumber: 1,
		})
	})
	if err != nil {
		t.Fatal(err)
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
	input.WAFID = waf.ID
	input.WAFVersionNumber = 2
	input.WAFVersionID = wafVer.ID
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
}

func verifyWAFVersionUpdate(t *testing.T, i *UpdateWAFVersionInput, o *WAFVersion) {

	if i.WAFVersionID != o.ID {
		t.Errorf("expected %s waf: got %s", i.WAFVersionID, o.ID)
	}
	if i.AllowedHTTPVersions != o.AllowedHTTPVersions {
		t.Errorf("expected %s waf: got %s", i.AllowedHTTPVersions, o.AllowedHTTPVersions)
	}
	if i.AllowedMethods != o.AllowedMethods {
		t.Errorf("expected %s waf: got %s", i.AllowedMethods, o.AllowedMethods)
	}
	if i.AllowedRequestContentType != o.AllowedRequestContentType {
		t.Errorf("expected %s waf: got %s", i.AllowedRequestContentType, o.AllowedRequestContentType)
	}
	if i.AllowedRequestContentTypeCharset != o.AllowedRequestContentTypeCharset {
		t.Errorf("expected %s waf: got %s", i.AllowedRequestContentTypeCharset, o.AllowedRequestContentTypeCharset)
	}
	if i.ArgLength != o.ArgLength {
		t.Errorf("expected %d waf: got %d", i.ArgLength, o.ArgLength)
	}
	if i.ArgNameLength != o.ArgNameLength {
		t.Errorf("expected %d waf: got %d", i.ArgNameLength, o.ArgNameLength)
	}
	if i.CombinedFileSizes != o.CombinedFileSizes {
		t.Errorf("expected %d waf: got %d", i.CombinedFileSizes, o.CombinedFileSizes)
	}
	if i.CriticalAnomalyScore != o.CriticalAnomalyScore {
		t.Errorf("expected %d waf: got %d", i.CriticalAnomalyScore, o.CriticalAnomalyScore)
	}
	if i.CRSValidateUTF8Encoding != o.CRSValidateUTF8Encoding {
		t.Errorf("expected %v waf: got %v", i.CRSValidateUTF8Encoding, o.CRSValidateUTF8Encoding)
	}
	if i.ErrorAnomalyScore != o.ErrorAnomalyScore {
		t.Errorf("expected %d waf: got %d", i.ErrorAnomalyScore, o.ErrorAnomalyScore)
	}
	if i.HighRiskCountryCodes != o.HighRiskCountryCodes {
		t.Errorf("expected %s waf: got %s", i.HighRiskCountryCodes, o.HighRiskCountryCodes)
	}
	if i.HTTPViolationScoreThreshold != o.HTTPViolationScoreThreshold {
		t.Errorf("expected %d waf: got %d", i.HTTPViolationScoreThreshold, o.HTTPViolationScoreThreshold)
	}
	if i.InboundAnomalyScoreThreshold != o.InboundAnomalyScoreThreshold {
		t.Errorf("expected %d waf: got %d", i.InboundAnomalyScoreThreshold, o.InboundAnomalyScoreThreshold)
	}
	if i.LFIScoreThreshold != o.LFIScoreThreshold {
		t.Errorf("expected %d waf: got %d", i.LFIScoreThreshold, o.LFIScoreThreshold)
	}
	if i.MaxFileSize != o.MaxFileSize {
		t.Errorf("expected %d waf: got %d", i.MaxFileSize, o.MaxFileSize)
	}
	if i.MaxNumArgs != o.MaxNumArgs {
		t.Errorf("expected %d waf: got %d", i.MaxNumArgs, o.MaxNumArgs)
	}
	if i.NoticeAnomalyScore != o.NoticeAnomalyScore {
		t.Errorf("expected %d waf: got %d", i.NoticeAnomalyScore, o.NoticeAnomalyScore)
	}
	if i.ParanoiaLevel != o.ParanoiaLevel {
		t.Errorf("expected %d waf: got %d", i.ParanoiaLevel, o.ParanoiaLevel)
	}
	if i.PHPInjectionScoreThreshold != o.PHPInjectionScoreThreshold {
		t.Errorf("expected %d waf: got %d", i.PHPInjectionScoreThreshold, o.PHPInjectionScoreThreshold)
	}
	if i.RCEScoreThreshold != o.RCEScoreThreshold {
		t.Errorf("expected %d waf: got %d", i.RCEScoreThreshold, o.RCEScoreThreshold)
	}
	if i.RestrictedExtensions != o.RestrictedExtensions {
		t.Errorf("expected %s waf: got %s", i.RestrictedExtensions, o.RestrictedExtensions)
	}
	if i.RestrictedHeaders != o.RestrictedHeaders {
		t.Errorf("expected %s waf: got %s", i.RestrictedHeaders, o.RestrictedHeaders)
	}
	if i.RFIScoreThreshold != o.RFIScoreThreshold {
		t.Errorf("expected %d waf: got %d", i.RFIScoreThreshold, o.RFIScoreThreshold)
	}
	if i.SessionFixationScoreThreshold != o.SessionFixationScoreThreshold {
		t.Errorf("expected %d waf: got %d", i.SessionFixationScoreThreshold, o.SessionFixationScoreThreshold)
	}
	if i.SQLInjectionScoreThreshold != o.SQLInjectionScoreThreshold {
		t.Errorf("expected %d waf: got %d", i.SQLInjectionScoreThreshold, o.SQLInjectionScoreThreshold)
	}
	if i.TotalArgLength != o.TotalArgLength {
		t.Errorf("expected %d waf: got %d", i.TotalArgLength, o.TotalArgLength)
	}
	if i.WarningAnomalyScore != o.WarningAnomalyScore {
		t.Errorf("expected %d waf: got %d", i.WarningAnomalyScore, o.WarningAnomalyScore)
	}
	if i.XSSScoreThreshold != o.XSSScoreThreshold {
		t.Errorf("expected %d waf: got %d", i.XSSScoreThreshold, o.XSSScoreThreshold)
	}

}

func buildUpdateInput() *UpdateWAFVersionInput {
	return &UpdateWAFVersionInput{
		Comment:                          "my comment",
		AllowedHTTPVersions:              "HTTP/1.0 HTTP/1.1",
		AllowedMethods:                   "GET HEAD POST",
		AllowedRequestContentType:        "application/x-www-form-urlencoded|multipart/form-data|text/xml|application/xml",
		AllowedRequestContentTypeCharset: "utf-8|iso-8859-1",
		ArgLength:                        800,
		ArgNameLength:                    200,
		CombinedFileSizes:                20000000,
		CriticalAnomalyScore:             12,
		CRSValidateUTF8Encoding:          true,
		ErrorAnomalyScore:                10,
		HighRiskCountryCodes:             "gb",
		HTTPViolationScoreThreshold:      20,
		InboundAnomalyScoreThreshold:     20,
		LFIScoreThreshold:                20,
		MaxFileSize:                      20000000,
		MaxNumArgs:                       510,
		NoticeAnomalyScore:               8,
		ParanoiaLevel:                    2,
		PHPInjectionScoreThreshold:       20,
		RCEScoreThreshold:                20,
		RestrictedExtensions:             ".asa/ .asax/ .ascx/ .axd/ .backup/ .bak/ .bat/ .cdx/ .cer/ .cfg/ .cmd/ .com/",
		RestrictedHeaders:                "/proxy/ /lock-token/",
		RFIScoreThreshold:                20,
		SessionFixationScoreThreshold:    20,
		SQLInjectionScoreThreshold:       20,
		TotalArgLength:                   12800,
		WarningAnomalyScore:              20,
		XSSScoreThreshold:                20,
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

	resp, err := testClient.ListAllWAFVersions(&ListAllWAFVersionsInput{
		WAFID: "4QXAURauMXa4KHQ3kRn5Yr",
	})
	if err != nil {
		t.Errorf("bad error: %s", err)
	}

	print(resp)
}

func TestClient_GetWAFVersion_validation(t *testing.T) {
	var err error
	_, err = testClient.GetWAFVersion(&GetWAFVersionInput{
		WAFID: "",
	})
	if err != ErrMissingWAFID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateWAFVersion_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateWAFVersion(&UpdateWAFVersionInput{
		WAFID: "",
	})
	if err != ErrMissingWAFID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateWAFVersion(&UpdateWAFVersionInput{
		WAFID:            "1",
		WAFVersionNumber: 0,
	})
	if err != ErrMissingWAFVersionNumber {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateWAFVersion(&UpdateWAFVersionInput{
		WAFID:            "1",
		WAFVersionNumber: 1,
		WAFVersionID:     "",
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
