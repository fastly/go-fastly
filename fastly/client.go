package fastly

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/go-querystring/query"
	"github.com/google/jsonapi"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/mitchellh/mapstructure"
)

// APIKeyEnvVar is the name of the environment variable where the Fastly API
// key should be read from.
const APIKeyEnvVar = "FASTLY_API_KEY" // #nosec G101

// APIKeyHeader is the name of the header that contains the Fastly API key.
const APIKeyHeader = "Fastly-Key" // #nosec G101

// EndpointEnvVar is the name of an environment variable that can be used
// to change the URL of API requests.
const EndpointEnvVar = "FASTLY_API_URL"

// DebugEnvVar is the name of an environment variable that can be used to switch
// the API client into debug mode.
const DebugEnvVar = "FASTLY_DEBUG_MODE"

// DefaultEndpoint is the default endpoint for Fastly. Since Fastly does not
// support an on-premise solution, this is likely to always be the default.
const DefaultEndpoint = "https://api.fastly.com"

// RealtimeStatsEndpointEnvVar is the name of an environment variable that can be used
// to change the URL of realtime stats requests.
const RealtimeStatsEndpointEnvVar = "FASTLY_RTS_URL"

// DefaultRealtimeStatsEndpoint is the realtime stats endpoint for Fastly.
const DefaultRealtimeStatsEndpoint = "https://rt.fastly.com"

// JSONMimeType is the MIME type for the JSON data format.
const JSONMimeType = "application/json"

// UserAgentEnvVar is the name of an environment variable that can be used
// to change the User-Agent of the http requests
const UserAgentEnvVar = "FASTLY_USER_AGENT"

// ProjectURL is the url for this library.
var ProjectURL = "github.com/fastly/go-fastly"

// ProjectVersion is the version of this library.
var ProjectVersion = "9.13.1"

// UserAgent is the user agent for this particular client.
var UserAgent = fmt.Sprintf("FastlyGo/%s (+%s; %s)",
	ProjectVersion, ProjectURL, runtime.Version())

// Client is the main entrypoint to the Fastly golang API library.
type Client struct {
	// Address is the address of Fastly's API endpoint.
	Address string
	// DebugMode enables HTTP request/response dumps.
	DebugMode bool
	// HTTPClient is the HTTP client to use. If one is not provided, a default
	// client will be used.
	HTTPClient *http.Client

	// apiKey is the Fastly API key to authenticate requests.
	apiKey string
	// remaining is last observed value of http header Fastly-RateLimit-Remaining
	remaining int
	// reset is last observed value of http header Fastly-RateLimit-Reset
	reset int64
	// updateLock forces serialization of calls that modify a service.
	// Concurrent modifications have undefined semantics.
	updateLock sync.Mutex
	// url is the parsed URL from Address
	url *url.URL
}

// RTSClient is the entrypoint to the Fastly's Realtime Stats API.
type RTSClient struct {
	client *Client
}

// DefaultClient instantiates a new Fastly API client. This function requires
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
	endpoint, ok := os.LookupEnv(EndpointEnvVar)

	if !ok {
		endpoint = DefaultEndpoint
	}

	return NewClientForEndpoint(key, endpoint)
}

// NewClientForEndpoint creates a new API client with the given key and API
// endpoint. Because Fastly allows some requests without an API key, this
// function will not error if the API token is not supplied. Attempts to make a
// request that requires an API key will return a 403 response.
func NewClientForEndpoint(key, endpoint string) (*Client, error) {
	client := &Client{apiKey: key, Address: endpoint}

	if endpoint, ok := os.LookupEnv(DebugEnvVar); ok && endpoint == "true" {
		client.DebugMode = true
	}

	if customUserAgent, ok := os.LookupEnv(UserAgentEnvVar); ok {
		UserAgent = fmt.Sprintf("%s, %s", customUserAgent, UserAgent)
	}

	return client.init()
}

