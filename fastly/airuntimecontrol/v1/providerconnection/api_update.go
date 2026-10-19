package providerconnection

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v17/fastly"
)

// UpdateInput specifies the information needed for the Update() function to
// perform the operation.
type UpdateInput struct {
	// ID is the ID of the provider connection (required).
	ID *string `json:"-"`
	// Models is the AI model identifiers this connection is allowed to use.
	Models []string `json:"models,omitempty"`
	// BaseURL is the base URL of the provider's API. Present only for `api-key`
	// connections.
	BaseURL *string `json:"base_url,omitempty"`
	// APIKey is the provider's secret API key.
	APIKey *string `json:"api_key,omitempty"`
	// CustomerRoleARN is the ARN of the IAM role in the customer's AWS account
	// that AI Runtime Control assumes to access Bedrock. Applies only to `aws-iam`
	// connections.
	CustomerRoleARN *string `json:"customer_role_arn,omitempty"`
	// Region is the cloud region for API requests. Applies only to `aws-iam`
	// and `gcp-iam` connections.
	Region *string `json:"region,omitempty"`
	// PoolID is the GCP workload identity pool ID. Applies only to `gcp-iam`
	// connections.
	PoolID *string `json:"pool_id,omitempty"`
	// WorkloadProviderID is the GCP workload identity provider ID. Applies only
	// to `gcp-iam` connections.
	WorkloadProviderID *string `json:"workload_provider_id,omitempty"`
	// ServiceAccountEmail is the email of the GCP service account to
	// impersonate. Applies only to `gcp-iam` connections.
	ServiceAccountEmail *string `json:"service_account_email,omitempty"`
	// ProjectID is the GCP project ID. Applies only to `gcp-iam` connections.
	ProjectID *string `json:"project_id,omitempty"`
	// ProjectNumber is the GCP project number. Applies only to `gcp-iam`
	// connections.
	ProjectNumber *string `json:"project_number,omitempty"`
}

// Update updates the attributes of a provider connection. Only the fields you
// include are changed.
func Update(ctx context.Context, c *fastly.Client, i *UpdateInput) (*ProviderConnection, error) {
	if i.ID == nil {
		return nil, fastly.ErrMissingID
	}

	path := fastly.ToSafeURL("ai-runtime-control", "v1", "provider-connections", *i.ID)

	resp, err := c.PatchJSON(ctx, path, i, fastly.CreateRequestOptions())
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
