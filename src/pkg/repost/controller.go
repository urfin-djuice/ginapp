package repost

import (
	"oko/pkg/account"
	"oko/pkg/db"
	"oko/pkg/ginapp/controller"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	View(c *gin.Context)
	List(c *gin.Context)
	NewRequest(c *gin.Context)
	Export(c *gin.Context)
}

func NewController() controller.Ctrl {
	repository := NewRequestRepository(db.GetDB())
	handler := NewHandler(repository)
	return controller.Ctrl{
		Name:     "repost",
		Handlers: controller.HandlerList{account.Auth(true, []int{})},
		Acts: []controller.Act{
			{Method: "POST", Route: "/", Handlers: []gin.HandlerFunc{handler.NewRequest}},
			{Method: "GET", Route: "/view", Handlers: []gin.HandlerFunc{handler.View}},
			{Method: "GET", Route: "/", Handlers: []gin.HandlerFunc{handler.List}},
			{Method: "GET", Route: "/export", Handlers: controller.HandlerList{handler.Export}},
		},
	}
}
