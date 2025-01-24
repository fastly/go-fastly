package products

import "github.com/fastly/go-fastly/v9/fastly"

// ProductOutput is an interface used to constrain the 'O' type
// parameters of API operation functions. Use of this interface allows
// the FunctionalTest constructors to apply common validation steps to
// the output received from an API operation, eliminating the need to
// duplicate that validation in every FunctionalTest case, and ensures
// that the API operation functions only accept types which would also
// be accepted by the FunctionalTest constructors.
//
// This interface matches the methods of the NullOutput, EnableOutput
// and ConfigureOutput structs below.
type ProductOutput interface {
	ProductID() string
	ServiceID() string
}

// NullInput is used to indicate to the generic FunctionalTest
// constructors that the operation being tested does not accept an
// Input struct
type NullInput struct{}

// NullOutput is used to indicate to the generic FunctionalTest
// constructors that the operation being tested does not produce an
// Output struct
type NullOutput struct{}

func (o NullOutput) ProductID() string {
	return ""
}

func (o NullOutput) ServiceID() string {
	return ""
}

// EnableOutput represents an enablement response from the Fastly
// API. Products will embed this structure into their own.
type EnableOutput struct {
	Product *EnableOutputNested `mapstructure:"product"`
	Service *EnableOutputNested `mapstructure:"service"`
}

type EnableOutputNested struct {
	Object *string `mapstructure:"object,omitempty"`
	ID     *string `mapstructure:"id,omitempty"`
}

// ProductID returns the ProductID inside an EnableOutput structure.
//
// This method is required, even though the field in the structure is
// exported, because the ProductOutput interface specifies it.
func (o EnableOutput) ProductID() string {
	if o.Product != nil && o.Product.ID != nil {
		return *o.Product.ID
	}
	return ""
}

// ServiceID returns the ServiceID inside an EnableOutput structure.
//
// This method is required, even though the field in the structure is
// exported, because the ProductOutput interface specifies it.
func (o EnableOutput) ServiceID() string {
	if o.Service != nil && o.Service.ID != nil {
		return *o.Service.ID
	}
	return ""
}

func NewEnableOutput(productID, serviceID string) (result EnableOutput) {
	result.Product = &EnableOutputNested{Object: fastly.ToPointer("product"), ID: &productID}
	result.Service = &EnableOutputNested{Object: fastly.ToPointer("service"), ID: &serviceID}
	return
}

// ConfigureOutput represents a configuration response from the Fastly
// API. Products will embed this structure into their own.
type ConfigureOutput struct {
	EnableOutput `mapstructure:",squash"`
}
