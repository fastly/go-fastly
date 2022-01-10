package fastly

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"testing"

	rnd "math/rand"
	"strings"
	"time"
)

func TestClient_PrivateKey(t *testing.T) {
	t.Parallel()

	fixtureBase := "tls/"

	// Create
	var err error
	var pk *PrivateKey
	record(t, fixtureBase+"create", func(c *Client) {
		pk, err = c.CreatePrivateKey(&CreatePrivateKeyInput{
			Key:  "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n",
			Name: "My private key",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, fixtureBase+"cleanup", func(c *Client) {
			c.DeletePrivateKey(&DeletePrivateKeyInput{
				ID: pk.ID,
			})
		})
	}()

	// List
	var lpk []*PrivateKey
	record(t, fixtureBase+"list", func(c *Client) {
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
	record(t, fixtureBase+"get", func(c *Client) {
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
	record(t, fixtureBase+"delete", func(c *Client) {
		err = c.DeletePrivateKey(&DeletePrivateKeyInput{
			ID: pk.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListPrivateKeys_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "tls/list", func(c *Client) {
		_, err = c.ListPrivateKeys(&ListPrivateKeysInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetPrivateKey_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "tls/get", func(c *Client) {
		_, err = c.GetPrivateKey(&GetPrivateKeyInput{
			ID: "PRIVATE_KEY_ID",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_CreatePrivateKey_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "tls/create", func(c *Client) {
		_, err = c.CreatePrivateKey(&CreatePrivateKeyInput{
			Key:  "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----\n",
			Name: "My private key",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_DeletePrivateKey_validation(t *testing.T) {
	t.Parallel()

	var err error
	record(t, "tls/delete", func(c *Client) {
		err = c.DeletePrivateKey(&DeletePrivateKeyInput{
			ID: "PRIVATE_KEY_ID",
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
	serialNumber := new(big.Int).SetInt64(rnd.Int63())
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

func formatCertificate(certificate *x509.Certificate, parent *x509.Certificate, publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) (string, error) {
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
