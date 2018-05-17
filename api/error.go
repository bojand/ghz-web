package api

// apiError is an API error message
type apiError struct {
	Message string `json:"message`
}

// newAPIError returns new API error from an err
func newAPIError(err error) *apiError {
	return &apiError{Message: err.Error()}
}
