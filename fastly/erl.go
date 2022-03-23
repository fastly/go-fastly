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
	HttpMethods        []string         `mapstructure:"http_methods"`
	ID                 string           `mapstructure:"id"`
	LoggerType         ERLLogger        `mapstructure:"logger_type"`
	Name               string           `mapstructure:"name"`
	PenaltyBoxDuration int              `mapstructure:"penalty_box_duration"` // 1..60
	Response           *ERLResponseType `mapstructure:"response"`             // required if Action != Log
	ResponseObjectName string           `mapstructure:"response_object_name"`
	RpsLimit           int              `mapstructure:"rps_limit"` // 10..10000
	ServiceId          string           `mapstructure:"service_id"`
	UpdatedAt          *time.Time       `mapstructure:"updated_at"`
	UriDictionaryName  string           `mapstructure:"uri_dictionary_name"`
	Version            int              `mapstructure:"version"` // 1..
	WindowSize         ERLWindowSize    `mapstructure:"window_size"`
}

// ERLResponseType models the response from the Fastly API.
type ERLResponseType struct {
	ERLStatus      int    `url:"status,omitempty"`
	ERLContentType string `url:"content_type,omitempty"`
	ERLContent     string `url:"content,omitempty"`
}

// ERLAction represents the action variants for when a rate limiter
// violation is detected.
type ERLAction string

const (
	ERLActionLogOnly        ERLAction = "log_only"
	ERLActionResponse       ERLAction = "response"
	ERLActionResponseObject ERLAction = "response_object"
)

// ERLLogger represents the supported log provider variants.
type ERLLogger string

const (
	ERLLogAzureBlob       ERLLogger = "azureblob"
	ERLLogBigQuery        ERLLogger = "bigquery"
	ERLLogCloudFiles      ERLLogger = "cloudfiles"
	ERLLogDataDog         ERLLogger = "datadog"
	ERLLogDigitalOcean    ERLLogger = "digitalocean"
	ERLLogElasticSearch   ERLLogger = "elasticsearch"
	ERLLogFtp             ERLLogger = "ftp"
	ERLLogGcs             ERLLogger = "gcs"
	ERLLogGoogleAnalytics ERLLogger = "googleanalytics"
	ERLLogHeroku          ERLLogger = "heroku"
	ERLLogHoneycomb       ERLLogger = "honeycomb"
	ERLLogHttp            ERLLogger = "http"
	ERLLogHttps           ERLLogger = "https"
	ERLLogKafta           ERLLogger = "kafka"
	ERLLogKinesis         ERLLogger = "kinesis"
	ERLLogLogEntries      ERLLogger = "logentries"
	ERLLogLoggly          ERLLogger = "loggly"
	ERLLogLogShuttle      ERLLogger = "logshuttle"
	ERLLogNewRelic        ERLLogger = "newrelic"
	ERLLogOpenStack       ERLLogger = "openstack"
	ERLLogPaperTrail      ERLLogger = "papertrail"
	ERLLogPubSub          ERLLogger = "pubsub"
	ERLLogS3              ERLLogger = "s3"
	ERLLogScalyr          ERLLogger = "scalyr"
	ERLLogSftp            ERLLogger = "sftp"
	ERLLogSplunk          ERLLogger = "splunk"
	ERLLogStackDriver     ERLLogger = "stackdriver"
	ERLLogSumoLogiuc      ERLLogger = "sumologic"
	ERLLogSysLog          ERLLogger = "syslog"
)

// ERLWindowSize represents the duration variants for when the RPS limit is
// exceeded.
type ERLWindowSize int

const (
	ERLSize1  ERLWindowSize = 1
	ERLSize10 ERLWindowSize = 10
	ERLSize60 ERLWindowSize = 60
)

// ERLsByName is a sortable list of ERLs
type ERLsByName []*ERL

// Len, Swap, and Less implement the sortable interface
func (s ERLsByName) Len() int      { return len(s) }
func (s ERLsByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
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
	HttpMethods        []string         `url:"http_methods,brackets"`
	Name               string           `url:"name"`
	PenaltyBoxDuration int              `url:"penalty_box_duration"`
	Response           *ERLResponseType `url:"response,omitempty"`
	RpsLimit           int              `url:"rps_limit"`
	ServiceID          string           `url:"-"`
	ServiceVersion     int              `url:"-"`
	WindowSize         ERLWindowSize    `url:"window_size"`
}

// CreateERL returns a new ERL.
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

	var erl *ERL
	if err := decodeBodyMap(resp.Body, &erl); err != nil {
		return nil, err
	}

	return erl, nil
}

// DeleteERLInput is used as input to the DeleteERL function.
type DeleteERLInput struct {
	ServiceID      string `form:"service_id"`
	ServiceVersion int    `form:"version"`
	ERLID          string `form:"id"`
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
	ServiceID      string `form:"service_id"`
	ServiceVersion int    `form:"version"`
	ERLID          string `form:"id"`
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
	HttpMethods        []string         `url:"http_methods,omitempty,brackets"`
	ID                 string           `url:"id"`
	Name               string           `url:"name,omitempty"`
	PenaltyBoxDuration int              `url:"penalty_box_duration,omitempty"`
	Response           *ERLResponseType `url:"response,omitempty"`
	RpsLimit           int              `url:"rps_limit,omitempty"`
	ServiceID          string           `url:"-"`
	ServiceVersion     int              `url:"-"`
	WindowSize         ERLWindowSize    `url:"window_size,omitempty"`
}

// UpdateERLInput updates the specified ERL.
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

	var erl *ERL
	if err := decodeBodyMap(resp.Body, &erl); err != nil {
		return nil, err
	}

	return erl, nil
}
