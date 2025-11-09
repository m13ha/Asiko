package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type apiErrorPayload struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Fields  []struct {
		Field   string `json:"field"`
		Message string `json:"message"`
		Rule    string `json:"rule"`
	} `json:"fields,omitempty"`
	RequestID string `json:"request_id"`
}

func decodeAPIError(t *testing.T, body []byte) apiErrorPayload {
	t.Helper()
	var payload apiErrorPayload
	require.NoError(t, json.Unmarshal(body, &payload), "response body: %s", string(body))
	return payload
}
