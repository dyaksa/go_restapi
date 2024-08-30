package errors

type AppError struct {
	Code    int
	Err     error
	Details ErrorDetail
}

func NewAppError(code int, details ErrorDetail) *AppError {
	appError := &AppError{}
	appError.Code = code
	appError.Details = details
	return appError
}

func (appErr AppError) WithError(err error) *AppError {
	appErr.Err = err
	return &appErr
}