// NewRealtimeStatsClient instantiates a new Fastly API client for the realtime stats.
// This function requires the environment variable `FASTLY_API_KEY` is set and contains
// a valid API key to authenticate with Fastly.
func NewRealtimeStatsClient() *RTSClient {
	endpoint, ok := os.LookupEnv(RealtimeStatsEndpointEnvVar)

	if !ok {
		endpoint = DefaultRealtimeStatsEndpoint
	}

	c, err := NewClientForEndpoint(os.Getenv(APIKeyEnvVar), endpoint)
	if err != nil {
		panic(err)
	}
	return &RTSClient{client: c}
}

// NewRealtimeStatsClientForEndpoint creates an RTSClient from a token and endpoint url.
// `token` is a Fastly API token and `endpoint` is RealtimeStatsEndpoint for the production
// realtime stats API.
func NewRealtimeStatsClientForEndpoint(token, endpoint string) (*RTSClient, error) {
	c, err := NewClientForEndpoint(token, endpoint)
	if err != nil {
		return nil, err
	}
	return &RTSClient{client: c}, nil
}

func (c *Client) init() (*Client, error) {
	// Until we do a request, we don't know how many are left.
	// Use the default limit as a first guess:
	// https://developer.fastly.com/reference/api/#rate-limiting
	c.remaining = 1000

	u, err := url.Parse(c.Address)
	if err != nil {
		return nil, err
	}
	c.url = u

	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{
			// IMPORTANT: Avoid cleanhttp.DefaultTransport() which disables keepalive.
			Transport: cleanhttp.DefaultPooledTransport(),
		}
	}

	return c, nil
}

// RateLimitRemaining returns the number of non-read requests left before
// rate limiting causes a 429 Too Many Requests error.
func (c *Client) RateLimitRemaining() int {
	return c.remaining
}

// RateLimitReset returns the next time the rate limiter's counter will be
// reset.
func (c *Client) RateLimitReset() time.Time {
	return time.Unix(c.reset, 0)
}

// Get issues an HTTP GET request.
func (c *Client) Get(p string, ro *RequestOptions) (*http.Response, error) {
	if ro == nil {
		ro = new(RequestOptions)
	}
	ro.Parallel = true
	return c.Request(http.MethodGet, p, ro)
}

// GetJSON issues an HTTP GET request and indicates that the response
// should be JSON encoded.
func (c *Client) GetJSON(p string, ro *RequestOptions) (*http.Response, error) {
	if ro == nil {
		ro = new(RequestOptions)
	}
	if ro.Headers == nil {
		ro.Headers = make(map[string]string)
	}
	ro.Parallel = true
	ro.Headers["Accept"] = JSONMimeType
	return c.Request(http.MethodGet, p, ro)
}

// Head issues an HTTP HEAD request.
func (c *Client) Head(p string, ro *RequestOptions) (*http.Response, error) {
	if ro == nil {
		ro = new(RequestOptions)
	}
	ro.Parallel = true
	return c.Request(http.MethodHead, p, ro)
}

// Patch issues an HTTP PATCH request.
func (c *Client) Patch(p string, ro *RequestOptions) (*http.Response, error) {
	return c.Request(http.MethodPatch, p, ro)
}

// PatchForm issues an HTTP PUT request with the given interface form-encoded.
func (c *Client) PatchForm(p string, i any, ro *RequestOptions) (*http.Response, error) {
	return c.RequestForm(http.MethodPatch, p, i, ro)
}

// PatchJSON issues an HTTP PUT request with the given interface json-encoded.
func (c *Client) PatchJSON(p string, i any, ro *RequestOptions) (*http.Response, error) {
	return c.RequestJSON(http.MethodPatch, p, i, ro)
}

// PatchJSONAPI issues an HTTP PUT request with the given interface json-encoded.
func (c *Client) PatchJSONAPI(p string, i any, ro *RequestOptions) (*http.Response, error) {
	return c.RequestJSONAPI(http.MethodPatch, p, i, ro)
}

// Post issues an HTTP POST request.
func (c *Client) Post(p string, ro *RequestOptions) (*http.Response, error) {
	return c.Request(http.MethodPost, p, ro)
}

