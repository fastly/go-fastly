package fastly

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestClient_Bigqueries(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "bigqueries/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	secretKey := privateKey()

	// Create
	var bq *BigQuery
	Record(t, "bigqueries/create", func(c *Client) {
		bq, err = c.CreateBigQuery(context.TODO(), &CreateBigQueryInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           ToPointer("test-bigquery"),
			ProjectID:      ToPointer("example-fastly-log"),
			Dataset:        ToPointer("fastly_log_test"),
			Table:          ToPointer("fastly_logs"),
			Template:       ToPointer(""),
			User:           ToPointer("fastly-bigquery-log@example-fastly-log.iam.gserviceaccount.com"),
			AccountName:    ToPointer("service-account"),
			SecretKey:      ToPointer(secretKey),
			Format:         ToPointer("{\n \"timestamp\":\"%{begin:%Y-%m-%dT%H:%M:%S}t\",\n  \"time_elapsed\":%{time.elapsed.usec}V,\n  \"is_tls\":%{if(req.is_ssl, \"true\", \"false\")}V,\n  \"client_ip\":\"%{req.http.Fastly-Client-IP}V\",\n  \"geo_city\":\"%{client.geo.city}V\",\n  \"geo_country_code\":\"%{client.geo.country_code}V\",\n  \"request\":\"%{req.request}V\",\n  \"host\":\"%{req.http.Fastly-Orig-Host}V\",\n  \"url\":\"%{json.escape(req.url)}V\",\n  \"request_referer\":\"%{json.escape(req.http.Referer)}V\",\n  \"request_user_agent\":\"%{json.escape(req.http.User-Agent)}V\",\n  \"request_accept_language\":\"%{json.escape(req.http.Accept-Language)}V\",\n  \"request_accept_charset\":\"%{json.escape(req.http.Accept-Charset)}V\",\n  \"cache_status\":\"%{regsub(fastly_info.state, \"^(HIT-(SYNTH)|(HITPASS|HIT|MISS|PASS|ERROR|PIPE)).*\", \"\\\\2\\\\3\") }V\"\n}"),
			Placement:      ToPointer("waf_debug"),
			FormatVersion:  ToPointer(2),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "bigqueries/cleanup", func(c *Client) {
			_ = c.DeleteBigQuery(context.TODO(), &DeleteBigQueryInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-bigquery",
			})

			_ = c.DeleteBigQuery(context.TODO(), &DeleteBigQueryInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-bigquery",
			})
		})
	}()

	if *bq.Name != "test-bigquery" {
		t.Errorf("bad name: %q", *bq.Name)
	}
	if *bq.ProjectID != "example-fastly-log" {
		t.Errorf("bad project_id: %q", *bq.ProjectID)
	}
	if *bq.Dataset != "fastly_log_test" {
		t.Errorf("bad dataset: %q", *bq.Dataset)
	}
	if *bq.Table != "fastly_logs" {
		t.Errorf("bad table: %q", *bq.Table)
	}
	if *bq.Template != "" {
		t.Errorf("bad template_suffix: %q", *bq.Template)
	}
	if *bq.User != "fastly-bigquery-log@example-fastly-log.iam.gserviceaccount.com" {
		t.Errorf("bad user: %q", *bq.User)
	}
	if *bq.AccountName != "service-account" {
		t.Errorf("bad account name: %q", *bq.AccountName)
	}
	if strings.TrimSpace(*bq.SecretKey) != strings.TrimSpace(secretKey) {
		t.Errorf("bad secret_key: %q", *bq.SecretKey)
	}
	if *bq.Format != "{\n \"timestamp\":\"%{begin:%Y-%m-%dT%H:%M:%S}t\",\n  \"time_elapsed\":%{time.elapsed.usec}V,\n  \"is_tls\":%{if(req.is_ssl, \"true\", \"false\")}V,\n  \"client_ip\":\"%{req.http.Fastly-Client-IP}V\",\n  \"geo_city\":\"%{client.geo.city}V\",\n  \"geo_country_code\":\"%{client.geo.country_code}V\",\n  \"request\":\"%{req.request}V\",\n  \"host\":\"%{req.http.Fastly-Orig-Host}V\",\n  \"url\":\"%{json.escape(req.url)}V\",\n  \"request_referer\":\"%{json.escape(req.http.Referer)}V\",\n  \"request_user_agent\":\"%{json.escape(req.http.User-Agent)}V\",\n  \"request_accept_language\":\"%{json.escape(req.http.Accept-Language)}V\",\n  \"request_accept_charset\":\"%{json.escape(req.http.Accept-Charset)}V\",\n  \"cache_status\":\"%{regsub(fastly_info.state, \"^(HIT-(SYNTH)|(HITPASS|HIT|MISS|PASS|ERROR|PIPE)).*\", \"\\\\2\\\\3\") }V\"\n}" {
		t.Errorf("bad format: %q", *bq.Format)
	}
	if *bq.ResponseCondition != "" {
		t.Errorf("bad response_condition: %q", *bq.ResponseCondition)
	}
	if *bq.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", *bq.Placement)
	}
	if *bq.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", *bq.FormatVersion)
	}

	// List
	var bqs []*BigQuery
	Record(t, "bigqueries/list", func(c *Client) {
		bqs, err = c.ListBigQueries(context.TODO(), &ListBigQueriesInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(bqs) < 1 {
		t.Errorf("bad bigqueries: %v", bqs)
	}

	// Get
	var nbq *BigQuery
	Record(t, "bigqueries/get", func(c *Client) {
		nbq, err = c.GetBigQuery(context.TODO(), &GetBigQueryInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-bigquery",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *bq.Name != *nbq.Name {
		t.Errorf("bad name: %q", *bq.Name)
	}
	if *bq.ProjectID != *nbq.ProjectID {
		t.Errorf("bad project_id: %q", *bq.ProjectID)
	}
	if *bq.Dataset != *nbq.Dataset {
		t.Errorf("bad dataset: %q", *bq.Dataset)
	}
	if *bq.Table != *nbq.Table {
		t.Errorf("bad table: %q", *bq.Table)
	}
	if *bq.Template != *nbq.Template {
		t.Errorf("bad template_suffix: %q", *bq.Template)
	}
	if *bq.User != *nbq.User {
		t.Errorf("bad user: %q", *bq.User)
	}
	if *bq.SecretKey != *nbq.SecretKey {
		t.Errorf("bad secret_key: %q", *bq.SecretKey)
	}
	if *bq.Format != *nbq.Format {
		t.Errorf("bad format: %q", *bq.Format)
	}
	if *bq.ResponseCondition != *nbq.ResponseCondition {
		t.Errorf("bad response_condition: %q", *bq.ResponseCondition)
	}
	if *bq.Placement != *nbq.Placement {
		t.Errorf("bad placement: %q", *bq.Placement)
	}
	if *bq.FormatVersion != *nbq.FormatVersion {
		t.Errorf("bad format_version: %q", *bq.FormatVersion)
	}

	// Update
	var ubq *BigQuery
	Record(t, "bigqueries/update", func(c *Client) {
		ubq, err = c.UpdateBigQuery(context.TODO(), &UpdateBigQueryInput{
			ServiceID:        TestDeliveryServiceID,
			ServiceVersion:   *tv.Number,
			Name:             "test-bigquery",
			NewName:          ToPointer("new-test-bigquery"),
			ProcessingRegion: ToPointer("eu"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *ubq.Name != "new-test-bigquery" {
		t.Errorf("bad name: %q", *ubq.Name)
	}
	if *ubq.ProcessingRegion != "eu" {
		t.Errorf("bad log_processing_region: %q", *ubq.ProcessingRegion)
	}

	// Delete
	Record(t, "bigqueries/delete", func(c *Client) {
		err = c.DeleteBigQuery(context.TODO(), &DeleteBigQueryInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-bigquery",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListBigQueries_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListBigQueries(context.TODO(), &ListBigQueriesInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListBigQueries(context.TODO(), &ListBigQueriesInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateBigQuery_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateBigQuery(context.TODO(), &CreateBigQueryInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateBigQuery(context.TODO(), &CreateBigQueryInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetBigQuery_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetBigQuery(context.TODO(), &GetBigQueryInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetBigQuery(context.TODO(), &GetBigQueryInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetBigQuery(context.TODO(), &GetBigQueryInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateBigQuery_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateBigQuery(context.TODO(), &UpdateBigQueryInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateBigQuery(context.TODO(), &UpdateBigQueryInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateBigQuery(context.TODO(), &UpdateBigQueryInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteBigQuery_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteBigQuery(context.TODO(), &DeleteBigQueryInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteBigQuery(context.TODO(), &DeleteBigQueryInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteBigQuery(context.TODO(), &DeleteBigQueryInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
