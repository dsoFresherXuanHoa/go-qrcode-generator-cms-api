package entity

type standardResponse struct {
	StatusCode int         `json:"statusCode"`
	StatusText string      `json:"statusText"`
	Data       interface{} `json:"data"`
	Error      string      `json:"error"`
	Message    string      `json:"message"`
}

func NewStandardResponse(data interface{}, statusCode int, statusText string, err string, message string) standardResponse {
	return standardResponse{Data: data, StatusCode: statusCode, StatusText: statusText, Error: err, Message: message}
}
