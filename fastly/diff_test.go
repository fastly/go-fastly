package fastly

import "testing"

func TestClient_Diff(t *testing.T) {
	t.Parallel()

	tv1 := testVersion(t)
	tv2 := testVersion(t)

	// Diff should be empty
	d, err := testClient.GetDiff(&GetDiffInput{
		Service: testServiceID,
		From:    tv1.Number,
		To:      tv2.Number,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a diff
	_, err = testClient.CreateBackend(&CreateBackendInput{
		Service: testServiceID,
		Version: tv2.Number,
		Name:    "test-backend",
		Address: "integ-test.hashicorp.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure we delete the backend we just created
	defer func() {
		testClient.DeleteBackend(&DeleteBackendInput{
			Service: testServiceID,
			Version: tv2.Number,
			Name:    "test-backend",
		})
	}()

	// Diff should mot be empty
	d, err = testClient.GetDiff(&GetDiffInput{
		Service: testServiceID,
		From:    tv1.Number,
		To:      tv2.Number,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(d.Diff) == 0 {
		t.Errorf("bad diff: %s", d.Diff)
	}
}

func TestClient_Diff_validation(t *testing.T) {
	var err error
	_, err = testClient.GetDiff(&GetDiffInput{
		Service: "",
	})
	if err != ErrMissingService {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDiff(&GetDiffInput{
		Service: "foo",
		From:    "",
	})
	if err != ErrMissingFrom {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetDiff(&GetDiffInput{
		Service: "foo",
		From:    "1",
		To:      "",
	})
	if err != ErrMissingTo {
		t.Errorf("bad error: %s", err)
	}
}
