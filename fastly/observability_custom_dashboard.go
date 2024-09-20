package fastly

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const path = "/observability/dashboards"

const (
	ErrInvalidDataSourceType    = ""
	ErrInvalidMetric            = ""
	ErrInvalidVisualizationType = ""
	ErrInvalidPlotType          = ""
	ErrInvalidCalculationMethod = ""
)

// ObservabilityCustomDashboard is a named container for a custom dashboard configuration
type ObservabilityCustomDashboard struct {
	// The date and time the dashboard was created
	CreatedAt time.Time
	// The ID of the user who created the dashboard
	CreatedBy string
	// A short description of the dashboard
	Description string
	// The unique identifier of the dashboard
	ID string
	// A list of DashboardItems
	Items []DashboardItem
	// A human-readable name
	Name string
	// The date and time the dashboard was last updated
	UpdatedAt time.Time
	// The ID of the user who last modified the dashboard
	UpdatedBy string
}

// DashboardItem describes an item (or "widget") of a dashboard
type DashboardItem struct {
	// DataSource describes the source of the metrics to be displayed (required).
	DataSource DataSource `json:"data_source"`
	// ID is a unique identifier for the DashboardItem (read-only).
	ID string `json:"id,omitempty"`
	// Span is the number of columns (1-12) for the DashboardItem to span (default: 4).
	Span uint8 `json:"span"`
	// Subtitle is a human-readable subtitle to display, often a description of the visualization (optional).
	Subtitle string `json:"subtitle"`
	// Title is a human-readable title to display (optional).
	Title string `json:"title"`
	// Visualization describes the way the DashboardItem should display data (required).
	Visualization Visualization `json:"visualization"`
}

type SourceType int

const (
	SourceTypeUndefined SourceType = iota
	SourceTypeStatsEdge
	SourceTypeStatsDomain
	SourceTypeStatsOrigin
)

var stringToSourceType = map[string]SourceType{
	"stats.edge":   SourceTypeStatsEdge,
	"stats.domain": SourceTypeStatsDomain,
	"stats.origin": SourceTypeStatsOrigin,
}

var sourceTypeToString = map[SourceType]string{
	SourceTypeStatsEdge:   "stats.edge",
	SourceTypeStatsDomain: "stats.domain",
	SourceTypeStatsOrigin: "stats.origin",
}

func (st SourceType) String() string {
	return sourceTypeToString[st]
}
func (st SourceType) MarshalJSON() ([]byte, error) {
	if st == SourceTypeUndefined {
		return nil, errors.New("cannot marshal undefined SourceType")
	}
	return json.Marshal(st.String())
}
func (st *SourceType) UnmarshalJSON(data []byte) (err error) {
	var str string
	if err = json.Unmarshal(data, &str); err != nil {
		return err
	}
	var ok bool
	if *st, ok = stringToSourceType[str]; !ok {
		return NewFieldError("DataSource.Type").Message(fmt.Sprintf("Invalid value \"%s\" for DataSource.Type", data))
	}
	return nil
}

type DataSource struct {
	// Config describes configuration options for the selected data source (required).
	Config SourceConfig `json:"config"`
	// Type is the source of the data to display (required).
	Type SourceType `json:"type"`
}

type Metric string

type SourceConfig struct {
	Metrics []Metric `json:"metrics"`
}

type VisualizationType string

const VisualizationTypeChart VisualizationType = "chart"

type Visualization struct {
	// Config describes configuration options for the given visualization (required).
	Config VisualizationConfig `json:"config"`
	// Type is type of visualization to display. Currently only "chart" is supported (required).
	Type VisualizationType `json:"type"`
}

type PlotType string

const (
	PlotTypeLine         PlotType = "line"
	PlotTypeBar                   = "bar"
	PlotTypeDonut                 = "donut"
	PlotTypeSingleMetric          = "single-metric"
)

type Format string

const (
	FormatNumber       Format = "number"
	FormatBytes               = "bytes"
	FormatPercent             = "percent"
	FormatRequests            = "requests"
	FormatResponses           = "responses"
	FormatSeconds             = "seconds"
	FormatMilliseconds        = "milliseconds"
	FormatRatio               = "ratio"
	FormatBitrate             = "bitrate"
)

type CalculationMethod string

const (
	CalculationMethodAvg    CalculationMethod = "avg"
	CalculationMethodSum                      = "sum"
	CalculationMethodMin                      = "min"
	CalculationMethodMax                      = "max"
	CalculationMethodLatest                   = "latest"
)

