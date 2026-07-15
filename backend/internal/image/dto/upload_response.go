package dto

type UploadResponse struct {
	ID        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ObjectKey string `json:"object_key" example:"images/550e8400-e29b-41d4-a716-446655440000.jpg"`
	Status    string `json:"status" example:"queued"`
}

type UploadSuccessResponse struct {
	Success bool           `json:"success" example:"true"`
	Message string         `json:"message" example:"Image uploaded and queued successfully"`
	Data    UploadResponse `json:"data"`
	Meta    any            `json:"meta,omitempty"`
}
