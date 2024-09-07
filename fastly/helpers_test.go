package fastly

import (
	"testing"
)

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

		expected := "/detail"

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
