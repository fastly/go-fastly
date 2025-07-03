package workspacealerts

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fastly/go-fastly/v10/fastly"
	"github.com/fastly/go-fastly/v10/fastly/ngwaf/v1/workspacealerts/datadog"
	"github.com/fastly/go-fastly/v10/fastly/ngwaf/v1/workspacealerts/jira"
	"github.com/fastly/go-fastly/v10/fastly/ngwaf/v1/workspacealerts/mailinglist"
	"github.com/fastly/go-fastly/v10/fastly/ngwaf/v1/workspacealerts/microsoftteams"
	"github.com/fastly/go-fastly/v10/fastly/ngwaf/v1/workspacealerts/opsgenie"
	"github.com/fastly/go-fastly/v10/fastly/ngwaf/v1/workspacealerts/pagerduty"
	"github.com/fastly/go-fastly/v10/fastly/ngwaf/v1/workspacealerts/slack"
	"github.com/fastly/go-fastly/v10/fastly/ngwaf/v1/workspacealerts/webhook"
)

// ConfigIntegrations specifies the input objects for the config field for Update() operations.
type ConfigIntegrations struct {
	// Datadog is the Datadog integration configuration.
	Datadog *datadog.Integration `json:"Datadog,omitempty"`
	// Jira is the Jira integration configuration.
	Jira *jira.Integration `json:"jira,omitempty"`
	// MailingList is the mailing list integration configuration.
	MailingList *mailinglist.Integration `json:"mailinglist,omitempty"`
	// MicrosoftTeams is the Microsoft Teams integration configuration.
	MicrosoftTeams *microsoftteams.Integration `json:"microsoftteams,omitempty"`
	// Opsgenie is the Opsgenie integration configuration.
	Opsgenie *opsgenie.Integration `json:"Opsgenie,omitempty"`
	// PagerDuty is the PagerDuty integration configuration.
	PagerDuty *pagerduty.Integration `json:"pagerduty,omitempty"`
	// Slack is the Slack integration configuration.
	Slack *slack.Integration `json:"slack,omitempty"`
	// Webook is the webhook integration configuration.
	Webook *webhook.Integration `json:"webhook,omitempty"`
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
	if i.Config.Jira != nil {
		if i.Config.Jira.Host == nil {
			return nil, fastly.ErrMissingHost
		}
		if i.Config.Jira.Key == nil {
			return nil, fastly.ErrMissingKey
		}
		if i.Config.Jira.Project == nil {
			return nil, fastly.ErrMissingProject
			}
		if i.Config.Jira.Username == nil {
			return nil, fastly.ErrMissingUserName
		}
	}
	if i.Config.MailingList != nil {
		if i.Config.MailingList.Address == nil {
			return nil, fastly.ErrMissingAddress
		}
	}
	if i.Config.MicrosoftTeams != nil {
		if i.Config.MicrosoftTeams.Webhook == nil {
			return nil, fastly.ErrMissingWebhook
		}
	}
	if i.Config.Opsgenie != nil {
		if i.Config.Opsgenie.Key == nil {
			return nil, fastly.ErrMissingKey
		}
	}

	// To-Do - update rest of the below code
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
