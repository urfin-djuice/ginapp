package action

import (
	"context"
	"net/http"
	"oko/pkg/e"
	"oko/pkg/env"
	"oko/pkg/ginapp"
	"oko/pkg/ginapp/types"
	pb "oko/srv/proxy/proto"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"

	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/etcd"
)

type handler struct {
	srv pb.ActionService
}

func NewHandler() *handler { //nolint
	reg := etcd.NewRegistry(
		registry.Addrs(env.GetEnvOrPanic("ETCD_ADDRESS")),
	)
	service := micro.NewService(micro.Registry(reg))

	service.Init()

	cl := service.Client()

	_ = cl.Init(client.RequestTimeout(time.Second * 30))

	actionSrv := pb.NewActionService("go.micro.srv.proxy", cl)

	return &handler{
		srv: actionSrv,
	}
}

// Exist godoc
// @Summary Exist action item
// @Description Exist action item
// @ID exist-action
// @Tags Action
// @Accept json
// @Produce json
// @Param id path int true "Exist action item id"
// @Success 200
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /action/{id} [head]
// @Security ApiKeyAuth
func (h handler) Exist(c *gin.Context) {
	actionID, err := ginapp.GetUint32PathParam(c, "id")
	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	exist, err := h.srv.Exist(context.Background(), &pb.ActionRequest{
		Id: actionID,
	})
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
// @Summary List action item
// @Description List action item
// @ID list-action
// @Tags Action
// @Accept json
// @Produce json
// @Param rule_id query int false "Rule id for filtering list"
// @Param object query types.PaginationRequest true "Pagination"
// @Success 200 {object} action.ListResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /action/ [get]
// @Security ApiKeyAuth
func (h handler) List(c *gin.Context) {
	filter := &ListFilter{}
	err := c.ShouldBindQuery(filter)

	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	var form types.PaginationRequest
	if err = c.ShouldBind(&form); err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	req := &pb.ActionListRequest{
		CurrentPage: form.CurrentPage,
		PerPage:     form.PerPage,
	}

	if filter.RuleID != nil {
		req.RuleId = &wrappers.UInt32Value{Value: *filter.RuleID}
	}
	list, err := h.srv.List(context.Background(), req)
	if err != nil || list == nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}
	views := []View{}
	if list.Data != nil && len(list.Data) > 0 {
		views = make([]View, 0, len(list.Data))
		for _, item := range list.Data {
			if item != nil {
				views = append(views, toView(*item))
			}
		}
	}

	meta := &types.PaginationResponse{
		PaginationRequest: types.PaginationRequest{
			CurrentPage: list.Meta.CurrentPage,
			PerPage:     list.Meta.PerPage,
		},
		TotalPages:   list.Meta.TotalPage,
		TotalRecords: list.Meta.TotalRecord,
	}

	result := types.Response{
		Data: views,
		Meta: meta,
	}

	result.Success(c)
}

// Get godoc
// @Summary Get action item
// @Description Get action item
// @ID get-action
// @Tags Action
// @Accept json
// @Produce json
// @Param id path int true "Get action item id"
// @Success 200 {object} action.GetResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /action/{id} [get]
// @Security ApiKeyAuth
func (h handler) Get(c *gin.Context) {
	actID, err := ginapp.GetUint32PathParam(c, "id")
	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	act, err := h.srv.Get(context.Background(), &pb.ActionRequest{Id: actID})
	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}
	if act == nil || act.Data == nil {
		e.ErrorResponse(c, http.StatusNotFound, "Something went wrong")
	}

	view := toView(*act.Data)
	types.SuccessResponse(c, view)
}

// Create godoc
// @Summary Create action item
// @Description Create action item
// @ID create-action
// @Tags Action
// @Accept json
// @Produce json
// @Param object body action.CreateRequest true "Create action item body"
// @Success 200 {object} action.GetResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /action/ [post]
// @Security ApiKeyAuth
func (h handler) Create(c *gin.Context) {
	act := CreateRequest{}
	err := c.BindJSON(&act)
	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	res, err := h.srv.Create(context.Background(), &pb.ActionCreateRequest{
		RuleId: act.RuleID,
		Type:   pb.ActionType(act.Type),
		Params: act.Params,
	})
	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}
	types.SuccessResponse(c, toView(*res.Data))
}

