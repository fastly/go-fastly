package fastly

import (
	"testing"
)

func TestNullString(t *testing.T) {
	t.Parallel()

	t.Run("non-empty string returns pointer", func(t *testing.T) {
		v := "hello"
		result := NullString(v)
		if result == nil {
			t.Fatal("expected non-nil pointer")
		}
		if *result != v {
			t.Errorf("expected %q, got %q", v, *result)
		}
	})

	t.Run("empty string returns nil", func(t *testing.T) {
		result := NullString("")
		if result != nil {
			t.Errorf("expected nil, got %q", *result)
		}
	})
}

func TestNullInt(t *testing.T) {
	t.Parallel()

	t.Run("non-zero int returns pointer", func(t *testing.T) {
		v := 42
		result := NullInt(v)
		if result == nil {
			t.Fatal("expected non-nil pointer")
		}
		if *result != v {
			t.Errorf("expected %d, got %d", v, *result)
		}
	})

	t.Run("zero returns nil", func(t *testing.T) {
		result := NullInt(0)
		if result != nil {
			t.Errorf("expected nil, got %d", *result)
		}
	})
}

func TestToSafeURL(t *testing.T) {
	t.Parallel()

	t.Run("normal", func(t *testing.T) {
		path := ToSafeURL("services", "1234", "detail")

		expected := "/services/1234/detail"

		if path != expected {
			t.Errorf("expected \n\n%q\n\n to be \n\n%q\n\n", path, expected)
		}
	})

	t.Run("suppress '.'", func(t *testing.T) {
		path := ToSafeURL("services", ".", "detail")

		expected := "/services/detail"

		if path != expected {
			t.Errorf("expected \n\n%q\n\n to be \n\n%q\n\n", path, expected)
		}
	})

	t.Run("suppress '..'", func(t *testing.T) {
		path := ToSafeURL("services", "..", "detail")

		expected := "/services/detail"

		if path != expected {
			t.Errorf("expected \n\n%q\n\n to be \n\n%q\n\n", path, expected)
		}
	})

	t.Run("encode '/'", func(t *testing.T) {
		path := ToSafeURL("services", "1234/detail")

		expected := "/services/1234%2Fdetail"

		if path != expected {
			t.Errorf("expected \n\n%q\n\n to be \n\n%q\n\n", path, expected)
		}
	})

	t.Run("suppress empty components", func(t *testing.T) {
		path := ToSafeURL("services", "1234", "", "detail")

		expected := "/services/1234/detail"

		if path != expected {
			t.Errorf("expected \n\n%q\n\n to be \n\n%q\n\n", path, expected)
		}
	})
}
