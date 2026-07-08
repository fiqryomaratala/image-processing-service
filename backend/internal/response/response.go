package response

type SuccessBody struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

type ErrorBody struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Errors  []ErrorItem `json:"errors,omitempty"`
}

type ErrorItem struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}
