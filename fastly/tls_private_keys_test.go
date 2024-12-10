package fastly

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"
)

func TestClient_PrivateKey(t *testing.T) {
	t.Parallel()

	fixtureBase := "tls/"

	// prepare test key
	_, key, err := buildPrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	// Create
	var pk *PrivateKey
	Record(t, fixtureBase+"create", func(c *Client) {
		pk, err = c.CreatePrivateKey(&CreatePrivateKeyInput{
			Key:  key,
			Name: "My private key",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		_ = TestClient.DeletePrivateKey(&DeletePrivateKeyInput{
			ID: pk.ID,
		})
	}()

	// List
	var lpk []*PrivateKey
	Record(t, fixtureBase+"list", func(c *Client) {
		lpk, err = c.ListPrivateKeys(&ListPrivateKeysInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(lpk) < 1 {
		t.Errorf("bad privatekeys: %v", lpk)
	}

	// Get
	var gpk *PrivateKey
	Record(t, fixtureBase+"get", func(c *Client) {
		gpk, err = c.GetPrivateKey(&GetPrivateKeyInput{
			ID: pk.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if pk.Name != gpk.Name {
		t.Errorf("bad name: %q (%q)", pk.Name, gpk.Name)
	}

	// Delete
	Record(t, fixtureBase+"delete", func(c *Client) {
		err = c.DeletePrivateKey(&DeletePrivateKeyInput{
			ID: pk.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func buildPrivateKey() (*rsa.PrivateKey, string, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate private key: %s", err)
	}
	privKeyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return nil, "", fmt.Errorf("unable to marshal private key: %v", err)
	}
	keyBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privKeyBytes,
	}
	return key, strings.TrimSpace(string(pem.EncodeToMemory(keyBlock))), nil
}

func buildCertificate(privateKey *rsa.PrivateKey, domains ...string) (string, error) {
	now := time.Now()

	n, err := rand.Int(rand.Reader, big.NewInt(1000))
	if err != nil {
		return "", err
	}
	serialNumber := new(big.Int).SetInt64(n.Int64())
	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			SerialNumber: fmt.Sprintf("%d", serialNumber),
		},
		NotBefore:             now,
		NotAfter:              now.Add(24 * 90 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		SignatureAlgorithm:    x509.SHA256WithRSA,
		IsCA:                  true,
		DNSNames:              domains,
	}

	if len(domains) != 0 {
		template.Subject.CommonName = domains[0]
	}

	c, err := formatCertificate(template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", err
	}
	return c, nil
}

func formatCertificate(certificate, parent *x509.Certificate, publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) (string, error) {
	derBytes, err := x509.CreateCertificate(
		rand.Reader,
		certificate,
		parent,
		publicKey,
		privateKey,
	)
	if err != nil {
		return "", err
	}
	certBlock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	}
	return strings.TrimSpace(string(pem.EncodeToMemory(certBlock))), nil
}
