package transactiontypes

type Response struct {
	Status          string            `json:"status"`
	Message         string            `json:"message"`
	TransactionType []TransactionType `json:"data"`
}

// NewResponse is a utility function to create a new Response
func NewResponse(status, message string, transactionType []TransactionType) *Response {
	return &Response{
		Status:          status,
		Message:         message,
		TransactionType: transactionType,
	}
}

// NewSuccessResponse creates a success response with a list of types
func NewSuccessResponse(transactionType []TransactionType) *Response {
	return NewResponse("200", "Success", transactionType)
}

// NewErrorResponse creates an error response with a custom status and message
func NewErrorResponse(status, message string) *Response {
	return NewResponse(status, message, nil)
}

// NewCreatedResponse creates a response for a newly created resource
func NewCreatedResponse(transactionType []TransactionType) *Response {
	return NewResponse("201", "transactionType created successfully", transactionType)
}

// NewNoContentResponse creates a response for no content (for delete or update)
func NewNoContentResponse() *Response {
	return NewResponse("204", "No content", nil)
}
