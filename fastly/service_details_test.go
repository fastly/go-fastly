package fastly

import (
	"context"
	"errors"
	"testing"
)

func TestClient_Services(t *testing.T) {
	t.Parallel()

	var err error

	// Create
	var s *Service
	Record(t, "services/create", func(c *Client) {
		s, err = c.CreateService(context.TODO(), &CreateServiceInput{
			Name:    ToPointer("test-service"),
			Comment: ToPointer("comment"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "services/cleanup", func(c *Client) {
			_ = c.DeleteService(context.TODO(), &DeleteServiceInput{
				ServiceID: *s.ServiceID,
			})

			_ = c.DeleteService(context.TODO(), &DeleteServiceInput{
				ServiceID: *s.ServiceID,
			})
		})
	}()

	if *s.Name != "test-service" {
		t.Errorf("bad name: %q", *s.Name)
	}
	if *s.Comment != "comment" {
		t.Errorf("bad comment: %q", *s.Comment)
	}

	// List
	var ss []*Service
	Record(t, "services/list", func(c *Client) {
		ss, err = c.ListServices(context.TODO(), &ListServicesInput{
			Direction: ToPointer("descend"),
			Sort:      ToPointer("created"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ss) < 1 {
		t.Errorf("bad services: %v", ss)
	}

	// List with paginator
	var ss2 []*Service
	var paginator *ListPaginator[Service]
	Record(t, "services/list_paginator", func(c *Client) {
		paginator = c.GetServices(context.TODO(), &GetServicesInput{
			Direction: ToPointer("descend"),
			PerPage:   ToPointer(200),
			Sort:      ToPointer("created"),
		})

		for paginator.HasNext() {
			data, err := paginator.GetNext()
			if err != nil {
				t.Errorf("Bad paginator (remaining: %d): %s", paginator.Remaining(), err)
				return
			}
			ss2 = append(ss2, data...)
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ss2) != len(ss) {
		t.Errorf("expected %d services but got: %d", len(ss), len(ss2))
	}

	// Get
	var ns *Service
	Record(t, "services/get", func(c *Client) {
		ns, err = c.GetService(context.TODO(), &GetServiceInput{
			ServiceID: *s.ServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *s.Name != *ns.Name {
		t.Errorf("bad name: %q (%q)", *s.Name, *ns.Name)
	}
	if *s.Comment != *ns.Comment {
		t.Errorf("bad comment: %q (%q)", *s.Comment, *ns.Comment)
	}
	if ns.CreatedAt == nil {
		t.Errorf("Bad created at: empty")
	}
	if ns.UpdatedAt == nil {
		t.Errorf("Bad updated at: empty")
	}
	if s.DeletedAt != nil {
		t.Errorf("Bad deleted at: %s", *ns.DeletedAt)
	}

	// Get Details
	var nsd *ServiceDetail
	Record(t, "services/details", func(c *Client) {
		nsd, err = c.GetServiceDetails(context.TODO(), &GetServiceInput{
			ServiceID: *s.ServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *s.Name != *nsd.Name {
		t.Errorf("bad name: %q (%q)", *s.Name, *nsd.Name)
	}
	if *s.Comment != *nsd.Comment {
		t.Errorf("bad comment: %q (%q)", *s.Comment, *nsd.Comment)
	}
	if *nsd.Version.Number == 0 {
		t.Errorf("Service Detail Version is empty: (%#v)", nsd)
	}

	// Search
	var fs *Service
	Record(t, "services/search", func(c *Client) {
		fs, err = c.SearchService(context.TODO(), &SearchServiceInput{
			Name: "test-service",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *s.Name != *fs.Name {
		t.Errorf("bad name: %q (%q)", *s.Name, *fs.Name)
	}
	if *s.Comment != *fs.Comment {
		t.Errorf("bad comment: %q (%q)", *s.Comment, *fs.Comment)
	}

	// Update
	var us *Service
	Record(t, "services/update", func(c *Client) {
		us, err = c.UpdateService(context.TODO(), &UpdateServiceInput{
			ServiceID: *s.ServiceID,
			Name:      ToPointer("new-test-service"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *us.Name != "new-test-service" {
		t.Errorf("bad name: %q", *us.Name)
	}

	// Delete
	Record(t, "services/delete", func(c *Client) {
		err = c.DeleteService(context.TODO(), &DeleteServiceInput{
			ServiceID: *s.ServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	//	List Domains
	var ds ServiceDomainsList
	Record(t, "services/domain", func(c *Client) {
		ds, err = c.ListServiceDomains(context.TODO(), &ListServiceDomainInput{
			ServiceID: *s.ServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ds) != 0 {
		t.Fatalf("bad services: %v", ds)
	}
}

func TestClient_Services_Compute(t *testing.T) {
	t.Parallel()

	var err error

	// Create
	var s *Service
	Record(t, "services/compute/create", func(c *Client) {
		s, err = c.CreateService(context.TODO(), &CreateServiceInput{
			Name:    ToPointer("test-service"),
			Comment: ToPointer("comment"),
			Type:    ToPointer("wasm"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "services/compute/cleanup", func(c *Client) {
			_ = c.DeleteService(context.TODO(), &DeleteServiceInput{
				ServiceID: *s.ServiceID,
			})

			_ = c.DeleteService(context.TODO(), &DeleteServiceInput{
				ServiceID: *s.ServiceID,
			})
		})
	}()

	if *s.Name != "test-service" {
		t.Errorf("bad name: %q", *s.Name)
	}
	if *s.Comment != "comment" {
		t.Errorf("bad comment: %q", *s.Comment)
	}
	if *s.Type != "wasm" {
		t.Errorf("bad type: %q", *s.Type)
	}

	// List
	var ss []*Service
	Record(t, "services/compute/list", func(c *Client) {
		ss, err = c.ListServices(context.TODO(), &ListServicesInput{
			Direction: ToPointer("descend"),
			Sort:      ToPointer("created"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ss) < 1 {
		t.Errorf("bad services: %v", ss)
	}

	// List with paginator
	var ss2 []*Service
	var paginator *ListPaginator[Service]
	Record(t, "services/compute/list_paginator", func(c *Client) {
		paginator = c.GetServices(context.TODO(), &GetServicesInput{
			Direction: ToPointer("descend"),
			PerPage:   ToPointer(200),
			Sort:      ToPointer("created"),
		})

		for paginator.HasNext() {
			data, err := paginator.GetNext()
			if err != nil {
				t.Errorf("Bad paginator (remaining: %d): %s", paginator.Remaining(), err)
				return
			}
			ss2 = append(ss2, data...)
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(ss2) != len(ss) {
		t.Errorf("expected %d services but got: %d", len(ss), len(ss2))
	}

	// Get
	var ns *Service
	Record(t, "services/compute/get", func(c *Client) {
		ns, err = c.GetService(context.TODO(), &GetServiceInput{
			ServiceID: *s.ServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *s.Name != *ns.Name {
		t.Errorf("bad name: %q (%q)", *s.Name, *ns.Name)
	}
	if *s.Comment != *ns.Comment {
		t.Errorf("bad comment: %q (%q)", *s.Comment, *ns.Comment)
	}
	if ns.CreatedAt == nil {
		t.Errorf("Bad created at: empty")
	}
	if ns.UpdatedAt == nil {
		t.Errorf("Bad updated at: empty")
	}
	if s.DeletedAt != nil {
		t.Errorf("Bad deleted at: %s", *ns.DeletedAt)
	}

	// Get Details
	var nsd *ServiceDetail
	Record(t, "services/compute/details", func(c *Client) {
		nsd, err = c.GetServiceDetails(context.TODO(), &GetServiceInput{
			ServiceID: *s.ServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *s.Name != *nsd.Name {
		t.Errorf("bad name: %q (%q)", *s.Name, *nsd.Name)
	}
	if *s.Comment != *nsd.Comment {
		t.Errorf("bad comment: %q (%q)", *s.Comment, *nsd.Comment)
	}
	if *nsd.Version.Number == 0 {
		t.Errorf("Service Detail Version is empty: (%#v)", nsd)
	}

	// Search
	var fs *Service
	Record(t, "services/compute/search", func(c *Client) {
		fs, err = c.SearchService(context.TODO(), &SearchServiceInput{
			Name: "test-service",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *s.Name != *fs.Name {
		t.Errorf("bad name: %q (%q)", *s.Name, *fs.Name)
	}
	if *s.Comment != *fs.Comment {
		t.Errorf("bad comment: %q (%q)", *s.Comment, *fs.Comment)
	}

	// Update
	var us *Service
	Record(t, "services/compute/update", func(c *Client) {
		us, err = c.UpdateService(context.TODO(), &UpdateServiceInput{
			ServiceID: *s.ServiceID,
			Name:      ToPointer("new-test-service"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if *us.Name != "new-test-service" {
		t.Errorf("bad name: %q", *us.Name)
	}

	// Delete
	Record(t, "services/compute/delete", func(c *Client) {
		err = c.DeleteService(context.TODO(), &DeleteServiceInput{
			ServiceID: *s.ServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	//	List Domains
	var ds ServiceDomainsList
	Record(t, "services/compute/domain", func(c *Client) {
		ds, err = c.ListServiceDomains(context.TODO(), &ListServiceDomainInput{
			ServiceID: *s.ServiceID,
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
	_, err := TestClient.GetService(context.TODO(), &GetServiceInput{})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateService_validation(t *testing.T) {
	_, err := TestClient.UpdateService(context.TODO(), &UpdateServiceInput{})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteService_validation(t *testing.T) {
	err := TestClient.DeleteService(context.TODO(), &DeleteServiceInput{})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}
