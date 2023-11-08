package entity

type standardResponse struct {
	StatusCode int         `json:"statusCode"`
	StatusText string      `json:"statusText"`
	Data       interface{} `json:"data"`
	Error      string      `json:"error"`
	Message    string      `json:"message"`
}

type standardWithPaginateResponse struct {
	StatusCode int         `json:"statusCode"`
	StatusText string      `json:"statusText"`
	Data       interface{} `json:"data"`
	Error      string      `json:"error"`
	Message    string      `json:"message"`
	Paging     Paginate    `json:"paginate"`
}

func NewStandardResponse(data interface{}, statusCode int, statusText string, err string, message string) standardResponse {
	return standardResponse{Data: data, StatusCode: statusCode, StatusText: statusText, Error: err, Message: message}
}

func NewStandardWithPaginateResponse(data interface{}, statusCode int, statusText string, err string, message string, paging Paginate) standardWithPaginateResponse {
	return standardWithPaginateResponse{Data: data, StatusCode: statusCode, StatusText: statusText, Error: err, Message: message, Paging: paging}
}
