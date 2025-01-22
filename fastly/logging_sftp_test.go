package fastly

import (
	"errors"
	"strings"
	"testing"
)

func TestClient_SFTPs(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "sftps/version", func(c *Client) {
		tv = testVersion(t, c)
	})
	knownHosts := "example.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEmuNYPrqg9tjqfR14ye3Pvsm9sjIx6EJwm5tMXIMaN1"

	// Create
	var sftpCreateResp1, sftpCreateResp2, sftpCreateResp3 *SFTP
	Record(t, "sftps/create", func(c *Client) {
		sftpCreateResp1, err = c.CreateSFTP(&CreateSFTPInput{
			ServiceID:        TestDeliveryServiceID,
			ServiceVersion:   *tv.Number,
			Name:             ToPointer("test-sftp"),
			Address:          ToPointer("example.com"),
			Port:             ToPointer(1234),
			PublicKey:        ToPointer(pgpPublicKey()),
			SecretKey:        ToPointer(privateKey()),
			SSHKnownHosts:    ToPointer(knownHosts),
			User:             ToPointer("username"),
			Password:         ToPointer("password"),
			Path:             ToPointer("/dir"),
			Period:           ToPointer(12),
			CompressionCodec: ToPointer("snappy"),
			FormatVersion:    ToPointer(2),
			Format:           ToPointer("format"),
			MessageType:      ToPointer("blank"),
			TimestampFormat:  ToPointer("%Y"),
			Placement:        ToPointer("waf_debug"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	Record(t, "sftps/create2", func(c *Client) {
		sftpCreateResp2, err = c.CreateSFTP(&CreateSFTPInput{
			ServiceID:       TestDeliveryServiceID,
			ServiceVersion:  *tv.Number,
			Name:            ToPointer("test-sftp-2"),
			Address:         ToPointer("example.com"),
			Port:            ToPointer(1234),
			PublicKey:       ToPointer(pgpPublicKey()),
			SecretKey:       ToPointer(privateKey()),
			SSHKnownHosts:   ToPointer(knownHosts),
			User:            ToPointer("username"),
			Password:        ToPointer("password"),
			Path:            ToPointer("/dir"),
			Period:          ToPointer(12),
			GzipLevel:       ToPointer(8),
			FormatVersion:   ToPointer(2),
			Format:          ToPointer("format"),
			MessageType:     ToPointer("blank"),
			TimestampFormat: ToPointer("%Y"),
			Placement:       ToPointer("waf_debug"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	Record(t, "sftps/create3", func(c *Client) {
		sftpCreateResp3, err = c.CreateSFTP(&CreateSFTPInput{
			ServiceID:        TestDeliveryServiceID,
			ServiceVersion:   *tv.Number,
			Name:             ToPointer("test-sftp-3"),
			Address:          ToPointer("example.com"),
			Port:             ToPointer(1234),
			PublicKey:        ToPointer(pgpPublicKey()),
			SecretKey:        ToPointer(privateKey()),
			SSHKnownHosts:    ToPointer(knownHosts),
			User:             ToPointer("username"),
			Password:         ToPointer("password"),
			Path:             ToPointer("/dir"),
			Period:           ToPointer(12),
			CompressionCodec: ToPointer("snappy"),
			FormatVersion:    ToPointer(2),
			Format:           ToPointer("format"),
			MessageType:      ToPointer("blank"),
			TimestampFormat:  ToPointer("%Y"),
			Placement:        ToPointer("waf_debug"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// This case is expected to fail because both CompressionCodec and
	// GzipLevel are present.
	Record(t, "sftps/create4", func(c *Client) {
		_, err = c.CreateSFTP(&CreateSFTPInput{
			ServiceID:        TestDeliveryServiceID,
			ServiceVersion:   *tv.Number,
			Name:             ToPointer("test-sftp-4"),
			Address:          ToPointer("example.com"),
			Port:             ToPointer(1234),
			PublicKey:        ToPointer(pgpPublicKey()),
			SecretKey:        ToPointer(privateKey()),
			SSHKnownHosts:    ToPointer(knownHosts),
			User:             ToPointer("username"),
			Password:         ToPointer("password"),
			Path:             ToPointer("/dir"),
			Period:           ToPointer(12),
			CompressionCodec: ToPointer("snappy"),
			GzipLevel:        ToPointer(8),
			FormatVersion:    ToPointer(2),
			Format:           ToPointer("format"),
			MessageType:      ToPointer("blank"),
			TimestampFormat:  ToPointer("%Y"),
			Placement:        ToPointer("waf_debug"),
		})
	})
	if err == nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "sftps/cleanup", func(c *Client) {
			_ = c.DeleteSFTP(&DeleteSFTPInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-sftp",
			})

			_ = c.DeleteSFTP(&DeleteSFTPInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-sftp-2",
			})

			_ = c.DeleteSFTP(&DeleteSFTPInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-sftp-3",
			})

			_ = c.DeleteSFTP(&DeleteSFTPInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-sftp",
			})
		})
	}()

	if *sftpCreateResp1.Name != "test-sftp" {
		t.Errorf("bad name: %q", *sftpCreateResp1.Name)
	}
	if *sftpCreateResp1.Address != "example.com" {
		t.Errorf("bad address: %q", *sftpCreateResp1.Address)
	}
	if *sftpCreateResp1.Port != 1234 {
		t.Errorf("bad port: %q", *sftpCreateResp1.Port)
	}
	if *sftpCreateResp1.PublicKey != pgpPublicKey() {
		t.Errorf("bad public_key: %q", *sftpCreateResp1.PublicKey)
	}
	if strings.TrimSpace(*sftpCreateResp1.SecretKey) != strings.TrimSpace(privateKey()) {
		t.Errorf("bad secret_key: %q", *sftpCreateResp1.SecretKey)
	}
	if *sftpCreateResp1.SSHKnownHosts != knownHosts {
		t.Errorf("bad ssh_known_hosts: %q", *sftpCreateResp1.SSHKnownHosts)
	}
	if *sftpCreateResp1.User != "username" {
		t.Errorf("bad user: %q", *sftpCreateResp1.User)
	}
	if *sftpCreateResp1.Password != "password" {
		t.Errorf("bad password: %q", *sftpCreateResp1.Password)
	}
	if *sftpCreateResp1.Path != "/dir" {
		t.Errorf("bad path: %q", *sftpCreateResp1.Path)
	}
	if *sftpCreateResp1.Period != 12 {
		t.Errorf("bad period: %q", *sftpCreateResp1.Period)
	}
	if *sftpCreateResp1.CompressionCodec != "snappy" {
		t.Errorf("bad compression_codec: %q", *sftpCreateResp1.CompressionCodec)
	}
	if *sftpCreateResp1.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", *sftpCreateResp1.GzipLevel)
	}
	if *sftpCreateResp1.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", *sftpCreateResp1.FormatVersion)
	}
	if *sftpCreateResp1.Format != "format" {
		t.Errorf("bad format: %q", *sftpCreateResp1.Format)
	}
	if *sftpCreateResp1.TimestampFormat != "%Y" {
		t.Errorf("bad timestamp_format: %q", *sftpCreateResp1.TimestampFormat)
	}
	if *sftpCreateResp1.MessageType != "blank" {
		t.Errorf("bad message_type: %q", *sftpCreateResp1.MessageType)
	}
	if *sftpCreateResp1.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", *sftpCreateResp1.Placement)
	}
	if sftpCreateResp2.CompressionCodec != nil {
		t.Errorf("bad compression_codec: %q", *sftpCreateResp2.CompressionCodec)
	}
	if *sftpCreateResp2.GzipLevel != 8 {
		t.Errorf("bad gzip_level: %q", *sftpCreateResp2.GzipLevel)
	}
	if *sftpCreateResp3.CompressionCodec != "snappy" {
		t.Errorf("bad compression_codec: %q", *sftpCreateResp3.CompressionCodec)
	}
	if *sftpCreateResp3.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", *sftpCreateResp3.GzipLevel)
	}

	// List
	var sftps []*SFTP
	Record(t, "sftps/list", func(c *Client) {
		sftps, err = c.ListSFTPs(&ListSFTPsInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(sftps) < 1 {
		t.Errorf("bad sftps: %v", sftps)
	}

	// Get
	var sftpGetResp *SFTP
	Record(t, "sftps/get", func(c *Client) {
		sftpGetResp, err = c.GetSFTP(&GetSFTPInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-sftp",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *sftpCreateResp1.Name != *sftpGetResp.Name {
		t.Errorf("bad name: %q", *sftpCreateResp1.Name)
	}
	if *sftpCreateResp1.Address != *sftpGetResp.Address {
		t.Errorf("bad address: %q", *sftpCreateResp1.Address)
	}
	if *sftpCreateResp1.Port != *sftpGetResp.Port {
		t.Errorf("bad port: %q", *sftpCreateResp1.Port)
	}
	if *sftpCreateResp1.PublicKey != *sftpGetResp.PublicKey {
		t.Errorf("bad public_key: %q", *sftpCreateResp1.PublicKey)
	}
	if *sftpCreateResp1.SecretKey != *sftpGetResp.SecretKey {
		t.Errorf("bad secret_key: %q", *sftpCreateResp1.SecretKey)
	}
	if *sftpCreateResp1.SSHKnownHosts != *sftpGetResp.SSHKnownHosts {
		t.Errorf("bad ssh_known_hosts: %q", *sftpCreateResp1.SSHKnownHosts)
	}
	if *sftpCreateResp1.User != *sftpGetResp.User {
		t.Errorf("bad user: %q", *sftpCreateResp1.User)
	}
	if *sftpCreateResp1.Password != *sftpGetResp.Password {
		t.Errorf("bad password: %q", *sftpCreateResp1.Password)
	}
	if *sftpCreateResp1.Path != *sftpGetResp.Path {
		t.Errorf("bad path: %q", *sftpCreateResp1.Path)
	}
	if *sftpCreateResp1.Period != *sftpGetResp.Period {
		t.Errorf("bad period: %q", *sftpCreateResp1.Period)
	}
	if *sftpCreateResp1.CompressionCodec != *sftpGetResp.CompressionCodec {
		t.Errorf("bad compression_codec: %q", *sftpCreateResp1.CompressionCodec)
	}
	if *sftpCreateResp1.GzipLevel != *sftpGetResp.GzipLevel {
		t.Errorf("bad gzip_level: %q", *sftpCreateResp1.GzipLevel)
	}
	if *sftpCreateResp1.FormatVersion != *sftpGetResp.FormatVersion {
		t.Errorf("bad format_version: %q", *sftpCreateResp1.FormatVersion)
	}
	if *sftpCreateResp1.Format != *sftpGetResp.Format {
		t.Errorf("bad format: %q", *sftpCreateResp1.Format)
	}
	if *sftpCreateResp1.TimestampFormat != *sftpGetResp.TimestampFormat {
		t.Errorf("bad timestamp_format: %q", *sftpCreateResp1.TimestampFormat)
	}
	if *sftpCreateResp1.MessageType != "blank" {
		t.Errorf("bad message_type: %q", *sftpCreateResp1.MessageType)
	}
	if *sftpCreateResp1.Placement != *sftpGetResp.Placement {
		t.Errorf("bad placement: %q", *sftpCreateResp1.Placement)
	}

	// Update
	var sftpUpdateResp1, sftpUpdateResp2, sftpUpdateResp3 *SFTP
	Record(t, "sftps/update", func(c *Client) {
		sftpUpdateResp1, err = c.UpdateSFTP(&UpdateSFTPInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-sftp",
			NewName:        ToPointer("new-test-sftp"),
			GzipLevel:      ToPointer(8),
			MessageType:    ToPointer("classic"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	Record(t, "sftps/update2", func(c *Client) {
		sftpUpdateResp2, err = c.UpdateSFTP(&UpdateSFTPInput{
			ServiceID:        TestDeliveryServiceID,
			ServiceVersion:   *tv.Number,
			Name:             "test-sftp-2",
			CompressionCodec: ToPointer("zstd"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	Record(t, "sftps/update3", func(c *Client) {
		sftpUpdateResp3, err = c.UpdateSFTP(&UpdateSFTPInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-sftp-3",
			GzipLevel:      ToPointer(9),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *sftpUpdateResp1.Name != "new-test-sftp" {
		t.Errorf("bad name: %q", *sftpUpdateResp1.Name)
	}
	if *sftpUpdateResp1.MessageType != "classic" {
		t.Errorf("bad message_type: %q", *sftpUpdateResp1.MessageType)
	}
	if sftpUpdateResp1.CompressionCodec != nil {
		t.Errorf("bad compression_codec: %q", *sftpUpdateResp1.CompressionCodec)
	}
	if *sftpUpdateResp1.GzipLevel != 8 {
		t.Errorf("bad gzip_level: %q", *sftpUpdateResp1.GzipLevel)
	}
	if *sftpUpdateResp2.CompressionCodec != "zstd" {
		t.Errorf("bad compression_codec: %q", *sftpUpdateResp2.CompressionCodec)
	}
	if *sftpUpdateResp2.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", *sftpUpdateResp2.GzipLevel)
	}
	if sftpUpdateResp3.CompressionCodec != nil {
		t.Errorf("bad compression_codec: %q", *sftpUpdateResp3.CompressionCodec)
	}
	if *sftpUpdateResp3.GzipLevel != 9 {
		t.Errorf("bad gzip_level: %q", *sftpUpdateResp3.GzipLevel)
	}

	// Delete
	Record(t, "sftps/delete", func(c *Client) {
		err = c.DeleteSFTP(&DeleteSFTPInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-sftp",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListSFTPs_validation(t *testing.T) {
	var err error

	_, err = TestClient.ListSFTPs(&ListSFTPsInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListSFTPs(&ListSFTPsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateSFTP_validation(t *testing.T) {
	var err error

	_, err = TestClient.CreateSFTP(&CreateSFTPInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateSFTP(&CreateSFTPInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetSFTP_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetSFTP(&GetSFTPInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetSFTP(&GetSFTPInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetSFTP(&GetSFTPInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateSFTP_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateSFTP(&UpdateSFTPInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateSFTP(&UpdateSFTPInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateSFTP(&UpdateSFTPInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteSFTP_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteSFTP(&DeleteSFTPInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteSFTP(&DeleteSFTPInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteSFTP(&DeleteSFTPInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
