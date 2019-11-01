package fastly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/ajg/form"
	"github.com/google/jsonapi"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/mitchellh/mapstructure"
)

// APIKeyEnvVar is the name of the environment variable where the Fastly API
// key should be read from.
const APIKeyEnvVar = "FASTLY_API_KEY"

// APIKeyHeader is the name of the header that contains the Fastly API key.
const APIKeyHeader = "Fastly-Key"

// DefaultEndpoint is the default endpoint for Fastly. Since Fastly does not
// support an on-premise solution, this is likely to always be the default.
const DefaultEndpoint = "https://api.fastly.com"

// RealtimeStatsEndpoint is the realtime stats endpoint for Fastly.
const RealtimeStatsEndpoint = "https://rt.fastly.com"

// ProjectURL is the url for this library.
var ProjectURL = "github.com/fastly/go-fastly"

// ProjectVersion is the version of this library.
var ProjectVersion = "1.3.0"

// UserAgent is the user agent for this particular client.
var UserAgent = fmt.Sprintf("FastlyGo/%s (+%s; %s)",
	ProjectVersion, ProjectURL, runtime.Version())

type ACLClient interface {
	ListACLs(i *ListACLsInput) ([]*ACL, error)
	CreateACL(i *CreateACLInput) (*ACL, error)
	DeleteACL(i *DeleteACLInput) error
	GetACL(i *GetACLInput) (*ACL, error)
	UpdateACL(i *UpdateACLInput) (*ACL, error)
	ListACLEntries(i *ListACLEntriesInput) ([]*ACLEntry, error)
	GetACLEntry(i *GetACLEntryInput) (*ACLEntry, error)
	CreateACLEntry(i *CreateACLEntryInput) (*ACLEntry, error)
	DeleteACLEntry(i *DeleteACLEntryInput) error
	UpdateACLEntry(i *UpdateACLEntryInput) (*ACLEntry, error)
	BatchModifyACLEntries(i *BatchModifyACLEntriesInput) error
}

type BackendClient interface {
	ListBackends(i *ListBackendsInput) ([]*Backend, error)
	CreateBackend(i *CreateBackendInput) (*Backend, error)
	GetBackend(i *GetBackendInput) (*Backend, error)
	UpdateBackend(i *UpdateBackendInput) (*Backend, error)
	DeleteBackend(i *DeleteBackendInput) error
}

