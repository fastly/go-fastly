package workspacealerts

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
)

// ConfigIntegrations specifies the input objects for the config field for Update() operations.
type ConfigIntegrations struct {
	// Datadog is the Datadog integration configuration.
	Datadog *IntegrationDatadog `json:"Datadog,omitempty"`
	// Jira is the Jira integration configuration.
	Jira *IntegrationJira `json:"jira,omitempty"`
	// MailingList is the mailing list integration configuration.
	MailingList *IntegrationMailingList `json:"mailinglist,omitempty"`
	// MicrosoftTeams is the Microsoft Teams integration configuration.
	MicrosoftTeams *IntegrationMicrosoftTeams `json:"microsoftteams,omitempty"`
	// Opsgenie is the Opsgenie integration configuration.
	Opsgenie *IntegrationOpsgenie `json:"Opsgenie,omitempty"`
	// PagerDuty is the PagerDuty integration configuration.
	PagerDuty *IntegrationPagerduty `json:"pagerduty,omitempty"`
	// Slack is the Slack integration configuration.
	Slack *IntegrationSlack `json:"slack,omitempty"`
	// Webook is the webhook integration configuration.
	Webook *IntegrationWebhook `json:"webhook,omitempty"`
}

// IntegrationDatadog is the config object for integration type datadog.
type IntegrationDatadog struct {
	// Key is the Datadog integration key (required).
	Key *string `json:"key,omitempty"`
	// Site is the Datadog site (required).
	Site *string `json:"site,omitempty"`
}

// IntegrationJira is the config object for integration type jira.
type IntegrationJira struct {
	// Host is the name of the Jira instnace (required).
	Host *string `json:"host,omitempty"`
	// IssueType is the Jira issue type associated with the ticket (optional).
	IssueType *string `json:"issue_type,omitempty"`
	// Key is the Jira API key / secret field (required).
	Key *string `json:"key,omitempty"`
	// Project specifies the Jira project where the issue will be created (required).
	Project *string `json:"project,omitempty"`
	// Username is the Jira username of the user who created the ticket (required).
	Username *string `json:"username,omitempty"`
}

// IntegrationMailingList is the config object for integration type mailinglist.
type IntegrationMailingList struct {
	// Address An email address (required).
	Address *string `json:"address,omitempty"`
}

// IntegrationMicrosoftTeams is the config object for integration type microsoftteams.
type IntegrationMicrosoftTeams struct {
	// Webhook is the Microsoft Teams webhook (required).
	Webhook *string `json:"webhook,omitempty"`
}

// IntegrationOpsgenie is the config object for integration type opsgenie.
type IntegrationOpsgenie struct {
	// Key is the Opsgenie integration key (required).
	Key *string `json:"key,omitempty"`
}

// IntegrationPagerduty is the config object for integration type pagerduty.
type IntegrationPagerduty struct {
	// Key is the PagerDuty integration key (required).
	Key *string `json:"key,omitempty"`
}

// IntegrationSlack is the config object for integration type slack.
type IntegrationSlack struct {
	// Webhook is the Slack webhook (required).
	Webhook *string `json:"webhook,omitempty"`
}

// IntegrationWebhook is the config object for integration type webhook.
type IntegrationWebhook struct {
	// Webhook is the Webhook URL (required).
	Webhook *string `json:"webhook,omitempty"`
}

// UpdateInput specifies the information needed for the Put() function to perform
// the operation.
type UpdateInput struct {
	// AlertID is The unique identifier of the workspace alert (required).
	AlertID *string
	// Config is the configuration associated with the workspace integration (required).
	Config ConfigIntegrations `json:"-"`
	// Context, if supplied, will be used as the Request's context.
	Context *context.Context `json:"-"`
	// Events is a list of event types (required).
	Events *string `json:"-"`
	// Type is the type of the workspace integration.
	Type string `json:"-"`
	// WorkspaceID is the workspace identifier (required).
	WorkspaceID *string `json:"-"`
}

// Update updates the specified virtual patch.
func Update(c *fastly.Client, i *UpdateInput) (*WorkspaceAlert, error) {
	if i.AlertID == nil {
		return nil, fastly.ErrMissingAlertID
	}

	if (i.Config == ConfigIntegrations{}) {
		return nil, fastly.ErrMissingConfig
	}

	if i.WorkspaceID == nil {
		return nil, fastly.ErrMissingWorkspaceID
	}

	// Validate integration configuration
	if i.Config.Datadog != nil {
		if i.Config.Datadog.Key == nil {
			return nil, fastly.ErrMissingKey
		}
		if i.Config.Datadog.Site == nil {
			return nil, fastly.ErrMissingSite
		}
	}

	path := fastly.ToSafeURL("ngwaf", "v1", "workspaces", *i.WorkspaceID, "virtual-patches", *i.VirtualPatchID)

	resp, err := c.PatchJSON(path, i, fastly.CreateRequestOptions(i.Context))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var vp *WorkspaceAlert
	if err := json.NewDecoder(resp.Body).Decode(&vp); err != nil {
		return nil, fmt.Errorf("failed to decode json response: %w", err)
	}

	return vp, nil
}
