package account

import (
	"oko/pkg/ginapp/controller"

	"github.com/gin-gonic/gin"
)

func NewController() controller.Ctrl {
	return controller.Ctrl{
		Name:     "account",
		Handlers: nil,
		Acts: []controller.Act{
			{Method: "POST", Route: "/sign-up/", Handlers: []gin.HandlerFunc{SignUp}},
			{Method: "POST", Route: "/sign-in/", Handlers: []gin.HandlerFunc{SignIn}},
			{Method: "POST", Route: "/sign-out/", Handlers: []gin.HandlerFunc{Auth(false, []int{}), SignOut}},
			{Method: "GET", Route: "/profile/", Handlers: []gin.HandlerFunc{Auth(true, []int{}), Profile}},
		},
	}
}