type LoggingClient interface {
	ListBigQueries(i *ListBigQueriesInput) ([]*BigQuery, error)
	CreateBigQuery(i *CreateBigQueryInput) (*BigQuery, error)
	GetBigQuery(i *GetBigQueryInput) (*BigQuery, error)
	UpdateBigQuery(i *UpdateBigQueryInput) (*BigQuery, error)
	DeleteBigQuery(i *DeleteBigQueryInput) error
	ListFTPs(i *ListFTPsInput) ([]*FTP, error)
	CreateFTP(i *CreateFTPInput) (*FTP, error)
	GetFTP(i *GetFTPInput) (*FTP, error)
	UpdateFTP(i *UpdateFTPInput) (*FTP, error)
	DeleteFTP(i *DeleteFTPInput) error
	ListGCSs(i *ListGCSsInput) ([]*GCS, error)
	CreateGCS(i *CreateGCSInput) (*GCS, error)
	GetGCS(i *GetGCSInput) (*GCS, error)
	UpdateGCS(i *UpdateGCSInput) (*GCS, error)
	DeleteGCS(i *DeleteGCSInput) error
	ListBlobStorages(i *ListBlobStoragesInput) ([]*BlobStorage, error)
	CreateBlobStorage(i *CreateBlobStorageInput) (*BlobStorage, error)
	GetBlobStorage(i *GetBlobStorageInput) (*BlobStorage, error)
	UpdateBlobStorage(i *UpdateBlobStorageInput) (*BlobStorage, error)
	DeleteBlobStorage(i *DeleteBlobStorageInput) error
	ListLogentries(i *ListLogentriesInput) ([]*Logentries, error)
	CreateLogentries(i *CreateLogentriesInput) (*Logentries, error)
	GetLogentries(i *GetLogentriesInput) (*Logentries, error)
	UpdateLogentries(i *UpdateLogentriesInput) (*Logentries, error)
	DeleteLogentries(i *DeleteLogentriesInput) error
	ListPapertrails(i *ListPapertrailsInput) ([]*Papertrail, error)
	CreatePapertrail(i *CreatePapertrailInput) (*Papertrail, error)
	GetPapertrail(i *GetPapertrailInput) (*Papertrail, error)
	UpdatePapertrail(i *UpdatePapertrailInput) (*Papertrail, error)
	DeletePapertrail(i *DeletePapertrailInput) error
	ListS3s(i *ListS3sInput) ([]*S3, error)
	CreateS3(i *CreateS3Input) (*S3, error)
	GetS3(i *GetS3Input) (*S3, error)
	UpdateS3(i *UpdateS3Input) (*S3, error)
	DeleteS3(i *DeleteS3Input) error
	ListSplunks(i *ListSplunksInput) ([]*Splunk, error)
	CreateSplunk(i *CreateSplunkInput) (*Splunk, error)
	GetSplunk(i *GetSplunkInput) (*Splunk, error)
	UpdateSplunk(i *UpdateSplunkInput) (*Splunk, error)
	DeleteSplunk(i *DeleteSplunkInput) error
	ListSumologics(i *ListSumologicsInput) ([]*Sumologic, error)
	CreateSumologic(i *CreateSumologicInput) (*Sumologic, error)
	GetSumologic(i *GetSumologicInput) (*Sumologic, error)
	UpdateSumologic(i *UpdateSumologicInput) (*Sumologic, error)
	DeleteSumologic(i *DeleteSumologicInput) error
	ListSyslogs(i *ListSyslogsInput) ([]*Syslog, error)
	CreateSyslog(i *CreateSyslogInput) (*Syslog, error)
	GetSyslog(i *GetSyslogInput) (*Syslog, error)
	UpdateSyslog(i *UpdateSyslogInput) (*Syslog, error)
	DeleteSyslog(i *DeleteSyslogInput) error
}

type GzipClient interface {
	ListGzips(i *ListGzipsInput) ([]*Gzip, error)
	CreateGzip(i *CreateGzipInput) (*Gzip, error)
	GetGzip(i *GetGzipInput) (*Gzip, error)
	UpdateGzip(i *UpdateGzipInput) (*Gzip, error)
	DeleteGzip(i *DeleteGzipInput) error
}

type CacheSettingsClient interface {
	ListCacheSettings(i *ListCacheSettingsInput) ([]*CacheSetting, error)
	CreateCacheSetting(i *CreateCacheSettingInput) (*CacheSetting, error)
	GetCacheSetting(i *GetCacheSettingInput) (*CacheSetting, error)
	UpdateCacheSetting(i *UpdateCacheSettingInput) (*CacheSetting, error)
	DeleteCacheSetting(i *DeleteCacheSettingInput) error
}

type ConditionsClient interface {
	ListConditions(i *ListConditionsInput) ([]*Condition, error)
	CreateCondition(i *CreateConditionInput) (*Condition, error)
	GetCondition(i *GetConditionInput) (*Condition, error)
	UpdateCondition(i *UpdateConditionInput) (*Condition, error)
	DeleteCondition(i *DeleteConditionInput) error
}

type AccountClient interface {
	GetBilling(i *GetBillingInput) (*Billing, error)
	GetAPIEvents(i *GetAPIEventsFilterInput) (GetAPIEventsResponse, error)
	GetAPIEvent(i *GetAPIEventInput) (*Event, error)
}

type UtilityClient interface {
	EdgeCheck(i *EdgeCheckInput) ([]*EdgeCheck, error)
}

type DiffClient interface {
	GetDiff(i *GetDiffInput) (*Diff, error)
}

