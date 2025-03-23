package models

// ErrorResponse represents a standard error response
// @Description Standard error response object
type ErrorResponse struct {
	Error string `json:"error"`
}

// InternalServerErrorResponse represents a internal error response
// @Description Standard error response object
type InternalServerErrorResponse struct {
	Error string `json:"internal server error"`
}

// BadRequestResponse represents a internal error response
// @Description Standard error response object
type BadRequestResponse struct {
	Error string `json:"bad request"`
}

// UnauthorizedRequestResponse represents a internal error response
// @Description Standard error response object
type UnauthorizedRequestResponse struct {
	Error string `json:"unauthorized request"`
}

// EmptyResponse represents a internal error response
// @Description Standard error response object
type EmptyResponse struct {
	Error string `json:"ok"`
}

// NotFoundResponse represents a internal error response
// @Description Standard error response object
type NotFoundResponse struct {
	Error string `json:"not found"`
}
