package fastly

import "testing"

func TestEncodeFormValues(t *testing.T) {
	i := struct {
		A string `form:"a"`
		B string `form:"b.c"`
		C uint   `form:"c,omitempty"`
	}{
		`a\.b.c\.d`,
		`a\.b.c\.d`,
		0,
	}

	s, err := encodeFormValues(i)
	if err != nil {
		t.Fatal(err)
	}

	expected := `a=a%5C.b.c%5C.d&b.c=a%5C.b.c%5C.d`
	if s != expected {
		t.Fatalf("expected %q to be %q", s, expected)
	}
}