type DirectorClient interface {
	ListDirectors(i *ListDirectorsInput) ([]*Director, error)
	CreateDirector(i *CreateDirectorInput) (*Director, error)
	GetDirector(i *GetDirectorInput) (*Director, error)
	UpdateDirector(i *UpdateDirectorInput) (*Director, error)
	DeleteDirector(i *DeleteDirectorInput) error
	CreateDirectorBackend(i *CreateDirectorBackendInput) (*DirectorBackend, error)
	GetDirectorBackend(i *GetDirectorBackendInput) (*DirectorBackend, error)
	DeleteDirectorBackend(i *DeleteDirectorBackendInput) error
}

type DomainClient interface {
	ListDomains(i *ListDomainsInput) ([]*Domain, error)
	CreateDomain(i *CreateDomainInput) (*Domain, error)
	GetDomain(i *GetDomainInput) (*Domain, error)
	UpdateDomain(i *UpdateDomainInput) (*Domain, error)
	DeleteDomain(i *DeleteDomainInput) error
}

type HeaderClient interface {
	ListHeaders(i *ListHeadersInput) ([]*Header, error)
	CreateHeader(i *CreateHeaderInput) (*Header, error)
	GetHeader(i *GetHeaderInput) (*Header, error)
	UpdateHeader(i *UpdateHeaderInput) (*Header, error)
	DeleteHeader(i *DeleteHeaderInput) error
}

type HealthCheckClient interface {
	ListHealthChecks(i *ListHealthChecksInput) ([]*HealthCheck, error)
	CreateHealthCheck(i *CreateHealthCheckInput) (*HealthCheck, error)
	GetHealthCheck(i *GetHealthCheckInput) (*HealthCheck, error)
	UpdateHealthCheck(i *UpdateHealthCheckInput) (*HealthCheck, error)
	DeleteHealthCheck(i *DeleteHealthCheckInput) error
}

type PurgeClient interface {
	Purge(i *PurgeInput) (*Purge, error)
	PurgeKey(i *PurgeKeyInput) (*Purge, error)
	PurgeAll(i *PurgeAllInput) (*Purge, error)
}

type RequestSettingsClient interface {
	ListRequestSettings(i *ListRequestSettingsInput) ([]*RequestSetting, error)
	CreateRequestSetting(i *CreateRequestSettingInput) (*RequestSetting, error)
	GetRequestSetting(i *GetRequestSettingInput) (*RequestSetting, error)
	UpdateRequestSetting(i *UpdateRequestSettingInput) (*RequestSetting, error)
	DeleteRequestSetting(i *DeleteRequestSettingInput) error
}

type ResponseObjectClient interface {
	ListResponseObjects(i *ListResponseObjectsInput) ([]*ResponseObject, error)
	CreateResponseObject(i *CreateResponseObjectInput) (*ResponseObject, error)
	GetResponseObject(i *GetResponseObjectInput) (*ResponseObject, error)
	UpdateResponseObject(i *UpdateResponseObjectInput) (*ResponseObject, error)
	DeleteResponseObject(i *DeleteResponseObjectInput) error
}

type ServiceClient interface {
	ListServices(i *ListServicesInput) ([]*Service, error)
	CreateService(i *CreateServiceInput) (*Service, error)
	GetService(i *GetServiceInput) (*Service, error)
	GetServiceDetails(i *GetServiceInput) (*ServiceDetail, error)
	UpdateService(i *UpdateServiceInput) (*Service, error)
	DeleteService(i *DeleteServiceInput) error
	SearchService(i *SearchServiceInput) (*Service, error)
	ListServiceDomains(i *ListServiceDomainInput) (ServiceDomainsList, error)
}

type SettingsClient interface {
	GetSettings(i *GetSettingsInput) (*Settings, error)
	UpdateSettings(i *UpdateSettingsInput) (*Settings, error)
}

type StatsClient interface {
	GetStats(i *GetStatsInput) (*StatsResponse, error)
	GetUsage(i *GetUsageInput) (*UsageResponse, error)
	GetUsageByService(i *GetUsageInput) (*UsageByServiceResponse, error)
	GetRegions() (*RegionsResponse, error)
}

