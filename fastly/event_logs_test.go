package fastly

import (
	"bytes"
	"io/ioutil"
	"testing"
)

var testEventID = "3OMewexIMbzrQj77xxxxxx"

func TestClient_APIEvents(t *testing.T) {
	t.Parallel()

	var err error
	var events GetAPIEventsResponse
	record(t, "events/get_events", func(c *Client) {
		events, err = c.GetAPIEvents(&GetAPIEventsFilterInput{
			PageNumber: 1,
			MaxResults: 1,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(events.Events) < 1 {
		t.Errorf("bad events: %v", events.Events)
	}

	var event *Event
	record(t, "events/get_event", func(c *Client) {
		event, err = c.GetAPIEvent(&GetAPIEventInput{
			EventID: events.Events[0].ID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(event.ID) < 1 {
		t.Errorf("bad event: %v", event)
	}

}

func TestClient_GetAPIEvent_validation(t *testing.T) {
	var err error
	_, err = testClient.GetAPIEvent(&GetAPIEventInput{
		EventID: "",
	})
	if err != ErrMissingEventID {
		t.Errorf("bad error: %s", err)
	}

}

func TestGetAPIEventsFilterInput_formatFilters(t *testing.T) {
	tests := []struct {
		description string
		filters     GetAPIEventsFilterInput
		expected    map[string]string
	}{
		{
			description: "converts both strings and ints to strings",
			filters: GetAPIEventsFilterInput{
				CustomerID: "65135846153687547",
				ServiceID:  "5343548168357658",
				EventType:  "version.activate",
				UserID:     "654681384354746951",
				MaxResults: 1,
				PageNumber: 2,
			},
			expected: map[string]string{
				"filter[customer_id]": "65135846153687547",
				"filter[service_id]":  "5343548168357658",
				"filter[event_type]":  "version.activate",
				"filter[user_id]":     "654681384354746951",
				"page[size]":          "1",
				"page[number]":        "2",
			},
		},
	}
	for _, testcase := range tests {
		answer := testcase.filters.formatEventFilters()
		if len(answer) != len(testcase.expected) {
			t.Errorf("In test %s: Expected map with %d entries,got one with %d", testcase.description, len(testcase.expected), len(answer))
		}
		for key, value := range testcase.expected {
			if answer[key] != value {
				t.Errorf("In test %s: Expected %s key to have value %s, got %s", testcase.description, key, value, answer[key])
			}
		}
	}
}

func TestGetEventsPages(t *testing.T) {
	tests := []struct {
		description   string
		input         string
		expectedPages EventsPaginationInfo
		expectedErr   error
	}{
		{
			description: "returns the next page",
			input:       `{"links": {"next": "https://google.com/2"}, "data": []}`,
			expectedPages: EventsPaginationInfo{
				Next: "https://google.com/2",
			},
		},
		{
			description: "returns multiple pages",
			input:       `{"links": {"next": "https://google.com/2", "first": "https://google.com/1"}, "data": []}`,
			expectedPages: EventsPaginationInfo{
				First: "https://google.com/1",
				Next:  "https://google.com/2",
			},
		},
		{
			description:   "returns no pages",
			input:         `{"data": []}`,
			expectedPages: EventsPaginationInfo{},
		},
	}
	for _, testcase := range tests {
		pages, reader, err := getEventsPages(bytes.NewReader([]byte(testcase.input)))
		if pages != testcase.expectedPages {
			t.Errorf("Test %s: Expected pages %+v, got %+v", testcase.description, testcase.expectedPages, pages)
		}

		// we expect to be able to get the original input out again
		resultBytes, _ := ioutil.ReadAll(reader)
		if string(resultBytes) != testcase.input {
			t.Errorf("Test %s: Expected body %s, got %s", testcase.description, testcase.input, string(resultBytes))
		}
		if err != testcase.expectedErr {
			t.Errorf("Test %s: Expected error %v, got %v", testcase.description, testcase.expectedErr, err)
		}
	}
}