// PostForm issues an HTTP POST request with the given interface form-encoded.
func (c *Client) PostForm(p string, i any, ro *RequestOptions) (*http.Response, error) {
	return c.RequestForm(http.MethodPost, p, i, ro)
}

// PostJSON issues an HTTP POST request with the given interface json-encoded.
func (c *Client) PostJSON(p string, i any, ro *RequestOptions) (*http.Response, error) {
	return c.RequestJSON(http.MethodPost, p, i, ro)
}

// PostJSONAPI issues an HTTP POST request with the given interface json-encoded.
func (c *Client) PostJSONAPI(p string, i any, ro *RequestOptions) (*http.Response, error) {
	return c.RequestJSONAPI(http.MethodPost, p, i, ro)
}

// PostJSONAPIBulk issues an HTTP POST request with the given interface json-encoded and bulk requests.
func (c *Client) PostJSONAPIBulk(p string, i any, ro *RequestOptions) (*http.Response, error) {
	return c.RequestJSONAPIBulk(http.MethodPost, p, i, ro)
}

// Put issues an HTTP PUT request.
func (c *Client) Put(p string, ro *RequestOptions) (*http.Response, error) {
	return c.Request(http.MethodPut, p, ro)
}

// PutForm issues an HTTP PUT request with the given interface form-encoded.
func (c *Client) PutForm(p string, i any, ro *RequestOptions) (*http.Response, error) {
	return c.RequestForm(http.MethodPut, p, i, ro)
}

// PutFormFile issues an HTTP PUT request (multipart/form-encoded) to put a file to an endpoint.
func (c *Client) PutFormFile(urlPath, filePath, fieldName string, ro *RequestOptions) (*http.Response, error) {
	return c.RequestFormFile(http.MethodPut, urlPath, filePath, fieldName, ro)
}

// PutFormFileFromReader issues an HTTP PUT request (multipart/form-encoded) to put a file to an endpoint.
func (c *Client) PutFormFileFromReader(urlPath, fileName string, fileBytes io.Reader, fieldName string, ro *RequestOptions) (*http.Response, error) {
	return c.RequestFormFileFromReader(http.MethodPut, urlPath, fileName, fileBytes, fieldName, ro)
}

// PutJSON issues an HTTP PUT request with the given interface json-encoded.
func (c *Client) PutJSON(p string, i any, ro *RequestOptions) (*http.Response, error) {
	return c.RequestJSON(http.MethodPut, p, i, ro)
}

// PutJSONAPI issues an HTTP PUT request with the given interface json-encoded.
func (c *Client) PutJSONAPI(p string, i any, ro *RequestOptions) (*http.Response, error) {
	return c.RequestJSONAPI(http.MethodPut, p, i, ro)
}

// Delete issues an HTTP DELETE request.
func (c *Client) Delete(p string, ro *RequestOptions) (*http.Response, error) {
	return c.Request(http.MethodDelete, p, ro)
}

// DeleteJSONAPI issues an HTTP DELETE request with the given interface json-encoded.
func (c *Client) DeleteJSONAPI(p string, i any, ro *RequestOptions) (*http.Response, error) {
	return c.RequestJSONAPI(http.MethodDelete, p, i, ro)
}

// DeleteJSONAPIBulk issues an HTTP DELETE request with the given interface json-encoded and bulk requests.
func (c *Client) DeleteJSONAPIBulk(p string, i any, ro *RequestOptions) (*http.Response, error) {
	return c.RequestJSONAPIBulk(http.MethodDelete, p, i, ro)
}

