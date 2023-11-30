package fastly

// ToPointer converts T to *T.
func ToPointer[T ~string | ~int | ~int32 | ~int64 | uint | uint8 | ~bool](v T) *T {
	return &v
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
