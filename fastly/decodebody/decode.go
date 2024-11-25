package decodebody

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
)

// DecodeBodyMap is used to decode an HTTP response body into a mapstructure struct.
func DecodeBodyMap(body io.Reader, out any) error {
	var parsed any
	dec := json.NewDecoder(body)
	if err := dec.Decode(&parsed); err != nil {
		return err
	}

	return DecodeMap(parsed, out)
}

// DecodeMap decodes an `in` struct or map to a mapstructure tagged `out`.
// It applies the decoder defaults used throughout go-fastly.
// Note that this uses opposite argument order from Go's copy().
func DecodeMap(in, out any) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapToHTTPHeaderHookFunc(),
			stringToTimeHookFunc(),
		),
		WeaklyTypedInput: true,
		Result:           out,
	})
	if err != nil {
		return err
	}
	return decoder.Decode(in)
}

// mapToHTTPHeaderHookFunc returns a function that converts maps into an
// http.Header value.
func mapToHTTPHeaderHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data any,
	) (any, error) {
		if f.Kind() != reflect.Map {
			return data, nil
		}
		if t != reflect.TypeOf(new(http.Header)) {
			return data, nil
		}

		typed, ok := data.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("cannot convert %T to http.Header", data)
		}

		n := map[string][]string{}
		for k, v := range typed {
			switch tv := v.(type) {
			case string:
				n[k] = []string{tv}
			case []string:
				n[k] = tv
			case int, int8, int16, int32, int64:
				n[k] = []string{fmt.Sprintf("%d", tv)}
			case float32, float64:
				n[k] = []string{fmt.Sprintf("%f", tv)}
			default:
				return nil, fmt.Errorf("cannot convert %T to http.Header", v)
			}
		}

		return n, nil
	}
}

// stringToTimeHookFunc returns a function that converts strings to a time.Time
// value.
func stringToTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data any,
	) (any, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(time.Now()) {
			return data, nil
		}

		// Convert it by parsing
		v, err := time.Parse(time.RFC3339, data.(string))
		if err != nil {
			// DictionaryInfo#get uses it's own special time format for now.
			v, _ := data.(string) // type assert to avoid runtime panic (v will have zero value for its type)
			return time.Parse("2006-01-02 15:04:05", v)
		}
		return v, err
	}
}
