package rss

import (
	"net/http"
	"oko/pkg/e"
	"oko/pkg/ginapp"
	"oko/pkg/ginapp/types"

	"github.com/gin-gonic/gin"
)

type rssHandler struct {
	rep *rssRepository
}

// Create godoc
// @Summary Create
// @Description Create
// @ID create-rss
// @Tags Rss
// @Accept json
// @Produce json
// @Param object body rss.CreateForm true "Rss update request fields"
// @Success 200 {object} types.StdResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /rss [post]
func (r rssHandler) Create(c *gin.Context) {
	req := CreateForm{}
	if err := c.BindJSON(&req); err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, "Something went wrong")
		return
	}

	if err := r.rep.CreateRss(req.DomainID, req.RssLink); err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}
	types.SuccessEmptyResponse(c)
}

// Update godoc
// @Summary Update
// @Description Update
// @ID update-rss
// @Tags Rss
// @Accept json
// @Produce json
// @Param object body rss.UpdateForm true "Rss update request fields"
// @Success 200 {object} types.StdResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /rss [put]
func (r rssHandler) Update(c *gin.Context) {
	req := UpdateForm{}
	if err := c.BindJSON(&req); err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, "Something went wrong")
		return
	}

	if err := r.rep.UpdateRss(req.DomainID, req.RssLink, req.RssID); err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	types.SuccessEmptyResponse(c)
}

// Delete godoc
// @Summary Delete
// @Description Delete
// @ID delete-rss
// @Tags Rss
// @Accept json
// @Produce json
// @Param id path int true "Rss ID"
// @Success 200 {object} types.StdResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /rss/{id} [delete]
func (r rssHandler) Delete(c *gin.Context) {
	rssID, err := ginapp.GetUint32PathParam(c, "id")
	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, "Something went wrong")
		return
	}

	if err = r.rep.DeleteRss(uint(rssID)); err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, "Something went wrong")
		return
	}

	types.SuccessEmptyResponse(c)
}

// List godoc
// @Summary List
// @Description List
// @ID get-rss-list
// @Tags Rss
// @Accept json
// @Produce json
// @Success 200 {object} types.StdResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /rss [get]
func (r rssHandler) List(c *gin.Context) {
	list, err := r.rep.List()
	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	types.SuccessResponse(c, list)
}

// Get godoc
// @Summary Get
// @Description Get
// @ID get-rss
// @Tags Rss
// @Accept json
// @Produce json
// @Param id path int true "Rss ID"
// @Success 200 {object} types.StdResponse
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /rss/{id} [get]
func (r rssHandler) Get(c *gin.Context) {
	id, err := ginapp.GetUint32PathParam(c, "id")
	if err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, "Something went wrong")
		return
	}
	rss, err := r.rep.Get(uint(id))
	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	types.SuccessResponse(c, rss)
}
