package fastly

import (
	"fmt"
	"time"
)

// DictionaryInfo represents a dictionary metadata response from the Fastly API.
type DictionaryInfo struct {
	// LastUpdated is the Time-stamp (GMT) when the dictionary was last updated.
	LastUpdated *time.Time `mapstructure:"last_updated"`

	// Digest is the hash of the dictionary content.
	Digest string `mapstructure:"digest"`

	// ItemCount is the number of items belonging to the dictionary.
	ItemCount int `mapstructure:"item_count"`
}

// GetDictionaryInfoInput is used as input to the GetDictionary function.
type GetDictionaryInfoInput struct {
	// ServiceID is the ID of the service Dictionary belongs to (required).
	ServiceID string

	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int

	// ID is the alphanumeric string identifying a dictionary.
	ID string
}

// GetDictionaryInfo gets the dictionary metadata with the given parameters.
func (c *Client) GetDictionaryInfo(i *GetDictionaryInfoInput) (*DictionaryInfo, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}

	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := fmt.Sprintf("/service/%s/version/%d/dictionary/%s/info", i.ServiceID, i.ServiceVersion, i.ID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var b *DictionaryInfo
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}