// Update godoc
// @Summary Update action item
// @Description Update action item
// @ID update-action
// @Tags Action
// @Accept json
// @Produce json
// @Param id path int true "Action ID"
// @Param object body action.UpdateRequestBody true "Update action item body"
// @Success 200 {object} action.GetResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /action/{id} [put]
// @Security ApiKeyAuth
func (h handler) Update(c *gin.Context) {
	act := &UpdateRequest{}
	err := c.BindJSON(act)
	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	act.ID, err = ginapp.GetUint32PathParam(c, "id")
	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, e.CustomFieldError{
			Name:    "id",
			Tag:     "id",
			Param:   "",
			Value:   0,
			Message: "Action ID is required",
		})
		return
	}

	request := pb.ActionUpdateRequest{
		Id: &wrappers.UInt32Value{Value: act.ID},
	}

	if act.Params != nil {
		request.Params = &wrappers.StringValue{Value: *act.Params}
	}
	if act.Type != nil {
		request.Type = pb.ActionType(*act.Type)
	}

	res, err := h.srv.Update(context.Background(), &request)

	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}
	types.SuccessResponse(c, res)
}

// Delete godoc
// @Summary Delete action item
// @Description Delete action item
// @ID delete-action
// @Tags Action
// @Accept json
// @Produce json
// @Param id path int true "Update action item id"
// @Success 200 {object} types.Response
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /action/{id} [delete]
// @Security ApiKeyAuth
func (h handler) Delete(c *gin.Context) {
	actID, err := ginapp.GetUint32PathParam(c, "id")
	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	_, err = h.srv.Delete(context.Background(), &pb.ActionRequest{
		Id: actID,
	})

	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}
	types.SuccessEmptyResponse(c)
}

// AddRule godoc
// @Summary Add rule to action
// @Description Add rule to action
// @ID add-rule-action
// @Tags Action
// @Accept json
// @Produce json
// @Param id path int true "Action ID"
// @Param object body action.AddRuleRequestBody true "Add rule body"
// @Success 200
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /action/{id}/rule [post]
// @Security ApiKeyAuth
func (h handler) AddRuleToAction(c *gin.Context) {
	req := &AddRuleRequest{}
	err := c.Bind(req)
	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	req.ID, err = ginapp.GetUint32PathParam(c, "id")
	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, e.CustomFieldError{
			Name:    "id",
			Tag:     "id",
			Param:   "",
			Value:   0,
			Message: "Action ID is required",
		})
		return
	}

	resp, err := h.srv.AddRule(context.Background(), &pb.AddRuleRequest{
		ActionId: req.ID,
		RuleId:   req.RuleID,
	})
	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	types.SuccessResponse(c, toView(*resp.Data))
}

// DelRuleAction godoc
// @Summary Delete rule from action
// @Description Delete rule from action
// @ID del-rule-action
// @Tags Action
// @Accept json
// @Produce json
// @Param id path int true "Action ID"
// @Param object body action.AddRuleRequestBody true "Del rule body"
// @Success 200
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /action/{id}/rule [delete]
// @Security ApiKeyAuth
func (h handler) DelRuleFormAction(c *gin.Context) {
	req := &AddRuleRequest{}
	err := c.Bind(req)
	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, "Something went wrong")
		return
	}
	req.ID, err = ginapp.GetUint32PathParam(c, "id")
	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, e.CustomFieldError{
			Name:    "id",
			Tag:     "id",
			Param:   "",
			Value:   0,
			Message: "Action ID is required",
		})
		return
	}

	resp, err := h.srv.DelRule(context.Background(), &pb.AddRuleRequest{
		ActionId: req.ID,
		RuleId:   req.RuleID,
	})
	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	types.SuccessResponse(c, toView(*resp.Data))
}
