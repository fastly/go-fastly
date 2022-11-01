//
// API reference:
// https://developer.fastly.com/reference/api/vcl-services/rate-limiter/
// NB: ERL is an optional feature that must be enabled before use
package fastly

import (
	"fmt"
	"sort"
	"time"
)

// ERL models the response from the Fastly API.
type ERL struct {
	Action             ERLAction        `mapstructure:"action"`
	ClientKey          []string         `mapstructure:"client_key"`
	CreatedAt          *time.Time       `mapstructure:"created_at"`
	DeletedAt          *time.Time       `mapstructure:"deleted_at"`
	FeatureRevision    int              `mapstructure:"feature_revision"` // 1..
	HTTPMethods        []string         `mapstructure:"http_methods"`
	ID                 string           `mapstructure:"id"`
	LoggerType         ERLLogger        `mapstructure:"logger_type"`
	Name               string           `mapstructure:"name"`
	PenaltyBoxDuration int              `mapstructure:"penalty_box_duration"` // 1..60
	Response           *ERLResponseType `mapstructure:"response"`             // required if Action != Log
	ResponseObjectName string           `mapstructure:"response_object_name"`
	RpsLimit           int              `mapstructure:"rps_limit"` // 10..10000
	ServiceID          string           `mapstructure:"service_id"`
	UpdatedAt          *time.Time       `mapstructure:"updated_at"`
	URIDictionaryName  string           `mapstructure:"uri_dictionary_name"`
	Version            int              `mapstructure:"version"` // 1..
	WindowSize         ERLWindowSize    `mapstructure:"window_size"`
}

// ERLResponseType models the response from the Fastly API.
type ERLResponseType struct {
	ERLContent     string `url:"content,omitempty"`
	ERLContentType string `url:"content_type,omitempty"`
	ERLStatus      int    `url:"status,omitempty"`
}

// ERLAction represents the action variants for when a rate limiter
// violation is detected.
type ERLAction string

const (
	// ERLActionLogOnly represents an action variant.
	ERLActionLogOnly ERLAction = "log_only"
	// ERLActionResponse represents an action variant.
	ERLActionResponse ERLAction = "response"
	// ERLActionResponseObject represents an action variant.
	ERLActionResponseObject ERLAction = "response_object"
)

// ERLLogger represents the supported log provider variants.
type ERLLogger string

const (
	// ERLLogAzureBlob represents a log provider variant.
	ERLLogAzureBlob ERLLogger = "azureblob"
	// ERLLogBigQuery represents a log provider variant.
	ERLLogBigQuery ERLLogger = "bigquery"
	// ERLLogCloudFiles represents a log provider variant.
	ERLLogCloudFiles ERLLogger = "cloudfiles"
	// ERLLogDataDog represents a log provider variant.
	ERLLogDataDog ERLLogger = "datadog"
	// ERLLogDigitalOcean represents a log provider variant.
	ERLLogDigitalOcean ERLLogger = "digitalocean"
	// ERLLogElasticSearch represents a log provider variant.
	ERLLogElasticSearch ERLLogger = "elasticsearch"
	// ERLLogFtp represents a log provider variant.
	ERLLogFtp ERLLogger = "ftp"
	// ERLLogGcs represents a log provider variant.
	ERLLogGcs ERLLogger = "gcs"
	// ERLLogGoogleAnalytics represents a log provider variant.
	ERLLogGoogleAnalytics ERLLogger = "googleanalytics"
	// ERLLogHeroku represents a log provider variant.
	ERLLogHeroku ERLLogger = "heroku"
	// ERLLogHoneycomb represents a log provider variant.
	ERLLogHoneycomb ERLLogger = "honeycomb"
	// ERLLogHTTP represents a log provider variant.
	ERLLogHTTP ERLLogger = "http"
	// ERLLogHTTPS represents a log provider variant.
	ERLLogHTTPS ERLLogger = "https"
	// ERLLogKafta represents a log provider variant.
	ERLLogKafta ERLLogger = "kafka"
	// ERLLogKinesis represents a log provider variant.
	ERLLogKinesis ERLLogger = "kinesis"
	// ERLLogLogEntries represents a log provider variant.
	ERLLogLogEntries ERLLogger = "logentries"
	// ERLLogLoggly represents a log provider variant.
	ERLLogLoggly ERLLogger = "loggly"
	// ERLLogLogShuttle represents a log provider variant.
	ERLLogLogShuttle ERLLogger = "logshuttle"
	// ERLLogNewRelic represents a log provider variant.
	ERLLogNewRelic ERLLogger = "newrelic"
	// ERLLogOpenStack represents a log provider variant.
	ERLLogOpenStack ERLLogger = "openstack"
	// ERLLogPaperTrail represents a log provider variant.
	ERLLogPaperTrail ERLLogger = "papertrail"
	// ERLLogPubSub represents a log provider variant.
	ERLLogPubSub ERLLogger = "pubsub"
	// ERLLogS3 represents a log provider variant.
	ERLLogS3 ERLLogger = "s3"
	// ERLLogScalyr represents a log provider variant.
	ERLLogScalyr ERLLogger = "scalyr"
	// ERLLogSftp represents a log provider variant.
	ERLLogSftp ERLLogger = "sftp"
	// ERLLogSplunk represents a log provider variant.
	ERLLogSplunk ERLLogger = "splunk"
	// ERLLogStackDriver represents a log provider variant.
	ERLLogStackDriver ERLLogger = "stackdriver"
	// ERLLogSumoLogic represents a log provider variant.
	ERLLogSumoLogic ERLLogger = "sumologic"
	// ERLLogSysLog represents a log provider variant.
	ERLLogSysLog ERLLogger = "syslog"
)

