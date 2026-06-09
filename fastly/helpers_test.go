package fastly

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestToPointer(t *testing.T) {
	t.Parallel()

	t.Run("string", func(t *testing.T) {
		t.Parallel()

		input := "hello"
		got := ToPointer(input)
		if got == nil {
			t.Fatal("expected non-nil pointer")
		}
		if *got != input {
			t.Errorf("expected %q, got %q", input, *got)
		}
	})

	t.Run("int", func(t *testing.T) {
		t.Parallel()

		input := 42
		got := ToPointer(input)
		if got == nil {
			t.Fatal("expected non-nil pointer")
		}
		if *got != input {
			t.Errorf("expected %d, got %d", input, *got)
		}
	})

	t.Run("bool", func(t *testing.T) {
		t.Parallel()

		input := true
		got := ToPointer(input)
		if got == nil {
			t.Fatal("expected non-nil pointer")
		}
		if *got != input {
			t.Errorf("expected %t, got %t", input, *got)
		}
	})
}

func TestToValue(t *testing.T) {
	t.Parallel()

	t.Run("string", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name  string
			input *string
			want  string
		}{
			{
				name:  "non-nil pointer returns value",
				input: ToPointer("hello"),
				want:  "hello",
			},
			{
				name:  "nil pointer returns zero value",
				input: nil,
				want:  "",
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := ToValue(tt.input)
				if got != tt.want {
					t.Errorf("expected %q, got %q", tt.want, got)
				}
			})
		}
	})

	t.Run("int", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name  string
			input *int
			want  int
		}{
			{
				name:  "non-nil pointer returns value",
				input: ToPointer(42),
				want:  42,
			},
			{
				name:  "nil pointer returns zero value",
				input: nil,
				want:  0,
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := ToValue(tt.input)
				if got != tt.want {
					t.Errorf("expected %d, got %d", tt.want, got)
				}
			})
		}
	})

	t.Run("bool", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name  string
			input *bool
			want  bool
		}{
			{
				name:  "non-nil pointer returns value",
				input: ToPointer(true),
				want:  true,
			},
			{
				name:  "nil pointer returns zero value",
				input: nil,
				want:  false,
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := ToValue(tt.input)
				if got != tt.want {
					t.Errorf("expected %t, got %t", tt.want, got)
				}
			})
		}
	})
}

func TestNullString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantNil bool
	}{
		{
			name:  "non-empty string returns pointer",
			input: "hello",
			want:  "hello",
		},
		{
			name:    "empty string returns nil",
			input:   "",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := NullString(tt.input)

			if tt.wantNil {
				if got != nil {
					t.Fatalf("expected nil, got %q", *got)
				}
				return
			}

			if got == nil {
				t.Fatal("expected non-nil pointer")
			}
			if *got != tt.want {
				t.Errorf("expected %q, got %q", tt.want, *got)
			}
		})
	}
}

func TestNullInt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   int
		want    int
		wantNil bool
	}{
		{
			name:  "positive int returns pointer",
			input: 42,
			want:  42,
		},
		{
			name:  "negative int returns pointer",
			input: -1,
			want:  -1,
		},
		{
			name:    "zero returns nil",
			input:   0,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := NullInt(tt.input)

			if tt.wantNil {
				if got != nil {
					t.Fatalf("expected nil, got %d", *got)
				}
				return
			}

			if got == nil {
				t.Fatal("expected non-nil pointer")
			}
			if *got != tt.want {
				t.Errorf("expected %d, got %d", tt.want, *got)
			}
		})
	}
}

func TestNullableMarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("string", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name  string
			input *Nullable[string]
			want  string
		}{
			{
				name:  "null value serializes to null",
				input: NullValue[string](),
				want:  "null",
			},
			{
				name:  "non-null value serializes to string",
				input: NewNullable("hello"),
				want:  `"hello"`,
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got, err := json.Marshal(tt.input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if string(got) != tt.want {
					t.Errorf("expected %s, got %s", tt.want, string(got))
				}
			})
		}
	})

	t.Run("int", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name  string
			input *Nullable[int]
			want  string
		}{
			{
				name:  "null value serializes to null",
				input: NullValue[int](),
				want:  "null",
			},
			{
				name:  "non-null value serializes to int",
				input: NewNullable(42),
				want:  "42",
			},
		}

		for _, tt := range tests {
			tt := tt
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got, err := json.Marshal(tt.input)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if string(got) != tt.want {
					t.Errorf("expected %s, got %s", tt.want, string(got))
				}
			})
		}
	})
}

func TestToSafeURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input []string
		want  string
	}{
		{
			name:  "normal",
			input: []string{"services", "1234", "detail"},
			want:  "/services/1234/detail",
		},
		{
			name:  "suppress '.'",
			input: []string{"services", ".", "detail"},
			want:  "/services/detail",
		},
		{
			name:  "suppress '..'",
			input: []string{"services", "..", "detail"},
			want:  "/services/detail",
		},
		{
			name:  "encode '/'",
			input: []string{"services", "1234/detail"},
			want:  "/services/1234%2Fdetail",
		},
		{
			name:  "suppress empty components",
			input: []string{"services", "1234", "", "detail"},
			want:  "/services/1234/detail",
		},
		{
			name:  "encode spaces",
			input: []string{"services", "hello world"},
			want:  "/services/hello%20world",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ToSafeURL(tt.input...)
			if got != tt.want {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestGetResponseInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		reader  strings.Reader
		want    infoResponse
		wantErr bool
	}{
		{
			name:   "valid response info",
			reader: *strings.NewReader(`{"links":{"first":"/page/1","last":"/page/3","next":"/page/2"},"meta":{"current_page":1,"per_page":20,"record_count":50,"total_pages":3}}`),
			want: infoResponse{
				Links: paginationInfo{
					First: "/page/1",
					Last:  "/page/3",
					Next:  "/page/2",
				},
				Meta: metaInfo{
					CurrentPage: 1,
					PerPage:     20,
					RecordCount: 50,
					TotalPages:  3,
				},
			},
		},
		{
			name:    "invalid json returns error",
			reader:  *strings.NewReader(`{`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := getResponseInfo(&tt.reader)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Errorf("expected %+v, got %+v", tt.want, got)
			}
		})
	}

	t.Run("reader error returns error", func(t *testing.T) {
		t.Parallel()

		_, err := getResponseInfo(errorReader{})
		if err == nil {
			t.Fatal("expected error")
		}
	})
}

type errorReader struct{}

func (errorReader) Read(_ []byte) (int, error) {
	return 0, errors.New("read error")
}
