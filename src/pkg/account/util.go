package account

import (
	"net/http"
	"oko/pkg/e"

	"github.com/gin-gonic/gin"
)

func loginErrorResponse(c *gin.Context, email string) {
	e.ErrorResponse(c, http.StatusBadRequest, e.CustomFieldError{
		Name:    "email",
		Tag:     "email",
		Param:   "",
		Value:   email,
		Message: "Invalid username or password",
	})
}
