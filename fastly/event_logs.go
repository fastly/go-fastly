package fastly

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/jsonapi"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Events represents an event_logs item response from the Fastly API.
// type Event struct {
// 	Data  []*Data `mapstructure:"data"`
// 	Links Links   `mapstructure:"links"`
// }

type Event struct {
	ID          string                  `jsonapi:"primary,event"`
	CustomerID  string                  `jsonapi:"attr,customer_id"`
	Description string                  `jsonapi:"attr,description"`
	EventType   string                  `jsonapi:"attr,event_type"`
	IP          string                  `jsonapi:"attr,ip"`
	Metadata    *map[string]interface{} `jsonapi:"attr,metadata"`
	ServiceID   string                  `jsonapi:"attr,service_id"`
	UserID      string                  `jsonapi:"attr,user_id"`
	CreatedAt   time.Time               `jsonapi:"attr,created_at"`
	Admin       bool                    `jsonapi:"attr,admin"`
	// Type       string     `mapstructure:"type"`
	// Attributes Attributes `mapstructure:"attributes"`
}

// eventType is used for reflection because JSONAPI wants to know what it's
// decoding into.
// var eventType = reflect.TypeOf(new(Event))

// type Attributes struct {
// 	CustomerID  string                  `mapstructure:"customer_id"`
// 	Description string                  `mapstructure:"description"`
// 	EventType   string                  `mapstructure:"event_type"`
// 	IP          string                  `mapstructure:"ip"`
// 	Metadata    *map[string]interface{} `mapstructure:"metadata"`
// 	ServiceID   string                  `mapstructure:"service_id"`
// 	UserID      string                  `mapstructure:"user_id"`
// 	CreatedAt   time.Time               `mapstructure:"created_at"`
// 	Admin       bool                    `mapstructure:"admin"`
// }

type GetAPIEventsInput struct {
	CustomerID string
	ServiceId  string
	Filters    GetAPIEventsFilter
}

// GetAPIEventsFilter is used as input to the GetAPIEvents function.
type GetAPIEventsFilter struct {
	// CustomerID to Limit the returned events to a specific customer.
	CustomerID string

	// ServiceID to Limit the returned events to a specific service.
	ServiceID string

	// EventType to Limit the returned events to a specific event type. See above for event codes.
	EventType string

	// UserID to Limit the returned events to a specific user.
	UserID string

	// Number is the Pagination page number.
	PageNumber int

	// Size is the Number of items to return on each paginated page.
	MaxResults int
}

// type Links struct {
// 	Next string `mapstructure:"next"`
// 	Last string `mapstructure:"last"`
// }

// eventLinksResponse is used to pull the "Links" pagination fields from
// a call to Fastly; these are excluded from the results of the jsonapi
// call to `UnmarshalManyPayload()`, so we have to fetch them separately.
type eventLinksResponse struct {
	Links eventsPaginationInfo `json:"links"`
}

// eventsPaginationInfo stores links to searches related to the current one, showing
// any information about additional results being stored on another page
type eventsPaginationInfo struct {
	First string `json:"first,omitempty"`
	Last  string `json:"last,omitempty"`
	Next  string `json:"next,omitempty"`
}

// GetWAFRuleStatusesResponse is the data returned to the user from a GetAPIEvents call
type GetAPIEventsResponse struct {
	Events []*Event
}

func (i *GetAPIEventsInput) formatFilters() map[string]string {
	input := i.Filters
	result := map[string]string{}
	pairings := map[string]interface{}{
		"filter[customer_id]": input.CustomerID,
		"filter[service_id]":  input.ServiceID,
		"filter[event_type]":  input.EventType,
		"filter[user_id]":     input.UserID,
		"page[size]":          input.MaxResults,
		"page[number]":        input.PageNumber, // starts at 1, not 0
	}
	// NOTE: This setup means we will not be able to send the zero value
	// of any of these filters. It doesn't appear we would need to at present.
	for key, value := range pairings {
		switch t := reflect.TypeOf(value).String(); t {
		case "string":
			if value != "" {
				result[key] = value.(string)
			}
		case "int":
			if value != 0 {
				result[key] = strconv.Itoa(value.(int))
			}
		case "[]int":
			// convert ints to strings
			toStrings := []string{}
			values := value.([]int)
			for _, i := range values {
				toStrings = append(toStrings, strconv.Itoa(i))
			}
			// concat strings
			if len(values) > 0 {
				result[key] = strings.Join(toStrings, ",")
			}
		}
	}
	return result

}

// GetAPIEvents gets the events for a particular customer
func (c *Client) GetAPIEvents(i *GetAPIEventsInput) (GetAPIEventsResponse, error) {
	eventsResponse := GetAPIEventsResponse{Events: []*Event{}}
	//
	// if i.CustomerID == "" {
	// 	return eventsResponse, ErrMissingCustomerID
	// }
	// if i.ServiceId == "" {
	// 	return eventsResponse, ErrMissingService
	// }

	path := fmt.Sprintln("/events")
	filters := &RequestOptions{Params: i.formatFilters()}

	resp, err := c.Get(path, filters)
	if err != nil {
		return eventsResponse, err
	}
	err = c.interpretAPIEventsPage(&eventsResponse, resp)
	// NOTE: It's possible for statusResponse to be partially completed before an error
	// was encountered, so the presence of a statusResponse doesn't preclude the presence of
	// an error.
	return eventsResponse, err
}

