package fastly

// String is a helper that returns a pointer to the string value passed in.
func String(v string) *string {
	return &v
}

// Int is a helper that returns a pointer to the int value passed in.
func Int(v int) *int {
	return &v
}

// Uint is a helper that returns a pointer to the uint value passed in.
func Uint(v uint) *uint {
	return &v
}

// Uint8 is a helper that returns a pointer to the uint8 value passed in.
func Uint8(v uint8) *uint8 {
	return &v
}

// Bool is a helper that returns a pointer to the bool value passed in.
func Bool(v bool) *bool {
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
