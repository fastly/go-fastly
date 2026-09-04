package fastly

import (
	"testing"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/stretchr/testify/require"
)

func TestRedactResponseFields(t *testing.T) {
	i := &cassette.Interaction{
		Response: cassette.Response{
			Body: `{"id":"abc123","access_token":"arc_sk_realsecret","user_id":"realuser","other":"keep-me"}`,
		},
	}

	filter := redactResponseFields([]string{"access_token", "user_id"})
	require.NoError(t, filter(i))

	require.Equal(t,
		`{"id":"abc123","access_token":"`+RedactedFixturePlaceholder+`","user_id":"`+RedactedFixturePlaceholder+`","other":"keep-me"}`,
		i.Response.Body,
	)
}
