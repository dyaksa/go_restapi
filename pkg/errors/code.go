package errors

type ErrorCode string

const (
	ServerError     ErrorCode = "INTERNAL_SERVER_ERROR"
	CommitError     ErrorCode = "COMMIT_ERROR"
	NotFound        ErrorCode = "NOT_FOUND"
	BadRequest      ErrorCode = "BAD_REQUEST"
	InvalidJson     ErrorCode = "INVALID_JSON"
	BindHeap        ErrorCode = "BIND_HEAP"
	InsertData      ErrorCode = "INSERT_DATA"
	FetchData       ErrorCode = "FETCH_DATA"
	UpdateData      ErrorCode = "UPDATE_DATA"
	ValidationError ErrorCode = "VALIDATION_ERROR"
)
