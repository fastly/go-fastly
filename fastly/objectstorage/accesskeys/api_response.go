package accesskeys

import "time"

// AccessKey is the API response structure for the create and describe operations.
type AccessKey struct {
	// AccessKey is an AccessKey identifier.
	AccessKeyID string `json:"access_key"`
	// SecretKey is the secret for the access key
	SecretKey string `json:"secret_key"`
	// Description is human readable description for the access key.
	Description string `json:"description"`
	// Permission is the permissions the key has.
	Permission string `json:"permission"`
	// Buckets is the list of buckets associated with the access key.
	Buckets []string `json:"buckets"`
	// CreatedAt is the timestamp associated with the creation of the access key.
	CreatedAt time.Time `json:"created_at"`
}

// AccessKeys is the API response structure for the list access keys operation.
type AccessKeys struct {
	// Data is the list of returned AccessKeys.
	Data []AccessKey `json:"data"`
	// Meta is additional information about the request
	Meta MetaAccessKeys `json:"meta"`
}

type MetaAccessKeys struct{}
