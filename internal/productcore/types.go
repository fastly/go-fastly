package productcore

// ProductOutput is an interface used to constrain the 'O' type
// parameter of the functions in this package. Use of this interface
// allows the FunctionalTest constructors to apply common validation
// steps to the output received from an API operation, eliminating the
// need to duplicate that validation in every FunctionalTest case, and
// ensures that the API operation functions only accept types which
// would also be accepted by the FunctionalTest constructors.
//
// This interface matches the methods of the NullOutput, EnableOutput
// and ConfigureOutput structs below.
type ProductOutput interface {
	ProductID() *string
	ServiceID() *string
}

// NullInput is used to indicate to the generic FunctionalTest
// constructors that the operation being tested does not accept an
// Input struct
type NullInput struct{}

// NullOutput is used to indicate to the generic FunctionalTest
// constructors that the operation being tested does not produce an
// Output struct
type NullOutput struct{}

func (o *NullOutput) ProductID() *string {
	return nil
}

func (o *NullOutput) ServiceID() *string {
	return nil
}

// EnableOutput represents an enablement response from the Fastly
// API. Products will embed this structure into their own.
type EnableOutput struct {
	Product *enableOutputNested `mapstructure:"product"`
	Service *enableOutputNested `mapstructure:"service"`
}

type enableOutputNested struct {
	Object *string `mapstructure:"object,omitempty"`
	ID     *string `mapstructure:"id,omitempty"`
}

// GetProductID return the ProductID inside an EnableOutput structure.
//
// This method is required, even though the field in the structure is
// exported, because there is an interface in an internal package
// which expects 'output' types to provide this method.
func (o *EnableOutput) ProductID() *string {
	return o.Product.ID
}

// GetServiceID return the ServiceID inside an EnableOutput structure.
//
// This method is required, even though the field in the structure is
// exported, because there is an interface in an internal package
// which expects 'output' types to provide this method.
func (o *EnableOutput) ServiceID() *string {
	return o.Service.ID
}

// ConfigureOutput represents a configuration response from the Fastly
// API. Products will embed this structure into their own.
type ConfigureOutput struct {
	EnableOutput `mapstructure:",squash"`
}
