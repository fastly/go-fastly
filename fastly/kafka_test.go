package fastly

import (
	"strings"
	"testing"
)

func TestClient_Kafkas(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "kafkas/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	caCert := strings.TrimSpace(caCert())
	clientCert := strings.TrimSpace(certificate())
	clientKey := strings.TrimSpace(privateKey())

	// Create
	var k *Kafka
	record(t, "kafkas/create", func(c *Client) {
		k, err = c.CreateKafka(&CreateKafkaInput{
			Service:          testServiceID,
			Version:          tv.Number,
			Name:             String("test-kafka"),
			Brokers:          String("192.168.1.1,192.168.1.2"),
			Topic:            String("kafka-topic"),
			RequiredACKs:     String("-1"),
			UseTLS:           CBool(true),
			CompressionCodec: String("lz4"),
			Format:           String("%h %l %u %t \"%r\" %>s %b"),
			FormatVersion:    Uint(2),
			Placement:        String("waf_debug"),
			TLSCACert:        String(caCert),
			TLSHostname:      String("example.com"),
			TLSClientCert:    String(clientCert),
			TLSClientKey:     String(clientKey),
			ParseLogKeyvals:  CBool(true),
			RequestMaxBytes:  Uint(1024),
			AuthMethod:       String("scram-sha-512"),
			User:             String("foobar"),
			Password:         String("deadbeef"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "kafkas/cleanup", func(c *Client) {
			c.DeleteKafka(&DeleteKafkaInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "test-kafka",
			})

			c.DeleteKafka(&DeleteKafkaInput{
				Service: testServiceID,
				Version: tv.Number,
				Name:    "new-test-kafka",
			})
		})
	}()

	if k.Name != "test-kafka" {
		t.Errorf("bad name: %q", k.Name)
	}
	if k.Brokers != "192.168.1.1,192.168.1.2" {
		t.Errorf("bad url: %q", k.Brokers)
	}
	if k.Topic != "kafka-topic" {
		t.Errorf("bad topic: %q", k.Topic)
	}
	if k.RequiredACKs != "-1" {
		t.Errorf("bad required_acks: %q", k.RequiredACKs)
	}
	if k.UseTLS != true {
		t.Errorf("bad use_tls: %t", k.UseTLS)
	}
	if k.CompressionCodec != "lz4" {
		t.Errorf("bad compression_codec: %q", k.CompressionCodec)
	}
	if k.Format != "%h %l %u %t \"%r\" %>s %b" {
		t.Errorf("bad format: %q", k.Format)
	}
	if k.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", k.FormatVersion)
	}
	if k.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", k.Placement)
	}
	if k.TLSCACert != caCert {
		t.Errorf("bad tls_ca_cert: %q", k.TLSCACert)
	}
	if k.TLSHostname != "example.com" {
		t.Errorf("bad tls_hostname: %q", k.TLSHostname)
	}
	if k.TLSClientCert != clientCert {
		t.Errorf("bad tls_client_cert: %q", k.TLSClientCert)
	}
	if k.TLSClientKey != clientKey {
		t.Errorf("bad tls_client_key: %q", k.TLSClientKey)
	}
	if k.ParseLogKeyvals != true {
		t.Errorf("bad parse_log_keyvals: %t", k.ParseLogKeyvals)
	}
	if k.RequestMaxBytes != 1024 {
		t.Errorf("bad request_max_bytes: %q", k.RequestMaxBytes)
	}
	if k.AuthMethod != "scram-sha-512" {
		t.Errorf("bad auth_method: %q", k.AuthMethod)
	}
	if k.User != "foobar" {
		t.Errorf("bad user: %q", k.User)
	}
	if k.Password != "deadbeef" {
		t.Errorf("bad password: %q", k.Password)
	}

	// List
	var ks []*Kafka
	record(t, "kafkas/list", func(c *Client) {
		ks, err = c.ListKafkas(&ListKafkasInput{
			Service: testServiceID,
			Version: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ks) < 1 {
		t.Errorf("bad kafkas: %v", ks)
	}

	// Get
	var nk *Kafka
	record(t, "kafkas/get", func(c *Client) {
		nk, err = c.GetKafka(&GetKafkaInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-kafka",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if k.Name != nk.Name {
		t.Errorf("bad name: %q", k.Name)
	}
	if k.Brokers != nk.Brokers {
		t.Errorf("bad url: %q", k.Brokers)
	}
	if k.Topic != nk.Topic {
		t.Errorf("bad topic: %q", k.Topic)
	}
	if k.RequiredACKs != nk.RequiredACKs {
		t.Errorf("bad required_acks: %q", k.RequiredACKs)
	}
	if k.UseTLS != nk.UseTLS {
		t.Errorf("bad use_tls: %t", k.UseTLS)
	}
	if k.CompressionCodec != nk.CompressionCodec {
		t.Errorf("bad compression_codec: %q", k.CompressionCodec)
	}
	if k.Format != nk.Format {
		t.Errorf("bad format: %q", k.Format)
	}
	if k.FormatVersion != nk.FormatVersion {
		t.Errorf("bad format_version: %q", k.FormatVersion)
	}
	if k.Placement != nk.Placement {
		t.Errorf("bad placement: %q", k.Placement)
	}
	if k.TLSCACert != nk.TLSCACert {
		t.Errorf("bad tls_ca_cert: %q", k.TLSCACert)
	}
	if k.TLSHostname != nk.TLSHostname {
		t.Errorf("bad tls_hostname: %q", k.TLSHostname)
	}
	if k.TLSClientCert != nk.TLSClientCert {
		t.Errorf("bad tls_client_cert: %q", k.TLSClientCert)
	}
	if k.TLSClientKey != nk.TLSClientKey {
		t.Errorf("bad tls_client_key: %q", k.TLSClientKey)
	}
	if k.ParseLogKeyvals != true {
		t.Errorf("bad parse_log_keyvals: %t", k.ParseLogKeyvals)
	}
	if k.RequestMaxBytes != 1024 {
		t.Errorf("bad request_max_bytes: %q", k.RequestMaxBytes)
	}
	if k.AuthMethod != "scram-sha-512" {
		t.Errorf("bad auth_method: %q", k.AuthMethod)
	}
	if k.User != "foobar" {
		t.Errorf("bad user: %q", k.User)
	}
	if k.Password != "deadbeef" {
		t.Errorf("bad password: %q", k.Password)
	}

	// Update
	var uk *Kafka
	record(t, "kafkas/update", func(c *Client) {
		uk, err = c.UpdateKafka(&UpdateKafkaInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "test-kafka",
			NewName: String("new-test-kafka"),
			Topic:   String("new-kafka-topic"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if uk.Name != "new-test-kafka" {
		t.Errorf("bad name: %q", uk.Name)
	}
	if uk.Topic != "new-kafka-topic" {
		t.Errorf("bad topic: %q", uk.Topic)
	}

	// Delete
	record(t, "kafkas/delete", func(c *Client) {
		err = c.DeleteKafka(&DeleteKafkaInput{
			Service: testServiceID,
			Version: tv.Number,
			Name:    "new-test-kafka",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListKafkas_validation(t *testing.T) {
	var err error
	_, err = testClient.ListKafkas(&ListKafkasInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListKafkas(&ListKafkasInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateKafka_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateKafka(&CreateKafkaInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateKafka(&CreateKafkaInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetKafka_validation(t *testing.T) {
	var err error
	_, err = testClient.GetKafka(&GetKafkaInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetKafka(&GetKafkaInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetKafka(&GetKafkaInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateKafka_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateKafka(&UpdateKafkaInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateKafka(&UpdateKafkaInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateKafka(&UpdateKafkaInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteKafka_validation(t *testing.T) {
	var err error
	err = testClient.DeleteKafka(&DeleteKafkaInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteKafka(&DeleteKafkaInput{
		Service: "foo",
		Version: 0,
	})
	if err != ErrMissingVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteKafka(&DeleteKafkaInput{
		Service: "foo",
		Version: 1,
		Name:    "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
