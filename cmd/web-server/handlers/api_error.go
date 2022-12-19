package handlers

type ApiError = struct {
	ErrorCode    int    `json:"code"`
	ErrorMessage string `json:"message"`
}