type AuthenticationClient interface {
	ListTokens() ([]*Token, error)
	ListCustomerTokens(i *ListCustomerTokensInput) ([]*Token, error)
	GetTokenSelf() (*Token, error)
	CreateToken(i *CreateTokenInput) (*Token, error)
	DeleteToken(i *DeleteTokenInput) error
	DeleteTokenSelf() error
}

type VCLClient interface {
	ListVCLs(i *ListVCLsInput) ([]*VCL, error)
	GetVCL(i *GetVCLInput) (*VCL, error)
	GetGeneratedVCL(i *GetGeneratedVCLInput) (*VCL, error)
	CreateVCL(i *CreateVCLInput) (*VCL, error)
	UpdateVCL(i *UpdateVCLInput) (*VCL, error)
	ActivateVCL(i *ActivateVCLInput) (*VCL, error)
	DeleteVCL(i *DeleteVCLInput) error
}

type SnippetClient interface {
	CreateSnippet(i *CreateSnippetInput) (*Snippet, error)
	UpdateSnippet(i *UpdateSnippetInput) (*Snippet, error)
	UpdateDynamicSnippet(i *UpdateDynamicSnippetInput) (*DynamicSnippet, error)
	DeleteSnippet(i *DeleteSnippetInput) error
	ListSnippets(i *ListSnippetsInput) ([]*Snippet, error)
	GetSnippet(i *GetSnippetInput) (*Snippet, error)
	GetDynamicSnippet(i *GetDynamicSnippetInput) (*DynamicSnippet, error)
}

type ConfigVersionClient interface {
	ListVersions(i *ListVersionsInput) ([]*Version, error)
	LatestVersion(i *LatestVersionInput) (*Version, error)
	CreateVersion(i *CreateVersionInput) (*Version, error)
	GetVersion(i *GetVersionInput) (*Version, error)
	UpdateVersion(i *UpdateVersionInput) (*Version, error)
	ActivateVersion(i *ActivateVersionInput) (*Version, error)
	DeactivateVersion(i *DeactivateVersionInput) (*Version, error)
	CloneVersion(i *CloneVersionInput) (*Version, error)
	ValidateVersion(i *ValidateVersionInput) (bool, string, error)
	LockVersion(i *LockVersionInput) (*Version, error)
}

type WafClient interface {
	ListWAFs(i *ListWAFsInput) ([]*WAF, error)
	CreateWAF(i *CreateWAFInput) (*WAF, error)
	GetWAF(i *GetWAFInput) (*WAF, error)
	UpdateWAF(i *UpdateWAFInput) (*WAF, error)
	DeleteWAF(i *DeleteWAFInput) error
	GetOWASP(i *GetOWASPInput) (*OWASP, error)
	CreateOWASP(i *CreateOWASPInput) (*OWASP, error)
	UpdateOWASP(i *UpdateOWASPInput) (*OWASP, error)
	GetRules() ([]*Rule, error)
	GetRule(i *GetRuleInput) (*Rule, error)
	GetRuleVCL(i *GetRuleInput) (*RuleVCL, error)
	GetWAFRuleVCL(i *GetWAFRuleVCLInput) (*RuleVCL, error)
	GetWAFRuleRuleSets(i *GetWAFRuleRuleSetsInput) (*Ruleset, error)
	UpdateWAFRuleSets(i *UpdateWAFRuleRuleSetsInput) (*Ruleset, error)
	GetWAFRuleStatuses(i *GetWAFRuleStatusesInput) (GetWAFRuleStatusesResponse, error)
	GetWAFRuleStatus(i *GetWAFRuleStatusInput) (WAFRuleStatus, error)
	UpdateWAFRuleStatus(i *UpdateWAFRuleStatusInput) (WAFRuleStatus, error)
	UpdateWAFRuleTagStatus(input *UpdateWAFRuleTagStatusInput) (GetWAFRuleStatusesResponse, error)
	UpdateWAFConfigSet(i *UpdateWAFConfigSetInput) (UpdateWAFConfigSetResponse, error)
}

