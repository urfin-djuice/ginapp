package e

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/iancoleman/strcase"
	"gopkg.in/go-playground/validator.v9"
)

const DefaultValidatorMessage = "Incorrect field"

type CustomFieldError struct {
	Name    string      `json:"name"`
	Tag     string      `json:"tag"`
	Param   string      `json:"param"`
	Value   interface{} `json:"value"`
	Message string      `json:"message"`
}

var MsgFlags = map[int]string{
	Success:                    "ok",
	Created:                    "created successfully",
	Error:                      "fail",
	InvalidParams:              "invalid params",
	Unauthorized:               "unauthorized",
	NotFound:                   "not found",
	NotAcceptable:              "not acceptable",
	ErrorAuthCheckTokenFail:    "token check failed",
	ErrorAuthCheckTokenTimeout: "token expired",
	ErrorAuthToken:             "token authentication failed",
	ErrorAuth:                  "authentication failed",
	ErrorAuthUnconfirmed:       "unconfirmed account",
	ErrorAuthBanned:            "banned account",
	ErrorAuthRole:              "no role needed",
	ErrorAuthAppKeyNotAllowed:  "use of the method with the application key is not allowed",
	Processing:                 "link is being processed",
}

var ValidatorMessages = map[string]string{}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[Error]
}

func ErrorResponse(c *gin.Context, code int, err interface{}) {
	ginH := gin.H{
		"code":    code,
		"message": GetMsg(code),
	}

	errStruct := make([]CustomFieldError, 0)

	const validationErrorMsg = "Some fields are incorrect"
	switch v := err.(type) {
	case validator.ValidationErrors:
		for _, fe := range v {
			message := ValidatorMessages[fe.Tag()]
			if message == "" {
				message = DefaultValidatorMessage
			}
			errStruct = append(errStruct, CustomFieldError{
				Name:    strcase.ToSnake(fe.Field()),
				Tag:     strcase.ToSnake(fe.Tag()),
				Param:   strcase.ToSnake(fe.Param()),
				Value:   fmt.Sprintf("%v", fe.Value()),
				Message: message,
			})
		}
		ginH["err"] = validationErrorMsg
	case CustomFieldError:
		errStruct = append(errStruct, v)
		ginH["err"] = validationErrorMsg
	case error:
		ginH["err"] = v.Error()
	case string:
		ginH["err"] = v
	default:
		ginH["err"] = fmt.Sprintf("%T", v)
	}

	ginH["err_struct"] = errStruct

	c.JSON(CodeStatuses[code], ginH)
}
