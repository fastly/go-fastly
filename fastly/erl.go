//
// API reference: 
// https://developer.fastly.com/reference/api/vcl-services/rate-limiter/
package fastly

import (
	"fmt"
	"sort"
	"time"
)

// ERL models the response from the Fastly API.
type ERL struct {
	Action             ERLAction        `mapstructure:"action"`
	ClientKey          []string             `mapstructure:"client_key"`
	CreatedAt          *time.Time           `mapstructure:"created_at"`
	DeletedAt          *time.Time           `mapstructure:"deleted_at"`
	FeatureRevision    int                  `mapstructure:"feature_revision"` // 1..
	HttpMethods        []ERLHttpMethods `mapstructure:"http_methods"`
	ID                 string               `mapstructure:"id"`
	LoggerType         ERLLogger        `mapstructure:"logger_type"`
	Name               string               `mapstructure:"name"`
	PenaltyBoxDuration int                  `mapstructure:"version"`  // 1..60
	Response           *ERLResponseType     `mapstructure:"response"` // required if Action != Log
	ResponseObjectName string               `mapstructure:"response_object_name"`
	RpsLimit           int                  `mapstructure:"rps_limit"` // 10..10000
	ServiceId          string               `mapstructure:"service_id"`
	UpdatedAt          *time.Time           `mapstructure:"updated_at"`
	UriDictionaryName  string               `mapstructure:"uri_dictionary_name"`
	Version            int                  `mapstructure:"version"` // 1..
	WindowSize         ERLWindowSize    `mapstructure:"window_size"`
}

// Type for ERL
type ERLResponseType struct {
	ERLStatus      int    `json:"status,omitempty"`
	ERLContentType string `json:"content_type,omitempty"`
	ERLContent     string `json:"content,omitempty"`
}

// Enum for ERL
type ERLActionEnum string

const (
	ERLActionResponse       ERLActionEnum = "response"
	ERLActionResponseObject ERLActionEnum = "response_object"
	ERLActionLogOnly        ERLActionEnum = "log_only"
)

// Enum for ERL
type ERLHttpMethodsEnum string

const (
	ERLGet     ERLHttpMethodsEnum = "GET"
	ERLPut     ERLHttpMethodsEnum = "PUT"
	ERLTrace   ERLHttpMethodsEnum = "TRACE"
	ERLPost    ERLHttpMethodsEnum = "POST"
	ERLHead    ERLHttpMethodsEnum = "HEAD"
	ERLDelete  ERLHttpMethodsEnum = "DELETE"
	ERLPatch   ERLHttpMethodsEnum = "PATCH"
	ERLOptions ERLHttpMethodsEnum = "OPTIONS"
)

// Enum for ERL
type ERLLoggerEnum string

const (
	ERLLogAzureBlob       ERLLoggerEnum = "azureblob"
	ERLLogBigQuery        ERLLoggerEnum = "bigquery"
	ERLLogCloudFiles      ERLLoggerEnum = "cloudfiles"
	ERLLogDataDog         ERLLoggerEnum = "datadog"
	ERLLogDigitalOcean    ERLLoggerEnum = "digitalocean"
	ERLLogElasticSearch   ERLLoggerEnum = "elasticsearch"
	ERLLogFtp             ERLLoggerEnum = "ftp"
	ERLLogGcs             ERLLoggerEnum = "gcs"
	ERLLogGoogleAnalytics ERLLoggerEnum = "googleanalytics"
	ERLLogHeroku          ERLLoggerEnum = "heroku"
	ERLLogHoneycomb       ERLLoggerEnum = "honeycomb"
	ERLLogHttp            ERLLoggerEnum = "http"
	ERLLogHttps           ERLLoggerEnum = "https"
	ERLLogKafta           ERLLoggerEnum = "kafka"
	ERLLogKinesis         ERLLoggerEnum = "kinesis"
	ERLLogLogEntries      ERLLoggerEnum = "logentries"
	ERLLogLoggly          ERLLoggerEnum = "loggly"
	ERLLogLogShuttle      ERLLoggerEnum = "logshuttle"
	ERLLogNewRelic        ERLLoggerEnum = "newrelic"
	ERLLogOpenStack       ERLLoggerEnum = "openstack"
	ERLLogPaperTrail      ERLLoggerEnum = "papertrail"
	ERLLogPubSub          ERLLoggerEnum = "pubsub"
	ERLLogS3              ERLLoggerEnum = "s3"
	ERLLogScalyr          ERLLoggerEnum = "scalyr"
	ERLLogSftp            ERLLoggerEnum = "sftp"
	ERLLogSplunk          ERLLoggerEnum = "splunk"
	ERLLogStackDriver     ERLLoggerEnum = "stackdriver"
	ERLLogSumoLogiuc      ERLLoggerEnum = "sumologic"
	ERLLogSysLog          ERLLoggerEnum = "syslog"
)

