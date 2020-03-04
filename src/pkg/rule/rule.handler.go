package rule

import (
	"context"
	"net/http"
	"oko/pkg/e"
	"oko/pkg/env"
	"oko/pkg/ginapp"
	"oko/pkg/ginapp/types"
	pb "oko/srv/proxy/proto"
	"time"

	"github.com/thoas/go-funk"

	"github.com/gin-gonic/gin"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
)

type handler struct {
	srv pb.RuleService
}

func NewHandler() *handler { //nolint
	reg := etcd.NewRegistry(
		registry.Addrs(env.GetEnvOrPanic("ETCD_ADDRESS")),
	)
	service := micro.NewService(micro.Registry(reg))

	service.Init()

	sdCl := service.Client()

	_ = sdCl.Init(client.RequestTimeout(time.Second * 30))

	ruleService := pb.NewRuleService("go.micro.srv.proxy", sdCl)

	return &handler{srv: ruleService}
}

// Exist godoc
// @Summary Exist rule item
// @Description Exist rule item
// @ID exist-rule
// @Tags Rule
// @Accept json
// @Produce json
// @Param id query int true "Rule item id"
// @Success 200
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /rule/{id} [head]
// @Security ApiKeyAuth
func (h handler) Exist(c *gin.Context) {
	id, err := ginapp.GetUint32PathParam(c, "id")

	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	exist, err := h.srv.Exist(context.Background(), &pb.RuleRequest{Id: id})
	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if exist.Data {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusNotFound)
	}
}

// Create godoc
// @Summary Create rule item
// @Description Create rule item
// @ID create-rule
// @Tags Rule
// @Accept json
// @Produce json
// @Param object body rule.CreateRequest true "Create rule item body"
// @Success 200 {object} rule.GetResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /rule/ [post]
// @Security ApiKeyAuth
func (h handler) Create(c *gin.Context) {
	req := &CreateRequest{}
	err := c.BindJSON(req)
	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	res, err := h.srv.Create(context.Background(), &pb.RuleCreateRequest{
		Host:   req.Host,
		Status: req.Status,
	})
	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	types.SuccessResponse(c, toView(res.Data))
}

// List godoc
// @Summary List rules
// @Description List rules
// @ID list-rule
// @Tags Rule
// @Accept json
// @Produce json
// @Param host query string false "Rule item id"
// @Param object query types.PaginationRequest true "Pagination"
// @Success 200 {object} rule.ListResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /rule/ [get]
// @Security ApiKeyAuth
func (h handler) List(c *gin.Context) {
	host := c.Query("host")

	var form types.PaginationRequest
	if err := c.ShouldBind(&form); err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	list, err := h.srv.List(context.Background(), &pb.RuleListRequest{
		CurrentPage: form.CurrentPage,
		PerPage:     form.PerPage,
		Host:        host,
	})
	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	meta := &types.PaginationResponse{
		PaginationRequest: types.PaginationRequest{
			CurrentPage: list.Meta.CurrentPage,
			PerPage:     list.Meta.PerPage,
		},
		TotalPages:   list.Meta.TotalPage,
		TotalRecords: list.Meta.TotalRecord,
	}

	data := funk.Map(list.Data, toView)
	result := types.Response{
		Data: data,
		Meta: meta,
	}

	result.Success(c)
}

// Get godoc
// @Summary Get rule item
// @Description Get rule item
// @ID get-rule
// @Tags Rule
// @Accept json
// @Produce json
// @Param id path string true "Rule item id"
// @Success 200 {object} rule.GetResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /rule/{id} [get]
// @Security ApiKeyAuth
func (h handler) Get(c *gin.Context) {
	id, err := ginapp.GetUint32PathParam(c, "id")

	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.srv.Get(context.Background(), &pb.RuleRequest{Id: id})
	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	if res == nil || res.Data == nil {
		e.ErrorResponse(c, http.StatusNotFound, "Something went wrong")
		return
	}
	types.SuccessResponse(c, toView(res.Data))
}

// Update godoc
// @Summary Update rule item
// @Description Update rule item
// @ID update-rule
// @Tags Rule
// @Accept json
// @Produce json
// @Param id path int true "Update rule id"
// @Param object body rule.UpdateRequest true "Update rule item body"
// @Success 200 {object} rule.GetResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /rule/{id} [put]
// @Security ApiKeyAuth
func (h handler) Update(c *gin.Context) {
	id, err := ginapp.GetUint32PathParam(c, "id")
	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	req := &CreateRequest{}
	if err = c.BindJSON(req); err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.srv.Update(context.Background(), &pb.RuleCreateRequest{
		Id:     id,
		Host:   req.Host,
		Status: req.Status,
	})
	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	types.SuccessResponse(c, toView(res.Data))
}

// Delete godoc
// @Summary Delete rule item
// @Description Delete rule item
// @ID Delete-rule
// @Tags Rule
// @Accept json
// @Produce json
// @Param id query string true "Delete rule item id"
// @Success 200 {object} types.Response
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /rule/{id} [delete]
// @Security ApiKeyAuth
func (h handler) Delete(c *gin.Context) {
	id, err := ginapp.GetUint32PathParam(c, "id")

	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if _, err := h.srv.Delete(context.Background(), &pb.RuleRequest{Id: id}); err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	types.SuccessEmptyResponse(c)
}
