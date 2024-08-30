package web

type ApiResponse struct {
	Code    int `json:"code"`
	Data    any `json:"data,omitempty"`
	Details any `json:"details,omitempty"`
}
