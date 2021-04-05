package engine

type Response struct {
	Data interface{} `json:"data"`
	Meta MetaData    `json:"meta"`
}

type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Meta       MetaData    `json:"meta"`
	Pagination Pagination  `json:"pagination"`
}

func NewResponse(data interface{}, Type, message string) Response {
	return Response{
		Data: data,
		Meta: NewMetaData(Type, message),
	}
}

func NewResponsePaginated(data interface{}, pagination Pagination, message string) PaginationResponse {
	return PaginationResponse{
		Data:       data,
		Pagination: pagination,
		Meta:       NewMetaData("OK", message),
	}
}

func NewResponseOK(data interface{}, message string) Response {
	return NewResponse(data, "OK", message)
}

func NewResponseCreated(data interface{}, message string) Response {
	return NewResponse(data, "CREATED", message)
}

func NewResponseError(err *Error) Response {
	return NewResponse(err.Extra, err.Meta.Type, err.Meta.Message)
}

func NewResponseRemove(message string) Response {
	return NewResponse(message, "OK", message)
}
