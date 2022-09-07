package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// BigQuery represents a BigQuery response from the Fastly API.
type BigQuery struct {
	ServiceID      string `mapstructure:"service_id"`
	ServiceVersion int    `mapstructure:"version"`

	Name              string     `mapstructure:"name"`
	Format            string     `mapstructure:"format"`
	User              string     `mapstructure:"user"`
	ProjectID         string     `mapstructure:"project_id"`
	Dataset           string     `mapstructure:"dataset"`
	Table             string     `mapstructure:"table"`
	Template          string     `mapstructure:"template_suffix"`
	SecretKey         string     `mapstructure:"secret_key"`
	ResponseCondition string     `mapstructure:"response_condition"`
	Placement         string     `mapstructure:"placement"`
	FormatVersion     uint       `mapstructure:"format_version"`
	CreatedAt         *time.Time `mapstructure:"created_at"`
	UpdatedAt         *time.Time `mapstructure:"updated_at"`
	DeletedAt         *time.Time `mapstructure:"deleted_at"`
}

// bigQueriesByName is a sortable list of BigQueries.
type bigQueriesByName []*BigQuery

// Len, Swap, and Less implement the sortable interface.
func (s bigQueriesByName) Len() int      { return len(s) }
func (s bigQueriesByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
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

// ListBigQueries returns the list of BigQueries for the configuration version.
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
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	Name              string `url:"name,omitempty"`
	ProjectID         string `url:"project_id,omitempty"`
	Dataset           string `url:"dataset,omitempty"`
	Table             string `url:"table,omitempty"`
	Template          string `url:"template_suffix,omitempty"`
	User              string `url:"user,omitempty"`
	SecretKey         string `url:"secret_key,omitempty"`
	Format            string `url:"format,omitempty"`
	ResponseCondition string `url:"response_condition,omitempty"`
	Placement         string `url:"placement,omitempty"`
	FormatVersion     uint   `url:"format_version,omitempty"`
}

// CreateBigQuery creates a new Fastly BigQuery.
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
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the BigQuery to fetch.
	Name string
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
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the BigQuery to update.
	Name string

	NewName           *string `url:"name,omitempty"`
	ProjectID         *string `url:"project_id,omitempty"`
	Dataset           *string `url:"dataset,omitempty"`
	Table             *string `url:"table,omitempty"`
	Template          *string `url:"template_suffix,omitempty"`
	User              *string `url:"user,omitempty"`
	SecretKey         *string `url:"secret_key,omitempty"`
	Format            *string `url:"format,omitempty"`
	ResponseCondition *string `url:"response_condition,omitempty"`
	Placement         *string `url:"placement,omitempty"`
	FormatVersion     *uint   `url:"format_version,omitempty"`
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
	// ServiceID is the ID of the service (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// Name is the name of the BigQuery to delete (required).
	Name string
}

// DeleteBigQuery deletes the given BigQuery version.
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
