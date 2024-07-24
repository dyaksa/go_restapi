package web

type ResponseCategory struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type Response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
	Meta   interface{} `json:"meta,omitempty"`
}

type Meta struct {
	Page int `json:"page"`
	Size int `json:"size"`
}
