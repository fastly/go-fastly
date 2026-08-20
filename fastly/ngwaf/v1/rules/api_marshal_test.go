package rules

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConditionItem_Marshal_RoundTrip(t *testing.T) {
	assert := require.New(t)

	orig := []byte(`[
	  {"type":"single","field":"ip","operator":"equals","value":"1.2.3.4"},
	  {"type":"group","group_operator":"all","conditions":[
	    {"type":"single","field":"method","operator":"equals","value":"POST"},
	    {"type":"multival","field":"request_header","operator":"exists","group_operator":"any","conditions":[
	      {"type":"single","field":"name","operator":"equals","value":"X-Test"},
	      {"type":"single","field":"value_string","operator":"contains","value":"abc"}
	    ]}
	  ]},
	  {"type":"multival","field":"query_parameter","operator":"does_not_exist","group_operator":"all","conditions":[
	    {"type":"single","field":"name","operator":"equals","value":"debug"}
	  ]}
	]`)

	var parsed []ConditionItem
	assert.NoError(json.Unmarshal(orig, &parsed))

	out, err := json.Marshal(parsed)
	assert.NoError(err)

	// Compare as generic JSON (maps/slices), not as raw bytes.
	var want any
	var got any
	assert.NoError(json.Unmarshal(orig, &want))
	assert.NoError(json.Unmarshal(out, &got))

	assert.Equal(want, got, "round-trip JSON should be identical")
}

func TestGroupConditionItem_Marshal_RoundTrip(t *testing.T) {
	assert := require.New(t)

	orig := []byte(`[
	  {"type":"single","field":"country","operator":"equals","value":"AD"},
	  {"type":"multival","field":"request_cookie","operator":"exists","group_operator":"all","conditions":[
	    {"type":"single","field":"name","operator":"equals","value":"session"},
	    {"type":"single","field":"value_string","operator":"contains","value":"xyz"}
	  ]}
	]`)

	var parsed []GroupConditionItem
	assert.NoError(json.Unmarshal(orig, &parsed))

	out, err := json.Marshal(parsed)
	assert.NoError(err)

	var want any
	var got any
	assert.NoError(json.Unmarshal(orig, &want))
	assert.NoError(json.Unmarshal(out, &got))

	assert.Equal(want, got, "round-trip JSON should be identical")
}

func TestConditionItem_Marshal_DoesNotEmitFieldsWrapper(t *testing.T) {
	assert := require.New(t)

	ci := ConditionItem{
		Type: "single",
		Fields: SingleCondition{
			Field:    "ip",
			Operator: "equals",
			Value:    "1.2.3.4",
		},
	}

	b, err := json.Marshal(ci)
	assert.NoError(err)

	var m map[string]any
	assert.NoError(json.Unmarshal(b, &m))

	_, hasFields := m["Fields"]
	assert.False(hasFields, "marshaled output must not contain Fields")
	assert.Equal("single", m["type"])
	assert.Equal("ip", m["field"])
	assert.Equal("equals", m["operator"])
	assert.Equal("1.2.3.4", m["value"])
}

func TestCreateClientIdentifier_Marshal_SignalPayload(t *testing.T) {
	assert := require.New(t)

	typ := "signal_payload"
	ci := &CreateClientIdentifier{
		Type:   &typ,
		Signal: "site.test",
	}

	b, err := json.Marshal(ci)
	assert.NoError(err)

	var m map[string]any
	assert.NoError(json.Unmarshal(b, &m))

	assert.Equal("signal_payload", m["type"])
	assert.Equal("site.test", m["signal"])
	_, hasKey := m["key"]
	_, hasName := m["name"]
	assert.False(hasKey, "signal_payload identifier must not emit key")
	assert.False(hasName, "signal_payload identifier must not emit name")
}

func TestCreateClientIdentifier_Marshal_OmitsSignalForOtherTypes(t *testing.T) {
	assert := require.New(t)

	typ := "ip"
	ci := &CreateClientIdentifier{Type: &typ}

	b, err := json.Marshal(ci)
	assert.NoError(err)

	var m map[string]any
	assert.NoError(json.Unmarshal(b, &m))

	_, hasSignal := m["signal"]
	assert.False(hasSignal, "ip identifier must not emit an empty signal field")
}

func TestUpdateClientIdentifier_Marshal_SignalPayload(t *testing.T) {
	assert := require.New(t)

	typ := "signal_payload"
	ci := &UpdateClientIdentifier{
		Type:   &typ,
		Signal: "site.test",
	}

	b, err := json.Marshal(ci)
	assert.NoError(err)

	var m map[string]any
	assert.NoError(json.Unmarshal(b, &m))

	assert.Equal("signal_payload", m["type"])
	assert.Equal("site.test", m["signal"])
}

func TestUpdateClientIdentifier_Marshal_OmitsSignalForOtherTypes(t *testing.T) {
	assert := require.New(t)

	typ := "request_header"
	name := "X-Test"
	ci := &UpdateClientIdentifier{Type: &typ, Name: &name}

	b, err := json.Marshal(ci)
	assert.NoError(err)

	var m map[string]any
	assert.NoError(json.Unmarshal(b, &m))

	_, hasSignal := m["signal"]
	assert.False(hasSignal, "request_header identifier must not emit an empty signal field")
}