// Client is the main entrypoint to the Fastly golang API library.
type Client struct {
	// Address is the address of Fastly's API endpoint.
	Address string

	// HTTPClient is the HTTP client to use. If one is not provided, a default
	// client will be used.
	HTTPClient *http.Client

	// apiKey is the Fastly API key to authenticate requests.
	apiKey string

	// url is the parsed URL from Address
	url *url.URL
}

// RTSClient is the entrypoint to the Fastly's Realtime Stats API.
type RTSClient struct {
	client *Client
}

// DefaultClient instantiates a new Fastly API client. This function requires
// the environment variable `FASTLY_API_KEY` is set and contains a valid API key
// the environment variable `FASTLY_API_KEY` is set and contains a valid API key
// to authenticate with Fastly.
func DefaultClient() *Client {
	client, err := NewClient(os.Getenv(APIKeyEnvVar))
	if err != nil {
		panic(err)
	}
	return client
}

// NewClient creates a new API client with the given key and the default API
// endpoint. Because Fastly allows some requests without an API key, this
// function will not error if the API token is not supplied. Attempts to make a
// request that requires an API key will return a 403 response.
func NewClient(key string) (*Client, error) {
	return NewClientForEndpoint(key, DefaultEndpoint)
}

// NewClientForEndpoint creates a new API client with the given key and API
// endpoint. Because Fastly allows some requests without an API key, this
// function will not error if the API token is not supplied. Attempts to make a
// request that requires an API key will return a 403 response.
func NewClientForEndpoint(key string, endpoint string) (*Client, error) {
	client := &Client{apiKey: key, Address: endpoint}
	client, err := client.init()
	return client, err
}

// NewRealtimeStatsClient instantiates a new Fastly API client for the realtime stats.
// This function requires the environment variable `FASTLY_API_KEY` is set and contains
// a valid API key to authenticate with Fastly.
func NewRealtimeStatsClient() *RTSClient {
	c, err := NewClientForEndpoint(os.Getenv(APIKeyEnvVar), RealtimeStatsEndpoint)
	if err != nil {
		panic(err)
	}
	return &RTSClient{client: c}
}

func (c *Client) init() (*Client, error) {
	u, err := url.Parse(c.Address)
	if err != nil {
		return nil, err
	}
	c.url = u

	if c.HTTPClient == nil {
		c.HTTPClient = cleanhttp.DefaultClient()
	}

	return c, nil
}

// Get issues an HTTP GET request.
func (c *Client) Get(p string, ro *RequestOptions) (*http.Response, error) {
	return c.Request("GET", p, ro)
}

// Head issues an HTTP HEAD request.
func (c *Client) Head(p string, ro *RequestOptions) (*http.Response, error) {
	return c.Request("HEAD", p, ro)
}

// Patch issues an HTTP PATCH request.
func (c *Client) Patch(p string, ro *RequestOptions) (*http.Response, error) {
	return c.Request("PATCH", p, ro)
}

// PatchForm issues an HTTP PUT request with the given interface form-encoded.
func (c *Client) PatchForm(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestForm("PATCH", p, i, ro)
}

// PatchJSON issues an HTTP PUT request with the given interface json-encoded.
func (c *Client) PatchJSON(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestJSON("PATCH", p, i, ro)
}

// PatchJSONAPI issues an HTTP PUT request with the given interface json-encoded.
func (c *Client) PatchJSONAPI(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestJSONAPI("PATCH", p, i, ro)
}

// Post issues an HTTP POST request.
func (c *Client) Post(p string, ro *RequestOptions) (*http.Response, error) {
	return c.Request("POST", p, ro)
}

// PostForm issues an HTTP POST request with the given interface form-encoded.
func (c *Client) PostForm(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestForm("POST", p, i, ro)
}

// PostJSON issues an HTTP POST request with the given interface json-encoded.
func (c *Client) PostJSON(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestJSON("POST", p, i, ro)
}

// PostJSONAPI issues an HTTP POST request with the given interface json-encoded.
func (c *Client) PostJSONAPI(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestJSONAPI("POST", p, i, ro)
}

