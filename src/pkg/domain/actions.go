package domain

import (
	"net/http"
	"oko/pkg/e"
	"oko/pkg/rss"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type domainHandler struct {
	repository Repository
}

func NewHandler(repo Repository) Handler {
	return &domainHandler{
		repository: repo,
	}
}

// List godoc
// @Summary List
// @Description List
// @ID get-domains-list
// @Tags Domain
// @Accept json
// @Produce json
// @Success 200 {object} types.StdResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /domain [get]
func (h *domainHandler) List(c *gin.Context) {
	models, err := h.repository.List()
	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	serializer := ListSerializer{Domains: models}

	c.JSON(
		http.StatusOK,
		gin.H{
			"data": serializer.To(),
		})
}

// Get godoc
// @Summary Get
// @Description Get
// @ID get-domains-get
// @Tags Domain
// @Accept json
// @Produce json
// @Param id path int true "Domain ID"
// @Success 200 {object} types.StdResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /domain/{id} [get]
func (h *domainHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, "Something went wrong")
		return
	}

	model, err := h.repository.Get(uint(id))
	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, "Something went wrong")
		return
	}

	serializer := Serializer{*model}

	c.JSON(
		http.StatusOK,
		gin.H{
			"data": serializer.To(),
		})
}

// Create godoc
// @Summary Create
// @Description Create
// @ID get-domains-create
// @Tags Domain
// @Accept json
// @Produce json
// @Param object body domain.CreateForm true "Domain create request fields"
// @Success 200 {object} types.StdResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /domain [post]
func (h *domainHandler) Create(c *gin.Context) {
	var form CreateForm

	if err := c.ShouldBind(&form); err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	rsses := []*rss.Rss{{
		Link: form.RssLink,
	}}
	model := &Domain{
		Name:             form.Name,
		Rss:              rsses,
		TelegramUsername: form.TelegramUsername,
	}

	if err := h.repository.Create(model); err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, "Something went wrong")
		return
	}

	serializer := Serializer{*model}

	c.JSON(
		http.StatusOK,
		gin.H{
			"data": serializer.To(),
		})
}

// Update godoc
// @Summary Update
// @Description Update
// @ID get-domains-update
// @Tags Domain
// @Accept json
// @Produce json
// @Param id path int true "Domain ID"
// @Param object body domain.UpdateForm true "Domain update request fields"
// @Success 200 {object} types.StdResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /domain/{id} [put]
func (h *domainHandler) Update(c *gin.Context) {
	var form UpdateForm

	if err := c.ShouldBind(&form); err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, "Something went wrong")
		return
	}

	model := &Domain{
		Name:             form.Name,
		TelegramUsername: form.TelegramUsername,
	}

	if err := h.repository.Update(uint(id), model); err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, "Something went wrong")
	}
}

// Delete godoc
// @Summary Delete
// @Description Delete
// @ID get-domains-delete
// @Tags Domain
// @Accept json
// @Produce json
// @Param id path int true "Domain ID"
// @Success 200 {object} types.StdResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /domain/{id} [delete]
// @Security ApiKeyAuth
func (h *domainHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, "Something went wrong")
		return
	}

	model := &Domain{
		Model: gorm.Model{
			ID: uint(id),
		},
	}

	if err := h.repository.Delete(model); err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, "Something went wrong")
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"data": "ok",
		})
}
