package dto

import "github.com/google/uuid"

type CheckJob struct {
	APIID   uuid.UUID `json:"api_id"`
	ApiName string    `json:"api_name"`

	URL    string `json:"url"`
	Method string `json:"method"`

	// Auth
	AuthType  string  `json:"auth_type"`            // none, bearer, api-key
	AuthIn    *string `json:"auth_in,omitempty"`    // header, query
	AuthKey   *string `json:"auth_key,omitempty"`   // Authorization, x-api-key
	AuthValue *string `json:"auth_value,omitempty"` // token / key

	// Request
	Headers  map[string]string `json:"headers,omitempty"`
	BodyType *string           `json:"body_type,omitempty"` // json, form-data, none
	Body     any               `json:"body,omitempty"`

	// Execution
	TimeoutMs int `json:"timeout_ms"`

	// Expectations
	ExpectedStatusCodes    []int   `json:"expected_status_codes,omitempty"`
	ExpectedResponseTimeMs *int    `json:"expected_response_time_ms,omitempty"`
	ExpectedBodyContains   *string `json:"expected_body_contains,omitempty"`
}
