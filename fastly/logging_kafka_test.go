package fastly

import (
	"errors"
	"strings"
	"testing"
)

const requestMaxBytes = 2048

func TestClient_Kafkas(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	Record(t, "kafkas/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	caCert := strings.TrimSpace(caCert())
	clientCert := strings.TrimSpace(certificate())
	clientKey := strings.TrimSpace(privateKey())

	// Create
	var k *Kafka
	Record(t, "kafkas/create", func(c *Client) {
		k, err = c.CreateKafka(&CreateKafkaInput{
			AuthMethod:       ToPointer("scram-sha-512"),
			Brokers:          ToPointer("192.168.1.1,192.168.1.2"),
			CompressionCodec: ToPointer("lz4"),
			Format:           ToPointer("%h %l %u %t \"%r\" %>s %b"),
			FormatVersion:    ToPointer(2),
			Name:             ToPointer("test-kafka"),
			ParseLogKeyvals:  ToPointer(Compatibool(true)),
			Password:         ToPointer("deadbeef"),
			Placement:        ToPointer("waf_debug"),
			RequestMaxBytes:  ToPointer(requestMaxBytes),
			RequiredACKs:     ToPointer("-1"),
			ServiceID:        TestDeliveryServiceID,
			ServiceVersion:   *tv.Number,
			TLSCACert:        ToPointer(caCert),
			TLSClientCert:    ToPointer(clientCert),
			TLSClientKey:     ToPointer(clientKey),
			TLSHostname:      ToPointer("example.com"),
			Topic:            ToPointer("kafka-topic"),
			UseTLS:           ToPointer(Compatibool(true)),
			User:             ToPointer("foobar"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "kafkas/cleanup", func(c *Client) {
			_ = c.DeleteKafka(&DeleteKafkaInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "test-kafka",
			})

			_ = c.DeleteKafka(&DeleteKafkaInput{
				ServiceID:      TestDeliveryServiceID,
				ServiceVersion: *tv.Number,
				Name:           "new-test-kafka",
			})
		})
	}()

	if *k.Name != "test-kafka" {
		t.Errorf("bad name: %q", *k.Name)
	}
	if *k.Brokers != "192.168.1.1,192.168.1.2" {
		t.Errorf("bad url: %q", *k.Brokers)
	}
	if *k.Topic != "kafka-topic" {
		t.Errorf("bad topic: %q", *k.Topic)
	}
	if *k.RequiredACKs != "-1" {
		t.Errorf("bad required_acks: %q", *k.RequiredACKs)
	}
	if !*k.UseTLS {
		t.Errorf("bad use_tls: %t", *k.UseTLS)
	}
	if *k.CompressionCodec != "lz4" {
		t.Errorf("bad compression_codec: %q", *k.CompressionCodec)
	}
	if *k.Format != "%h %l %u %t \"%r\" %>s %b" {
		t.Errorf("bad format: %q", *k.Format)
	}
	if *k.FormatVersion != 2 {
		t.Errorf("bad format_version: %q", *k.FormatVersion)
	}
	if *k.Placement != "waf_debug" {
		t.Errorf("bad placement: %q", *k.Placement)
	}
	if *k.TLSCACert != caCert {
		t.Errorf("bad tls_ca_cert: %q", *k.TLSCACert)
	}
	if *k.TLSHostname != "example.com" {
		t.Errorf("bad tls_hostname: %q", *k.TLSHostname)
	}
	if *k.TLSClientCert != clientCert {
		t.Errorf("bad tls_client_cert: %q", *k.TLSClientCert)
	}
	if *k.TLSClientKey != clientKey {
		t.Errorf("bad tls_client_key: %q", *k.TLSClientKey)
	}
	if !*k.ParseLogKeyvals {
		t.Errorf("bad parse_log_keyvals: %t", *k.ParseLogKeyvals)
	}
	if *k.RequestMaxBytes != requestMaxBytes {
		t.Errorf("bad request_max_bytes: %q", *k.RequestMaxBytes)
	}
	if *k.AuthMethod != "scram-sha-512" {
		t.Errorf("bad auth_method: %q", *k.AuthMethod)
	}
	if *k.User != "foobar" {
		t.Errorf("bad user: %q", *k.User)
	}
	if *k.Password != "deadbeef" {
		t.Errorf("bad password: %q", *k.Password)
	}

	// List
	var ks []*Kafka
	Record(t, "kafkas/list", func(c *Client) {
		ks, err = c.ListKafkas(&ListKafkasInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
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
	Record(t, "kafkas/get", func(c *Client) {
		nk, err = c.GetKafka(&GetKafkaInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-kafka",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *k.Name != *nk.Name {
		t.Errorf("bad name: %q", *k.Name)
	}
	if *k.Brokers != *nk.Brokers {
		t.Errorf("bad url: %q", *k.Brokers)
	}
	if *k.Topic != *nk.Topic {
		t.Errorf("bad topic: %q", *k.Topic)
	}
	if *k.RequiredACKs != *nk.RequiredACKs {
		t.Errorf("bad required_acks: %q", *k.RequiredACKs)
	}
	if *k.UseTLS != *nk.UseTLS {
		t.Errorf("bad use_tls: %t", *k.UseTLS)
	}
	if *k.CompressionCodec != *nk.CompressionCodec {
		t.Errorf("bad compression_codec: %q", *k.CompressionCodec)
	}
	if *k.Format != *nk.Format {
		t.Errorf("bad format: %q", *k.Format)
	}
	if *k.FormatVersion != *nk.FormatVersion {
		t.Errorf("bad format_version: %q", *k.FormatVersion)
	}
	if *k.Placement != *nk.Placement {
		t.Errorf("bad placement: %q", *k.Placement)
	}
	if *k.TLSCACert != *nk.TLSCACert {
		t.Errorf("bad tls_ca_cert: %q", *k.TLSCACert)
	}
	if *k.TLSHostname != *nk.TLSHostname {
		t.Errorf("bad tls_hostname: %q", *k.TLSHostname)
	}
	if *k.TLSClientCert != *nk.TLSClientCert {
		t.Errorf("bad tls_client_cert: %q", *k.TLSClientCert)
	}
	if *k.TLSClientKey != *nk.TLSClientKey {
		t.Errorf("bad tls_client_key: %q", *k.TLSClientKey)
	}
	if !*k.ParseLogKeyvals {
		t.Errorf("bad parse_log_keyvals: %t", *k.ParseLogKeyvals)
	}
	if *k.RequestMaxBytes != requestMaxBytes {
		t.Errorf("bad request_max_bytes: %q", *k.RequestMaxBytes)
	}
	if *k.AuthMethod != "scram-sha-512" {
		t.Errorf("bad auth_method: %q", *k.AuthMethod)
	}
	if *k.User != "foobar" {
		t.Errorf("bad user: %q", *k.User)
	}
	if *k.Password != "deadbeef" {
		t.Errorf("bad password: %q", *k.Password)
	}

	// Update
	var uk *Kafka
	Record(t, "kafkas/update", func(c *Client) {
		uk, err = c.UpdateKafka(&UpdateKafkaInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "test-kafka",
			NewName:        ToPointer("new-test-kafka"),
			Topic:          ToPointer("new-kafka-topic"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *uk.Name != "new-test-kafka" {
		t.Errorf("bad name: %q", *uk.Name)
	}
	if *uk.Topic != "new-kafka-topic" {
		t.Errorf("bad topic: %q", *uk.Topic)
	}

	// Delete
	Record(t, "kafkas/delete", func(c *Client) {
		err = c.DeleteKafka(&DeleteKafkaInput{
			ServiceID:      TestDeliveryServiceID,
			ServiceVersion: *tv.Number,
			Name:           "new-test-kafka",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListKafkas_validation(t *testing.T) {
	var err error
	_, err = TestClient.ListKafkas(&ListKafkasInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.ListKafkas(&ListKafkasInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateKafka_validation(t *testing.T) {
	var err error
	_, err = TestClient.CreateKafka(&CreateKafkaInput{
		ServiceID: "",
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.CreateKafka(&CreateKafkaInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetKafka_validation(t *testing.T) {
	var err error

	_, err = TestClient.GetKafka(&GetKafkaInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetKafka(&GetKafkaInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.GetKafka(&GetKafkaInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateKafka_validation(t *testing.T) {
	var err error

	_, err = TestClient.UpdateKafka(&UpdateKafkaInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateKafka(&UpdateKafkaInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	_, err = TestClient.UpdateKafka(&UpdateKafkaInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteKafka_validation(t *testing.T) {
	var err error

	err = TestClient.DeleteKafka(&DeleteKafkaInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingName) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteKafka(&DeleteKafkaInput{
		Name:           "test",
		ServiceVersion: 1,
	})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}

	err = TestClient.DeleteKafka(&DeleteKafkaInput{
		Name:      "test",
		ServiceID: "foo",
	})
	if !errors.Is(err, ErrMissingServiceVersion) {
		t.Errorf("bad error: %s", err)
	}
}
