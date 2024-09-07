package fastly

// MultiConstraint is a generic constraint for ToPointer/ToValue.
type MultiConstraint interface {
	[]string | ~string | ~int | int32 | ~int64 | uint | uint8 | uint32 | uint64 | float64 | ~bool
}

// ToPointer converts T to *T.
func ToPointer[T MultiConstraint](v T) *T {
	return &v
}

// ToValue converts *T to T.
// If v is nil, then return T's zero value.
func ToValue[T MultiConstraint](v *T) T {
	if v != nil {
		return *v
	}
	var zero T
	return zero
}

// NullString is a helper that returns a pointer to the string value passed in
// or nil if the string is empty.
//
// NOTE: historically this has only been utilized by
// https://github.com/fastly/terraform-provider-fastly
func NullString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
