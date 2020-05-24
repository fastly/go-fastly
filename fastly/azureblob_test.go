package fastly

import "testing"

func TestClient_AzureBlobs(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "azureblobs/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var azureblob *AzureBlob
	record(t, "azureblobs/create", func(c *Client) {
		azureblob, err = c.CreateAzureBlob(&CreateAzureBlobInput{
			Service:         testServiceID,
			Version:         tv.Number,
			Name:            "test-azureblob",
			Container:       "container",
			AccountName:     "account_name",
			SASToken:        "sv=2018-04-05&st=2018-04-29T22%3A18%3A26Z&sr=b&se=2020-04-30T02%3A23%3A26Z&sp=w&sig=Z%2FRHIX5Xcg0Mq2rqI3OlWTjEg2tYkboXr1P9ZUXDtkk%3D",
			Path:            "/path",
			Period:          12,
			GzipLevel:       9,
			FormatVersion:   2,
			Format:          "format",
			MessageType:     "blank",
			TimestampFormat: "%Y",
			Placement:       "waf_debug",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "azureblobs/cleanup", func(c *Client) {
			c.DeleteAzureBlob(&DeleteAzureBlobInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test-azureblob",
			})

			c.DeleteAzureBlob(&DeleteAzureBlobInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "new-test-azureblob",
			})
		})
	}()

	if azureblob.Name != "test-azureblob" {
		t.Errorf("bad name: %q", azureblob.Name)
	}
	if azureblob.Container != "container" {
		t.Errorf("bad container: %q", azureblob.Container)
	}
	if azureblob.AccountName != "account_name" {
		t.Errorf("bad account_name: %q", azureblob.AccountName)
	}
	if azureblob.SASToken != "sv=2018-04-05&st=2018-04-29T22%3A18%3A26Z&sr=b&se=2020-04-30T02%3A23%3A26Z&sp=w&sig=Z%2FRHIX5Xcg0Mq2rqI3OlWTjEg2tYkboXr1P9ZUXDtkk%3D" {
		t.Errorf("bad sas_token: %q", azureblob.SASToken)
	}
	if azureblob.Path != "/path" {
		t.Errorf("bad path: %q", azureblob.Path)
	}
	if azureblob.Period != 12 {
		t.Errorf("bad period: %q", azureblob.Period)
	}
	if azureblob.GzipLevel != 9 {
		t.Errorf("bad gzip_level: %q", azureblob.GzipLevel)
	}
	if azureblob.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", azureblob.FormatVersion)
	}
	if azureblob.Format != "format" {
		t.Errorf("bad format: %q", azureblob.Format)
	}
	if azureblob.TimestampFormat != "%Y" {
		t.Errorf("bad timestamp_format: %q", azureblob.TimestampFormat)
	}
	if azureblob.MessageType != "blank" {
		t.Errorf("bad message_type: %q", azureblob.MessageType)
	}
	if azureblob.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", azureblob.Placement)
	}

	// List
	var azureblobs []*AzureBlob
	record(t, "azureblobs/list", func(c *Client) {
		azureblobs, err = c.ListAzureBlobs(&ListAzureBlobsInput{
			Service: testServiceID,
			Version: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(azureblobs) < 1 {
		t.Errorf("bad azureblobs: %v", azureblobs)
	}

	// Get
	var nazureblob *AzureBlob
	record(t, "azureblobs/get", func(c *Client) {
		nazureblob, err = c.GetAzureBlob(&GetAzureBlobInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-azureblob",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if azureblob.Name != nazureblob.Name {
		t.Errorf("bad name: %q", azureblob.Name)
	}
	if azureblob.Container != nazureblob.Container {
		t.Errorf("bad container: %q", azureblob.Container)
	}
	if azureblob.AccountName != nazureblob.AccountName {
		t.Errorf("bad account_name: %q", azureblob.AccountName)
	}
	if azureblob.SASToken != nazureblob.SASToken {
		t.Errorf("bad sas_token: %q", azureblob.SASToken)
	}
	if azureblob.Path != nazureblob.Path {
		t.Errorf("bad path: %q", azureblob.Path)
	}
	if azureblob.Period != nazureblob.Period {
		t.Errorf("bad period: %q", azureblob.Period)
	}
	if azureblob.GzipLevel != nazureblob.GzipLevel {
		t.Errorf("bad gzip_level: %q", azureblob.GzipLevel)
	}
	if azureblob.FormatVersion != nazureblob.FormatVersion {
		t.Errorf("bad format_version: %q", azureblob.FormatVersion)
	}
	if azureblob.Format != nazureblob.Format {
		t.Errorf("bad format: %q", azureblob.Format)
	}
	if azureblob.TimestampFormat != nazureblob.TimestampFormat {
		t.Errorf("bad timestamp_format: %q", azureblob.TimestampFormat)
	}
	if azureblob.MessageType != nazureblob.MessageType {
		t.Errorf("bad message_type: %q", azureblob.MessageType)
	}
	if azureblob.Placement != nazureblob.Placement {
		t.Errorf("bad placement: %q", azureblob.Placement)
	}

	// Update
	var uazureblob *AzureBlob
	record(t, "azureblobs/update", func(c *Client) {
		uazureblob, err = c.UpdateAzureBlob(&UpdateAzureBlobInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-azureblob",
			NewName: "new-test-azureblob",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if uazureblob.Name != "new-test-azureblob" {
		t.Errorf("bad name: %q", uazureblob.Name)
	}

	// Delete
	record(t, "azureblobs/delete", func(c *Client) {
		err = c.DeleteAzureBlob(&DeleteAzureBlobInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-azureblob",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListAzureBlobs_validation(t *testing.T) {
	var err error
	_, err = testClient.ListAzureBlobs(&ListAzureBlobsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListAzureBlobs(&ListAzureBlobsInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateAzureBlob_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateAzureBlob(&CreateAzureBlobInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateAzureBlob(&CreateAzureBlobInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetAzureBlob_validation(t *testing.T) {
	var err error
	_, err = testClient.GetAzureBlob(&GetAzureBlobInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetAzureBlob(&GetAzureBlobInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetAzureBlob(&GetAzureBlobInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateAzureBlob_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateAzureBlob(&UpdateAzureBlobInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateAzureBlob(&UpdateAzureBlobInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateAzureBlob(&UpdateAzureBlobInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteAzureBlob_validation(t *testing.T) {
	var err error
	err = testClient.DeleteAzureBlob(&DeleteAzureBlobInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteAzureBlob(&DeleteAzureBlobInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteAzureBlob(&DeleteAzureBlobInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