// Request makes an HTTP request against the HTTPClient using the given verb,
// Path, and request options.
func (c *Client) Request(verb, p string, ro *RequestOptions) (*http.Response, error) {
	req, err := c.RawRequest(verb, p, ro)
	if err != nil {
		return nil, err
	}

	if ro == nil || !ro.Parallel {
		c.updateLock.Lock()
		defer c.updateLock.Unlock()
	}

	if c.DebugMode {
		r := req.Clone(context.Background())
		// 'r' and 'req' both have a reference to a Body that
		// is an io.Reader, but only one of them can read its
		// contents since io.Reader is not seekable and cannot
		// be rewound

		r.Header.Del(APIKeyHeader)
		dump, _ := httputil.DumpRequest(r, true)

		// httputil.DumpRequest has read the Body from 'r',
		// and set r.Body to an io.ReadCloser that will return
		// the same bytes as the original Body, so we can
		// stuff that into 'req' for when the request is
		// actually sent; it can't be read in 'r', but the
		// lifetime of 'r' ends at the end of this block
		req.Body = r.Body

		fmt.Printf("http.Request (dump): %q\n", dump)
	}

	// nosemgrep: trailofbits.go.invalid-usage-of-modified-variable.invalid-usage-of-modified-variable
	resp, err := checkResp(c.HTTPClient.Do(req))

	if c.DebugMode && resp != nil {
		dump, _ := httputil.DumpResponse(resp, true)
		fmt.Printf("http.Response (dump): %q\n", dump)
	}

	if err != nil {
		return resp, err
	}

	if verb != http.MethodGet && verb != http.MethodHead {
		remaining := resp.Header.Get("Fastly-RateLimit-Remaining")
		if remaining != "" {
			if val, err := strconv.Atoi(remaining); err == nil {
				c.remaining = val
			}
		}
		reset := resp.Header.Get("Fastly-RateLimit-Reset")
		if reset != "" {
			if val, err := strconv.ParseInt(reset, 10, 64); err == nil {
				c.reset = val
			}
		}
	}

	return resp, nil
}

// RequestOptions is the list of options to pass to the request.
type RequestOptions struct {
	// Body is an io.Reader object that will be streamed or uploaded with the
	// Request.
	Body io.Reader
	// BodyLength is the final size of the Body.
	BodyLength int64
	// Headers is a map of key-value pairs that will be added to the Request.
	Headers map[string]string
	// HealthCheckHeaders indicates if there is any special parsing required to
	// support the health check API endpoint (refer to client.RequestForm).
	//
	// TODO: Lookout for this when it comes to the future code-generated API
	// client world, as this special case might get omitted accidentally.
	HealthCheckHeaders bool
	// Can this request run in parallel
	Parallel bool
	// Params is a map of key-value pairs that will be added to the Request.
	Params map[string]string
}

// RawRequest accepts a verb, URL, and RequestOptions struct and returns the
// constructed http.Request and any errors that occurred.
func (c *Client) RawRequest(verb, p string, ro *RequestOptions) (*http.Request, error) {
	// Ensure we have request options.
	if ro == nil {
		ro = new(RequestOptions)
	}

	// Append the path to the URL.
	u := strings.TrimRight(c.url.String(), "/") + "/" + strings.TrimLeft(p, "/")

	// Create the request object.
	request, err := http.NewRequest(verb, u, ro.Body)
	if err != nil {
		return nil, err
	}

	params := make(url.Values)
	for k, v := range ro.Params {
		params.Add(k, v)
	}
	request.URL.RawQuery = params.Encode()

	// Set the API key.
	if len(c.apiKey) > 0 {
		request.Header.Set(APIKeyHeader, c.apiKey)
	}

	// Set the User-Agent.
	request.Header.Set("User-Agent", UserAgent)

	// Add any custom headers.
	for k, v := range ro.Headers {
		request.Header.Add(k, v)
	}

	// Add Content-Length if we have it.
	if ro.BodyLength > 0 {
		request.ContentLength = ro.BodyLength
	}

	return request, nil
}

// SimpleGet combines the RawRequest and Request methods,
// but doesn't add any parameters or change any encoding in the URL
// passed to it. It's mostly for calling the URLs given to us
// directly from Fastly without mangling them.
func (c *Client) SimpleGet(target string) (*http.Response, error) {
	// We parse the URL and then convert it right back to a string
	// later; this just acts as a check that Fastly isn't sending
	// us nonsense.
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	if len(c.apiKey) > 0 {
		request.Header.Set(APIKeyHeader, c.apiKey)
	}
	request.Header.Set("User-Agent", UserAgent)

	// nosemgrep: trailofbits.go.invalid-usage-of-modified-variable.invalid-usage-of-modified-variable
	return checkResp(c.HTTPClient.Do(request))
}

