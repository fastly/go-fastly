package fastly

// ProductConfiguration represents the base of a response from the
// Fastly API. Configurable products will embed this structure into
// their own product-specific structure.
type ProductConfiguration struct {
	Product *ProductConfigurationNested `mapstructure:"product"`
	Service *ProductConfigurationNested `mapstructure:"service"`
}

type ProductConfigurationNested struct {
	Object    *string `mapstructure:"object,omitempty"`
	ProductID *string `mapstructure:"id,omitempty"`
}
