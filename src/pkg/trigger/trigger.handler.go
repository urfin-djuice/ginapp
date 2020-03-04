package trigger

import (
	"context"
	"net/http"
	"oko/pkg/e"
	"oko/pkg/env"
	"oko/pkg/ginapp"
	"oko/pkg/ginapp/types"
	"oko/pkg/log"
	pb "oko/srv/proxy/proto"
	"time"

	"github.com/thoas/go-funk"

	"github.com/golang/protobuf/ptypes/wrappers"

	"github.com/gin-gonic/gin"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
)

type handler struct {
	srv pb.TriggerService
}

func NewHandler() *handler { //nolint
	reg := etcd.NewRegistry(
		registry.Addrs(env.GetEnvOrPanic("ETCD_ADDRESS")),
	)
	service := micro.NewService(micro.Registry(reg))

	service.Init()

	cl := service.Client()

	_ = cl.Init(client.RequestTimeout(time.Second * 30))
	return &handler{srv: pb.NewTriggerService("go.micro.srv.proxy", cl)}
}

func ctx() context.Context {
	return context.Background()
}

// Exist godoc
// @Summary Exist trigger item
// @Description Exist trigger item
// @ID exist-trigger
// @Tags Trigger
// @Accept json
// @Produce json
// @Param id path int true "Trigger id"
// @Success 200
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /trigger/{id} [head]
// @Security ApiKeyAuth
func (h handler) Exist(c *gin.Context) {
	id, err := ginapp.GetUint32PathParam(c, "id")

	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	exist, err := h.srv.Exist(ctx(), &pb.TriggerRequest{Id: id})
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

// List godoc
// @Summary List trigger item
// @Description List trigger item
// @ID list-trigger
// @Tags Trigger
// @Accept json
// @Produce json
// @Param rule_id query int false "Rule id for filter triggers"
// @Param object query types.PaginationRequest true "Pagination"
// @Success 200 {object} trigger.ListResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /trigger/ [get]
// @Security ApiKeyAuth
func (h handler) List(c *gin.Context) {
	filter := &ListFilter{}
	err := c.BindQuery(filter)

	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	var form types.PaginationRequest
	if err = c.ShouldBind(&form); err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	req := &pb.TriggerListRequest{
		CurrentPage: form.CurrentPage,
		PerPage:     form.PerPage,
	}

	if filter.RuleID != nil {
		req.RuleId = &wrappers.UInt32Value{Value: *filter.RuleID}
	}
	list, err := h.srv.List(ctx(), req)
	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
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
// @Summary Get trigger item
// @Description Get trigger item
// @ID get-trigger
// @Tags Trigger
// @Accept json
// @Produce json
// @Param id path int true "Trigger id"
// @Success 200 {object} trigger.GetResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /trigger/{id} [get]
// @Security ApiKeyAuth
func (h handler) Get(c *gin.Context) {
	id, err := ginapp.GetUint32PathParam(c, "id")

	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	get, err := h.srv.Get(ctx(), &pb.TriggerRequest{Id: id})
	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}
	types.SuccessResponse(c, toView(get.Data))
}

// Update godoc
// @Summary Update trigger item
// @Description Update trigger item
// @ID update-trigger
// @Tags Trigger
// @Accept json
// @Produce json
// @Param id path int true "Trigger id"
// @Param object body trigger.UpdateRequest true "Update trigger item body"
// @Success 200 {object} trigger.GetResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /trigger/{id} [put]
// @Security ApiKeyAuth
func (h handler) Update(c *gin.Context) {
	id, err := ginapp.GetUint32PathParam(c, "id")
	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, e.CustomFieldError{
			Name:    "id",
			Tag:     "id",
			Param:   "",
			Value:   c.Param("id"),
			Message: "Invalid id",
		})
		return
	}
	trigger := &UpdateRequest{}
	err = c.BindJSON(&trigger)
	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	request := pb.TriggerUpdateRequest{
		Id: &wrappers.UInt32Value{Value: id},
	}

	if trigger.URL != nil {
		request.Url = &wrappers.StringValue{Value: *trigger.URL}
	}

	if trigger.Params != nil {
		request.Params = &wrappers.StringValue{Value: *trigger.Params}
	}
	if trigger.Type != nil {
		request.Type = ToTriggerType(*trigger.Type)
	}

	if trigger.RuleID != nil {
		request.RuleId = &wrappers.UInt32Value{Value: *trigger.RuleID}
	}

	res, err := h.srv.Update(context.Background(), &request)

	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}
	types.SuccessResponse(c, toView(res.Data))
}

// Create godoc
// @Summary Create trigger item
// @Description Create trigger item
// @ID create-trigger
// @Tags Trigger
// @Accept json
// @Produce json
// @Param object body trigger.CreateRequest true "Create trigger item body"
// @Success 200 {object} trigger.GetResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /trigger/ [post]
// @Security ApiKeyAuth
func (h handler) Create(c *gin.Context) {
	trigger := CreateRequest{}
	err := c.BindJSON(&trigger)
	if err != nil {
		log.Println(err)
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.srv.Create(context.Background(), &pb.TriggerCreateRequest{
		Url:    trigger.URL,
		RuleId: trigger.RuleID,
		Type:   ToTriggerType(trigger.Type),
		Params: trigger.Params,
	})
	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	types.SuccessResponse(c, toView(res.Data))
}

// Delete godoc
// @Summary Delete trigger item
// @Description Delete trigger item
// @ID delete-trigger
// @Tags Trigger
// @Accept json
// @Produce json
// @Param id path int true "Trigger id"
// @Success 200 {object} types.StdResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /trigger/{id} [delete]
// @Security ApiKeyAuth
func (h handler) Delete(c *gin.Context) {
	id, err := ginapp.GetUint32PathParam(c, "id")

	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	_, err = h.srv.Delete(ctx(), &pb.TriggerRequest{Id: id})
	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, err)
		return
	}
	types.SuccessEmptyResponse(c)
}
