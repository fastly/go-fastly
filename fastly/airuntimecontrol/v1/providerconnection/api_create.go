package providerconnection

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v18/fastly"
)

// CreateInput specifies the information needed for the Create() function to
// perform the operation.
type CreateInput struct {
	// Name is the display name of the provider connection (required).
	Name *string `json:"name"`
	// Models is the AI model identifiers this connection is allowed to use
	// (required).
	Models []string `json:"models"`
	// AuthType is the authentication type for this connection. One of
	// `api-key`, `aws-iam`, `gcp-iam`. Defaults to `api-key`.
	AuthType *string `json:"auth_type,omitempty"`
	// BaseURL is the base URL of the provider's API (required for `api-key`).
	BaseURL *string `json:"base_url,omitempty"`
	// APIKey is the provider's secret API key (required for `api-key`).
	APIKey *string `json:"api_key,omitempty"`
	// CustomerRoleARN is the ARN of the IAM role in the customer's AWS account
	// that AI Runtime Control assumes to access Bedrock (required for
	// `aws-iam`).
	CustomerRoleARN *string `json:"customer_role_arn,omitempty"`
	// Region is the cloud region for API requests: the AWS region for Bedrock
	// (`aws-iam`) or the GCP region for Vertex AI (`gcp-iam`). Required for
	// `aws-iam` and `gcp-iam`.
	Region *string `json:"region,omitempty"`
	// PoolID is the GCP workload identity pool ID (required for `gcp-iam`).
	// 4-32 characters: lowercase letters, digits, and hyphens, starting with a
	// letter.
	PoolID *string `json:"pool_id,omitempty"`
	// WorkloadProviderID is the GCP workload identity provider ID (required for
	// `gcp-iam`). 4-32 characters: lowercase letters, digits, and hyphens,
	// starting with a letter.
	WorkloadProviderID *string `json:"workload_provider_id,omitempty"`
	// ServiceAccountEmail is the email of the GCP service account to
	// impersonate (required for `gcp-iam`).
	ServiceAccountEmail *string `json:"service_account_email,omitempty"`
	// ProjectID is the GCP project ID (required for `gcp-iam`).
	ProjectID *string `json:"project_id,omitempty"`
	// ProjectNumber is the GCP project number (required for `gcp-iam`).
	ProjectNumber *string `json:"project_number,omitempty"`
}

// Create creates a provider connection. Accounts are limited to 25 provider
// connections. Creating a 26th connection returns a 409 response.
func Create(ctx context.Context, c *fastly.Client, i *CreateInput) (*ProviderConnection, error) {
	if i.Name == nil {
		return nil, fastly.ErrMissingName
	}
	if len(i.Models) == 0 {
		return nil, fastly.ErrMissingModels
	}
	var authType string
	if i.AuthType != nil {
		authType = *i.AuthType
	}
	switch authType {
	case "aws-iam":
		if i.CustomerRoleARN == nil {
			return nil, fastly.ErrMissingCustomerRoleARN
		}
		if i.Region == nil {
			return nil, fastly.ErrMissingRegion
		}
	case "gcp-iam":
		if i.PoolID == nil {
			return nil, fastly.ErrMissingPoolID
		}
		if i.WorkloadProviderID == nil {
			return nil, fastly.ErrMissingWorkloadProviderID
		}
		if i.ServiceAccountEmail == nil {
			return nil, fastly.ErrMissingServiceAccountEmail
		}
		if i.ProjectID == nil {
			return nil, fastly.ErrMissingProjectID
		}
		if i.ProjectNumber == nil {
			return nil, fastly.ErrMissingProjectNumber
		}
		if i.Region == nil {
			return nil, fastly.ErrMissingRegion
		}
	default:
		if i.BaseURL == nil {
			return nil, fastly.ErrMissingBaseURL
		}
		if i.APIKey == nil {
			return nil, fastly.ErrMissingAPIKey
		}
	}

	path := fastly.ToSafeURL("ai-runtime-control", "v1", "provider-connections")

	resp, err := c.PostJSON(ctx, path, i, fastly.CreateRequestOptions())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pc *ProviderConnection
	if err := json.NewDecoder(resp.Body).Decode(&pc); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return pc, nil
}
