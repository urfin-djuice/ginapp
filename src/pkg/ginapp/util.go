package ginapp

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUint32PathParam(c *gin.Context, pathParamName string) (uint32, error) {
	param := c.Param(pathParamName)
	parseInt, err := strconv.ParseInt(param, 10, 32)
	return uint32(parseInt), err
}
