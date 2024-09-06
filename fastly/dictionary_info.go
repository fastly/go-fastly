package fastly

import (
	"strconv"
	"time"
)

// DictionaryInfo represents a dictionary metadata response from the Fastly API.
type DictionaryInfo struct {
	// Digest is the hash of the dictionary content.
	Digest *string `mapstructure:"digest"`
	// ItemCount is the number of items belonging to the dictionary.
	ItemCount *int `mapstructure:"item_count"`
	// LastUpdated is the Time-stamp (GMT) when the dictionary was last updated.
	LastUpdated *time.Time `mapstructure:"last_updated"`
}

// GetDictionaryInfoInput is used as input to the GetDictionary function.
type GetDictionaryInfoInput struct {
	// DictionaryID is the alphanumeric string identifying a dictionary (required).
	DictionaryID string
	// ServiceID is the ID of the service Dictionary belongs to (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetDictionaryInfo retrieves the specified resource.
func (c *Client) GetDictionaryInfo(i *GetDictionaryInfoInput) (*DictionaryInfo, error) {
	if i.DictionaryID == "" {
		return nil, ErrMissingDictionaryID
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := ToSafeURL("service", i.ServiceID, "version", strconv.Itoa(i.ServiceVersion), "dictionary", i.DictionaryID, "info")

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var b *DictionaryInfo
	if err := decodeBodyMap(resp.Body, &b); err != nil {
		return nil, err
	}
	return b, nil
}
