package repost

import (
	"bytes"
	"encoding/csv"
	"net/http"
	"oko/pkg/e"

	"github.com/gin-gonic/gin"
)

// Export godoc
// @Summary Export all respost to csv table
// @Description Export all respost to csv table
// @ID export-repost
// @Tags Repost
// @Accept json
// @Produce json
// @Param object query repost.RequestForm true "Repost find request (date format ex.: 2006-01-02T15:04:05Z or 2006-01-02T15:04:05-00:00 with time zone)" //nolint
// @Success 200
// @Failure 400 {object} types.ResponseErrorSwg
// @Failure 404 {object} types.ResponseErrorSwg
// @Failure 500 {object} types.ResponseErrorSwg
// @Router /repost/export [get]
// @Security ApiKeyAuth
func (h *repostHandler) Export(c *gin.Context) {
	r := &RequestForm{}
	var err error

	if err = c.ShouldBind(r); err != nil {
		e.ErrorResponse(c, http.StatusBadRequest, "Fail to write header")
		return
	}

	accID, _ := c.Get("account_id")

	export, err := h.repository.GetForExport(r.URL, accID.(int), r.DateFrom, r.DateTo)
	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Error fetch datas for export")
		return
	}
	buff, err := writeToCsv(export)
	if err != nil {
		e.ErrorResponse(c, http.StatusInternalServerError, "Fail get export records")
		return
	}
	extraHeaders := map[string]string{
		"Content-Disposition": `attachment; filename="export.csv"`,
	}
	c.DataFromReader(http.StatusOK, int64(buff.Len()), "application/octet-stream", buff, extraHeaders)
}

func writeToCsv(exportRecs []RecordForExport) (buffer *bytes.Buffer, err error) {
	header := []string{"заголовок репоста", "ссылка репоста", "родительская ссылка", "уровень вложенности репоста", "дата публикации"}
	buffer = bytes.NewBuffer([]byte{})
	writer := csv.NewWriter(buffer)
	err = writer.Write(header)
	if err != nil {
		return
	}

	csvRec := make([]string, 5)
	for _, exp := range exportRecs {
		exp.fillCsvRec(csvRec)
		if err = writer.Write(csvRec); err != nil {
			return
		}
	}

	writer.Flush()
	return
}
