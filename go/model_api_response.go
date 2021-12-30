package openapi

// ApiResponse - A generic API response
type ApiResponse struct {
	Code int32 `bson:"code,omitempty"`

	Type string `bson:"type,omitempty"`

	Message string `bson:"message,omitempty"`
}
