package fastly

import "testing"

func TestClient_Services(t *testing.T) {
	t.Parallel()

	// Create
	s, err := testClient.CreateService(&CreateServiceInput{
		Name:    "test-service",
		Comment: "comment",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		testClient.DeleteService(&DeleteServiceInput{
			ID: s.ID,
		})

		testClient.DeleteService(&DeleteServiceInput{
			ID: s.ID,
		})
	}()

	if s.Name != "test-service" {
		t.Errorf("bad name: %q", s.Name)
	}
	if s.Comment != "comment" {
		t.Errorf("bad comment: %q", s.Comment)
	}

	// List
	bs, err := testClient.ListServices(&ListServicesInput{})
	if err != nil {
		t.Fatal(err)
	}
	if len(bs) < 1 {
		t.Errorf("bad services: %v", bs)
	}

	// Get
	ns, err := testClient.GetService(&GetServiceInput{
		ID: s.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if s.Name != ns.Name {
		t.Errorf("bad name: %q (%q)", s.Name, ns.Name)
	}
	if s.Comment != ns.Comment {
		t.Errorf("bad comment: %q (%q)", s.Comment, ns.Comment)
	}

	// Get Details
	nsd, err := testClient.GetServiceDetails(&GetServiceInput{
		ID: s.ID,
	})
	if err != nil {
		t.Fatal(err)
	}

	if s.Name != nsd.Name {
		t.Errorf("bad name: %q (%q)", s.Name, nsd.Name)
	}
	if s.Comment != nsd.Comment {
		t.Errorf("bad comment: %q (%q)", s.Comment, nsd.Comment)
	}
	if nsd.Version.Number == "" {
		t.Errorf("Service Detail Version is empty: (%#v)", nsd)
	}

	// Search
	nb, err := testClient.SearchService(&SearchServiceInput{
		Name: "test-service",
	})
	if err != nil {
		t.Fatal(err)
	}
	if s.Name != nb.Name {
		t.Errorf("bad name: %q (%q)", s.Name, nb.Name)
	}
	if s.Comment != nb.Comment {
		t.Errorf("bad comment: %q (%q)", s.Comment, nb.Comment)
	}

	// Update
	us, err := testClient.UpdateService(&UpdateServiceInput{
		ID:   s.ID,
		Name: "new-test-service",
	})
	if err != nil {
		t.Fatal(err)
	}
	if us.Name != "new-test-service" {
		t.Errorf("bad name: %q", us.Name)
	}

	// Delete
	if err := testClient.DeleteService(&DeleteServiceInput{
		ID: s.ID,
	}); err != nil {
		t.Fatal(err)
	}
}

func TestClient_GetService_validation(t *testing.T) {
	var err error
	_, err = testClient.GetService(&GetServiceInput{})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateService_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateService(&UpdateServiceInput{})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteService_validation(t *testing.T) {
	var err error
	err = testClient.DeleteService(&DeleteServiceInput{})
	if err != ErrMissingID {
		t.Errorf("bad error: %s", err)
	}
}