// interpretWAFRuleStatusesPage accepts a Fastly response for a set of WAF rule statuses
// and unmarshals the results. If there are more pages of results, it fetches the next
// page, adds that response to the array of results, and repeats until all results have
// been fetched.
func (c *Client) interpretAPIEventsPage(answer *GetAPIEventsResponse, received *http.Response) error {
	// before we pull the status info out of the response body, fetch
	// pagination info from it:
	pages, body, err := getEventsPages(received.Body)
	if err != nil {
		return err
	}
	data, err := jsonapi.UnmarshalManyPayload(body, reflect.TypeOf(new(Event)))
	if err != nil {
		return err
	}

	for i := range data {
		typed, ok := data[i].(*Event)
		if !ok {
			return fmt.Errorf("got back response of unexpected type")
		}
		answer.Events = append(answer.Events, typed)
	}
	if pages.Next != "" {
		// NOTE: pages.Next URL includes filters already
		resp, err := c.SimpleGet(pages.Next)
		if err != nil {
			return err
		}
		c.interpretAPIEventsPage(answer, resp)
	}
	return nil
}

// getEventsPages parses a response to get the pagination data without destroying
// the reader we receive as "resp.Body"; this essentially copies resp.Body
// and returns it so we can use it again.
func getEventsPages(body io.Reader) (eventsPaginationInfo, io.Reader, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(body, &buf)

	bodyBytes, err := ioutil.ReadAll(tee)
	if err != nil {
		return eventsPaginationInfo{}, nil, err
	}

	var pages eventLinksResponse
	json.Unmarshal(bodyBytes, &pages)
	return pages.Links, bytes.NewReader(buf.Bytes()), nil
}

// // CreateDictionaryInput is used as input to the CreateDictionary function.
// type CreateDictionaryInput struct {
// 	// Service is the ID of the service. Version is the specific configuration
// 	// version. Both fields are required.
// 	Service string
// 	Version int
//
// 	Name string `form:"name,omitempty"`
// }
// // CreateDictionary creates a new Fastly dictionary.
// func (c *Client) CreateDictionary(i *CreateDictionaryInput) (*Dictionary, error) {
// 	if i.Service == "" {
// 		return nil, ErrMissingService
// 	}
//
// 	if i.Version == 0 {
// 		return nil, ErrMissingVersion
// 	}
//
// 	path := fmt.Sprintf("/service/%s/version/%d/dictionary", i.Service, i.Version)
// 	resp, err := c.PostForm(path, i, nil)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var b *Dictionary
// 	if err := decodeJSON(&b, resp.Body); err != nil {
// 		return nil, err
// 	}
// 	return b, nil
// }
//
// // GetDictionaryInput is used as input to the GetDictionary function.
// type GetDictionaryInput struct {
// 	// Service is the ID of the service. Version is the specific configuration
// 	// version. Both fields are required.
// 	Service string
// 	Version int
//
// 	// Name is the name of the dictionary to fetch.
// 	Name string
// }
//
// // GetDictionary gets the dictionary configuration with the given parameters.
// func (c *Client) GetDictionary(i *GetDictionaryInput) (*Dictionary, error) {
// 	if i.Service == "" {
// 		return nil, ErrMissingService
// 	}
//
// 	if i.Version == 0 {
// 		return nil, ErrMissingVersion
// 	}
//
// 	if i.Name == "" {
// 		return nil, ErrMissingName
// 	}
//
// 	path := fmt.Sprintf("/service/%s/version/%d/dictionary/%s", i.Service, i.Version, i.Name)
// 	resp, err := c.Get(path, nil)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var b *Dictionary
// 	if err := decodeJSON(&b, resp.Body); err != nil {
// 		return nil, err
// 	}
// 	return b, nil
// }
//
// // UpdateDictionaryInput is used as input to the UpdateDictionary function.
// type UpdateDictionaryInput struct {
// 	// Service is the ID of the service. Version is the specific configuration
// 	// version. Both fields are required.
// 	Service string
// 	Version int
//
// 	// Name is the name of the dictionary to update.
// 	Name string
//
// 	NewName string `form:"name,omitempty"`
// }
//
// // UpdateDictionary updates a specific dictionary.
// func (c *Client) UpdateDictionary(i *UpdateDictionaryInput) (*Dictionary, error) {
// 	if i.Service == "" {
// 		return nil, ErrMissingService
// 	}
//
// 	if i.Version == 0 {
// 		return nil, ErrMissingVersion
// 	}
//
// 	if i.Name == "" {
// 		return nil, ErrMissingName
// 	}
//
// 	path := fmt.Sprintf("/service/%s/version/%d/dictionary/%s", i.Service, i.Version, i.Name)
// 	resp, err := c.PutForm(path, i, nil)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	var b *Dictionary
// 	if err := decodeJSON(&b, resp.Body); err != nil {
// 		return nil, err
// 	}
// 	return b, nil
// }
//
// // DeleteDictionaryInput is the input parameter to DeleteDictionary.
// type DeleteDictionaryInput struct {
// 	// Service is the ID of the service. Version is the specific configuration
// 	// version. Both fields are required.
// 	Service string
// 	Version int
//
// 	// Name is the name of the dictionary to delete (required).
// 	Name string
// }
//
// // DeleteDictionary deletes the given dictionary version.
// func (c *Client) DeleteDictionary(i *DeleteDictionaryInput) error {
// 	if i.Service == "" {
// 		return ErrMissingService
// 	}
//
// 	if i.Version == 0 {
// 		return ErrMissingVersion
// 	}
//
// 	if i.Name == "" {
// 		return ErrMissingName
// 	}
//
// 	path := fmt.Sprintf("/service/%s/version/%d/dictionary/%s", i.Service, i.Version, i.Name)
// 	resp, err := c.Delete(path, nil)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()
//
// 	// Unlike other endpoints, the dictionary endpoint does not return a status
// 	// response - it just returns a 200 OK.
// 	return nil
// }
