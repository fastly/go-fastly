package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// BigQuery represents a BigQuery response from the Fastly API.
type BigQuery struct {
	AccountName       string     `mapstructure:"account_name"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	Dataset           string     `mapstructure:"dataset"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
	Format            string     `mapstructure:"format"`
	FormatVersion     uint       `mapstructure:"format_version"`
	Name              string     `mapstructure:"name"`
	Placement         string     `mapstructure:"placement"`
	ProjectID         string     `mapstructure:"project_id"`
	ResponseCondition string     `mapstructure:"response_condition"`
	SecretKey         string     `mapstructure:"secret_key"`
	ServiceID         string     `mapstructure:"service_id"`
	ServiceVersion    int        `mapstructure:"version"`
	Table             string     `mapstructure:"table"`
	Template          string     `mapstructure:"template_suffix"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	User              string     `mapstructure:"user"`
}

// bigQueriesByName is a sortable list of BigQueries.
type bigQueriesByName []*BigQuery

// Len implements the sortable interface.
func (s bigQueriesByName) Len() int {
	return len(s)
}

// Swap implements the sortable interface.
func (s bigQueriesByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implements the sortable interface.
func (s bigQueriesByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListBigQueriesInput is used as input to the ListBigQueries function.
type ListBigQueriesInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListBigQueries retrieves all resources.
func (c *Client) ListBigQueries(i *ListBigQueriesInput) ([]*BigQuery, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/bigquery", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bigQueries []*BigQuery
	if err := decodeBodyMap(resp.Body, &bigQueries); err != nil {
		return nil, err
	}
	sort.Stable(bigQueriesByName(bigQueries))
	return bigQueries, nil
}

// CreateBigQueryInput is used as input to the CreateBigQuery function.
type CreateBigQueryInput struct {
	AccountName       string `url:"account_name,omitempty"`
	Dataset           string `url:"dataset,omitempty"`
	Format            string `url:"format,omitempty"`
	FormatVersion     uint   `url:"format_version,omitempty"`
	Name              string `url:"name,omitempty"`
	Placement         string `url:"placement,omitempty"`
	ProjectID         string `url:"project_id,omitempty"`
	ResponseCondition string `url:"response_condition,omitempty"`
	SecretKey         string `url:"secret_key,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	Table          string `url:"table,omitempty"`
	Template       string `url:"template_suffix,omitempty"`
	User           string `url:"user,omitempty"`
}

// CreateBigQuery creates a new resource.
func (c *Client) CreateBigQuery(i *CreateBigQueryInput) (*BigQuery, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/bigquery", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bigQuery *BigQuery
	if err := decodeBodyMap(resp.Body, &bigQuery); err != nil {
		return nil, err
	}
	return bigQuery, nil
}

// GetBigQueryInput is used as input to the GetBigQuery function.
type GetBigQueryInput struct {
	// Name is the name of the BigQuery to fetch.
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetBigQuery gets the BigQuery configuration with the given parameters.
func (c *Client) GetBigQuery(i *GetBigQueryInput) (*BigQuery, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/bigquery/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bigQuery *BigQuery
	if err := decodeBodyMap(resp.Body, &bigQuery); err != nil {
		return nil, err
	}
	return bigQuery, nil
}

// UpdateBigQueryInput is used as input to the UpdateBigQuery function.
type UpdateBigQueryInput struct {
	AccountName   *string `url:"account_name,omitempty"`
	Dataset       *string `url:"dataset,omitempty"`
	Format        *string `url:"format,omitempty"`
	FormatVersion *uint   `url:"format_version,omitempty"`
	// Name is the name of the BigQuery to update.
	Name              string
	NewName           *string `url:"name,omitempty"`
	Placement         *string `url:"placement,omitempty"`
	ProjectID         *string `url:"project_id,omitempty"`
	ResponseCondition *string `url:"response_condition,omitempty"`
	SecretKey         *string `url:"secret_key,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
	Table          *string `url:"table,omitempty"`
	Template       *string `url:"template_suffix,omitempty"`
	User           *string `url:"user,omitempty"`
}

// UpdateBigQuery updates a specific BigQuery.
func (c *Client) UpdateBigQuery(i *UpdateBigQueryInput) (*BigQuery, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.Name == "" {
		return nil, ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/bigquery/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var bigQuery *BigQuery
	if err := decodeBodyMap(resp.Body, &bigQuery); err != nil {
		return nil, err
	}
	return bigQuery, nil
}

// DeleteBigQueryInput is the input parameter to DeleteBigQuery.
type DeleteBigQueryInput struct {
	// Name is the name of the BigQuery to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteBigQuery deletes the specified resource.
func (c *Client) DeleteBigQuery(i *DeleteBigQueryInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	if i.Name == "" {
		return ErrMissingName
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/bigquery/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r *statusResp
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return fmt.Errorf("not ok")
	}
	return nil
}
