package products

// ProductOutput is an interface used to constrain the 'O' type
// parameter of the FunctionalTest constructors in this package. Use
// of this interface allows those constructors to apply common
// validation steps to the output received from an API operation,
// eliminating the need to duplicate that validation in every
// FunctionalTest case.
//
// This interface matches the methods of the EnableOutput and
// ConfigureOutput structs in the fastly/products package.
type ProductOutput interface {
	ProductID() *string
	ServiceID() *string
}
