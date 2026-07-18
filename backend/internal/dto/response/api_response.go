package response

type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Result  interface{} `json:"result,omitempty"`
}

type PageResponse struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"size"`
}

func Success(result interface{}) ApiResponse {
	return ApiResponse{Code: 200, Message: "ok", Result: result}
}

func Error(code int, message string) ApiResponse {
	return ApiResponse{Code: code, Message: message}
}