type VisualizationConfig struct {
	// CalculationMethod is the aggregation function to apply to the dataset (optional).
	CalculationMethod *CalculationMethod `json:"calculation_method,omitempty"`
	// Format is the units to use to format the data (optional, default: number).
	Format *Format `json:"format,omitempty"`
	// PlotType is the type of chart to display (required).
	PlotType PlotType `json:"plot_type"`
}

type ListDashboardsResponse struct {
	Data []ObservabilityCustomDashboard `json:"data"`
	Meta DashboardMeta                  `json:"meta"`
}

// DashboardMeta holds metadata about a dashboards query.
type DashboardMeta struct {
	Limit      int    `json:"limit"`
	NextCursor string `json:"next_cursor"`
	Sort       string `json:"sort"`
	Total      int    `json:"total"`
}

// ListObservabilityCustomDashboardsInput is used as input to the ListObservabilityCustomDashboards function.
type ListObservabilityCustomDashboardsInput struct {
	// Cursor is the pagination cursor from a previous request's meta (optional).
	Cursor *string
	// Limit is the maximum number of items included in each response (optional).
	Limit *int
	// Sort is the field on which to sort dashboards (optional).
	Sort *string
}

func (c *Client) ListObservabilityCustomDashboards(i *ListObservabilityCustomDashboardsInput) (*ListDashboardsResponse, error) {
	ro := &RequestOptions{
		Params: map[string]string{},
	}
	if i.Cursor != nil {
		ro.Params["cursor"] = *i.Cursor
	}
	if i.Limit != nil {
		ro.Params["limit"] = strconv.Itoa(*i.Limit)
	}
	if i.Sort != nil {
		ro.Params["sort"] = *i.Sort
	}

	resp, err := c.Get(path, ro)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ldr *ListDashboardsResponse
	if err := json.NewDecoder(resp.Body).Decode(&ldr); err != nil {
		return nil, err
	}

	return ldr, nil
}

type CreateObservabilityCustomDashboardInput struct {
	Description *string         `json:"description,omitempty"`
	Items       []DashboardItem `json:"items"`
	Name        string          `json:"name"`
}

func (c *Client) CreateObservabilityCustomDashboard(i *CreateObservabilityCustomDashboardInput) (*ObservabilityCustomDashboard, error) {
	path := ToSafeURL("observability", "dashboards")
	resp, err := c.PostJSON(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ocd *ObservabilityCustomDashboard
	if err := json.NewDecoder(resp.Body).Decode(&ocd); err != nil {
		return nil, err
	}
	return ocd, nil
}

type GetObservabilityCustomDashboardInput struct {
	// ID of the dashboard to fetch (required).
	ID *string
}

func (c *Client) GetObservabilityCustomDashboard(i *GetObservabilityCustomDashboardInput) (*ObservabilityCustomDashboard, error) {
	if i.ID == nil {
		return nil, ErrMissingID
	}

	path := ToSafeURL("observability", "dashboards", *i.ID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ocd *ObservabilityCustomDashboard
	if err := json.NewDecoder(resp.Body).Decode(&ocd); err != nil {
		return nil, err
	}
	return ocd, nil
}

type UpdateObservabilityCustomDashboardInput struct {
	Description *string `json:"description,omitempty"`
	// ID of the dashboard to fetch (required).
	ID    *string          `json:"-"`
	Items *[]DashboardItem `json:"items,omitempty"`
	Name  *string          `json:"name,omitempty"`
}

func (c *Client) UpdateObservabilityCustomDashboard(i *UpdateObservabilityCustomDashboardInput) (*ObservabilityCustomDashboard, error) {
	if i.ID == nil {
		return nil, ErrMissingID
	}

	path := ToSafeURL("observability", "dashboards", *i.ID)
	resp, err := c.PatchJSON(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ocd *ObservabilityCustomDashboard
	if err := json.NewDecoder(resp.Body).Decode(&ocd); err != nil {
		return nil, err
	}
	return ocd, nil
}

type DeleteObservabilityCustomDashboardInput struct {
	// ID of the dashboard to delete (required).
	ID *string
}

func (c *Client) DeleteObservabilityCustomDashboard(i *DeleteObservabilityCustomDashboardInput) error {
	if i.ID == nil {
		return ErrMissingID
	}
	path := ToSafeURL("observability", "dashboards", *i.ID)
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return NewHTTPError(resp)
	}

	return nil
}
