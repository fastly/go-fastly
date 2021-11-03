package fastly

import (
	"net/url"
)

type BatchOperation string

const (
	CreateBatchOperation BatchOperation = "create"
	UpdateBatchOperation BatchOperation = "update"
	UpsertBatchOperation BatchOperation = "upsert"
	DeleteBatchOperation BatchOperation = "delete"

	// Represents the maximum number of operations that can be sent within a single batch request.
	// This is currently not documented in the API.
	BatchModifyMaximumOperations = 1000

	// Represents the maximum number of items that can be placed within an Edge Dictionary.
	MaximumDictionarySize = 10000

	// Represents the maximum number of entries that can be placed within an ACL.
	MaximumACLSize = 10000
)

type statusResp struct {
	Status string
	Msg    string
}

func (t *statusResp) Ok() bool {
	return t.Status == "ok"
}

// Helper function to get a pointer to bool
func CBool(b bool) *Compatibool {
	c := Compatibool(b)
	return &c
}

// Compatibool is a boolean value that marshalls to 0/1 instead of true/false
// for compatibility with Fastly's API.
type Compatibool bool

// EncodeValues implements github.com/google/go-querystring/query#Encoder interface.
func (b Compatibool) EncodeValues(key string, v *url.Values) error {
	if b {
		v.Add(key, "1")
		return nil
	}
	v.Add(key, "0")
	return nil
}
