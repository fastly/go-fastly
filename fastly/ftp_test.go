package fastly

import (
	"testing"
)

func TestClient_FTPs(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "ftps/version", func(c *Client) {
		tv = testVersion(t, c)
	})
	// Create
	var ftpCreateResp1, ftpCreateResp2, ftpCreateResp3 *FTP
	record(t, "ftps/create", func(c *Client) {
		ftpCreateResp1, err = c.CreateFTP(&CreateFTPInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-ftp",
			Address:          "example.com",
			Port:             1234,
			PublicKey:        pgpPublicKey(),
			Username:         "username",
			Password:         "password",
			Path:             "/dir",
			Period:           12,
			CompressionCodec: "snappy",
			FormatVersion:    2,
			Format:           "format",
			TimestampFormat:  "%Y",
			Placement:        "waf_debug",
			MessageType:      "classic",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "ftps/create2", func(c *Client) {
		ftpCreateResp2, err = c.CreateFTP(&CreateFTPInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-ftp-2",
			Address:          "example.com",
			Port:             1234,
			PublicKey:        pgpPublicKey(),
			Username:         "username",
			Password:         "password",
			Path:             "/dir",
			Period:           12,
			CompressionCodec: "snappy",
			GzipLevel:        8,
			FormatVersion:    2,
			Format:           "format",
			TimestampFormat:  "%Y",
			Placement:        "waf_debug",
			MessageType:      "classic",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "ftps/create3", func(c *Client) {
		ftpCreateResp3, err = c.CreateFTP(&CreateFTPInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-ftp-3",
			Address:          "example.com",
			Port:             1234,
			PublicKey:        pgpPublicKey(),
			Username:         "username",
			Password:         "password",
			Path:             "/dir",
			Period:           12,
			CompressionCodec: "snappy",
			FormatVersion:    2,
			Format:           "format",
			TimestampFormat:  "%Y",
			Placement:        "waf_debug",
			MessageType:      "classic",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "ftps/cleanup", func(c *Client) {
			c.DeleteFTP(&DeleteFTPInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-ftp",
			})

			c.DeleteFTP(&DeleteFTPInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-ftp-2",
			})

			c.DeleteFTP(&DeleteFTPInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-ftp-3",
			})

			c.DeleteFTP(&DeleteFTPInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "new-test-ftp",
			})
		})
	}()

	if ftpCreateResp1.Name != "test-ftp" {
		t.Errorf("bad name: %q", ftpCreateResp1.Name)
	}
	if ftpCreateResp1.Address != "example.com" {
		t.Errorf("bad address: %q", ftpCreateResp1.Address)
	}
	if ftpCreateResp1.Port != 1234 {
		t.Errorf("bad port: %q", ftpCreateResp1.Port)
	}
	if ftpCreateResp1.PublicKey != pgpPublicKey() {
		t.Errorf("bad public_key: %q", ftpCreateResp1.PublicKey)
	}
	if ftpCreateResp1.Username != "username" {
		t.Errorf("bad username: %q", ftpCreateResp1.Username)
	}
	if ftpCreateResp1.Password != "password" {
		t.Errorf("bad password: %q", ftpCreateResp1.Password)
	}
	if ftpCreateResp1.Path != "/dir" {
		t.Errorf("bad path: %q", ftpCreateResp1.Path)
	}
	if ftpCreateResp1.Period != 12 {
		t.Errorf("bad period: %q", ftpCreateResp1.Period)
	}
	if ftpCreateResp1.CompressionCodec != "snappy" {
		t.Errorf("bad compression_codec: %q", ftpCreateResp1.CompressionCodec)
	}
	if ftpCreateResp1.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", ftpCreateResp1.GzipLevel)
	}
	if ftpCreateResp1.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", ftpCreateResp1.FormatVersion)
	}
	if ftpCreateResp1.Format != "format" {
		t.Errorf("bad format: %q", ftpCreateResp1.Format)
	}
	if ftpCreateResp1.TimestampFormat != "%Y" {
		t.Errorf("bad timestamp_format: %q", ftpCreateResp1.TimestampFormat)
	}
	if ftpCreateResp1.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", ftpCreateResp1.Placement)
	}
	if ftpCreateResp1.MessageType != "classic" {
		t.Errorf("bad message type: %q", ftpCreateResp1.MessageType)
	}

	if ftpCreateResp2.CompressionCodec != "" {
		t.Errorf("bad compression_codec: %q", ftpCreateResp2.CompressionCodec)
	}
	if ftpCreateResp2.GzipLevel != 8 {
		t.Errorf("bad gzip_level: %q", ftpCreateResp2.GzipLevel)
	}

	if ftpCreateResp3.CompressionCodec != "snappy" {
		t.Errorf("bad compression_codec: %q", ftpCreateResp3.CompressionCodec)
	}
	if ftpCreateResp3.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", ftpCreateResp3.GzipLevel)
	}

	// List
	var ftps []*FTP
	record(t, "ftps/list", func(c *Client) {
		ftps, err = c.ListFTPs(&ListFTPsInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ftps) < 1 {
		t.Errorf("bad ftps: %v", ftps)
	}

	// Get
	var ftpGetResp *FTP
	record(t, "ftps/get", func(c *Client) {
		ftpGetResp, err = c.GetFTP(&GetFTPInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-ftp",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ftpCreateResp1.Name != ftpGetResp.Name {
		t.Errorf("bad name: %q", ftpCreateResp1.Name)
	}
	if ftpCreateResp1.Address != ftpGetResp.Address {
		t.Errorf("bad address: %q", ftpCreateResp1.Address)
	}
	if ftpCreateResp1.Port != ftpGetResp.Port {
		t.Errorf("bad port: %q", ftpCreateResp1.Port)
	}
	if ftpCreateResp1.PublicKey != ftpGetResp.PublicKey {
		t.Errorf("bad public_key: %q", ftpCreateResp1.PublicKey)
	}
	if ftpCreateResp1.Username != ftpGetResp.Username {
		t.Errorf("bad username: %q", ftpCreateResp1.Username)
	}
	if ftpCreateResp1.Password != ftpGetResp.Password {
		t.Errorf("bad password: %q", ftpCreateResp1.Password)
	}
	if ftpCreateResp1.Path != ftpGetResp.Path {
		t.Errorf("bad path: %q", ftpCreateResp1.Path)
	}
	if ftpCreateResp1.Period != ftpGetResp.Period {
		t.Errorf("bad period: %q", ftpCreateResp1.Period)
	}
	if ftpCreateResp1.CompressionCodec != ftpGetResp.CompressionCodec {
		t.Errorf("bad compression_codec: %q", ftpCreateResp1.CompressionCodec)
	}
	if ftpCreateResp1.GzipLevel != ftpGetResp.GzipLevel {
		t.Errorf("bad gzip_level: %q", ftpCreateResp1.GzipLevel)
	}
	if ftpCreateResp1.FormatVersion != ftpGetResp.FormatVersion {
		t.Errorf("bad format_version: %q", ftpCreateResp1.FormatVersion)
	}
	if ftpCreateResp1.Format != ftpGetResp.Format {
		t.Errorf("bad format: %q", ftpCreateResp1.Format)
	}
	if ftpCreateResp1.TimestampFormat != ftpGetResp.TimestampFormat {
		t.Errorf("bad timestamp_format: %q", ftpCreateResp1.TimestampFormat)
	}
	if ftpCreateResp1.Placement != ftpGetResp.Placement {
		t.Errorf("bad placement: %q", ftpCreateResp1.Placement)
	}
	if ftpCreateResp1.MessageType != ftpGetResp.MessageType {
		t.Errorf("bad message type: %q", ftpCreateResp1.MessageType)
	}

	// Update
	var ftpUpdateResp1, ftpUpdateResp2, ftpUpdateResp3 *FTP
	record(t, "ftps/update", func(c *Client) {
		ftpUpdateResp1, err = c.UpdateFTP(&UpdateFTPInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-ftp",
			NewName:          String("new-test-ftp"),
			CompressionCodec: String("zstd"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "ftps/update2", func(c *Client) {
		ftpUpdateResp2, err = c.UpdateFTP(&UpdateFTPInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-ftp-2",
			CompressionCodec: String("zstd"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	record(t, "ftps/update3", func(c *Client) {
		ftpUpdateResp3, err = c.UpdateFTP(&UpdateFTPInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-ftp-3",
			GzipLevel:      Uint8(9),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if ftpUpdateResp1.Name != "new-test-ftp" {
		t.Errorf("bad name: %q", ftpUpdateResp1.Name)
	}

	if ftpUpdateResp1.CompressionCodec != "zstd" {
		t.Errorf("bad compression_codec: %q", ftpUpdateResp1.CompressionCodec)
	}
	if ftpUpdateResp1.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", ftpUpdateResp1.GzipLevel)
	}

	if ftpUpdateResp2.CompressionCodec != "zstd" {
		t.Errorf("bad compression_codec: %q", ftpUpdateResp2.CompressionCodec)
	}
	if ftpUpdateResp2.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", ftpUpdateResp2.GzipLevel)
	}

	if ftpUpdateResp3.CompressionCodec != "" {
		t.Errorf("bad compression_codec: %q", ftpUpdateResp3.CompressionCodec)
	}
	if ftpUpdateResp3.GzipLevel != 9 {
		t.Errorf("bad gzip_level: %q", ftpUpdateResp3.GzipLevel)
	}

	// Delete
	record(t, "ftps/delete", func(c *Client) {
		err = c.DeleteFTP(&DeleteFTPInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "new-test-ftp",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListFTPs_validation(t *testing.T) {
	var err error
	_, err = testClient.ListFTPs(&ListFTPsInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListFTPs(&ListFTPsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateFTP_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateFTP(&CreateFTPInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateFTP(&CreateFTPInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetFTP_validation(t *testing.T) {
	var err error
	_, err = testClient.GetFTP(&GetFTPInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetFTP(&GetFTPInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetFTP(&GetFTPInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateFTP_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateFTP(&UpdateFTPInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateFTP(&UpdateFTPInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateFTP(&UpdateFTPInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteFTP_validation(t *testing.T) {
	var err error
	err = testClient.DeleteFTP(&DeleteFTPInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteFTP(&DeleteFTPInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteFTP(&DeleteFTPInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