// parseHealthCheckHeaders returns the serialised body with the custom health
// check headers appended.
//
// NOTE: The Google query library we use for parsing and encoding the provided
// struct values doesn't support the format `headers=["Foo: Bar"]` and so we
// have to manually construct this format.
func parseHealthCheckHeaders(s string) string {
	headers := []string{}
	result := []string{}
	segs := strings.Split(s, "&")
	for _, s := range segs {
		if strings.HasPrefix(strings.ToLower(s), "headers=") {
			v := strings.Split(s, "=")
			if len(v) == 2 {
				headers = append(headers, fmt.Sprintf("%q", strings.ReplaceAll(v[1], "%3A+", ":")))
			}
		} else {
			result = append(result, s)
		}
	}
	if len(headers) > 0 {
		result = append(result, "headers=%5B"+strings.Join(headers, ",")+"%5D")
	}
	return strings.Join(result, "&")
}

// RequestForm makes an HTTP request with the given interface being encoded as
// form data.
func (c *Client) RequestForm(verb, p string, i any, ro *RequestOptions) (*http.Response, error) {
	if ro == nil {
		ro = new(RequestOptions)
	}

	if ro.Headers == nil {
		ro.Headers = make(map[string]string)
	}
	ro.Headers["Content-Type"] = "application/x-www-form-urlencoded"

	v, err := query.Values(i)
	if err != nil {
		return nil, err
	}

	body := v.Encode()
	if ro.HealthCheckHeaders {
		body = parseHealthCheckHeaders(body)
	}

	ro.Body = strings.NewReader(body)
	ro.BodyLength = int64(len(body))

	return c.Request(verb, p, ro)
}

// RequestFormFile makes an HTTP request to upload a file to an endpoint.
func (c *Client) RequestFormFile(verb, urlPath, filePath, fieldName string, ro *RequestOptions) (*http.Response, error) {
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}
	defer file.Close() // #nosec G307

	return c.RequestFormFileFromReader(verb, urlPath, filepath.Base(filePath), file, fieldName, ro)
}

// RequestFormFileFromReader makes an HTTP request to upload a raw reader to an endpoint.
func (c *Client) RequestFormFileFromReader(verb, urlPath, fileName string, fileBytes io.Reader, fieldName string, ro *RequestOptions) (*http.Response, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		return nil, fmt.Errorf("error creating multipart form: %w", err)
	}

	_, err = io.Copy(part, fileBytes)
	if err != nil {
		return nil, fmt.Errorf("error copying file to multipart form: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("error closing multipart form: %w", err)
	}

	if ro == nil {
		ro = new(RequestOptions)
	}
	if ro.Headers == nil {
		ro.Headers = make(map[string]string)
	}
	ro.Headers["Content-Type"] = writer.FormDataContentType()
	ro.Headers["Accept"] = JSONMimeType
	ro.Body = &body
	ro.BodyLength = int64(body.Len())

	return c.Request(verb, urlPath, ro)
}

// setHeaders ensures that RequestOptions has headers set and applies the given content type and accept headers.
func setHeaders(ro *RequestOptions, contentType, accept string) *RequestOptions {
	if ro == nil {
		ro = new(RequestOptions)
	}
	if ro.Headers == nil {
		ro.Headers = make(map[string]string)
	}
	ro.Headers["Content-Type"] = contentType
	ro.Headers["Accept"] = accept
	return ro
}

// marshalBodyJSON marshals the input into JSON format and sets it in RequestOptions.
func marshalBodyJSON(i any, ro *RequestOptions) error {
	if i == nil {
		return nil
	}
	body, err := json.Marshal(i)
	if err != nil {
		return err
	}
	ro.Body = bytes.NewReader(body)
	ro.BodyLength = int64(len(body))
	return nil
}

