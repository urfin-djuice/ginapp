package action

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
	h := NewHandler()
	return controller.Ctrl{
		Name:     "action",
		Handlers: controller.HandlerList{account.Auth(true, []int{})},
		Acts: []controller.Act{
			{Method: "GET", Route: "/:id", Handlers: controller.HandlerList{h.Get}},
			{Method: "GET", Route: "/", Handlers: controller.HandlerList{h.List}},
			{Method: "DELETE", Route: "/:id", Handlers: controller.HandlerList{h.Delete}},
			{Method: "POST", Route: "/", Handlers: controller.HandlerList{h.Create}},
			{Method: "PUT", Route: "/:id", Handlers: controller.HandlerList{h.Update}},
			{Method: "HEAD", Route: "/:id", Handlers: controller.HandlerList{h.Exist}},
			{Method: "POST", Route: "/:id/rule", Handlers: controller.HandlerList{h.AddRuleToAction}},
			{Method: "DELETE", Route: "/:id/rule", Handlers: controller.HandlerList{h.DelRuleFormAction}},
		},
	}
}
