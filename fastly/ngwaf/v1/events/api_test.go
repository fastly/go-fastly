package events

import (
	"errors"
	"testing"
	"time"

	"github.com/fastly/go-fastly/v10/fastly"
	"github.com/fastly/go-fastly/v10/fastly/ngwaf/v1/signals"
)

const (
	TestEventID = "6841c2c07d3691b0f5b95130"
)

func TestClient_GetEvent(t *testing.T) {
	t.Parallel()

	getEventInput := new(GetInput)
	getEventInput.EventID = fastly.ToPointer(TestEventID)
	getEventInput.WorkspaceID = fastly.ToPointer(fastly.TestNGWAFWorkspaceID)

	var event *Event
	var err error
	createdAt, _ := time.Parse(time.RFC3339, "2025-06-05T16:15:58Z")
	detectedAt, _ := time.Parse(time.RFC3339, "2025-06-05T16:16:00Z")
	expiresAt, _ := time.Parse(time.RFC3339, "2025-06-05T16:20:02Z")
	testEvent := Event{
		Action:              "flagged",
		BlockSignals:        nil,
		BlockedRequestCount: 0,
		Country:             "US",
		CreatedAt:           createdAt,
		DetectedAt:          detectedAt,
		ExpiresAt:           expiresAt,
		EventID:             TestEventID,
		FlaggedRequestCount: 0,
		IsExpired:           true,
		Reasons: []signals.Reason{
			{
				SignalID: "CMDEXE",
				Count:    97,
			},
		},
		RemoteHostname: "pool-96-224-50-187.nycmny.fios.verizon.net",
		RequestCount:   97,
		Source:         "96.224.50.187",
		Type:           "attack",
		UserAgents: []string{
			"curl/8.7.1",
		},
		Window: 60,
	}

	// get an event
	fastly.Record(t, "get_event", func(c *fastly.Client) {
		event, err = Get(c, getEventInput)
	})
	if err != nil {
		t.Fatal(err)
	}
	if event.Action != testEvent.Action {
		t.Errorf("unexpected event Action: got %q, expected %q", event.Action, testEvent.Action)
	}
	if len(event.BlockSignals) != len(testEvent.BlockSignals) {
		t.Errorf("unexpected event Block Signals: got %q, expected %q", event.BlockSignals, testEvent.BlockSignals)
	}
	if event.BlockedRequestCount != testEvent.BlockedRequestCount {
		t.Errorf("unexpected event Blocked Request Count: got %d, expected %d", event.BlockedRequestCount, testEvent.BlockedRequestCount)
	}
	if event.CreatedAt != testEvent.CreatedAt {
		t.Errorf("unexpected event Created At: got %q, expected %q", event.CreatedAt, testEvent.CreatedAt)
	}
	if event.Country != testEvent.Country {
		t.Errorf("unexpected event Country: got %q, expected %q", event.Country, testEvent.Country)
	}
	if event.DetectedAt != testEvent.DetectedAt {
		t.Errorf("unexpected event Detected At: got %q, expected %q", event.DetectedAt, testEvent.DetectedAt)
	}
	if event.ExpiresAt != testEvent.ExpiresAt {
		t.Errorf("unexpected event Expires At: got %q, expected %q", event.ExpiresAt, testEvent.ExpiresAt)
	}
	if event.EventID != testEvent.EventID {
		t.Errorf("unexpected event Event ID: got %q, expected %q", event.EventID, testEvent.EventID)
	}
	if event.FlaggedRequestCount != testEvent.FlaggedRequestCount {
		t.Errorf("unexpected event Flagged Request Count: got %d, expected %d", event.FlaggedRequestCount, testEvent.FlaggedRequestCount)
	}
	if event.IsExpired != testEvent.IsExpired {
		t.Errorf("unexpected event Is Expired: got %t, expected %t", event.IsExpired, testEvent.IsExpired)
	}
	if event.Reasons[0].SignalID != testEvent.Reasons[0].SignalID {
		t.Errorf("unexpected event Reason id: got %q, expected %q", event.Reasons[0].SignalID, testEvent.Reasons[0].SignalID)
	}
	if event.Reasons[0].Count != testEvent.Reasons[0].Count {
		t.Errorf("unexpected event Reason id: got %d, expected %d", event.Reasons[0].Count, testEvent.Reasons[0].Count)
	}
	if event.RemoteHostname != testEvent.RemoteHostname {
		t.Errorf("unexpected event Remote Hostname: got %q, expected %q", event.RemoteHostname, testEvent.RemoteHostname)
	}
	if event.RequestCount != testEvent.RequestCount {
		t.Errorf("unexpected event Request Count: got %d, expected %d", event.RequestCount, testEvent.RequestCount)
	}
	if event.Source != testEvent.Source {
		t.Errorf("unexpected event Source: got %q, expected %q", event.Source, testEvent.Source)
	}
	if event.Type != testEvent.Type {
		t.Errorf("unexpected event Type: got %q, expected %q", event.Type, testEvent.Type)
	}
	if len(event.UserAgents) != len(testEvent.UserAgents) {
		t.Errorf("unexpected event User Agents: got %d, expected %d", len(event.UserAgents), len(testEvent.UserAgents))
	}
	if event.UserAgents[0] != testEvent.UserAgents[0] {
		t.Errorf("unexpected event User Agent: got %q, expected %q", event.UserAgents[0], testEvent.UserAgents[0])
	}
	if event.Window != testEvent.Window {
		t.Errorf("unexpected event Type: got %d, expected %d", event.Window, testEvent.Window)
	}

	var events *Events
	listEventInput := new(ListInput)
	listEventInput.WorkspaceID = fastly.ToPointer(fastly.TestNGWAFWorkspaceID)
	listEventInput.From = fastly.ToPointer("2024-05-27T14:08:03Z")

	// get a list of events
	fastly.Record(t, "list_event", func(c *fastly.Client) {
		events, err = List(c, listEventInput)
	})
	event = &events.Data[0]
	if err != nil {
		t.Fatal(err)
	}
	if event.Action != testEvent.Action {
		t.Errorf("unexpected event Action: got %q, expected %q", event.Action, testEvent.Action)
	}
	if len(event.BlockSignals) != len(testEvent.BlockSignals) {
		t.Errorf("unexpected event Block Signals: got %q, expected %q", event.BlockSignals, testEvent.BlockSignals)
	}
	if event.BlockedRequestCount != testEvent.BlockedRequestCount {
		t.Errorf("unexpected event Blocked Request Count: got %d, expected %d", event.BlockedRequestCount, testEvent.BlockedRequestCount)
	}
	if event.CreatedAt != testEvent.CreatedAt {
		t.Errorf("unexpected event Created At: got %q, expected %q", event.CreatedAt, testEvent.CreatedAt)
	}
	if event.Country != testEvent.Country {
		t.Errorf("unexpected event Country: got %q, expected %q", event.Country, testEvent.Country)
	}
	if event.DetectedAt != testEvent.DetectedAt {
		t.Errorf("unexpected event Detected At: got %q, expected %q", event.DetectedAt, testEvent.DetectedAt)
	}
	if event.ExpiresAt != testEvent.ExpiresAt {
		t.Errorf("unexpected event Expires At: got %q, expected %q", event.ExpiresAt, testEvent.ExpiresAt)
	}
	if event.EventID != testEvent.EventID {
		t.Errorf("unexpected event Event ID: got %q, expected %q", event.EventID, testEvent.EventID)
	}
	if event.FlaggedRequestCount != testEvent.FlaggedRequestCount {
		t.Errorf("unexpected event Flagged Request Count: got %d, expected %d", event.FlaggedRequestCount, testEvent.FlaggedRequestCount)
	}
	if event.IsExpired != testEvent.IsExpired {
		t.Errorf("unexpected event Is Expired: got %t, expected %t", event.IsExpired, testEvent.IsExpired)
	}
	if event.Reasons[0].SignalID != testEvent.Reasons[0].SignalID {
		t.Errorf("unexpected event Reason id: got %q, expected %q", event.Reasons[0].SignalID, testEvent.Reasons[0].SignalID)
	}
	if event.Reasons[0].Count != testEvent.Reasons[0].Count {
		t.Errorf("unexpected event Reason id: got %d, expected %d", event.Reasons[0].Count, testEvent.Reasons[0].Count)
	}
	if event.RemoteHostname != testEvent.RemoteHostname {
		t.Errorf("unexpected event Remote Hostname: got %q, expected %q", event.RemoteHostname, testEvent.RemoteHostname)
	}
	if event.RequestCount != testEvent.RequestCount {
		t.Errorf("unexpected event Request Count: got %d, expected %d", event.RequestCount, testEvent.RequestCount)
	}
	if event.Source != testEvent.Source {
		t.Errorf("unexpected event Source: got %q, expected %q", event.Source, testEvent.Source)
	}
	if event.Type != testEvent.Type {
		t.Errorf("unexpected event Type: got %q, expected %q", event.Type, testEvent.Type)
	}
	if len(event.UserAgents) != len(testEvent.UserAgents) {
		t.Errorf("unexpected event User Agents: got %d, expected %d", len(event.UserAgents), len(testEvent.UserAgents))
	}
	if event.UserAgents[0] != testEvent.UserAgents[0] {
		t.Errorf("unexpected event User Agent: got %q, expected %q", event.UserAgents[0], testEvent.UserAgents[0])
	}
	if event.Window != testEvent.Window {
		t.Errorf("unexpected event Type: got %d, expected %d", event.Window, testEvent.Window)
	}

	expireEventInput := new(ExpireInput)
	expireEventInput.WorkspaceID = fastly.ToPointer(fastly.TestNGWAFWorkspaceID)
	expireEventInput.EventID = fastly.ToPointer(TestEventID)
	expireEventInput.IsExpired = fastly.ToPointer(true)

	// expire an event.
	// this test relies on the fixture, in order to rerun this without the fixture you will have to have a new unexpired event.
	fastly.Record(t, "expire_event", func(c *fastly.Client) {
		event, err = Expire(c, expireEventInput)
	})
	if err != nil {
		t.Fatal(err)
	}
	if event.Action != testEvent.Action {
		t.Errorf("unexpected event Action: got %q, expected %q", event.Action, testEvent.Action)
	}
	if len(event.BlockSignals) != len(testEvent.BlockSignals) {
		t.Errorf("unexpected event Block Signals: got %q, expected %q", event.BlockSignals, testEvent.BlockSignals)
	}
	if event.BlockedRequestCount != testEvent.BlockedRequestCount {
		t.Errorf("unexpected event Blocked Request Count: got %d, expected %d", event.BlockedRequestCount, testEvent.BlockedRequestCount)
	}
	if event.CreatedAt != testEvent.CreatedAt {
		t.Errorf("unexpected event Created At: got %q, expected %q", event.CreatedAt, testEvent.CreatedAt)
	}
	if event.Country != testEvent.Country {
		t.Errorf("unexpected event Country: got %q, expected %q", event.Country, testEvent.Country)
	}
	if event.DetectedAt != testEvent.DetectedAt {
		t.Errorf("unexpected event Detected At: got %q, expected %q", event.DetectedAt, testEvent.DetectedAt)
	}
	if event.ExpiresAt != testEvent.ExpiresAt {
		t.Errorf("unexpected event Expires At: got %q, expected %q", event.ExpiresAt, testEvent.ExpiresAt)
	}
	if event.EventID != testEvent.EventID {
		t.Errorf("unexpected event Event ID: got %q, expected %q", event.EventID, testEvent.EventID)
	}
	if event.FlaggedRequestCount != testEvent.FlaggedRequestCount {
		t.Errorf("unexpected event Flagged Request Count: got %d, expected %d", event.FlaggedRequestCount, testEvent.FlaggedRequestCount)
	}
	if event.IsExpired != true {
		t.Errorf("unexpected event Is Expired: got %t, expected %t", event.IsExpired, true)
	}
	if event.Reasons[0].SignalID != testEvent.Reasons[0].SignalID {
		t.Errorf("unexpected event Reason id: got %q, expected %q", event.Reasons[0].SignalID, testEvent.Reasons[0].SignalID)
	}
	if event.Reasons[0].Count != testEvent.Reasons[0].Count {
		t.Errorf("unexpected event Reason id: got %d, expected %d", event.Reasons[0].Count, testEvent.Reasons[0].Count)
	}
	if event.RemoteHostname != testEvent.RemoteHostname {
		t.Errorf("unexpected event Remote Hostname: got %q, expected %q", event.RemoteHostname, testEvent.RemoteHostname)
	}
	if event.RequestCount != testEvent.RequestCount {
		t.Errorf("unexpected event Request Count: got %d, expected %d", event.RequestCount, testEvent.RequestCount)
	}
	if event.Source != testEvent.Source {
		t.Errorf("unexpected event Source: got %q, expected %q", event.Source, testEvent.Source)
	}
	if event.Type != testEvent.Type {
		t.Errorf("unexpected event Type: got %q, expected %q", event.Type, testEvent.Type)
	}
	if len(event.UserAgents) != len(testEvent.UserAgents) {
		t.Errorf("unexpected event User Agents: got %d, expected %d", len(event.UserAgents), len(testEvent.UserAgents))
	}
	if event.UserAgents[0] != testEvent.UserAgents[0] {
		t.Errorf("unexpected event User Agent: got %q, expected %q", event.UserAgents[0], testEvent.UserAgents[0])
	}
	if event.Window != testEvent.Window {
		t.Errorf("unexpected event Type: got %d, expected %d", event.Window, testEvent.Window)
	}
}

