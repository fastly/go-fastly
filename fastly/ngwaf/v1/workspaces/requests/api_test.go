package requests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/fastly/go-fastly/v10/fastly"
)

const (
	TestRequestID = "1b34db183e374ea7880f000000000001"
)

func TestClient_requests(t *testing.T) {
	getrequestInput := new(GetInput)
	getrequestInput.RequestID = fastly.ToPointer(string(TestRequestID))
	getrequestInput.WorkspaceID = fastly.ToPointer(fastly.TestNGWAFWorkspaceID)

	var request *Request
	var err error
	timestamp, _ := time.Parse(time.RFC3339, "2025-06-05T16:15:31Z")
	testRequest := Request{
		AgentResponseCode: 200,
		Country:           "US",
		ID:                "1b34db183e374ea7880f000000000001",
		Method:            "POST",
		Path:              "/",
		Protocol:          "HTTP/1.1",
		RemoteHostname:    "pool-96-224-50-187.nycmny.fios.verizon.net",
		RemoteIPAddress:   "96.224.50.187",
		RequestHeaders: []Header{
			{
				Name:  "X-Uat-Ip",
				Value: "10.2.3.4",
			},
			{
				Name:  "Fastly-Ff",
				Value: "7rhL8T7Um/0khntXIHgopW0fTXKYvcj2vPpd1I+60Rc=!EWR!cache-ewr-kewr1740025-EWR",
			},
			{
				Name:  "X-Sigsci-Client-Geo-City",
				Value: "ossining",
			},
			{
				Name:  "X-Timer",
				Value: "S1749140131.063312, VS0",
			},
			{
				Name:  "X-Sigsci-Requestid",
				Value: "1b34db183e374ea7880f000000000001",
			},
			{
				Name:  "Content-Length",
				Value: "23",
			},
			{
				Name:  "Content-Type",
				Value: "application/x-www-form-urlencoded",
			},
			{
				Name:  "Fastly-Client-Ip",
				Value: "96.224.50.187",
			},
			{
				Name:  "User-Agent",
				Value: "curl/8.7.1",
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
				Name:  "X-Sigsci-Serviceid-Prod",
				Value: "kKJb5bOFI47uHeBVluGfX1",
			},
			{
				Name:  "X-Varnish",
				Value: "3020601334",
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
				Name:  "Host",
				Value: "fastlydevtoolstesting.com",
			},
		},
		ResponseCode: 401,
		ResponseHeaders: []Header{
			{
				Name:  "X-Glitch-Proxy",
				Value: "true",
			},
			{
				Name:  "Date",
				Value: "Thu, 05 Jun 2025 16:15:31 GMT",
			},
			{
				Name:  "Content-Type",
				Value: "text/html; charset=utf-8",
			},
			{
				Name:  "X-Cache",
				Value: "MISS",
			},
			{
				Name:  "Content-Length",
				Value: "4696",
			},
			{
				Name:  "Connection",
				Value: "keep-alive",
			},
			{
				Name:  "X-Cache-Hits",
				Value: "0",
			},
			{
				Name:  "X-Served-By",
				Value: "cache-ewr-kewr1740079-EWR",
			},
			{
				Name:  "Cache-Control",
				Value: "no-cache",
			},
			{
				Name:  "Etag",
				Value: "W/\"1258-CRcGKaD3eZqNFrAdjmcK04t9o4g\"",
			},
		},
		ResponseSize:   4696,
		ResponseTime:   28,
		Scheme:         "http",
		ServerHostname: "fastlydevtoolstesting.com",
		ServerName:     "fastlydevtoolstesting.com",
		Signals: []Signal{
			{
				ID:       "CMDEXE",
				Location: "POST",
				Value:    "bar192=;cat /etc/passwd",
				Detector: "CmdExeRule",
			},
			{
				ID:       "HTTP4XX",
				Location: "",
				Value:    "401",
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
		request, err = Get(context.TODO(), c, getrequestInput)
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
	listRequestInput.WorkspaceID = fastly.ToPointer(fastly.TestNGWAFWorkspaceID)
	listRequestInput.Limit = fastly.ToPointer(100)
	listRequestInput.Query = fastly.ToPointer("from:-2d")

	// get a list of requests
	fastly.Record(t, "list_request", func(c *fastly.Client) {
		Requests, err = List(context.TODO(), c, listRequestInput)
	})
	request = &Requests.Data[7]
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
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = Get(context.TODO(), fastly.TestClient, &GetInput{
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		RequestID:   nil,
	})
	if !errors.Is(err, fastly.ErrMissingRequestID) {
		t.Errorf("expected ErrMissingrequestID: got %s", err)
	}
}

func TestClient_Listrequest_validation(t *testing.T) {
	var err error
	_, err = List(context.TODO(), fastly.TestClient, &ListInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
}
