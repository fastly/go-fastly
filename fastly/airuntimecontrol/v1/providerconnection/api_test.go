package providerconnection

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fastly/go-fastly/v17/fastly"
)

func TestClient_ProviderConnections(t *testing.T) {
	ctx := context.TODO()

	var err error

	// Create a provider connection.
	var created *ProviderConnection
	fastly.RecordRedacted(t, "create", []string{"api_key"}, func(c *fastly.Client) {
		created, err = Create(ctx, c, &CreateInput{
			Name:    new("go-fastly-test-connection"),
			Models:  []string{"claude-opus-4-7"},
			BaseURL: new("https://api.anthropic.com"),
			// `sk-go-fastly-test` is a placeholder. The API does not validate the provider key when the
			// connection is created, so a real Anthropic API key is not needed to re-record these tests.
			APIKey: new("sk-go-fastly-test"),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	require.NotEmpty(t, created.ID)
	require.NotEmpty(t, created.SecretID)
	require.Equal(t, "api-key", created.AuthType)

	connID := created.ID

	defer func() {
		fastly.Record(t, "delete", func(c *fastly.Client) {
			_ = Delete(ctx, c, &DeleteInput{
				ID: new(connID),
			})
		})
	}()

	// Get the provider connection.
	var fetched *ProviderConnection
	fastly.Record(t, "get", func(c *fastly.Client) {
		fetched, err = Get(ctx, c, &GetInput{
			ID: new(connID),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, fetched)
	require.Equal(t, connID, fetched.ID)
	require.NotEmpty(t, fetched.SecretID)
	require.Equal(t, "api-key", fetched.AuthType)
	require.Empty(t, fetched.VirtualKeyIDs)

	// List provider connections.
	var conns *ProviderConnections
	fastly.Record(t, "list", func(c *fastly.Client) {
		conns, err = List(ctx, c, &ListInput{
			IncludeDeleted: new(true),
		})
	})
	require.NoError(t, err)
	require.NotNil(t, conns)

	var foundCreated bool
	for _, pc := range conns.Data {
		if pc.ID == connID {
			foundCreated = true
			require.Nil(t, pc.DeletedAt)
			break
		}
	}
	require.True(t, foundCreated, "expected the created connection in the list results")

	// Update the provider connection.
	var updated *ProviderConnection
	fastly.Record(t, "update", func(c *fastly.Client) {
		updated, err = Update(ctx, c, &UpdateInput{
			ID:     new(connID),
			Models: []string{"claude-opus-4-6", "claude-sonnet-4-6"},
		})
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, connID, updated.ID)
	require.NotEmpty(t, updated.SecretID)
}

func TestClient_Create_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := Create(ctx, fastly.TestClient, &CreateInput{
		Models:  []string{"gpt-4"},
		BaseURL: new("https://api.openai.com/v1"),
		APIKey:  new("secret"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingName)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name:    new("OpenAI"),
		BaseURL: new("https://api.openai.com/v1"),
		APIKey:  new("secret"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingModels)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name:   new("OpenAI"),
		Models: []string{"gpt-4"},
		APIKey: new("secret"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingBaseURL)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name:    new("OpenAI"),
		Models:  []string{"gpt-4"},
		BaseURL: new("https://api.openai.com/v1"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingAPIKey)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name:     new("bedrock"),
		Models:   []string{"anthropic.claude-sonnet-4"},
		AuthType: new("aws-iam"),
		Region:   new("us-east-1"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingCustomerRoleARN)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name:            new("bedrock"),
		Models:          []string{"anthropic.claude-sonnet-4"},
		AuthType:        new("aws-iam"),
		CustomerRoleARN: new("arn:aws:iam::123456789012:role/arc"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingRegion)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name:                new("vertex"),
		Models:              []string{"gemini-2.5-pro"},
		AuthType:            new("gcp-iam"),
		WorkloadProviderID:  new("arc-provider"),
		ServiceAccountEmail: new("arc@my-project.iam.gserviceaccount.com"),
		ProjectID:           new("my-project"),
		ProjectNumber:       new("123456789012"),
		Region:              new("us-central1"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingPoolID)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name:                new("vertex"),
		Models:              []string{"gemini-2.5-pro"},
		AuthType:            new("gcp-iam"),
		PoolID:              new("arc-pool"),
		ServiceAccountEmail: new("arc@my-project.iam.gserviceaccount.com"),
		ProjectID:           new("my-project"),
		ProjectNumber:       new("123456789012"),
		Region:              new("us-central1"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingWorkloadProviderID)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name:               new("vertex"),
		Models:             []string{"gemini-2.5-pro"},
		AuthType:           new("gcp-iam"),
		PoolID:             new("arc-pool"),
		WorkloadProviderID: new("arc-provider"),
		ProjectID:          new("my-project"),
		ProjectNumber:      new("123456789012"),
		Region:             new("us-central1"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingServiceAccountEmail)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name:                new("vertex"),
		Models:              []string{"gemini-2.5-pro"},
		AuthType:            new("gcp-iam"),
		PoolID:              new("arc-pool"),
		WorkloadProviderID:  new("arc-provider"),
		ServiceAccountEmail: new("arc@my-project.iam.gserviceaccount.com"),
		ProjectNumber:       new("123456789012"),
		Region:              new("us-central1"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingProjectID)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name:                new("vertex"),
		Models:              []string{"gemini-2.5-pro"},
		AuthType:            new("gcp-iam"),
		PoolID:              new("arc-pool"),
		WorkloadProviderID:  new("arc-provider"),
		ServiceAccountEmail: new("arc@my-project.iam.gserviceaccount.com"),
		ProjectID:           new("my-project"),
		Region:              new("us-central1"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingProjectNumber)

	_, err = Create(ctx, fastly.TestClient, &CreateInput{
		Name:                new("vertex"),
		Models:              []string{"gemini-2.5-pro"},
		AuthType:            new("gcp-iam"),
		PoolID:              new("arc-pool"),
		WorkloadProviderID:  new("arc-provider"),
		ServiceAccountEmail: new("arc@my-project.iam.gserviceaccount.com"),
		ProjectID:           new("my-project"),
		ProjectNumber:       new("123456789012"),
	})
	require.ErrorIs(t, err, fastly.ErrMissingRegion)
}

func TestClient_GetBedrockSetup(t *testing.T) {
	ctx := context.TODO()

	var (
		setup *BedrockSetup
		err   error
	)
	fastly.RecordRedacted(t, "get_bedrock_setup", []string{"arc_role_arn", "external_id"}, func(c *fastly.Client) {
		setup, err = GetBedrockSetup(ctx, c)
	})
	require.NoError(t, err)
	require.NotNil(t, setup)
	require.NotEmpty(t, setup.ArcRoleARN)
	require.NotEmpty(t, setup.ExternalID)
}

func TestClient_GetVertexSetup(t *testing.T) {
	ctx := context.TODO()

	var (
		setup *VertexSetup
		err   error
	)
	fastly.RecordRedacted(t, "get_vertex_setup", []string{"oidc_issuer_url", "audience", "customer_id"}, func(c *fastly.Client) {
		setup, err = GetVertexSetup(ctx, c)
	})
	require.NoError(t, err)
	require.NotNil(t, setup)
	require.NotEmpty(t, setup.OIDCIssuerURL)
	require.NotEmpty(t, setup.Audience)
	require.NotEmpty(t, setup.CustomerID)
}

func TestClient_Get_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := Get(ctx, fastly.TestClient, &GetInput{ID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingID)
}

func TestClient_Update_validation(t *testing.T) {
	ctx := context.TODO()

	_, err := Update(ctx, fastly.TestClient, &UpdateInput{ID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingID)
}

func TestClient_Delete_validation(t *testing.T) {
	ctx := context.TODO()

	err := Delete(ctx, fastly.TestClient, &DeleteInput{ID: nil})
	require.ErrorIs(t, err, fastly.ErrMissingID)
}
