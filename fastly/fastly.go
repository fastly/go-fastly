package fastly

import (
	"bytes"
	"encoding"
	"net/url"
)

// BatchOperation represents batching variants.
type BatchOperation string

const (
	// CreateBatchOperation represents a batching variant.
	CreateBatchOperation BatchOperation = "create"
	// UpdateBatchOperation represents a batching variant.
	UpdateBatchOperation BatchOperation = "update"
	// UpsertBatchOperation represents a batching variant.
	UpsertBatchOperation BatchOperation = "upsert"
	// DeleteBatchOperation represents a batching variant.
	DeleteBatchOperation BatchOperation = "delete"

	// BatchModifyMaximumOperations represents the maximum number of operations
	// that can be sent within a single batch request. This is currently not
	// documented in the API.
	BatchModifyMaximumOperations = 1000

	// MaximumDictionarySize represents the maximum number of items that can be
	// placed within an Edge Dictionary.
	MaximumDictionarySize = 10000

	// MaximumACLSize represents the maximum number of entries that can be placed
	// within an ACL.
	MaximumACLSize = 10000
)

type statusResp struct {
	Msg    string
	Status string
}

func (t *statusResp) Ok() bool {
	return t.Status == "ok"
}

// Ensure Compatibool implements the proper interfaces.
var (
	_ encoding.TextMarshaler   = new(Compatibool)
	_ encoding.TextUnmarshaler = new(Compatibool)
)

// Compatibool is a boolean value that marshalls to 0/1 instead of true/false
// for compatibility with Fastly's API.
type Compatibool bool

// MarshalText implements the encoding.TextMarshaler interface.
func (b Compatibool) MarshalText() ([]byte, error) {
	if b {
		return []byte("1"), nil
	}
	return []byte("0"), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (b *Compatibool) UnmarshalText(t []byte) error {
	if bytes.Equal(t, []byte("1")) {
		*b = Compatibool(true)
	}
	return nil
}

// EncodeValues implements github.com/google/go-querystring/query#Encoder interface.
func (b Compatibool) EncodeValues(key string, v *url.Values) error {
	if b {
		v.Add(key, "1")
		return nil
	}
	v.Add(key, "0")
	return nil
}