// ERLWindowSize represents the duration variants for when the RPS limit is
// exceeded.
type ERLWindowSize int

const (
	// ERLSize1 represents a duration variant.
	ERLSize1 ERLWindowSize = 1
	// ERLSize10 represents a duration variant.
	ERLSize10 ERLWindowSize = 10
	// ERLSize60 represents a duration variant.
	ERLSize60 ERLWindowSize = 60
)

// ERLsByName is a sortable list of ERLs
type ERLsByName []*ERL

// Len implement the sortable interface.
func (s ERLsByName) Len() int {
	return len(s)
}

// Swap implement the sortable interface.
func (s ERLsByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implement the sortable interface.
func (s ERLsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListERLsInput is used as input to the ListERLs function.
type ListERLsInput struct {
	ServiceID      string // required
	ServiceVersion int    // required
}

// ListERLs returns the list of ERLs for the specified service version.
func (c *Client) ListERLs(i *ListERLsInput) ([]*ERL, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/rate-limiters", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var erls []*ERL
	if err := decodeBodyMap(resp.Body, &erls); err != nil {
		return nil, err
	}

	sort.Stable(ERLsByName(erls))
	return erls, nil
}

// CreateERLInput is used as input to the CreateERL function.
type CreateERLInput struct {
	Action             ERLAction        `url:"action"`
	ClientKey          []string         `url:"client_key,brackets"`
	HTTPMethods        []string         `url:"http_methods,brackets"`
	Name               string           `url:"name"`
	PenaltyBoxDuration int              `url:"penalty_box_duration"`
	Response           *ERLResponseType `url:"response,omitempty"`
	RpsLimit           int              `url:"rps_limit"`
	ServiceID          string           `url:"-"`
	ServiceVersion     int              `url:"-"`
	WindowSize         ERLWindowSize    `url:"window_size"`
}

// CreateERL creates a new resource.
func (c *Client) CreateERL(i *CreateERLInput) (*ERL, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/rate-limiters", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var erl *ERL
	if err := decodeBodyMap(resp.Body, &erl); err != nil {
		return nil, err
	}

	return erl, nil
}

// DeleteERLInput is used as input to the DeleteERL function.
type DeleteERLInput struct {
	ERLID          string `form:"id"`
	ServiceID      string `form:"service_id"`
	ServiceVersion int    `form:"version"`
}

// DeleteERL deletes the specified ERL.
func (c *Client) DeleteERL(i *DeleteERLInput) error {
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}
	if i.ERLID == "" {
		return ErrMissingID
	}

	path := fmt.Sprintf("/rate-limiters/%s", i.ERLID)
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

// GetERLInput is used as input to the GetERL function.
type GetERLInput struct {
	ERLID          string `form:"id"`
	ServiceID      string `form:"service_id"`
	ServiceVersion int    `form:"version"`
}

// GetERL returns the specified ERL.
func (c *Client) GetERL(i *GetERLInput) (*ERL, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}
	if i.ERLID == "" {
		return nil, ErrMissingID
	}

	path := fmt.Sprintf("/rate-limiters/%s", i.ERLID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var erl *ERL
	if err := decodeBodyMap(resp.Body, &erl); err != nil {
		return nil, err
	}

	return erl, nil
}

// UpdateERLInput is used as input to the UpdateERL function.
type UpdateERLInput struct {
	Action             ERLAction        `url:"action,omitempty"`
	ClientKey          []string         `url:"client_key,omitempty,brackets"`
	HTTPMethods        []string         `url:"http_methods,omitempty,brackets"`
	ID                 string           `url:"id"`
	Name               string           `url:"name,omitempty"`
	PenaltyBoxDuration int              `url:"penalty_box_duration,omitempty"`
	Response           *ERLResponseType `url:"response,omitempty"`
	RpsLimit           int              `url:"rps_limit,omitempty"`
	ServiceID          string           `url:"-"`
	ServiceVersion     int              `url:"-"`
	WindowSize         ERLWindowSize    `url:"window_size,omitempty"`
}

// UpdateERL updates the specified resource.
func (c *Client) UpdateERL(i *UpdateERLInput) (*ERL, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := fmt.Sprintf("/rate-limiters/%s", i.ID)
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var erl *ERL
	if err := decodeBodyMap(resp.Body, &erl); err != nil {
		return nil, err
	}

	return erl, nil
}
