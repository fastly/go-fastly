package fastly

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

// S3Redundancy represents the redundancy variants for S3.
type S3Redundancy string

// S3RedundancyPtr returns a pointer to a S3Redundancy.
func S3RedundancyPtr(v S3Redundancy) *S3Redundancy {
	return &v
}

// S3ServerSideEncryption represents the encryption variants for S3.
type S3ServerSideEncryption string

// S3ServerSideEncryptionPtr returns a pointer to a S3ServerSideEncryption.
func S3ServerSideEncryptionPtr(v S3ServerSideEncryption) *S3ServerSideEncryption {
	return &v
}

// S3AccessControlList represents the control list variants for S3.
type S3AccessControlList string

// S3AccessControlListPtr returns a pointer to a S3AccessControlList.
func S3AccessControlListPtr(v S3AccessControlList) *S3AccessControlList {
	return &v
}

const (
	// S3RedundancyStandard represents a redundancy variant.
	S3RedundancyStandard S3Redundancy = "standard"
	// S3RedundancyIntelligentTiering represents a redundancy variant.
	S3RedundancyIntelligentTiering S3Redundancy = "intelligent_tiering"
	// S3RedundancyStandardIA represents a redundancy variant.
	S3RedundancyStandardIA S3Redundancy = "standard_ia"
	// S3RedundancyOneZoneIA represents a redundancy variant.
	S3RedundancyOneZoneIA S3Redundancy = "onezone_ia"
	// S3RedundancyGlacierInstantRetrieval represents a redundancy variant.
	S3RedundancyGlacierInstantRetrieval S3Redundancy = "glacier_ir"
	// S3RedundancyGlacierFlexibleRetrieval represents a redundancy variant.
	S3RedundancyGlacierFlexibleRetrieval S3Redundancy = "glacier"
	// S3RedundancyGlacierDeepArchive represents a redundancy variant.
	S3RedundancyGlacierDeepArchive S3Redundancy = "deep_archive"
	// S3RedundancyReduced represents a redundancy variant.
	S3RedundancyReduced S3Redundancy = "reduced_redundancy"

	// S3ServerSideEncryptionAES represents an encryption variant.
	S3ServerSideEncryptionAES S3ServerSideEncryption = "AES256"
	// S3ServerSideEncryptionKMS represents an encryption variant.
	S3ServerSideEncryptionKMS S3ServerSideEncryption = "aws:kms"

	// S3AccessControlListPrivate represents a control list variant.
	S3AccessControlListPrivate S3AccessControlList = "private"
	// S3AccessControlListPublicRead represents a control list variant.
	S3AccessControlListPublicRead S3AccessControlList = "public-read"
	// S3AccessControlListPublicReadWrite represents a control list variant.
	S3AccessControlListPublicReadWrite S3AccessControlList = "public-read-write"
	// S3AccessControlListAWSExecRead represents a control list variant.
	S3AccessControlListAWSExecRead S3AccessControlList = "aws-exec-read"
	// S3AccessControlListAuthenticatedRead represents a control list variant.
	S3AccessControlListAuthenticatedRead S3AccessControlList = "authenticated-read"
	// S3AccessControlListBucketOwnerRead represents a control list variant.
	S3AccessControlListBucketOwnerRead S3AccessControlList = "bucket-owner-read"
	// S3AccessControlListBucketOwnerFullControl represents a control list variant.
	S3AccessControlListBucketOwnerFullControl S3AccessControlList = "bucket-owner-full-control"
)

// S3 represents a S3 response from the Fastly API.
type S3 struct {
	ACL                          S3AccessControlList    `mapstructure:"acl"`
	AccessKey                    string                 `mapstructure:"access_key"`
	BucketName                   string                 `mapstructure:"bucket_name"`
	CompressionCodec             string                 `mapstructure:"compression_codec"`
	CreatedAt                    *time.Time             `mapstructure:"created_at"`
	DeletedAt                    *time.Time             `mapstructure:"deleted_at"`
	Domain                       string                 `mapstructure:"domain"`
	Format                       string                 `mapstructure:"format"`
	FormatVersion                int                    `mapstructure:"format_version"`
	GzipLevel                    int                    `mapstructure:"gzip_level"`
	IAMRole                      string                 `mapstructure:"iam_role"`
	MessageType                  string                 `mapstructure:"message_type"`
	Name                         string                 `mapstructure:"name"`
	Path                         string                 `mapstructure:"path"`
	Period                       int                    `mapstructure:"period"`
	Placement                    string                 `mapstructure:"placement"`
	PublicKey                    string                 `mapstructure:"public_key"`
	Redundancy                   S3Redundancy           `mapstructure:"redundancy"`
	ResponseCondition            string                 `mapstructure:"response_condition"`
	SecretKey                    string                 `mapstructure:"secret_key"`
	ServerSideEncryption         S3ServerSideEncryption `mapstructure:"server_side_encryption"`
	ServerSideEncryptionKMSKeyID string                 `mapstructure:"server_side_encryption_kms_key_id"`
	ServiceID                    string                 `mapstructure:"service_id"`
	ServiceVersion               int                    `mapstructure:"version"`
	TimestampFormat              string                 `mapstructure:"timestamp_format"`
	UpdatedAt                    *time.Time             `mapstructure:"updated_at"`
}

