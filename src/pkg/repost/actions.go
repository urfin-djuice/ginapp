package repost

import (
	"math"
	"net/http"
	"oko/pkg/account"
	"oko/pkg/e"
	"oko/pkg/ginapp/types"
	"strings"

	"github.com/gin-gonic/gin"
)

type repostHandler struct {
	repository Repository
}

func NewHandler(repo Repository) Handler {
	return &repostHandler{
		repository: repo,
	}
}

// NewRequest godoc
// @Summary New Request
// @Description Create new repost request
// @ID post-repost-new-request
// @Tags Repost
// @Accept json
// @Produce json
// @Param object body repost.NewRequestForm true "Repost create request"
// @Success 200 {object} repost.NewRequestResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /repost/ [post]
// @Security ApiKeyAuth
func (h *repostHandler) NewRequest(c *gin.Context) {
	var form NewRequestForm

	if err := c.Bind(&form); err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	val, exists := c.Get("account_model")
	acc := val.(account.Account)
	model := &Request{
		URL: strings.TrimSpace(form.URL),
	}

	model, err := h.repository.GetOrNil(model)
	if model != nil {
		for _, a := range model.Accounts {
			if a.ID == acc.ID {
				e.ErrorResponse(c, http.StatusBadRequest, "Repost request already exist")
				return
			}
		}
	}

	if !exists {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Fail to create new repost request")
		return
	}

	if model == nil {
		model = &Request{
			URL:   strings.TrimSpace(form.URL),
			Level: 1,
		}
	}
	if err := h.repository.CreateAndAssign(model, acc); err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Fail to create new repost request")
		return
	}

	serializer := Serializer{*model}

	c.JSON(
		http.StatusCreated,
		NewRequestResponse{
			StdResponse: types.StdResponse{
				Status:  http.StatusCreated,
				Message: "Repost request created successfully",
			},
			Data: serializer.To(),
		})
}

// View godoc
// @Summary View repost request details and repost links
// @Description View repost request details and repost links
// @ID get-repost-view
// @Tags Repost
// @Accept json
// @Produce json
// @Param object query repost.RequestForm true "Repost find request (date format ex.: 2006-01-02T15:04:05Z (without time zone) or 2006-01-02T15:04:05-00:00 (with time zone))" //nolint
// @Success 200 {object} repost.NewRequestResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 404 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /repost/view [get]
// @Security ApiKeyAuth
func (h *repostHandler) View(c *gin.Context) {
	var form RequestForm

	if err := c.BindQuery(&form); err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}
	get, _ := c.Get("account_id")

	model, err := h.repository.GetWithDate(form.URL, get.(int), form.DateFrom, form.DateTo)
	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	serializer := Serializer{*model}

	c.JSON(
		http.StatusOK,
		NewRequestResponse{
			StdResponse: types.StdResponse{
				Status:  http.StatusOK,
				Message: "Ok",
			},
			Data: serializer.To(),
		})
}

// List godoc
// @Summary List repost requests
// @Description List repost request
// @ID get-repost-list
// @Tags Repost
// @Accept json
// @Produce json
// @Param object query repost.ListRequest true "Repost find request"
// @Success 200 {object} types.StdResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /repost [get]
// @Security ApiKeyAuth
func (h *repostHandler) List(c *gin.Context) {
	var form ListRequest

	if err := c.ShouldBind(&form); err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	accID, _ := c.Get("account_id")

	models, count, err := h.repository.List(form.PerPage, form.CurrentPage, 1, accID.(int), "")
	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, "Something went wrong")
		return
	}

	serializer := ListSerializer{Requests: models}

	c.JSON(
		http.StatusOK,
		gin.H{
			"data": serializer.To(),
			"meta": types.PaginationResponse{
				PaginationRequest: types.PaginationRequest{
					CurrentPage: form.CurrentPage,
					PerPage:     form.PerPage,
				},
				TotalRecords: count,
				TotalPages:   uint32(math.Ceil(float64(count) / float64(form.PerPage))),
			},
		})
}
