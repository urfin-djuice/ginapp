package types

import (
	"net/http"
	"oko/pkg/e"
	"time"

	"github.com/gin-gonic/gin"
)

type KeyVal struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type StdResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ResponseError struct {
	StdResponse
	Err error `json:"err"`
}

type ResponseAccessToken struct {
	AccessToken string `json:"access_token"`
}

type ResponseSignIn struct {
	StdResponse
	Data ResponseAccessToken `json:"data"`
}

type Profile struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Status    uint      `json:"status"`
	StatusStr string    `json:"status_str"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResponseProfile struct {
	StdResponse
	Data Profile `json:"data"`
}

type RoleItem struct {
	Role    int    `json:"role"`
	StrRole string `json:"str_role"`
}

type RolesList struct {
	StdResponse
	Data []RoleItem `json:"data"`
}

type StringArray struct {
	StdResponse
	Data []string `json:"data"`
}

//!!! JUST FOR SWAGGER

type ResponseErrorSwg struct {
	StdResponse
	Err       string               `json:"err"`
	ErrStruct []e.CustomFieldError `json:"err_struct"`
}

type SearchSwg struct {
	ID        int       `json:"id"`
	UID       string    `json:"uid"`
	Query     string    `json:"query"`
	Answer    string    `json:"answer"`
	AccountID int       `json:"account_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

type GetSearchResponseSwg struct {
	StdResponse
	Data SearchSwg `json:"data"`
}

type Response struct {
	StdResponse
	Data interface{} `json:"data"`
	Meta interface{} `json:"meta"`
}

func SuccessEmptyResponse(c *gin.Context) {
	c.JSON(http.StatusOK, &Response{
		StdResponse: StdResponse{
			Status:  e.Success,
			Message: e.GetMsg(e.Success),
		},
	})
}

func (r Response) Success(c *gin.Context) {
	r.StdResponse = StdResponse{
		Status:  e.Success,
		Message: e.GetMsg(e.Success),
	}

	c.JSON(http.StatusOK, r)
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &Response{
		StdResponse: StdResponse{
			Status:  e.Success,
			Message: e.GetMsg(e.Success),
		},
		Data: data,
	})
}