// Put issues an HTTP PUT request.
func (c *Client) Put(p string, ro *RequestOptions) (*http.Response, error) {
	return c.Request("PUT", p, ro)
}

// PutForm issues an HTTP PUT request with the given interface form-encoded.
func (c *Client) PutForm(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestForm("PUT", p, i, ro)
}

// PutJSON issues an HTTP PUT request with the given interface json-encoded.
func (c *Client) PutJSON(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestJSON("PUT", p, i, ro)
}

// PutJSONAPI issues an HTTP PUT request with the given interface json-encoded.
func (c *Client) PutJSONAPI(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestJSONAPI("PUT", p, i, ro)
}

// Delete issues an HTTP DELETE request.
func (c *Client) Delete(p string, ro *RequestOptions) (*http.Response, error) {
	return c.Request("DELETE", p, ro)
}

// Request makes an HTTP request against the HTTPClient using the given verb,
// Path, and request options.
func (c *Client) Request(verb, p string, ro *RequestOptions) (*http.Response, error) {
	req, err := c.RawRequest(verb, p, ro)
	if err != nil {
		return nil, err
	}

	resp, err := checkResp(c.HTTPClient.Do(req))
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RequestForm makes an HTTP request with the given interface being encoded as
// form data.
func (c *Client) RequestForm(verb, p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	if ro == nil {
		ro = new(RequestOptions)
	}

	if ro.Headers == nil {
		ro.Headers = make(map[string]string)
	}
	ro.Headers["Content-Type"] = "application/x-www-form-urlencoded"

	buf := new(bytes.Buffer)
	if err := form.NewEncoder(buf).KeepZeros(true).DelimitWith('|').Encode(i); err != nil {
		return nil, err
	}
	body := buf.String()

	ro.Body = strings.NewReader(body)
	ro.BodyLength = int64(len(body))

	return c.Request(verb, p, ro)
}

func (c *Client) RequestJSON(verb, p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	if ro == nil {
		ro = new(RequestOptions)
	}

	if ro.Headers == nil {
		ro.Headers = make(map[string]string)
	}
	ro.Headers["Content-Type"] = "application/json"
	ro.Headers["Accept"] = "application/json"

	body, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	ro.Body = bytes.NewReader(body)
	ro.BodyLength = int64(len(body))

	return c.Request(verb, p, ro)
}

func (c *Client) RequestJSONAPI(verb, p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	if ro == nil {
		ro = new(RequestOptions)
	}

	if ro.Headers == nil {
		ro.Headers = make(map[string]string)
	}
	ro.Headers["Content-Type"] = jsonapi.MediaType
	ro.Headers["Accept"] = jsonapi.MediaType

	var buf bytes.Buffer
	if err := jsonapi.MarshalPayload(&buf, i); err != nil {
		return nil, err
	}

	ro.Body = &buf
	ro.BodyLength = int64(buf.Len())

	return c.Request(verb, p, ro)
}

// checkResp wraps an HTTP request from the default client and verifies that the
// request was successful. A non-200 request returns an error formatted to
// included any validation problems or otherwise.
func checkResp(resp *http.Response, err error) (*http.Response, error) {
	// If the err is already there, there was an error higher up the chain, so
	// just return that.
	if err != nil {
		return resp, err
	}

	switch resp.StatusCode {
	case 200, 201, 202, 204, 205, 206:
		return resp, nil
	default:
		return resp, NewHTTPError(resp)
	}
}

// decodeJSON is used to decode an HTTP response body into an interface as JSON.
func decodeJSON(out interface{}, body io.ReadCloser) error {
	defer body.Close()

	var parsed interface{}
	dec := json.NewDecoder(body)
	if err := dec.Decode(&parsed); err != nil {
		return err
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapToHTTPHeaderHookFunc(),
			stringToTimeHookFunc(),
		),
		WeaklyTypedInput: true,
		Result:           out,
	})
	if err != nil {
		return err
	}
	return decoder.Decode(parsed)
}
