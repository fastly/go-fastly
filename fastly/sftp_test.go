package fastly

import "testing"

func TestClient_SFTPs(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "sftps/version", func(c *Client) {
		tv = testVersion(t, c)
	})
	knownHosts := "example.com ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIEmuNYPrqg9tjqfR14ye3Pvsm9sjIx6EJwm5tMXIMaN1"
	// Create
	var sftp *SFTP
	record(t, "sftps/create", func(c *Client) {
		sftp, err = c.CreateSFTP(&CreateSFTPInput{
			ServiceID:       testServiceID,
			ServiceVersion:  tv.Number,
			Name:            String("test-sftp"),
			Address:         String("example.com"),
			Port:            Uint(1234),
			PublicKey:       String(pgpPublicKey()),
			SecretKey:       String(privateKey()),
			SSHKnownHosts:   String(knownHosts),
			User:            String("username"),
			Password:        String("password"),
			Path:            String("/dir"),
			Period:          Uint(12),
			GzipLevel:       Uint(9),
			FormatVersion:   Uint(2),
			Format:          String("format"),
			MessageType:     String("blank"),
			TimestampFormat: String("%Y"),
			Placement:       String("waf_debug"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "sftps/cleanup", func(c *Client) {
			c.DeleteSFTP(&DeleteSFTPInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-sftp",
			})

			c.DeleteSFTP(&DeleteSFTPInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "new-test-sftp",
			})
		})
	}()

	if sftp.Name != "test-sftp" {
		t.Errorf("bad name: %q", sftp.Name)
	}
	if sftp.Address != "example.com" {
		t.Errorf("bad address: %q", sftp.Address)
	}
	if sftp.Port != 1234 {
		t.Errorf("bad port: %q", sftp.Port)
	}
	if sftp.PublicKey != pgpPublicKey() {
		t.Errorf("bad public_key: %q", sftp.PublicKey)
	}
	if sftp.SecretKey != privateKey() {
		t.Errorf("bad secret_key: %q", sftp.SecretKey)
	}
	if sftp.SSHKnownHosts != knownHosts {
		t.Errorf("bad ssh_known_hosts: %q", sftp.SSHKnownHosts)
	}
	if sftp.User != "username" {
		t.Errorf("bad user: %q", sftp.User)
	}
	if sftp.Password != "password" {
		t.Errorf("bad password: %q", sftp.Password)
	}
	if sftp.Path != "/dir" {
		t.Errorf("bad path: %q", sftp.Path)
	}
	if sftp.Period != 12 {
		t.Errorf("bad period: %q", sftp.Period)
	}
	if sftp.GzipLevel != 9 {
		t.Errorf("bad gzip_level: %q", sftp.GzipLevel)
	}
	if sftp.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", sftp.FormatVersion)
	}
	if sftp.Format != "format" {
		t.Errorf("bad format: %q", sftp.Format)
	}
	if sftp.TimestampFormat != "%Y" {
		t.Errorf("bad timestamp_format: %q", sftp.TimestampFormat)
	}
	if sftp.MessageType != "blank" {
		t.Errorf("bad message_type: %q", sftp.MessageType)
	}
	if sftp.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", sftp.Placement)
	}

	// List
	var sftps []*SFTP
	record(t, "sftps/list", func(c *Client) {
		sftps, err = c.ListSFTPs(&ListSFTPsInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(sftps) < 1 {
		t.Errorf("bad sftps: %v", sftps)
	}

	// Get
	var nsftp *SFTP
	record(t, "sftps/get", func(c *Client) {
		nsftp, err = c.GetSFTP(&GetSFTPInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-sftp",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if sftp.Name != nsftp.Name {
		t.Errorf("bad name: %q", sftp.Name)
	}
	if sftp.Address != nsftp.Address {
		t.Errorf("bad address: %q", sftp.Address)
	}
	if sftp.Port != nsftp.Port {
		t.Errorf("bad port: %q", sftp.Port)
	}
	if sftp.PublicKey != nsftp.PublicKey {
		t.Errorf("bad public_key: %q", sftp.PublicKey)
	}
	if sftp.SecretKey != nsftp.SecretKey {
		t.Errorf("bad secret_key: %q", sftp.SecretKey)
	}
	if sftp.SSHKnownHosts != nsftp.SSHKnownHosts {
		t.Errorf("bad ssh_known_hosts: %q", sftp.SSHKnownHosts)
	}
	if sftp.User != nsftp.User {
		t.Errorf("bad user: %q", sftp.User)
	}
	if sftp.Password != nsftp.Password {
		t.Errorf("bad password: %q", sftp.Password)
	}
	if sftp.Path != nsftp.Path {
		t.Errorf("bad path: %q", sftp.Path)
	}
	if sftp.Period != nsftp.Period {
		t.Errorf("bad period: %q", sftp.Period)
	}
	if sftp.GzipLevel != nsftp.GzipLevel {
		t.Errorf("bad gzip_level: %q", sftp.GzipLevel)
	}
	if sftp.FormatVersion != nsftp.FormatVersion {
		t.Errorf("bad format_version: %q", sftp.FormatVersion)
	}
	if sftp.Format != nsftp.Format {
		t.Errorf("bad format: %q", sftp.Format)
	}
	if sftp.TimestampFormat != nsftp.TimestampFormat {
		t.Errorf("bad timestamp_format: %q", sftp.TimestampFormat)
	}
	if sftp.MessageType != "blank" {
		t.Errorf("bad message_type: %q", sftp.MessageType)
	}
	if sftp.Placement != nsftp.Placement {
		t.Errorf("bad placement: %q", sftp.Placement)
	}

	// Update
	var usftp *SFTP
	record(t, "sftps/update", func(c *Client) {
		usftp, err = c.UpdateSFTP(&UpdateSFTPInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-sftp",
			NewName:        String("new-test-sftp"),
			GzipLevel:      Uint(0),
			MessageType:    String("classic"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if usftp.Name != "new-test-sftp" {
		t.Errorf("bad name: %q", usftp.Name)
	}
	if usftp.MessageType != "classic" {
		t.Errorf("bad message_type: %q", usftp.MessageType)
	}
	if usftp.GzipLevel != 0 {
		t.Errorf("bad gzip_level: %q", usftp.GzipLevel)
	}

	// Delete
	record(t, "sftps/delete", func(c *Client) {
		err = c.DeleteSFTP(&DeleteSFTPInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "new-test-sftp",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListSFTPs_validation(t *testing.T) {
	var err error
	_, err = testClient.ListSFTPs(&ListSFTPsInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListSFTPs(&ListSFTPsInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateSFTP_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateSFTP(&CreateSFTPInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateSFTP(&CreateSFTPInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetSFTP_validation(t *testing.T) {
	var err error
	_, err = testClient.GetSFTP(&GetSFTPInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetSFTP(&GetSFTPInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetSFTP(&GetSFTPInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateSFTP_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateSFTP(&UpdateSFTPInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateSFTP(&UpdateSFTPInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateSFTP(&UpdateSFTPInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteSFTP_validation(t *testing.T) {
	var err error
	err = testClient.DeleteSFTP(&DeleteSFTPInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteSFTP(&DeleteSFTPInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteSFTP(&DeleteSFTPInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
