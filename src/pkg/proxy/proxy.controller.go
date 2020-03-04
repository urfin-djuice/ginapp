package proxy

import (
	"oko/pkg/account"
	"oko/pkg/ginapp/controller"

	"github.com/gin-gonic/gin"
)

type iHandler interface { //nolint
	Delete(ctx *gin.Context)
	Create(ctx *gin.Context)
	Get(ctx *gin.Context)
	Update(ctx *gin.Context)
	Exist(ctx *gin.Context)
	List(ctx *gin.Context)
}

func NewController() controller.Ctrl {
	proxyHandler := NewHandler()
	return controller.Ctrl{
		Name:     "proxy",
		Handlers: controller.HandlerList{account.Auth(false, []int{})},
		Acts: []controller.Act{
			{Method: "GET", Route: "/:id", Handlers: controller.HandlerList{proxyHandler.Get}},
			{Method: "GET", Route: "/", Handlers: controller.HandlerList{proxyHandler.List}},
			{Method: "DELETE", Route: "/:id", Handlers: controller.HandlerList{proxyHandler.Delete}},
			{Method: "POST", Route: "/", Handlers: controller.HandlerList{proxyHandler.Create}},
			{Method: "PUT", Route: "/:id", Handlers: controller.HandlerList{proxyHandler.Update}},
			{Method: "HEAD", Route: "/:id", Handlers: controller.HandlerList{proxyHandler.Exist}},
		},
	}
}
