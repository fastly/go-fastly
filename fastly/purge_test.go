package fastly

import (
	"reflect"
	"testing"
)

func TestROParams(t *testing.T) {
	tests := map[string]struct {
		input  string
		output map[string]string
	}{
		"valid qs": {
			input: "https://www.example.com/?foo=bar&beep=boop",
			output: map[string]string{
				"foo":  "bar",
				"beep": "boop",
			},
		},
		"no qs": {
			input:  "https://www.example.com/",
			output: map[string]string{},
		},
		"key no value": {
			input: "https://www.example.com/?empty&foo=bar",
			output: map[string]string{
				"foo":   "bar",
				"empty": "",
			},
		},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			want := tc.output
			have, err := constructRequestOptionsParam(tc.input)
			if err != nil {
				t.Errorf("got an unexpected error: %s", err)
			}
			if !reflect.DeepEqual(want, have) {
				t.Errorf("want: %+v, have: %+v\n", want, have)
			}
		})
	}
}

// TestClient_Purge validates no runtime panics are raised by the Purge method.
//
// Specifically, we're ensuring that the setting of the `Soft` key to `true`
// (which will require assigning a header to the RequestOptions.Headers field)
// doesn't cause `panic: assignment to entry in nil map`.
func TestClient_Purge(t *testing.T) {
	t.Parallel()

	client := DefaultClient()
	url := "http://gofastly.fastly-testing.com/foo/bar"

	_, err := client.Purge(&PurgeInput{
		URL:  url,
		Soft: true,
	})
	if err == nil {
		t.Fatalf("we've accidentally purged a real URL: %s", url)
	}
}

func TestClient_PurgeKey(t *testing.T) {
	t.Parallel()

	var err error
	var purge *Purge
	Record(t, "purges/purge_by_key", func(c *Client) {
		purge, err = c.PurgeKey(&PurgeKeyInput{
			ServiceID: TestDeliveryServiceID,
			Key:       "foo",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *purge.Status != "ok" {
		t.Error("bad status")
	}
	if purge.PurgeID == nil {
		t.Error("bad id")
	}
}

func TestClient_PurgeKeys(t *testing.T) {
	t.Parallel()

	var err error
	var purges map[string]string
	Record(t, "purges/purge_by_keys", func(c *Client) {
		purges, err = c.PurgeKeys(&PurgeKeysInput{
			ServiceID: TestDeliveryServiceID,
			Keys:      []string{"foo", "bar", "baz"},
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(purges) != 3 {
		t.Error("bad length")
	}
}

func TestClient_PurgeAll(t *testing.T) {
	t.Parallel()

	var err error
	var purge *Purge
	Record(t, "purges/purge_all", func(c *Client) {
		purge, err = c.PurgeAll(&PurgeAllInput{
			ServiceID: TestDeliveryServiceID,
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	if *purge.Status != "ok" {
		t.Error("bad status")
	}
}