// Enum for ERL
type ERLWindowSizeEnum int

const (
	ERLSize1  ERLWindowSizeEnum = 1
	ERLSize10 ERLWindowSizeEnum = 10
	ERLSize60 ERLWindowSizeEnum = 60
)

// ERLsByName is a sortable list of ERLs
type ERLsByName []*ERL

// Len, Swap, and Less implement the sortable interface
func (s ERLsByName) Len() int      { return len(s) }
func (s ERLsByName) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ERLsByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ********************************************************************************************* //
//
// Function 1: List rate limiters
// Note(s):
//
// Return type
//   []*ERL

// Input to the ListERLs function
type ListERLsInput struct {
	ServiceID      string // required
	ServiceVersion int    // required
}

// Returns the list of ERLs for the specified version
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

// ********************************************************************************************* //
//
// Function 2: Create rate limiter
// Note(s):
//   VersionId must not be locked
//
// Return type
//   *ERL

// Input to the CreateERL function
type CreateERLInput struct {
	ServiceID          string               `json:"-"`
	ServiceVersion     int                  `json:"-"`
	Action             ERLActionEnum        `json:"action"`
	ClientKey          []string             `json:"client_key"`
	HttpMethods        []ERLHttpMethodsEnum `json:"http_methods"`
	Name               string               `json:"name"`
	PenaltyBoxDuration int                  `json:"penalty_box_duration"`
	Response           *ERLResponseType     `json:"response,omitempty"`
	RpsLimit           int                  `json:"rps_limit"`
	WindowSize         ERLWindowSizeEnum    `json:"window_size"`
}

// Returns the newly created ERL
func (c *Client) CreateERL(i *CreateERLInput) (*ERL, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/rate-limiters", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostJSON(path, i, nil)
	if err != nil {
		return nil, err
	}

	var erl *ERL
	if err := decodeBodyMap(resp.Body, &erl); err != nil {
		return nil, err
	}

	return erl, nil
}

// ********************************************************************************************* //
//
// Function 3: Delete rate limiter
// Note(s):
//
// Return type
//   error

// Input to the DeleteERL function
type DeleteERLInput struct {
	ServiceID      string `form:"service_id"`
	ServiceVersion int    `form:"version"`
	ERLID          string `form:"id"`
}

// Deletes the specified ERL
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

// ********************************************************************************************* //
//
// Function 4: Get rate limiter
// Note(s):
//
// Return type
//   error

// Input to the GetERL function
type GetERLInput struct {
	ServiceID      string `form:"service_id"`
	ServiceVersion int    `form:"version"`
	ERLID          string `form:"id"`
}

// Gets the specified ERL
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

// ********************************************************************************************* //
//
// Function 5: Update rate limiter
// Note(s):
//
// Return type
//   *ERL

// Input to the UpdateERL function
type UpdateERLInput struct {
	ServiceID          string               `json:"-"`
	ServiceVersion     int                  `json:"-"`
	ID                 string               `json:"id"`
	Action             ERLActionEnum        `json:"action,omitempty"`
	ClientKey          []string             `json:"client_key,omitempty"`
	HttpMethods        []ERLHttpMethodsEnum `json:"http_methods,omitempty"`
	Name               string               `json:"name,omitempty"`
	PenaltyBoxDuration int                  `json:"penalty_box_duration,omitempty"`
	Response           *ERLResponseType     `json:"response,omitempty"`
	RpsLimit           int                  `json:"rps_limit,omitempty"`
	WindowSize         ERLWindowSizeEnum    `json:"window_size,omitempty"`
}

// Updates the specified ERL
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
	// resp, err := c.PutForm(path, i, nil)
	resp, err := c.PutJSON(path, i, nil)
	if err != nil {
		return nil, err
	}

	var erl *ERL
	if err := decodeBodyMap(resp.Body, &erl); err != nil {
		return nil, err
	}

	return erl, nil
}
