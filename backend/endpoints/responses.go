package endpoints

// general response with error, success and data
type Response struct {
	Error   *ErrorResponse `json:"error"`
	Success bool           `json:"success"`
	Data    interface{}    `json:"data"`
}

type ErrorResponse struct {
	Error string                 `json:"code"`
	Msg   string                 `json:"msg"`
	Data  map[string]interface{} `json:"data"`
}

func newSuccessResponse(data interface{}) *Response {
	return &Response{
		Error:   nil,
		Success: true,
		Data:    data,
	}
}

func newErrorResponse(err string, msg string, data map[string]interface{}) *Response {
	return &Response{
		Error: &ErrorResponse{
			Error: err,
			Msg:   msg,
			Data:  data,
		},
		Success: false,
		Data:    nil,
	}
}
