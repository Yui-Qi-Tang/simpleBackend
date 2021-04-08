package recovery

import (
	"net/http"
)

const (
	NotFound             = "not_found"
	InternalError        = "server_internal_error"
	TooBusy              = "server_too_busy"
	DatabaseError        = "db_error"
	ThirdPartyServiceErr = "3-rd-party_service_error"
	InvalidRequest       = "invalid_request"
)

// APIError a form for api error
type APIError struct {
	Status int
	Msg    string
	Code   string
}

func (apiErr *APIError) Error() string {
	return apiErr.Msg
}

func (apiErr *APIError) response() response {
	return response{
		"code":  apiErr.Code,
		"error": apiErr.Msg,
	}
}

// Set500
func (apiErr *APIError) Set500(msg string) {
	apiErr.Status = http.StatusInternalServerError
	apiErr.Msg = msg
	apiErr.Code = InternalError
}

// Set400
func (apiErr *APIError) Set400(msg string) {
	apiErr.Status = http.StatusBadRequest
	apiErr.Msg = msg
	apiErr.Code = InvalidRequest
}
