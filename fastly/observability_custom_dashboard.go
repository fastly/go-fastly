package fastly

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// ObservabilityCustomDashboard is a named container for a custom dashboard configuration
type ObservabilityCustomDashboard struct {
	// The date and time the dashboard was created
	CreatedAt time.Time `json:"created_at"`
	// The ID of the user who created the dashboard
	CreatedBy string `json:"created_by"`
	// A short description of the dashboard
	Description string `json:"description"`
	// The unique identifier of the dashboard
	ID string `json:"id"`
	// A list of DashboardItems
	Items []DashboardItem `json:"items"`
	// A human-readable name
	Name string `json:"name"`
	// The date and time the dashboard was last updated
	UpdatedAt time.Time `json:"updated_at"`
	// The ID of the user who last modified the dashboard
	UpdatedBy string `json:"updated_by"`
}

// DashboardItem describes an item (or "widget") of a dashboard
type DashboardItem struct {
	// DataSource describes the source of the metrics to be displayed (required)
	DataSource DashboardDataSource `json:"data_source"`
	// ID is a unique identifier for the DashboardItem (read-only)
	ID string `json:"id,omitempty"`
	// Span is the number of columns (1-12) for the DashboardItem to span (default: 4)
	Span uint8 `json:"span"`
	// Subtitle is a human-readable subtitle to display, often a description of the visualization (optional)
	Subtitle string `json:"subtitle"`
	// Title is a human-readable title to display (optional)
	Title string `json:"title"`
	// Visualization describes the way the DashboardItem should display data (required)
	Visualization DashboardVisualization `json:"visualization"`
}

type DashboardSourceType string

const (
	SourceTypeStatsEdge   = "stats.edge"
	SourceTypeStatsDomain = "stats.domain"
	SourceTypeStatsOrigin = "stats.origin"
)

// DashboardDataSource describes the data to display in a DashboardItem
type DashboardDataSource struct {
	// Config describes configuration options for the selected data source (required)
	Config DashboardSourceConfig `json:"config"`
	// Type is the source of the data to display (required)
	Type DashboardSourceType `json:"type"`
}

type DashboardSourceConfig struct {
	// Metrics is the list metrics to visualize (required)
	// Valid options are defined by the selected SourceType. See https://www.fastly.com/documentation/reference/api/observability/custom-dashboards/#data-source
	Metrics []string `json:"metrics"`
}

type VisualizationType string

const VisualizationTypeChart VisualizationType = "chart"

type DashboardVisualization struct {
	// Config describes configuration options for the given visualization (required)
	Config VisualizationConfig `json:"config"`
	// Type is type of visualization to display. Currently only "chart" is supported (required)
	Type VisualizationType `json:"type"`
}

type PlotType string

const (
	PlotTypeLine         PlotType = "line"
	PlotTypeBar          PlotType = "bar"
	PlotTypeDonut        PlotType = "donut"
	PlotTypeSingleMetric PlotType = "single-metric"
)

type VisualizationFormat string

const (
	VisualizationFormatNumber       VisualizationFormat = "number"
	VisualizationFormatBytes        VisualizationFormat = "bytes"
	VisualizationFormatPercent      VisualizationFormat = "percent"
	VisualizationFormatRequests     VisualizationFormat = "requests"
	VisualizationFormatResponses    VisualizationFormat = "responses"
	VisualizationFormatSeconds      VisualizationFormat = "seconds"
	VisualizationFormatMilliseconds VisualizationFormat = "milliseconds"
	VisualizationFormatRatio        VisualizationFormat = "ratio"
	VisualizationFormatBitrate      VisualizationFormat = "bitrate"
)

type CalculationMethod string

const (
	CalculationMethodAvg    CalculationMethod = "avg"
	CalculationMethodSum    CalculationMethod = "sum"
	CalculationMethodMin    CalculationMethod = "min"
	CalculationMethodMax    CalculationMethod = "max"
	CalculationMethodLatest CalculationMethod = "latest"
)

type VisualizationConfig struct {
	// CalculationMethod is the aggregation function to apply to the dataset (optional)
	CalculationMethod *CalculationMethod `json:"calculation_method,omitempty"`
	// Format indicates the unit used to format the data (optional, default: number)
	Format *VisualizationFormat `json:"format,omitempty"`
	// PlotType is the type of chart to display (required)
	PlotType PlotType `json:"plot_type"`
}

type dashboardItemOption interface {
	apply(*DashboardItem)
}
type optionFunc func(*DashboardItem)

func (f optionFunc) apply(di *DashboardItem) {
	f(di)
}

func WithTitle(title string) dashboardItemOption {
	return optionFunc(func(di *DashboardItem) {
		di.Title = title
	})
}
func WithSubtitle(subtitle string) dashboardItemOption {
	return optionFunc(func(di *DashboardItem) {
		di.Subtitle = subtitle
	})
}
func WithSpan(span uint8) dashboardItemOption {
	return optionFunc(func(di *DashboardItem) {
		di.Span = span
	})
}
func WithCalculationMethod(calculationMethod CalculationMethod) dashboardItemOption {
	return optionFunc(func(di *DashboardItem) {
		di.Visualization.Config.CalculationMethod = &calculationMethod
	})
}
func WithFormat(format VisualizationFormat) dashboardItemOption {
	return optionFunc(func(di *DashboardItem) {
		di.Visualization.Config.Format = &format
	})
}

func NewDashboardItem(sourceType DashboardSourceType, metrics []string, plotType PlotType, options ...dashboardItemOption) DashboardItem {
	di := DashboardItem{
		DataSource: DashboardDataSource{
			Type: sourceType,
			Config: DashboardSourceConfig{
				Metrics: metrics,
			},
		},
		Visualization: DashboardVisualization{
			Type: VisualizationTypeChart,
			Config: VisualizationConfig{
				PlotType: plotType,
			},
		},
	}

	for _, o := range options {
		o.apply(&di)
	}

	return di
}

type ListDashboardsResponse struct {
	Data []ObservabilityCustomDashboard `json:"data"`
	Meta DashboardMeta                  `json:"meta"`
}

// DashboardMeta holds metadata about a dashboards query
type DashboardMeta struct {
	Limit      int    `json:"limit"`
	NextCursor string `json:"next_cursor"`
	Sort       string `json:"sort"`
	Total      int    `json:"total"`
}

// ListObservabilityCustomDashboardsInput is used as input to the ListObservabilityCustomDashboards function
type ListObservabilityCustomDashboardsInput struct {
	// Cursor is the pagination cursor from a previous request's meta (optional)
	Cursor *string
	// Limit is the maximum number of items included in each response (optional)
	Limit *int
	// Sort is the field on which to sort dashboards (optional)
	Sort *string
}

func (c *Client) ListObservabilityCustomDashboards(i *ListObservabilityCustomDashboardsInput) (*ListDashboardsResponse, error) {
	path := ToSafeURL("observability", "dashboards")
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
	// ID of the dashboard to fetch (required)
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
	// ID of the dashboard to fetch (required)
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
	// ID of the dashboard to delete (required)
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
