package products

// ServiceProductOutput is an interface used to constrain the 'O' type
// parameters of API operation functions. Use of this interface allows
// the FunctionalTestService constructors to apply common validation
// steps to the output received from an API operation, eliminating the
// need to duplicate that validation in every FunctionalTestService
// case, and ensures that the API operation functions only accept
// types which would also be accepted by the FunctionalTestService
// constructors.
//
// This interface matches the methods of the ServiceNullOutput,
// ServiceEnableOutput and ServiceConfigureOutput structs below.
type ServiceProductOutput interface {
	ProductID() string
	ServiceID() string
}

// ServiceNullInput is used to indicate to the generic FunctionalTestService
// constructors that the operation being tested does not accept an
// Input struct.
type ServiceNullInput struct{}

// ServiceNullOutput is used to indicate to the generic FunctionalTestService
// constructors that the operation being tested does not produce an
// Output struct.
type ServiceNullOutput struct{}

func (o ServiceNullOutput) ProductID() string {
	return ""
}

func (o ServiceNullOutput) ServiceID() string {
	return ""
}

// ServiceEnableOutput represents an enablement response from the Fastly
// API. Products will embed this structure into their own.
type ServiceEnableOutput struct {
	Product *ServiceEnableOutputNested `mapstructure:"product"`
	Service *ServiceEnableOutputNested `mapstructure:"service"`
}

type ServiceEnableOutputNested struct {
	Object *string `mapstructure:"object,omitempty"`
	ID     *string `mapstructure:"id,omitempty"`
}

// ProductID returns the ProductID inside an ServiceEnableOutput structure.
//
// This method is required, even though the field in the structure is
// exported, because the ServiceProductOutput interface specifies it.
func (o ServiceEnableOutput) ProductID() string {
	if o.Product != nil && o.Product.ID != nil {
		return *o.Product.ID
	}
	return ""
}

// ServiceID returns the ServiceID inside an ServiceEnableOutput structure.
//
// This method is required, even though the field in the structure is
// exported, because the ServiceProductOutput interface specifies it.
func (o ServiceEnableOutput) ServiceID() string {
	if o.Service != nil && o.Service.ID != nil {
		return *o.Service.ID
	}
	return ""
}

func NewServiceEnableOutput(productID, serviceID string) (result ServiceEnableOutput) {
	result.Product = &ServiceEnableOutputNested{Object: new("product"), ID: &productID}
	result.Service = &ServiceEnableOutputNested{Object: new("service"), ID: &serviceID}
	return
}

// ServiceConfigureOutput represents a configuration response from the Fastly
// API. Products will embed this structure into their own.
type ServiceConfigureOutput struct {
	ServiceEnableOutput `mapstructure:",squash"`
}
