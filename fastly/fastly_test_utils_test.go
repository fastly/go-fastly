package fastly

import (
	"testing"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/stretchr/testify/require"
)

func TestRedactBodyFields(t *testing.T) {
	i := &cassette.Interaction{
		Request: cassette.Request{
			Body: `{"name":"conn","api_key":"sk-realsecret"}`,
		},
		Response: cassette.Response{
			Body: `{"id":"abc123","access_token":"arc_sk_realsecret","user_id":"realuser","other":"keep-me"}`,
		},
	}

	filter := redactBodyFields([]string{"access_token", "user_id", "api_key"})
	require.NoError(t, filter(i))

	require.Equal(t,
		`{"name":"conn","api_key":"`+RedactedFixturePlaceholder+`"}`,
		i.Request.Body,
	)
	require.Equal(t,
		`{"id":"abc123","access_token":"`+RedactedFixturePlaceholder+`","user_id":"`+RedactedFixturePlaceholder+`","other":"keep-me"}`,
		i.Response.Body,
	)
}
