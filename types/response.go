package types

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Types   []Type `json:"data"`
}

// NewResponse is a utility function to create a new Response
func NewResponse(status, message string, types []Type) *Response {
	return &Response{
		Status:  status,
		Message: message,
		Types:   types,
	}
}

// NewSuccessResponse creates a success response with a list of types
func NewSuccessResponse(types []Type) *Response {
	return NewResponse("200", "Success", types)
}

// NewErrorResponse creates an error response with a custom status and message
func NewErrorResponse(status, message string) *Response {
	return NewResponse(status, message, nil)
}

// NewCreatedResponse creates a response for a newly created resource
func NewCreatedResponse(types []Type) *Response {
	return NewResponse("201", "Types created successfully", types)
}

// NewNoContentResponse creates a response for no content (for delete or update)
func NewNoContentResponse() *Response {
	return NewResponse("204", "No content", nil)
}
