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
		nsd, err = c.GetServiceDetails(context.TODO(), &GetServiceDetailsInput{
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
		nsd, err = c.GetServiceDetails(context.TODO(), &GetServiceDetailsInput{
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

func TestClient_GetServiceDetails_WithFilters(t *testing.T) {
	t.Parallel()

	var err error

	// Create
	var s *Service
	Record(t, "services/details_with_filters/create", func(c *Client) {
		s, err = c.CreateService(context.TODO(), &CreateServiceInput{
			Name:    ToPointer("test-service-filters"),
			Comment: ToPointer("test filters"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		Record(t, "services/details_with_filters/cleanup", func(c *Client) {
			_ = c.DeleteService(context.TODO(), &DeleteServiceInput{
				ServiceID: *s.ServiceID,
			})
		})
	}()

	// Get Details with versions.active filter
	var nsd *ServiceDetail
	Record(t, "services/details_with_filters/active", func(c *Client) {
		nsd, err = c.GetServiceDetails(context.TODO(), &GetServiceDetailsInput{
			ServiceID: *s.ServiceID,
			Filters: []ServiceDetailsFilter{
				{Key: "versions.active", Value: true},
			},
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
	if nsd.Version == nil {
		t.Fatal("Service Detail Version is nil")
	}
	if *nsd.Version.Number == 0 {
		t.Errorf("Service Detail Version is empty: (%#v)", nsd)
	}

	// Get Details with multiple filters
	var nsd2 *ServiceDetail
	Record(t, "services/details_with_filters/multiple", func(c *Client) {
		nsd2, err = c.GetServiceDetails(context.TODO(), &GetServiceDetailsInput{
			ServiceID: *s.ServiceID,
			Filters: []ServiceDetailsFilter{
				{Key: "versions.active", Value: true},
				{Key: "versions.staged", Value: true},
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *s.Name != *nsd2.Name {
		t.Errorf("bad name: %q (%q)", *s.Name, *nsd2.Name)
	}

	// Get Details with version parameter
	var nsd3 *ServiceDetail
	Record(t, "services/details_with_filters/version", func(c *Client) {
		nsd3, err = c.GetServiceDetails(context.TODO(), &GetServiceDetailsInput{
			ServiceID: *s.ServiceID,
			Version:   ToPointer(1),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *s.Name != *nsd3.Name {
		t.Errorf("bad name: %q (%q)", *s.Name, *nsd3.Name)
	}

	// Delete
	Record(t, "services/details_with_filters/delete", func(c *Client) {
		err = c.DeleteService(context.TODO(), &DeleteServiceInput{
			ServiceID: *s.ServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

// TestClient_GetServiceDetails_StagedFilter_ShouldNarrowVersions asserts the
// behavior filter[versions.staged]=true *should* have: narrowing the
// "versions" array down to just the staged version, symmetric with how
// filter[versions.active]=true already behaves.
func TestClient_GetServiceDetails_StagedFilter_ShouldNarrowVersions(t *testing.T) {
	var err error

	var s *Service
	Record(t, "services/details_staged_filter/create", func(c *Client) {
		s, err = c.CreateService(context.TODO(), &CreateServiceInput{
			Name:    ToPointer("test-service-staged-filter"),
			Comment: ToPointer("test staged filter"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		Record(t, "services/details_staged_filter/cleanup", func(c *Client) {
			_, _ = c.DeactivateVersion(context.TODO(), &DeactivateVersionInput{
				ServiceID:      *s.ServiceID,
				ServiceVersion: 1,
			})
			_ = c.DeleteService(context.TODO(), &DeleteServiceInput{
				ServiceID: *s.ServiceID,
			})
		})
	}()

	var v2 *Version
	Record(t, "services/details_staged_filter/create_version", func(c *Client) {
		v2, err = c.CreateVersion(context.TODO(), &CreateVersionInput{
			ServiceID: *s.ServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	Record(t, "services/details_staged_filter/activate_production", func(c *Client) {
		_, err = c.ActivateVersion(context.TODO(), &ActivateVersionInput{
			ServiceID:      *s.ServiceID,
			ServiceVersion: 1,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	Record(t, "services/details_staged_filter/activate_staging", func(c *Client) {
		_, err = c.ActivateVersion(context.TODO(), &ActivateVersionInput{
			ServiceID:      *s.ServiceID,
			ServiceVersion: *v2.Number,
			Environment:    "staging",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	var stagedFiltered *ServiceDetail
	Record(t, "services/details_staged_filter/staged_filter", func(c *Client) {
		stagedFiltered, err = c.GetServiceDetails(context.TODO(), &GetServiceDetailsInput{
			ServiceID: *s.ServiceID,
			Filters: []ServiceDetailsFilter{
				{Key: "versions.staged", Value: true},
			},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(stagedFiltered.Versions) != 1 {
		t.Errorf("expected filter[versions.staged]=true to narrow to 1 version, got %d: the API is not filtering on this key", len(stagedFiltered.Versions))
		return
	}
	if *stagedFiltered.Versions[0].Number != 2 || !*stagedFiltered.Versions[0].Staging {
		t.Errorf("expected only the staged version (2) back, got: %#v", stagedFiltered.Versions[0])
	}
}

func TestClient_GetServiceDetails_validation(t *testing.T) {
	_, err := TestClient.GetServiceDetails(context.TODO(), &GetServiceDetailsInput{})
	if !errors.Is(err, ErrMissingServiceID) {
		t.Errorf("bad error: %s", err)
	}
}
