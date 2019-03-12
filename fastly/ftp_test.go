package fastly

import "testing"

func TestClient_FTPs(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "ftps/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var ftp *FTP
	record(t, "ftps/create", func(c *Client) {
		ftp, err = c.CreateFTP(&CreateFTPInput{
			Service:         testServiceID,
			Version:         tv.Number,
			Name:            "test-ftp",
			Address:         "example.com",
			Port:            1234,
			Username:        "username",
			Password:        "password",
			Path:            "/dir",
			Period:          12,
			GzipLevel:       9,
			FormatVersion:   2,
			Format:          "format",
			TimestampFormat: "%Y",
			Placement:       "waf_debug",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "ftps/cleanup", func(c *Client) {
			c.DeleteFTP(&DeleteFTPInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test-ftp",
			})

			c.DeleteFTP(&DeleteFTPInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "new-test-ftp",
			})
		})
	}()

	if ftp.Name != "test-ftp" {
		t.Errorf("bad name: %q", ftp.Name)
	}
	if ftp.Address != "example.com" {
		t.Errorf("bad address: %q", ftp.Address)
	}
	if ftp.Port != 1234 {
		t.Errorf("bad port: %q", ftp.Port)
	}
	if ftp.Username != "username" {
		t.Errorf("bad username: %q", ftp.Username)
	}
	if ftp.Password != "password" {
		t.Errorf("bad password: %q", ftp.Password)
	}
	if ftp.Path != "/dir" {
		t.Errorf("bad path: %q", ftp.Path)
	}
	if ftp.Period != 12 {
		t.Errorf("bad period: %q", ftp.Period)
	}
	if ftp.GzipLevel != 9 {
		t.Errorf("bad gzip_level: %q", ftp.GzipLevel)
	}
	if ftp.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", ftp.FormatVersion)
	}
	if ftp.Format != "format" {
		t.Errorf("bad format: %q", ftp.Format)
	}
	if ftp.TimestampFormat != "%Y" {
		t.Errorf("bad timestamp_format: %q", ftp.TimestampFormat)
	}
	if ftp.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", ftp.Placement)
	}

	// List
	var ftps []*FTP
	record(t, "ftps/list", func(c *Client) {
		ftps, err = c.ListFTPs(&ListFTPsInput{
			Service: testServiceID,
			Version: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ftps) < 1 {
		t.Errorf("bad ftps: %v", ftps)
	}

	// Get
	var nftp *FTP
	record(t, "ftps/get", func(c *Client) {
		nftp, err = c.GetFTP(&GetFTPInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-ftp",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ftp.Name != nftp.Name {
		t.Errorf("bad name: %q", ftp.Name)
	}
	if ftp.Address != nftp.Address {
		t.Errorf("bad address: %q", ftp.Address)
	}
	if ftp.Port != nftp.Port {
		t.Errorf("bad port: %q", ftp.Port)
	}
	if ftp.Username != nftp.Username {
		t.Errorf("bad username: %q", ftp.Username)
	}
	if ftp.Password != nftp.Password {
		t.Errorf("bad password: %q", ftp.Password)
	}
	if ftp.Path != nftp.Path {
		t.Errorf("bad path: %q", ftp.Path)
	}
	if ftp.Period != nftp.Period {
		t.Errorf("bad period: %q", ftp.Period)
	}
	if ftp.GzipLevel != nftp.GzipLevel {
		t.Errorf("bad gzip_level: %q", ftp.GzipLevel)
	}
	if ftp.FormatVersion != nftp.FormatVersion {
		t.Errorf("bad format_version: %q", ftp.FormatVersion)
	}
	if ftp.Format != nftp.Format {
		t.Errorf("bad format: %q", ftp.Format)
	}
	if ftp.TimestampFormat != nftp.TimestampFormat {
		t.Errorf("bad timestamp_format: %q", ftp.TimestampFormat)
	}
	if ftp.Placement != nftp.Placement {
		t.Errorf("bad placement: %q", ftp.Placement)
	}

	// Update
	var uftp *FTP
	record(t, "ftps/update", func(c *Client) {
		uftp, err = c.UpdateFTP(&UpdateFTPInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-ftp",
			NewName: "new-test-ftp",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if uftp.Name != "new-test-ftp" {
		t.Errorf("bad name: %q", uftp.Name)
	}

	// Delete
	record(t, "ftps/delete", func(c *Client) {
		err = c.DeleteFTP(&DeleteFTPInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-ftp",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListFTPs_validation(t *testing.T) {
	var err error
	_, err = testClient.ListFTPs(&ListFTPsInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListFTPs(&ListFTPsInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateFTP_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateFTP(&CreateFTPInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateFTP(&CreateFTPInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetFTP_validation(t *testing.T) {
	var err error
	_, err = testClient.GetFTP(&GetFTPInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetFTP(&GetFTPInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetFTP(&GetFTPInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateFTP_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateFTP(&UpdateFTPInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateFTP(&UpdateFTPInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateFTP(&UpdateFTPInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteFTP_validation(t *testing.T) {
	var err error
	err = testClient.DeleteFTP(&DeleteFTPInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteFTP(&DeleteFTPInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteFTP(&DeleteFTPInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