// s3sByName is a sortable list of S3s.
type s3sByName []*S3

// Len implement the sortable interface.
func (s s3sByName) Len() int {
	return len(s)
}

// Swap implement the sortable interface.
func (s s3sByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less implement the sortable interface.
func (s s3sByName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// ListS3sInput is used as input to the ListS3s function.
type ListS3sInput struct {
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// ListS3s retrieves all resources.
func (c *Client) ListS3s(i *ListS3sInput) ([]*S3, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/s3", i.ServiceID, i.ServiceVersion)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s3s []*S3
	if err := decodeBodyMap(resp.Body, &s3s); err != nil {
		return nil, err
	}
	sort.Stable(s3sByName(s3s))
	return s3s, nil
}

// CreateS3Input is used as input to the CreateS3 function.
type CreateS3Input struct {
	// ACL is the access control list (ACL) specific request header.
	ACL *S3AccessControlList `url:"acl,omitempty"`
	//  AccessKey is the access key for your S3 account. Not required if iam_role is provided.
	AccessKey *string `url:"access_key,omitempty"`
	// BucketName is the bucket name for S3 account.
	BucketName *string `url:"bucket_name,omitempty"`
	// CompressionCodec is the codec used for compressing your logs. Valid values are zstd, snappy, and gzip.
	CompressionCodec *string `url:"compression_codec,omitempty"`
	// Domain is the domain of the Amazon S3 endpoint.
	Domain *string `url:"domain,omitempty"`
	// Format is a Fastly log format string.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// GzipLevel is the level of gzip encoding when sending logs (default 0, no compression).
	GzipLevel *int `url:"gzip_level,omitempty"`
	// IAMRole is the Amazon Resource Name (ARN) for the IAM role granting Fastly access to S3. Not required if access_key and secret_key are provided.
	IAMRole *string `url:"iam_role,omitempty"`
	// MessageType is how the message should be formatted (classic, loggly, logplex, blank).
	MessageType *string `url:"message_type,omitempty"`
	// Name is the name of the SFTP to update (required).
	Name *string `url:"name,omitempty"`
	// Path is the path to upload logs to.
	Path *string `url:"path,omitempty"`
	// Period is how frequently log files are finalized so they can be available for reading (in seconds).
	Period *int `url:"period,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// PublicKey is a PGP public key that Fastly will use to encrypt your log files before writing them to disk.
	PublicKey *string `url:"public_key,omitempty"`
	// Redundancy is the S3 redundancy level.
	Redundancy *S3Redundancy `url:"redundancy,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// SecretKey is the secret key for your S3 account. Not required if iam_role is provided.
	SecretKey *string `url:"secret_key,omitempty"`
	// ServerSideEncryption should be set to AES256 or aws:kms to enable S3 Server Side Encryption.
	ServerSideEncryption *S3ServerSideEncryption `url:"server_side_encryption,omitempty"`
	// ServerSideEncryptionKMSKeyID is an optional server-side KMS Key ID. Must be set if ServerSideEncryption is set to aws:kms or AES256.
	ServerSideEncryptionKMSKeyID *string `url:"server_side_encryption_kms_key_id,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// TimestampFormat is a timestamp format.
	TimestampFormat *string `url:"timestamp_format,omitempty"`
}

// CreateS3 creates a new resource.
func (c *Client) CreateS3(i *CreateS3Input) (*S3, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	isKMS := i.ServerSideEncryption != nil && *i.ServerSideEncryption == S3ServerSideEncryptionKMS
	hasNoKeyID := i.ServerSideEncryptionKMSKeyID != nil && *i.ServerSideEncryptionKMSKeyID == ""
	if isKMS && hasNoKeyID {
		return nil, ErrMissingServerSideEncryptionKMSKeyID
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/s3", i.ServiceID, i.ServiceVersion)
	resp, err := c.PostForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s3 *S3
	if err := decodeBodyMap(resp.Body, &s3); err != nil {
		return nil, err
	}
	return s3, nil
}

// GetS3Input is used as input to the GetS3 function.
type GetS3Input struct {
	// Name is the name of the S3 to fetch (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// GetS3 retrieves the specified resource.
func (c *Client) GetS3(i *GetS3Input) (*S3, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/s3/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s3 *S3
	if err := decodeBodyMap(resp.Body, &s3); err != nil {
		return nil, err
	}
	return s3, nil
}

// UpdateS3Input is used as input to the UpdateS3 function.
type UpdateS3Input struct {
	// ACL is the access control list (ACL) specific request header.
	ACL *S3AccessControlList `url:"acl,omitempty"`
	//  AccessKey is the access key for your S3 account. Not required if iam_role is provided.
	AccessKey *string `url:"access_key,omitempty"`
	// BucketName is the bucket name for S3 account.
	BucketName *string `url:"bucket_name,omitempty"`
	// CompressionCodec is the codec used for compressing your logs. Valid values are zstd, snappy, and gzip.
	CompressionCodec *string `url:"compression_codec,omitempty"`
	// Domain is the domain of the Amazon S3 endpoint.
	Domain *string `url:"domain,omitempty"`
	// Format is a Fastly log format string.
	Format *string `url:"format,omitempty"`
	// FormatVersion is the version of the custom logging format used for the configured endpoint.
	FormatVersion *int `url:"format_version,omitempty"`
	// GzipLevel is the level of gzip encoding when sending logs (default 0, no compression).
	GzipLevel *int `url:"gzip_level,omitempty"`
	// IAMRole is the Amazon Resource Name (ARN) for the IAM role granting Fastly access to S3. Not required if access_key and secret_key are provided.
	IAMRole *string `url:"iam_role,omitempty"`
	// MessageType is how the message should be formatted (classic, loggly, logplex, blank).
	MessageType *string `url:"message_type,omitempty"`
	// Name is the name of the S3 to update.
	Name string `url:"-"`
	// NewName is the new name for the resource.
	NewName *string `url:"name,omitempty"`
	// Path is the path to upload logs to.
	Path *string `url:"path,omitempty"`
	// Period is how frequently log files are finalized so they can be available for reading (in seconds).
	Period *int `url:"period,omitempty"`
	// Placement is where in the generated VCL the logging call should be placed.
	Placement *string `url:"placement,omitempty"`
	// PublicKey is a PGP public key that Fastly will use to encrypt your log files before writing them to disk.
	PublicKey *string `url:"public_key,omitempty"`
	// Redundancy is the S3 redundancy level.
	Redundancy *S3Redundancy `url:"redundancy,omitempty"`
	// ResponseCondition is the name of an existing condition in the configured endpoint, or leave blank to always execute.
	ResponseCondition *string `url:"response_condition,omitempty"`
	// SecretKey is the secret key for your S3 account. Not required if iam_role is provided.
	SecretKey *string `url:"secret_key,omitempty"`
	// ServerSideEncryption should be set to AES256 or aws:kms to enable S3 Server Side Encryption.
	ServerSideEncryption *S3ServerSideEncryption `url:"server_side_encryption,omitempty"`
	// ServerSideEncryptionKMSKeyID is an optional server-side KMS Key ID. Must be set if ServerSideEncryption is set to aws:kms or AES256.
	ServerSideEncryptionKMSKeyID *string `url:"server_side_encryption_kms_key_id,omitempty"`
	// ServiceID is the ID of the service (required).
	ServiceID string `url:"-"`
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int `url:"-"`
	// TimestampFormat is a timestamp format.
	TimestampFormat *string `url:"timestamp_format,omitempty"`
}

// UpdateS3 updates the specified resource.
func (c *Client) UpdateS3(i *UpdateS3Input) (*S3, error) {
	if i.Name == "" {
		return nil, ErrMissingName
	}
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return nil, ErrMissingServiceVersion
	}

	isKMS := i.ServerSideEncryption != nil && *i.ServerSideEncryption == S3ServerSideEncryptionKMS
	hasNoKeyID := i.ServerSideEncryptionKMSKeyID != nil && *i.ServerSideEncryptionKMSKeyID == ""
	if isKMS && hasNoKeyID {
		return nil, ErrMissingServerSideEncryptionKMSKeyID
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/s3/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.PutForm(path, i, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s3 *S3
	if err := decodeBodyMap(resp.Body, &s3); err != nil {
		return nil, err
	}
	return s3, nil
}

// DeleteS3Input is the input parameter to DeleteS3.
type DeleteS3Input struct {
	// Name is the name of the S3 to delete (required).
	Name string
	// ServiceID is the ID of the service (required).
	ServiceID string
	// ServiceVersion is the specific configuration version (required).
	ServiceVersion int
}

// DeleteS3 deletes the specified resource.
func (c *Client) DeleteS3(i *DeleteS3Input) error {
	if i.Name == "" {
		return ErrMissingName
	}
	if i.ServiceID == "" {
		return ErrMissingServiceID
	}
	if i.ServiceVersion == 0 {
		return ErrMissingServiceVersion
	}

	path := fmt.Sprintf("/service/%s/version/%d/logging/s3/%s", i.ServiceID, i.ServiceVersion, url.PathEscape(i.Name))
	resp, err := c.Delete(path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var r *statusResp
	if err := decodeBodyMap(resp.Body, &r); err != nil {
		return err
	}
	if !r.Ok() {
		return ErrNotOK
	}
	return nil
}
