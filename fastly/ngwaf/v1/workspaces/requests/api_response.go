package requests

import "time"

// Request is the API response to requests and request operations.
type Request struct {
	// AgentResponseCode is the Agent response code.
	AgentResponseCode int `json:"agent_response_code"`
	// Country is the origin country of the request.
	Country string `json:"country"`
	// ID is a base62-encoded representation of a UUID used to
	// uniquely identify a request.
	ID string `json:"id"`
	// Method is the HTTP method of the request.
	Method string `json:"method"`
	// Path is the request path of the request.
	Path string `json:"path"`
	// Protocol is the HTTP protocol of the request.
	Protocol string `json:"protocol"`
	// RemoteHostname is the Remote hostname of the request.
	RemoteHostname string `json:"remote_hostname"`
	// RemoteIPAddress is the remote IP address of the request.
	RemoteIPAddress string `json:"remote_ip"`
	// RequestHeaders are the request's headers.
	RequestHeaders []Header `json:"request_headers"`
	// ResponseCode is the response code of the request.
	ResponseCode int `json:"response_code"`
	// ResponseHeaders are the response's headers.
	ResponseHeaders []Header `json:"response_headers"`
	// ResponseSize is the HTTP response size.
	ResponseSize int `json:"response_size"`
	// ResponseTime is the HTTP response time in milliseconds.
	ResponseTime int `json:"response_time"`
	// Scheme is the request scheme.
	Scheme string `json:"scheme"`
	// ServerHostname is the hostname of the server.
	ServerHostname string `json:"server_hostname"`
	// ServerName is the server name.
	ServerName string `json:"server_name"`
	// Signals is the list of signals the request matched.
	Signals []Signal `json:"signals"`
	// Timestamp is the time when the request was made.
	Timestamp time.Time `json:"timestamp"`
	// TLSCipher is the TLS cipher of the request.
	TLSCipher string `json:"tls_cipher"`
	// TLSProtocol is the TLS protocol of the request.
	TLSProtocol string `json:"tls_protocol"`
	// URI is the request URI.
	URI string `json:"uri"`
	// UserAgent is the request's user agent.
	UserAgent string `json:"user_agent"`
}

// Header is a struct for request headers.
type Header struct {
	// Name is the name of the header.
	Name string `json:"name"`
	// Value is the value of the header.
	Value string `json:"value"`
}

// Signal is a struct for NGWAF signals within the scope of a request.
type Signal struct {
	// Detector is the dector of the signal.
	Detector string `json:"detector"`
	// ID is the signal ID.
	ID string `json:"id"`
	// Location is the location in the request that triggered the
	// signal
	Location string `json:"location"`
	// Value is the signal value
	Value string `json:"value"`
}

// Requests is the API response structure for the list requests
// operation.
type Requests struct {
	// Data is the list of returned workspaces.
	Data []Request `json:"data"`
	// Meta is the information for total workspaces.
	Meta MetaRequests `json:"meta"`
}

// MetaRequests is a subset of the Requests response structure.
type MetaRequests struct {
	// Limit is the limit of requests.
	Limit int `json:"limit"`
	// NextCursor is the next set of requests if paginated
	NextCursor string `json:"next_cursor"`
	// Total is the sum of requests.
	Total int `json:"total"`
}
