package domain

import (
	"oko/pkg/db"
	"oko/pkg/ginapp/controller"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	List(c *gin.Context)
	Get(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

func NewController() controller.Ctrl {
	repository := NewDomainRepository(db.GetDB())
	handler := NewHandler(repository)
	return controller.Ctrl{
		Name:     "domain",
		Handlers: nil,
		Acts: []controller.Act{
			{Method: "GET", Route: "/", Handlers: []gin.HandlerFunc{handler.List}},
			{Method: "GET", Route: "/:id", Handlers: []gin.HandlerFunc{handler.Get}},
			{Method: "PUT", Route: "/:id", Handlers: []gin.HandlerFunc{handler.Update}},
			{Method: "DELETE", Route: "/:id", Handlers: []gin.HandlerFunc{handler.Delete}},
			{Method: "POST", Route: "/", Handlers: []gin.HandlerFunc{handler.Create}},
		},
	}
}
