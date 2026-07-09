package dto

type UploadResponse struct {
	Filename    string `json:"filename" example:"sample.jpg"`
	ContentType string `json:"content_type" example:"image/jpeg"`
	Size        int64  `json:"size" example:"12345"`
}

type UploadSuccessResponse struct {
	Success bool           `json:"success" example:"true"`
	Message string         `json:"message" example:"File validation passed"`
	Data    UploadResponse `json:"data"`
	Meta    any            `json:"meta,omitempty"`
}
