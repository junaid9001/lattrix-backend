package dto

type ApiRegisterDto struct {
	Name        string  `json:"name" validate:"required,min=2,max=50"`
	Description *string `json:"description"`

	URL    string `json:"url" validate:"required,url"`
	Method string `json:"method" validate:"required,oneof=GET POST PUT DELETE PATCH"`

	AuthType  string  `json:"auth_type" validate:"required,oneof=NONE BEARER API_KEY"`
	AuthIn    *string `json:"auth_in" validate:"oneof=HEADER QUERY"`
	AuthKey   *string `json:"auth_key"`
	AuthValue *string `json:"auth_value"`

	Headers  map[string]string `json:"headers"`
	BodyType *string           `json:"body_type" validate:"oneof=JSON FORM NONE"`
	Body     map[string]any    `json:"body"`

	IntervalSeconds *int `json:"interval_seconds" validate:"min=10"`
	TimeoutMs       *int `json:"timeout_ms" validate:"min=1000"`

	ExpectedStatusCodes    []int   `json:"expected_status_codes"`
	ExpectedResponseTimeMs *int    `json:"expected_response_time_ms"`
	ExpectedBodyContains   *string `json:"expected_body_contains"`
}

type ApiUpdateDto struct {
	Name        *string `json:"name" validate:"omitempty,min=2,max=50"`
	Description *string `json:"description,omitempty"`

	URL    *string `json:"url" validate:"omitempty,url"`
	Method *string `json:"method" validate:"omitempty,oneof=GET POST PUT DELETE PATCH"`

	AuthType  *string `json:"auth_type" validate:"omitempty,oneof=NONE BEARER API_KEY"`
	AuthIn    *string `json:"auth_in" validate:"omitempty,oneof=HEADER QUERY"`
	AuthKey   *string `json:"auth_key,omitempty"`
	AuthValue *string `json:"auth_value,omitempty"`

	Headers  *map[string]string `json:"headers,omitempty"`
	BodyType *string            `json:"body_type" validate:"omitempty,oneof=JSON FORM NONE"`
	Body     *map[string]any    `json:"body,omitempty"`

	IntervalSeconds *int `json:"interval_seconds" validate:"omitempty,min=10"`
	TimeoutMs       *int `json:"timeout_ms" validate:"omitempty,min=1000"`

	ExpectedStatusCodes    *[]int  `json:"expected_status_codes,omitempty"`
	ExpectedResponseTimeMs *int    `json:"expected_response_time_ms,omitempty"`
	ExpectedBodyContains   *string `json:"expected_body_contains,omitempty"`
}
