package dto

type UploadResponse struct {
	ObjectKey   string `json:"object_key" example:"images/550e8400-e29b-41d4-a716-446655440000.jpg"`
	Bucket      string `json:"bucket" example:"image-processing"`
	ContentType string `json:"content_type" example:"image/jpeg"`
	Size        int64  `json:"size" example:"12345"`
}

type UploadSuccessResponse struct {
	Success bool           `json:"success" example:"true"`
	Message string         `json:"message" example:"File uploaded successfully"`
	Data    UploadResponse `json:"data"`
	Meta    any            `json:"meta,omitempty"`
}