func TestClient_GetEvent_validation(t *testing.T) {
	var err error
	_, err = Get(fastly.TestClient, &GetInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = Get(fastly.TestClient, &GetInput{
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		EventID:     nil,
	})
	if !errors.Is(err, fastly.ErrMissingEventID) {
		t.Errorf("expected ErrMissingEventID: got %s", err)
	}
}

func TestClient_ListEvent_validation(t *testing.T) {
	var err error
	_, err = List(fastly.TestClient, &ListInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = List(fastly.TestClient, &ListInput{
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		From:        nil,
	})
	if !errors.Is(err, fastly.ErrMissingFrom) {
		t.Errorf("expected ErrMissingFrom: got %s", err)
	}
}

func TestClient_ExpireEvent_validation(t *testing.T) {
	var err error
	_, err = Expire(fastly.TestClient, &ExpireInput{
		WorkspaceID: nil,
	})
	if !errors.Is(err, fastly.ErrMissingWorkspaceID) {
		t.Errorf("expected ErrMissingWorkspaceID: got %s", err)
	}
	_, err = Expire(fastly.TestClient, &ExpireInput{
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		EventID:     nil,
	})
	if !errors.Is(err, fastly.ErrMissingEventID) {
		t.Errorf("expected ErrMissingEventID: got %s", err)
	}

	_, err = Expire(fastly.TestClient, &ExpireInput{
		WorkspaceID: fastly.ToPointer(fastly.TestNGWAFWorkspaceID),
		EventID:     fastly.ToPointer(string(TestEventID)),
		IsExpired:   nil,
	})
	if !errors.Is(err, fastly.ErrMissingIsExpired) {
		t.Errorf("expected ErrMissingIsExpired: got %s", err)
	}
}
