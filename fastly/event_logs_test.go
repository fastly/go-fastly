package fastly

import (
	"testing"
)

var testEventID = "3OMewexIMbzrQj77xxxxxx"

// func TestClient_APIEvents(t *testing.T) {
// 	t.Parallel()
//
// 	var err error
//
// 	// Get Event
// 	var event *Event
// 	record(t *testing.T, fixture string, func(arg1 *Client) {
// 		[object Object]
// 	})
// 	record(t, "events/get", func(c *Client) {
// 		event, err = c.GetAPIEvent(&GetAPIEventInput{
// 			EventID: testEventID,
// 		})
// 	})
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if event.ID != testEventID {
// 		t.Errorf("expected %q to be %q", event.ID, testEventID)
// 	}
// }

// // Get Events
// var events []*Event
// record(t, "wafs/list", func(c *Client) {
// 	wafs, err = c.ListWAFs(&ListWAFsInput{
// 		Service: testServiceID,
// 		Version: tv.Number,
// 	})
// })
// if err != nil {
// 	t.Fatal(err)
// }
// if len(wafs) < 1 {
// 	t.Errorf("bad wafs: %v", wafs)
// }
//
// // Update
// // Create a new response object to attach
// var nro *ResponseObject
// record(t, "wafs/response_object/create_another", func(c *Client) {
// 	nro, err = c.CreateResponseObject(&CreateResponseObjectInput{
// 		Service:     testServiceID,
// 		Version:     tv.Number,
// 		Name:        "test-response-object-2",
// 		Status:      200,
// 		Response:    "Ok",
// 		Content:     "efgh",
// 		ContentType: "text/plain",
// 	})
// })

// }

func TestClient_GetAPIEvent_validation(t *testing.T) {
	var err error
	_, err = testClient.GetAPIEvent(&GetAPIEventInput{
		EventID: "",
	})
	if err != ErrMissingEventID {
		t.Errorf("bad error: %s", err)
	}
	_, err = testClient.GetAPIEvent(&GetAPIEventInput{
		EventID: testEventID,
	})
	if err != ErrMissingEventID {
		t.Errorf("bad error: %s", err)
	}
}

//
// func TestClient_UpdateWAF_validation(t *testing.T) {
// 	var err error
// 	_, err = testClient.UpdateWAF(&UpdateWAFInput{
// 		Service: "",
// 	})
// 	if err != ErrMissingService {
// 		t.Errorf("bad error: %s", err)
// 	}
//
// 	_, err = testClient.UpdateWAF(&UpdateWAFInput{
// 		Service: "foo",
// 		Version: 0,
// 	})
// 	if err != ErrMissingVersion {
// 		t.Errorf("bad error: %s", err)
// 	}
//
// 	_, err = testClient.UpdateWAF(&UpdateWAFInput{
// 		Service: "foo",
// 		Version: 1,
// 		WAFID:   "",
// 	})
// 	if err != ErrMissingWAFID {
// 		t.Errorf("bad error: %s", err)
// 	}
// }
//
// func TestClient_DeleteWAF_validation(t *testing.T) {
// 	var err error
// 	err = testClient.DeleteWAF(&DeleteWAFInput{
// 		Service: "",
// 	})
// 	if err != ErrMissingService {
// 		t.Errorf("bad error: %s", err)
// 	}
//
// 	err = testClient.DeleteWAF(&DeleteWAFInput{
// 		Service: "foo",
// 		Version: 0,
// 	})
// 	if err != ErrMissingVersion {
// 		t.Errorf("bad error: %s", err)
// 	}
//
// 	err = testClient.DeleteWAF(&DeleteWAFInput{
// 		Service: "foo",
// 		Version: 1,
// 		WAFID:   "",
// 	})
// 	if err != ErrMissingWAFID {
// 		t.Errorf("bad error: %s", err)
// 	}
// }

// func TestGetWAFRuleStatusesInput_formatFilters(t *testing.T) {
// 	tests := []struct {
// 		description string
// 		filters     GetWAFRuleStatusesFilters
// 		expected    map[string]string
// 	}{
// 		{
// 			description: "converts both strings and ints to strings",
// 			filters: GetWAFRuleStatusesFilters{
// 				Status:   "log",
// 				Accuracy: 10,
// 				Version:  "180ad",
// 			},
// 			expected: map[string]string{
// 				"filter[status]":         "log",
// 				"filter[rule][accuracy]": "10",
// 				"filter[rule][version]":  "180ad",
// 			},
// 		},
// 		{
// 			description: "converts arrays to strings",
// 			filters: GetWAFRuleStatusesFilters{
// 				Status:  "log",
// 				Version: "181ad",
// 				Tags:    []int{18, 1, 1093, 86308},
// 			},
// 			expected: map[string]string{
// 				"filter[status]":        "log",
// 				"filter[rule][version]": "181ad",
// 				"include":               "18,1,1093,86308",
// 			},
// 		},
// 	}
// 	for _, testcase := range tests {
// 		input := GetWAFRuleStatusesInput{
// 			Filters: testcase.filters,
// 		}
// 		answer := input.formatFilters()
// 		if len(answer) != len(testcase.expected) {
// 			t.Errorf("In test %s: Expected map with %d entries,got one with %d", testcase.description, len(testcase.expected), len(answer))
// 		}
// 		for key, value := range testcase.expected {
// 			if answer[key] != value {
// 				t.Errorf("In test %s: Expected %s key to have value %s, got %s", testcase.description, key, value, answer[key])
// 			}
// 		}
// 	}
// }
//
// func TestGetPages(t *testing.T) {
// 	tests := []struct {
// 		description   string
// 		input         string
// 		expectedPages paginationInfo
// 		expectedErr   error
// 	}{
// 		{
// 			description: "returns the next page",
// 			input:       `{"links": {"next": "https://google.com/2"}, "data": []}`,
// 			expectedPages: paginationInfo{
// 				Next: "https://google.com/2",
// 			},
// 		},
// 		{
// 			description: "returns multiple pages",
// 			input:       `{"links": {"next": "https://google.com/2", "first": "https://google.com/1"}, "data": []}`,
// 			expectedPages: paginationInfo{
// 				First: "https://google.com/1",
// 				Next:  "https://google.com/2",
// 			},
// 		},
// 		{
// 			description:   "returns no pages",
// 			input:         `{"data": []}`,
// 			expectedPages: paginationInfo{},
// 		},
// 	}
// 	for _, testcase := range tests {
// 		pages, reader, err := getPages(bytes.NewReader([]byte(testcase.input)))
// 		if pages != testcase.expectedPages {
// 			t.Errorf("Test %s: Expected pages %+v, got %+v", testcase.description, testcase.expectedPages, pages)
// 		}
//
// 		// we expect to be able to get the original input out again
// 		resultBytes, _ := ioutil.ReadAll(reader)
// 		if string(resultBytes) != testcase.input {
// 			t.Errorf("Test %s: Expected body %s, got %s", testcase.description, testcase.input, string(resultBytes))
// 		}
// 		if err != testcase.expectedErr {
// 			t.Errorf("Test %s: Expected error %v, got %v", testcase.description, testcase.expectedErr, err)
// 		}
// 	}
// }