// marshalBodyJSONAPI marshals the input into JSON API format and sets it in RequestOptions.
func marshalBodyJSONAPI(i any, ro *RequestOptions) error {
	if i == nil {
		return nil
	}
	var buf bytes.Buffer
	if err := jsonapi.MarshalPayload(&buf, i); err != nil {
		return err
	}
	ro.Body = &buf
	ro.BodyLength = int64(buf.Len())
	return nil
}

// RequestJSON constructs a JSON HTTP request.
func (c *Client) RequestJSON(verb, p string, i any, ro *RequestOptions) (*http.Response, error) {
	ro = setHeaders(ro, JSONMimeType, JSONMimeType)
	if err := marshalBodyJSON(i, ro); err != nil {
		return nil, err
	}
	return c.Request(verb, p, ro)
}

// RequestJSONAPI constructs a JSON API HTTP request.
func (c *Client) RequestJSONAPI(verb, p string, i any, ro *RequestOptions) (*http.Response, error) {
	ro = setHeaders(ro, jsonapi.MediaType, jsonapi.MediaType)
	if err := marshalBodyJSONAPI(i, ro); err != nil {
		return nil, err
	}
	return c.Request(verb, p, ro)
}

// RequestJSONAPIBulk constructs a bulk JSON API HTTP request.
func (c *Client) RequestJSONAPIBulk(verb, p string, i any, ro *RequestOptions) (*http.Response, error) {
	ro = setHeaders(ro, jsonapi.MediaType+"; ext=bulk", jsonapi.MediaType+"; ext=bulk")
	if err := marshalBodyJSONAPI(i, ro); err != nil {
		return nil, err
	}
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
	case http.StatusOK, http.StatusCreated, http.StatusAccepted, http.StatusNoContent, http.StatusResetContent, http.StatusPartialContent:
		return resp, nil
	default:
		return resp, NewHTTPError(resp)
	}
}

// DecodeBodyMap is used to decode an HTTP response body into a mapstructure struct.
func DecodeBodyMap(body io.Reader, out any) error {
	var parsed any
	dec := json.NewDecoder(body)
	if err := dec.Decode(&parsed); err != nil {
		return err
	}

	return decodeMap(parsed, out)
}

// decodeMap decodes an `in` struct or map to a mapstructure tagged `out`.
// It applies the decoder defaults used throughout go-fastly.
// Note that this uses opposite argument order from Go's copy().
func decodeMap(in, out any) error {
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
	return decoder.Decode(in)
}

// mapToHTTPHeaderHookFunc returns a function that converts maps into an
// http.Header value.
func mapToHTTPHeaderHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data any,
	) (any, error) {
		// Ensure source is a map and target is http.Header
		if f.Kind() != reflect.Map || t != reflect.TypeOf(new(http.Header)) {
			return data, nil
		}

		typed, ok := data.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("cannot convert %T to http.Header", data)
		}

		n := map[string][]string{}
		for k, v := range typed {
			switch tv := v.(type) {
			case string:
				n[k] = []string{tv}
			case []string:
				n[k] = tv
			case int, int8, int16, int32, int64:
				n[k] = []string{fmt.Sprintf("%d", tv)}
			case float32, float64:
				n[k] = []string{fmt.Sprintf("%f", tv)}
			default:
				return nil, fmt.Errorf("cannot convert %T to http.Header", v)
			}
		}

		return n, nil
	}
}

// stringToTimeHookFunc returns a function that converts strings to a time.Time
// value.
func stringToTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data any,
	) (any, error) {
		// Ensure the source is a string and the target is time.Time
		if f.Kind() != reflect.String || t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		// Attempt parsing in RFC3339 format
		str, _ := data.(string) // Safe type assertion; guaranteed by f.Kind()
		if v, err := time.Parse(time.RFC3339, str); err == nil {
			return v, nil
		}

		// Fallback to "2006-01-02 15:04:05" format
		if v, err := time.Parse("2006-01-02 15:04:05", str); err == nil {
			// DictionaryInfo#get uses it's own special time format for now.
			return v, nil
		}

		return nil, fmt.Errorf("unable to parse time string: %q", str)
	}
}
