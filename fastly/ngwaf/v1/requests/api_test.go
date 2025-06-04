package requests

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/fastly/go-fastly/v10/fastly"
)

const (
	TestRequestID = "1fb963ef0bae42889b00000000000001"
)

func TestClient_requests(t *testing.T) {
	t.Parallel()

	getrequestInput := new(GetInput)
	getrequestInput.RequestID = fastly.ToPointer(string(TestRequestID))
	getrequestInput.WorkspaceID = fastly.ToPointer(fastly.TestNGWAFEventsAndRequestsWorkspaceID)

	var request *Request
	var err error
	timestamp, _ := time.Parse(time.RFC3339, "2025-06-03T20:47:50Z")
	testRequest := Request{
		AgentResponseCode: 200,
		Country:           "US",
		ID:                "1fb963ef0bae42889b00000000000001",
		Method:            "POST",
		Path:              "/",
		Protocol:          "HTTP/1.1",
		RemoteHostname:    "pool-96-224-50-187.nycmny.fios.verizon.net",
		RemoteIPAddress:   "96.224.50.187",
		RequestHeaders: []Header{
			{
				Name:  "Content-Type",
				Value: "application/x-www-form-urlencoded",
			},
			{
				Name:  "User-Agent",
				Value: "curl/8.7.1",
			},
			{
				Name:  "X-Sigsci-Serviceid-Prod",
				Value: "5bHVNvR3QwYBtF59iFkL72",
			},
			{
				Name:  "X-Timer",
				Value: "S1748983670.229711, VS0",
			},
			{
				Name:  "Content-Length",
				Value: "23",
			},
			{
				Name:  "X-Sigsci-Client-Geo-Country-Code",
				Value: "US",
			},
			{
				Name:  "X-Sigsci-Edgemodule",
				Value: "vcl 3.1.0",
			},
			{
				Name:  "Accept",
				Value: "*/*",
			},
			{
				Name:  "Cdn-Loop",
				Value: "Fastly",
			},
			{
				Name:  "Fastly-Client-Ip",
				Value: "96.224.50.187",
			},
			{
				Name:  "Fastly-Ff",
				Value: "jjtNGp3COE7Cf2N9cZIJxo8uJocXdwD8AsyHr2Ir3HA=!NYC!cache-nyc-kteb1890034-NYC",
			},
			{
				Name:  "X-Sigsci-Client-Geo-City",
				Value: "ossining",
			},
			{
				Name:  "X-Sigsci-Requestid",
				Value: "1fb963ef0bae42889b00000000000001",
			},
			{
				Name:  "Host",
				Value: "ff3eb80ddda467d3d9fc2c6cbfb9fb95.com",
			},
			{
				Name:  "X-Uat-Ip",
				Value: "10.2.3.4",
			},
			{
				Name:  "X-Varnish",
				Value: "1146044231",
			},
		},
		ResponseCode: 503,
		ResponseHeaders: []Header{
			{
				Name:  "Connection",
				Value: "keep-alive",
			},
			{
				Name:  "Date",
				Value: "Tue, 03 Jun 2025 20:47:49 GMT",
			},
			{
				Name:  "Content-Type",
				Value: "text/plain",
			},
			{
				Name:  "Transfer-Encoding",
				Value: "chunked",
			},
			{
				Name:  "X-Cache",
				Value: "MISS",
			},
			{
				Name:  "X-Cache-Hits",
				Value: "0",
			},
			{
				Name:  "X-Served-By",
				Value: "cache-nyc-kteb1890043-NYC",
			},
		},
		ResponseSize:   19,
		ResponseTime:   7,
		Scheme:         "http",
		ServerHostname: "ff3eb80ddda467d3d9fc2c6cbfb9fb95.com",
		ServerName:     "ff3eb80ddda467d3d9fc2c6cbfb9fb95.com",
		Signals: []Signal{
			{
				ID:       "CMDEXE",
				Location: "POST",
				Value:    "bar200=;cat /etc/passwd",
				Detector: "CmdExeRule",
			},
			{
				ID:       "HTTP503",
				Location: "",
				Value:    "503",
				Detector: "HTTPErrorRule",
			},
		},
		Timestamp:   timestamp,
		TLSCipher:   "",
		TLSProtocol: "",
		URI:         "/",
		UserAgent:   "curl/8.7.1",
	}

	// get an request
	fastly.Record(t, "get_request", func(c *fastly.Client) {
		request, err = Get(c, getrequestInput)
	})
	if err != nil {
		t.Fatal(err)
	}
	if request.AgentResponseCode != testRequest.AgentResponseCode {
		t.Errorf("unexpected request AgentResponseCode: got %q, expected %q", request.AgentResponseCode, testRequest.AgentResponseCode)
	}
	if request.Country != testRequest.Country {
		t.Errorf("unexpected request Country: got %q, expected %q", request.Country, testRequest.Country)
	}
	if request.ID != testRequest.ID {
		t.Errorf("unexpected request ID: got %q, expected %q", request.ID, testRequest.ID)
	}
	if request.Method != testRequest.Method {
		t.Errorf("unexpected request Method: got %q, expected %q", request.Method, testRequest.Method)
	}
	if request.Path != testRequest.Path {
		t.Errorf("unexpected request Path: got %q, expected %q", request.Path, testRequest.Path)
	}
	if request.Protocol != testRequest.Protocol {
		t.Errorf("unexpected request Protocol: got %q, expected %q", request.Protocol, testRequest.Protocol)
	}
	if request.RemoteHostname != testRequest.RemoteHostname {
		t.Errorf("unexpected request RemoteHostname: got %q, expected %q", request.RemoteHostname, testRequest.RemoteHostname)
	}
	if request.RemoteIPAddress != testRequest.RemoteIPAddress {
		t.Errorf("unexpected request RemoteIPAddress: got %q, expected %q", request.RemoteIPAddress, testRequest.RemoteIPAddress)
	}
	assert.ElementsMatch(t, request.RequestHeaders, testRequest.RequestHeaders)
	if request.ResponseCode != testRequest.ResponseCode {
		t.Errorf("unexpected request ResponseCode: got %q, expected %q", request.ResponseCode, testRequest.ResponseCode)
	}
	assert.ElementsMatch(t, request.ResponseHeaders, testRequest.ResponseHeaders)
	if request.ResponseSize != testRequest.ResponseSize {
		t.Errorf("unexpected request ResponseSize: got %q, expected %q", request.ResponseSize, testRequest.ResponseSize)
	}
	if request.ResponseTime != testRequest.ResponseTime {
		t.Errorf("unexpected request ResponseTime: got %q, expected %q", request.ResponseTime, testRequest.ResponseTime)
	}
	if request.Scheme != testRequest.Scheme {
		t.Errorf("unexpected request Scheme: got %q, expected %q", request.Scheme, testRequest.Scheme)
	}
	if request.ServerHostname != testRequest.ServerHostname {
		t.Errorf("unexpected request ServerHostname: got %q, expected %q", request.ServerHostname, testRequest.ServerHostname)
	}
	if request.ServerName != testRequest.ServerName {
		t.Errorf("unexpected request ServerName: got %q, expected %q", request.ServerName, testRequest.ServerName)
	}
	assert.ElementsMatch(t, request.Signals, testRequest.Signals)
	if request.Timestamp != testRequest.Timestamp {
		t.Errorf("unexpected request Timestamp: got %q, expected %q", request.Timestamp, testRequest.Timestamp)
	}
	if request.TLSCipher != testRequest.TLSCipher {
		t.Errorf("unexpected request TLSCipher: got %q, expected %q", request.TLSCipher, testRequest.TLSCipher)
	}
	if request.TLSProtocol != testRequest.TLSProtocol {
		t.Errorf("unexpected request TLSProtocol: got %q, expected %q", request.TLSProtocol, testRequest.TLSProtocol)
	}
	if request.URI != testRequest.URI {
		t.Errorf("unexpected request URI: got %q, expected %q", request.URI, testRequest.URI)
	}
	if request.UserAgent != testRequest.UserAgent {
		t.Errorf("unexpected request UserAgent: got %q, expected %q", request.UserAgent, testRequest.UserAgent)
	}

	var Requests *Requests
	listRequestInput := new(ListInput)
	listRequestInput.WorkspaceID = fastly.ToPointer(fastly.TestNGWAFEventsAndRequestsWorkspaceID)
	listRequestInput.Limit = fastly.ToPointer(100)
	listRequestInput.Query = fastly.ToPointer("from:-2d")

	// get a list of requests
	fastly.Record(t, "list_request", func(c *fastly.Client) {
		Requests, err = List(c, listRequestInput)
	})
	request = &Requests.Data[0]
	if err != nil {
		t.Fatal(err)
	}
	if request.AgentResponseCode != testRequest.AgentResponseCode {
		t.Errorf("unexpected request AgentResponseCode: got %q, expected %q", request.AgentResponseCode, testRequest.AgentResponseCode)
	}
	if request.Country != testRequest.Country {
		t.Errorf("unexpected request Country: got %q, expected %q", request.Country, testRequest.Country)
	}
	if request.ID != testRequest.ID {
		t.Errorf("unexpected request ID: got %q, expected %q", request.ID, testRequest.ID)
	}
	if request.Method != testRequest.Method {
		t.Errorf("unexpected request Method: got %q, expected %q", request.Method, testRequest.Method)
	}
	if request.Path != testRequest.Path {
		t.Errorf("unexpected request Path: got %q, expected %q", request.Path, testRequest.Path)
	}
	if request.Protocol != testRequest.Protocol {
		t.Errorf("unexpected request Protocol: got %q, expected %q", request.Protocol, testRequest.Protocol)
	}
	if request.RemoteHostname != testRequest.RemoteHostname {
		t.Errorf("unexpected request RemoteHostname: got %q, expected %q", request.RemoteHostname, testRequest.RemoteHostname)
	}
	if request.RemoteIPAddress != testRequest.RemoteIPAddress {
		t.Errorf("unexpected request RemoteIPAddress: got %q, expected %q", request.RemoteIPAddress, testRequest.RemoteIPAddress)
	}
	assert.ElementsMatch(t, request.RequestHeaders, testRequest.RequestHeaders)
	if request.ResponseCode != testRequest.ResponseCode {
		t.Errorf("unexpected request ResponseCode: got %q, expected %q", request.ResponseCode, testRequest.ResponseCode)
	}
	assert.ElementsMatch(t, request.ResponseHeaders, testRequest.ResponseHeaders)
	if request.ResponseSize != testRequest.ResponseSize {
		t.Errorf("unexpected request ResponseSize: got %q, expected %q", request.ResponseSize, testRequest.ResponseSize)
	}
	if request.ResponseTime != testRequest.ResponseTime {
		t.Errorf("unexpected request ResponseTime: got %q, expected %q", request.ResponseTime, testRequest.ResponseTime)
	}
	if request.Scheme != testRequest.Scheme {
		t.Errorf("unexpected request Scheme: got %q, expected %q", request.Scheme, testRequest.Scheme)
	}
	if request.ServerHostname != testRequest.ServerHostname {
		t.Errorf("unexpected request ServerHostname: got %q, expected %q", request.ServerHostname, testRequest.ServerHostname)
	}
	if request.ServerName != testRequest.ServerName {
		t.Errorf("unexpected request ServerName: got %q, expected %q", request.ServerName, testRequest.ServerName)
	}
	assert.ElementsMatch(t, request.Signals, testRequest.Signals)
	if request.Timestamp != testRequest.Timestamp {
		t.Errorf("unexpected request Timestamp: got %q, expected %q", request.Timestamp, testRequest.Timestamp)
	}
	if request.TLSCipher != testRequest.TLSCipher {
		t.Errorf("unexpected request TLSCipher: got %q, expected %q", request.TLSCipher, testRequest.TLSCipher)
	}
	if request.TLSProtocol != testRequest.TLSProtocol {
		t.Errorf("unexpected request TLSProtocol: got %q, expected %q", request.TLSProtocol, testRequest.TLSProtocol)
	}
	if request.URI != testRequest.URI {
		t.Errorf("unexpected request URI: got %q, expected %q", request.URI, testRequest.URI)
	}
	if request.UserAgent != testRequest.UserAgent {
		t.Errorf("unexpected request UserAgent: got %q, expected %q", request.UserAgent, testRequest.UserAgent)
	}
}

func TestClient_Getrequest_validation(t *testing.T) {
	var err error
	_, err = Get(fastly.TestClient, &GetInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = Get(fastly.TestClient, &GetInput{
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFEventsAndRequestsWorkspaceID),
		RequestID:   nil,
	})
	if !errors.Is(err, fastly.ErrMissingRequestID) {
		t.Errorf("expected ErrMissingrequestID: got %s", err)
	}
}

func TestClient_Listrequest_validation(t *testing.T) {
	var err error
	_, err = List(fastly.TestClient, &ListInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = List(fastly.TestClient, &ListInput{
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFEventsAndRequestsWorkspaceID),
		Limit:       nil,
	})
	if !errors.Is(err, fastly.ErrMissingLimit) {
		t.Errorf("expected ErrMissingLimit: got %s", err)
	}
}
