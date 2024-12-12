package products

// EnableOutput represents an enablement response from the Fastly
// API. Some products will use this structure directly, others will
// embed it into their own structure.
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
func (o EnableOutput) ProductID() *string {
	return o.Product.ID
}

// GetServiceID return the ServiceID inside an EnableOutput structure.
//
// This method is required, even though the field in the structure is
// exported, because there is an interface in an internal package
// which expects 'output' types to provide this method.
func (o EnableOutput) ServiceID() *string {
	return o.Service.ID
}

// ConfigureOutput represents a configuration response from the Fastly
// API. Products will embed this into their own structure.
type ConfigureOutput struct {
	EnableOutput `mapstructure:",squash"`
}
