package accounts

type Response struct {
	Status     string       `json:"status"`
	Message    string       `json:"message"`
	AccountDTO []AccountDTO `json:"data"`
}

// NewResponse is a utility function to create a new Response
func NewResponse(status, message string, accounts []AccountDTO) *Response {
	return &Response{
		Status:     status,
		Message:    message,
		AccountDTO: accounts,
	}
}

func NewSuccessResponse(accounts []AccountDTO) *Response {
	return NewResponse("200", "Success", accounts)
}

func NewErrorResponse(status, message string) *Response {
	return NewResponse(status, message, nil)
}

func NewCreatedResponse(accounts []AccountDTO) *Response {
	return NewResponse("201", "Types created successfully", accounts)
}

func NewNoContentResponse() *Response {
	return NewResponse("204", "No content", nil)
}
