package e

import "net/http"

const (
	Success       = 200
	Created       = 201
	Error         = 500
	InvalidParams = 400
	Unauthorized  = 401
	NotFound      = 404
	NotAcceptable = 406

	ErrorAuthCheckTokenFail    = 20001
	ErrorAuthCheckTokenTimeout = 20002
	ErrorAuthToken             = 20003
	ErrorAuth                  = 20004
	ErrorAuthUnconfirmed       = 20005
	ErrorAuthBanned            = 20006
	ErrorAuthRole              = 20007
	ErrorAuthAppKeyNotAllowed  = 20008
	Processing                 = 20009
)

var CodeStatuses = map[int]int{
	Success:                    http.StatusOK,
	Created:                    http.StatusCreated,
	Error:                      http.StatusInternalServerError,
	InvalidParams:              http.StatusBadRequest,
	Unauthorized:               http.StatusUnauthorized,
	NotFound:                   http.StatusNotFound,
	NotAcceptable:              http.StatusNotAcceptable,
	ErrorAuthCheckTokenFail:    http.StatusUnauthorized,
	ErrorAuthCheckTokenTimeout: http.StatusUnauthorized,
	ErrorAuthToken:             http.StatusUnauthorized,
	ErrorAuth:                  http.StatusUnauthorized,
	ErrorAuthUnconfirmed:       http.StatusUnauthorized,
	ErrorAuthBanned:            http.StatusUnauthorized,
	ErrorAuthRole:              http.StatusUnauthorized,
	ErrorAuthAppKeyNotAllowed:  http.StatusUnauthorized,
}
