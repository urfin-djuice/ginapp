package links

import (
	"oko/pkg/db"
	"oko/pkg/ginapp/controller"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	List(c *gin.Context)
}

func NewController() controller.Ctrl {
	repository := NewLinkRepository(db.GetDB())
	handler := NewHandler(repository)

	return controller.Ctrl{
		Name:     "link",
		Handlers: nil,
		Acts: []controller.Act{
			{Method: "GET", Route: "/", Handlers: []gin.HandlerFunc{handler.List}},
		},
	}
}
