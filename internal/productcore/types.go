package productcore

// ProductOutput is an interface used to constrain the 'O' type
// parameter of the functions in this package. Use of this interface
// allows the FunctionalTest constructors to apply common validation
// steps to the output received from an API operation, eliminating the
// need to duplicate that validation in every FunctionalTest case, and
// ensures that the API operation functions only accept types which
// would also be accepted by the FunctionalTest constructors.
//
// This interface matches the methods of the EnableOutput and
// ConfigureOutput structs in the fastly/products package.
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
