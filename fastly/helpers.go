package fastly

import "net/url"

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

// ToSafeURL produces a safe (no path traversal, no unsafe characters) URL
// from the path components passed in.
func ToSafeURL(components ...string) string {
	for idx, component := range components {
		if component == ".." {
			components[idx] = ""
		} else {
			components[idx] = url.PathEscape(component)
		}
	}

	// it is safe to ignore the error returned from JoinPath
	// because the only time it will be non-nil is if parsing
	// the base path fails, but that will not fail since it is
	// a constant "/" string
	result, _ := url.JoinPath("/", components...)
	return result
}
