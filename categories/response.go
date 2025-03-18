package categories

type Response struct {
	Status     string     `json:"status"`
	Message    string     `json:"message"`
	Categories []Category `json:"data"`
}

// NewResponse is a utility function to create a new Response
func NewResponse(status, message string, categories []Category) *Response {
	return &Response{
		Status:     status,
		Message:    message,
		Categories: categories,
	}
}

// NewSuccessResponse creates a success response with a list of categories
func NewSuccessResponse(categories []Category) *Response {
	return NewResponse("200", "Success", categories)
}

// NewErrorResponse creates an error response with a custom status and message
func NewErrorResponse(status, message string) *Response {
	return NewResponse(status, message, nil)
}

// NewCreatedResponse creates a response for a newly created resource
func NewCreatedResponse(categories []Category) *Response {
	return NewResponse("201", "Category created successfully", categories)
}

// NewNoContentResponse creates a response for no content (for delete or update)
func NewNoContentResponse() *Response {
	return NewResponse("204", "No content", nil)
}
