package rule

import (
	"oko/pkg/account"
	"oko/pkg/ginapp/controller"

	"github.com/gin-gonic/gin"
)

type iHandler interface { //nolint
	Exist(c *gin.Context)
	List(c *gin.Context)
	Get(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

func NewController() controller.Ctrl {
	ruleHandler := NewHandler()
	return controller.Ctrl{
		Name:     "rule",
		Handlers: controller.HandlerList{account.Auth(true, []int{})},
		Acts: []controller.Act{
			{Method: "GET", Route: "/:id", Handlers: controller.HandlerList{ruleHandler.Get}},
			{Method: "GET", Route: "/", Handlers: controller.HandlerList{ruleHandler.List}},
			{Method: "DELETE", Route: "/:id", Handlers: controller.HandlerList{ruleHandler.Delete}},
			{Method: "POST", Route: "/", Handlers: controller.HandlerList{ruleHandler.Create}},
			{Method: "PUT", Route: "/:id", Handlers: controller.HandlerList{ruleHandler.Update}},
			{Method: "HEAD", Route: "/:id", Handlers: controller.HandlerList{ruleHandler.Exist}},
		},
	}
}
