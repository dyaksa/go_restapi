package errors

type ErrorDetail struct {
	Field          string `json:"field"`
	ValidationCode string `json:"validation_code"`
	Param          string `json:"param"`
	Message        string `json:"message"`
}
