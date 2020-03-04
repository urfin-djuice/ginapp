package links

import (
	"context"
	"net/http"
	"oko/pkg/e"
	"oko/pkg/env"
	"oko/pkg/ginapp/types"
	contentPB "oko/srv/content/proto"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
	"github.com/thoas/go-funk"
)

type linkHandler struct {
	repository Repository
}

func NewHandler(repo Repository) Handler {
	return &linkHandler{
		repository: repo,
	}
}

// List godoc
// @Summary List
// @Description List
// @ID get-links-list
// @Tags Link
// @Accept json
// @Produce json
// @Param object query links.ListRequest true "Links find request"
// @Success 200 {object} types.StdResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /link [get]
func (h *linkHandler) List(c *gin.Context) {
	var form ListRequest

	if err := c.ShouldBind(&form); err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	reg := etcd.NewRegistry(
		registry.Addrs(env.GetEnvOrPanic("ETCD_ADDRESS")),
	)
	service := micro.NewService(
		micro.Registry(reg),
	)
	service.Init()

	cl := service.Client()
	_ = cl.Init(
		client.RequestTimeout(time.Second * 30))

	ts := contentPB.NewContentService("go.micro.srv.content", cl)

	result, err := ts.List(context.Background(), &contentPB.ContentListRequest{
		Query:    form.Query,
		Page:     form.CurrentPage,
		Limit:    form.PerPage,
		DomainId: form.DomainID,
	})

	if err != nil || result == nil {
		e.ErrorResponse(c, http.StatusBadRequest, "Something went wrong")
		return
	}

	data := funk.Map(result.Data, func(model *contentPB.Content) *Response {
		id, _ := strconv.Atoi(model.Id[:len(model.Id)-8])
		if err != nil {
			e.ErrorResponse(c, http.StatusBadRequest, "Something went wrong")
		}

		link, _ := h.repository.Get(uint(id))
		if err != nil {
			e.ErrorResponse(c, http.StatusBadRequest, "Something went wrong")
		}

		serializer := Serializer{Link: link}

		return serializer.To(model.Data)
	}).([]*Response)

	c.JSON(
		http.StatusOK,
		gin.H{
			"data": data,
			"meta": MetaListResponse{
				PaginationResponse: types.PaginationResponse{
					PaginationRequest: types.PaginationRequest{
						CurrentPage: form.CurrentPage,
						PerPage:     form.PerPage,
					},
					TotalRecords: result.Meta.TotalHits,
					TotalPages:   result.Meta.TotalHits/form.PerPage + 1,
				},
				NegativeCount: result.Meta.NegativeCount,
				PositiveCount: result.Meta.PositiveCount,
				NeutralCount:  result.Meta.NeutralCount,
				Image:         nil,
			},
		})
}
