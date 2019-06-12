package fastly

import "testing"

func TestClient_Services(t *testing.T) {
	t.Parallel()

	var err error

	// Create
	var s *Service
	record(t, "services/create", func(c *Client) {
		s, err = c.CreateService(&CreateServiceInput{
			Name:    "test-service",
			Comment: "comment",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "services/cleanup", func(c *Client) {
			c.DeleteService(&DeleteServiceInput{
				ID: s.ID,
			})

			c.DeleteService(&DeleteServiceInput{
				ID: s.ID,
			})
		})
	}()

	if s.Name != "test-service" {
		t.Errorf("bad name: %q", s.Name)
	}
	if s.Comment != "comment" {
		t.Errorf("bad comment: %q", s.Comment)
	}

	// List
	var ss []*Service
	record(t, "services/list", func(c *Client) {
		ss, err = c.ListServices(&ListServicesInput{})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ss) < 1 {
		t.Errorf("bad services: %v", ss)
	}

	// Get
	var ns *Service
	record(t, "services/get", func(c *Client) {
		ns, err = c.GetService(&GetServiceInput{
			ID: s.ID,
		})
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

	if ns.CreatedAt == nil {
		t.Errorf("Bad created at: empty")
	}

	if ns.UpdatedAt == nil {
		t.Errorf("Bad updated at: empty")
	}

	if ns.DeletedAt != nil {
		t.Errorf("Bad deleted at: %s", ns.DeletedAt)
	}

	// Get Details
	var nsd *ServiceDetail
	record(t, "services/details", func(c *Client) {
		nsd, err = c.GetServiceDetails(&GetServiceInput{
			ID: s.ID,
		})
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
	if nsd.Version.Number == 0 {
		t.Errorf("Service Detail Version is empty: (%#v)", nsd)
	}

	// Search
	var fs *Service
	record(t, "services/search", func(c *Client) {
		fs, err = c.SearchService(&SearchServiceInput{
			Name: "test-service",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if s.Name != fs.Name {
		t.Errorf("bad name: %q (%q)", s.Name, fs.Name)
	}
	if s.Comment != fs.Comment {
		t.Errorf("bad comment: %q (%q)", s.Comment, fs.Comment)
	}

	// Update
	var us *Service
	record(t, "services/update", func(c *Client) {
		us, err = c.UpdateService(&UpdateServiceInput{
			ID:   s.ID,
			Name: "new-test-service",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if us.Name != "new-test-service" {
		t.Errorf("bad name: %q", us.Name)
	}

	// Delete
	record(t, "services/delete", func(c *Client) {
		err = c.DeleteService(&DeleteServiceInput{
			ID: s.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	//	List Domains
	var ds ServiceDomainsList
	record(t, "services/domain", func(c *Client) {
		ds, err = c.ListServiceDomains(&ListServiceDomainInput{
			ID: s.ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ds) != 0 {
		t.Fatalf("bad services: %v", ds)
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
