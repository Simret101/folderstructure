package response

type ErrorResponseFormat struct {
	Status   string   `json:"status"`
	Code     string   `json:"code"`
	Message  string   `json:"message,omitempty"`
	Data     any      `json:"data"`
	Metadata Metadata `json:"metadata"`
}
type ErrorLoginResponseFormat struct {
	Status   string      `json:"status"`
	Code     string      `json:"code"`
	Message  string      `json:"message,omitempty"`
	Data     any         `json:"data"`
	Metadata MetaPindata `json:"metadata"`
}
type ErrorAuthzResponseFormat struct {
	Status   string   `json:"status"`
	Code     string   `json:"code"`
	Message  any      `json:"message,omitempty"`
	Data     any      `json:"data"`
	Metadata Metadata `json:"metadata"`
}

type ErrorResponse struct {
	Code        int    `json:"code,omitempty"`
	Message     string `json:"message,omitempty"`
	Description string `json:"description,omitempty"`

	StackTrace string       `json:"stack_trace,omitempty"`
	FieldError []FieldError `json:"field_error,omitempty"`
	Data       interface{}  `json:"field_error,omitempty"`
}

type FieldError struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Metadata struct {
	ServerTimestamp string `json:"server_timestamp"`
	RequestID       string `json:"request_id"`
}
type MetaPindata struct {
	ServerTimestamp string                 `json:"server_timestamp"`
	RequestID       string                 `json:"request_id"`
	Extra           map[string]interface{} `json:"extra,omitempty"`
}

type SuccessResponse struct {
	Status   string      `json:"status"`
	Code     string      `json:"code"`
	Message  any         `json:"message,omitempty"`
	Data     interface{} `json:"data"`
	Metadata Metadata    `json:"metadata"`
}

