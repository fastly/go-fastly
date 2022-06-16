package fastly

import (
	"fmt"
)

type ServiceAuthorization struct {
	ID         string `mapstructure:"id"`
	Permission string `mapstructure:"permission"`
	UserID     string `mapstructure:"relationships.user.data.id"`
	ServiceID  string `mapstructure:"relationships.service.data.id"`
	Type       string `mapstructure:"type"`
}

type serviceAuthorizationAttributesModel struct {
	Permission string `mapstructure:"permission,omitempty" json:"permission,omitempty"`
}

type serviceAuthorizationRelationshipDataModel struct {
	ID   string `mapstructure:"id" json:"id,omitempty"`
	Type string `mapstructure:"type" json:"type,omitempty"`
}

type serviceAuthorizationRelationshipsAPIModel struct {
	Data serviceAuthorizationRelationshipDataModel `mapstructure:"data" json:"data,omitempty"`
}

type serviceAuthorizationRelationshipsModel struct {
	Service serviceAuthorizationRelationshipsAPIModel `mapstructure:"service" json:"service,omitempty"`
	User    serviceAuthorizationRelationshipsAPIModel `mapstructure:"user" json:"user,omitempty"`
}

type serviceAuthorizationDataModel struct {
	ID            string                                  `mapstructure:"id,omitempty" json:"id,omitempty"`
	Attributes    serviceAuthorizationAttributesModel     `mapstructure:"attributes" json:"attributes,omitempty"`
	Relationships *serviceAuthorizationRelationshipsModel `mapstructure:"relationships" json:"relationships,omitempty"`
	Type          string                                  `mapstructure:"type" json:"type,omitempty"`
}

type serviceAuthorizationAPIModel struct {
	Data serviceAuthorizationDataModel `mapstructure:"data" json:"data"`
}

func (a serviceAuthorizationAPIModel) ToServiceAuthorization() *ServiceAuthorization {
	return &ServiceAuthorization{
		ID:         a.Data.ID,
		Permission: a.Data.Attributes.Permission,
		UserID:     a.Data.Relationships.User.Data.ID,
		ServiceID:  a.Data.Relationships.Service.Data.ID,
		Type:       a.Data.Type,
	}
}

type GetServiceAuthorizationInput struct {
	ID string
}

func (c *Client) GetServiceAuthorization(i *GetServiceAuthorizationInput) (*ServiceAuthorization, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	path := fmt.Sprintf("/service-authorizations/%s", i.ID)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var a *serviceAuthorizationAPIModel
	if err := decodeBodyMap(resp.Body, &a); err != nil {
		return nil, err
	}

	return a.ToServiceAuthorization(), nil
}

type CreateServiceAuthorizationInput struct {
	ServiceID  string
	UserID     string
	Permission string
}

func (c *Client) CreateServiceAuthorization(i *CreateServiceAuthorizationInput) (*ServiceAuthorization, error) {
	if i.ServiceID == "" {
		return nil, ErrMissingServiceID
	}
	if i.UserID == "" {
		return nil, ErrMissingUserID
	}

	out := serviceAuthorizationAPIModel{
		Data: serviceAuthorizationDataModel{
			Attributes: serviceAuthorizationAttributesModel{
				Permission: i.Permission,
			},
			Relationships: &serviceAuthorizationRelationshipsModel{
				Service: serviceAuthorizationRelationshipsAPIModel{
					Data: serviceAuthorizationRelationshipDataModel{
						ID:   i.ServiceID,
						Type: "service",
					},
				},
				User: serviceAuthorizationRelationshipsAPIModel{
					Data: serviceAuthorizationRelationshipDataModel{
						ID:   i.UserID,
						Type: "user",
					},
				},
			},
			Type: "service_authorization",
		},
	}

	resp, err := c.PostJSON("/service-authorizations", &out, nil)
	if err != nil {
		return nil, err
	}

	var a *serviceAuthorizationAPIModel
	if err := decodeBodyMap(resp.Body, &a); err != nil {
		return nil, err
	}

	return a.ToServiceAuthorization(), nil
}

type UpdateServiceAuthorizationInput struct {
	ID          string
	Permissions string
}

func (c *Client) UpdateServiceAuthorization(i *UpdateServiceAuthorizationInput) (*ServiceAuthorization, error) {
	if i.ID == "" {
		return nil, ErrMissingID
	}

	if i.Permissions == "" {
		return nil, ErrMissingPermissions
	}

	out := &serviceAuthorizationAPIModel{
		Data: serviceAuthorizationDataModel{
			ID:   i.ID,
			Type: "service_authorization",
			Attributes: serviceAuthorizationAttributesModel{
				Permission: i.Permissions,
			},
		},
	}

	path := fmt.Sprintf("/service-authorizations/%s", i.ID)
	resp, err := c.PatchJSON(path, out, nil)
	if err != nil {
		return nil, err
	}

	var a *serviceAuthorizationAPIModel
	if err := decodeBodyMap(resp.Body, &a); err != nil {
		return nil, err
	}

	return a.ToServiceAuthorization(), nil
}

type DeleteServiceAuthorizationInput struct {
	ID string
}

func (c *Client) DeleteServiceAuthorization(i *DeleteServiceAuthorizationInput) error {
	if i.ID == "" {
		return ErrMissingID
	}

	path := fmt.Sprintf("/service-authorizations/%s", i.ID)
	_, err := c.Delete(path, nil)

	return err
}
