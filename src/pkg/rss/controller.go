package rss

import (
	"oko/pkg/account"
	"oko/pkg/db"
	"oko/pkg/ginapp/controller"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	List(c *gin.Context)
	Get(c *gin.Context)
}

func NewController() controller.Ctrl {
	handler := rssHandler{NewRssRepository(db.GetDB())}
	return controller.Ctrl{
		Name:     "rss",
		Handlers: controller.HandlerList{account.Auth(false, []int{})},
		Acts: []controller.Act{
			{Method: "GET", Route: "/", Handlers: []gin.HandlerFunc{handler.List}},
			{Method: "GET", Route: "/:id", Handlers: []gin.HandlerFunc{handler.Get}},
			{Method: "PUT", Route: "/", Handlers: []gin.HandlerFunc{handler.Update}},
			{Method: "DELETE", Route: "/:id", Handlers: []gin.HandlerFunc{handler.Delete}},
			{Method: "POST", Route: "/", Handlers: []gin.HandlerFunc{handler.Create}},
		},
	}
}
