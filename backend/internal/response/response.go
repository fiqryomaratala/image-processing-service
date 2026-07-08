package response

type SuccessBody struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Success"`
	Data    any    `json:"data,omitempty"`
	Meta    any    `json:"meta,omitempty"`
}

type ErrorBody struct {
	Success bool        `json:"success" example:"false"`
	Message string      `json:"message" example:"Validation failed"`
	Errors  []ErrorItem `json:"errors,omitempty"`
}

type ErrorItem struct {
	Field   string `json:"field,omitempty" example:"file"`
	Message string `json:"message" example:"file is required"`
	Code    string `json:"code,omitempty" example:"VALIDATION_ERROR"`
}

type HealthData struct {
	Status string `json:"status" example:"healthy"`
}

type HealthSuccessResponse struct {
	Success bool       `json:"success" example:"true"`
	Message string     `json:"message" example:"Image Processing Service API is running"`
	Data    HealthData `json:"data"`
	Meta    any        `json:"meta,omitempty"`
}
