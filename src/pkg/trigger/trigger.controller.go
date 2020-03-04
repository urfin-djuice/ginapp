package trigger

import (
	"oko/pkg/account"
	"oko/pkg/ginapp/controller"

	"github.com/gin-gonic/gin"
)

type iHandler interface { //nolint
	Exist(c *gin.Context)
	List(c *gin.Context)
	Get(c *gin.Context)
	CreateByStatus(c *gin.Context)
	UpdateByStatus(c *gin.Context)
	CreateByBody(c *gin.Context)
	UpdateByBody(c *gin.Context)
	Delete(c *gin.Context)
}

func NewController() controller.Ctrl {
	triggerHandler := NewHandler()
	return controller.Ctrl{
		Name:     "trigger",
		Handlers: controller.HandlerList{account.Auth(true, []int{})},
		Acts: []controller.Act{
			{Method: "GET", Route: "/", Handlers: controller.HandlerList{triggerHandler.List}},
			{Method: "GET", Route: "/:id", Handlers: controller.HandlerList{triggerHandler.Get}},
			{Method: "DELETE", Route: "/:id", Handlers: controller.HandlerList{triggerHandler.Delete}},
			{Method: "HEAD", Route: "/:id", Handlers: controller.HandlerList{triggerHandler.Exist}},
			{Method: "POST", Route: "/", Handlers: controller.HandlerList{triggerHandler.Create}},
			{Method: "PUT", Route: "/:id", Handlers: controller.HandlerList{triggerHandler.Update}},
		},
	}
}
