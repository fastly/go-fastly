package accesskeys

import (
	"testing"

	"github.com/fastly/go-fastly/v9/fastly"
)

func TestClient_AccessKey(t *testing.T) {
	t.Parallel()

	TestAccessKeyDescription := "THIS IS A TEST ACCESS KEY"
	TestAccessKeyPermission := "read-write-objects"
	TestAccessKeyBuckets := []string{"test-bucket"}

	var accessKeys *AccessKeys
	var err error

	// List all AccessKeys.
	fastly.Record(t, "list", func(c *fastly.Client) {
		accessKeys, err = ListAccessKeys(c)
	})
	if err != nil {
		t.Fatal(err)
	}

	// Make sure the test AccessKey we're going to create isn't among them.
	for _, ak := range accessKeys.Data {
		if ak.Description == TestAccessKeyDescription &&
			ak.Permission == TestAccessKeyPermission &&
			len(ak.Buckets) == len(TestAccessKeyBuckets) &&
			ak.Buckets[0] == TestAccessKeyBuckets[0] {
			t.Errorf("found test AccessKey %q, aborting", ak.AccessKeyID)
		}
	}

	// Create a AccessKey for testing.
	var accessKey *AccessKey
	fastly.Record(t, "create", func(c *fastly.Client) {
		accessKey, err = Create(c, &CreateInput{
			Description: fastly.ToPointer(TestAccessKeyDescription),
			Permission:  fastly.ToPointer(TestAccessKeyPermission),
			Buckets:     fastly.ToPointer(TestAccessKeyBuckets),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if accessKey.AccessKeyID == "" {
		t.Errorf("no access key ID was returned")
	}
	if accessKey.SecretKey == "" {
		t.Errorf("no access key secret was returned")
	}
	if accessKey.Description != TestAccessKeyDescription {
		t.Errorf("unexpected AccessKey name: got %q, expected %q", accessKey.Description, TestAccessKeyDescription)
	}
	if accessKey.Permission != TestAccessKeyPermission {
		t.Errorf("unexpected AccessKey permission: got %q, expected %q", accessKey.Permission, TestAccessKeyPermission)
	}
	if len(accessKey.Buckets) != len(TestAccessKeyBuckets) {
		t.Errorf("unexpected AccessKey buckets length: got %q, expected %q", len(accessKey.Buckets), len(TestAccessKeyBuckets))
	}
	if accessKey.Buckets[0] != TestAccessKeyBuckets[0] {
		t.Errorf("unexpected AccessKey bucket: got %q, expected %q", accessKey.Buckets[0], TestAccessKeyBuckets[0])
	}

	// Ensure we delete the test AccessKey at the end.
	defer func() {
		fastly.Record(t, "delete", func(c *fastly.Client) {
			err = Delete(c, &DeleteInput{
				AccessKeyID: fastly.ToPointer(accessKey.AccessKeyID),
			})
		})
		if err != nil {
			t.Errorf("error during AccessKey cleanup: %v", err)
		}
	}()

	// Get the test AccessKey.
	var ak *AccessKey
	fastly.Record(t, "get", func(c *fastly.Client) {
		ak, err = Get(c, &GetInput{
			AccessKeyID: fastly.ToPointer(accessKey.AccessKeyID),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if ak.Description != accessKey.Description {
		t.Errorf("unexpected AccessKey Description: got %q, expected %q", ak.Description, accessKey.Description)
	}
	if ak.Permission != accessKey.Permission {
		t.Errorf("unexpected AccessKey Permissions: got %q, expected %q", ak.Permission, accessKey.Permission)
	}
	if len(ak.Buckets) != len(accessKey.Buckets) {
		t.Errorf("unexpected AccessKey Buckets length: got %q, expected %q", len(ak.Buckets), len(accessKey.Buckets))
	}
	if ak.Buckets[0] != accessKey.Buckets[0] {
		t.Errorf("unexpected AccessKey Buckets contents: got %q, expected %q", ak.Buckets[0], accessKey.Buckets[0])
	}

	// List all entries of the test AccessKey and compare it to the input.
	var actualAccessKeys *AccessKeys
	fastly.Record(t, "list_with_new", func(c *fastly.Client) {
		actualAccessKeys, err = ListAccessKeys(c)
	})
	if err != nil {
		t.Errorf("error fetching list of AccessKey entries: %v", err)
	}

	actualNumberOfAccessKeyEntries := len(actualAccessKeys.Data)
	expectedNumberOfAccessKeyEntries := len(accessKeys.Data)
	// This checks the original number of access keys fetched in the creation check vs the number fetched after adding the test key
	if actualNumberOfAccessKeyEntries != expectedNumberOfAccessKeyEntries+1 {
		t.Errorf("incorrect number of AccessKeys returned, expected: %d, got %d", expectedNumberOfAccessKeyEntries, actualNumberOfAccessKeyEntries)
	}

	// Make sure the test AccessKey we've created is among them.
	var newKeyPresent = false
	for _, rak := range actualAccessKeys.Data {
		if rak.Description == TestAccessKeyDescription &&
			rak.Permission == TestAccessKeyPermission &&
			len(rak.Buckets) == len(TestAccessKeyBuckets) &&
			rak.Buckets[0] == TestAccessKeyBuckets[0] {
			newKeyPresent = true
		}
	}
	if !newKeyPresent {
		t.Errorf("missing test AccessKey %q, aborting", accessKey.AccessKeyID)
	}
}

func TestClient_Create_validation(t *testing.T) {
	_, err := Create(fastly.TestClient, &CreateInput{
		Description: nil,
	})
	if err != fastly.ErrMissingDescription {
		t.Errorf("expected ErrMissingDescription: got %s", err)
	}

	_, err = Create(fastly.TestClient, &CreateInput{
		Description: fastly.ToPointer("description"),
		Permission:  nil,
	})
	if err != fastly.ErrMissingPermission {
		t.Errorf("expected ErrMissingPermission: got %s", err)
	}

	_, err = Create(fastly.TestClient, &CreateInput{
		Description: fastly.ToPointer("description"),
		Permission:  fastly.ToPointer("bad-permission"),
	})
	if err != fastly.ErrInvalidPermission {
		t.Errorf("expected ErrInvalidPermission: got %s", err)
	}
}

func TestClient_Get_validation(t *testing.T) {
	_, err := Get(fastly.TestClient, &GetInput{
		AccessKeyID: nil,
	})
	if err != fastly.ErrMissingAccessKeyID {
		t.Errorf("expected ErrMissingAccessKeyID: got %s", err)
	}
}

func TestClient_Delete_validation(t *testing.T) {
	err := Delete(fastly.TestClient, &DeleteInput{
		AccessKeyID: nil,
	})
	if err != fastly.ErrMissingAccessKeyID {
		t.Errorf("expected ErrMissingAccessKeyID: got %s", err)
	}
}
