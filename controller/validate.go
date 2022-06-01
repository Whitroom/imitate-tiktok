package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func BindAndValid(ctx *gin.Context, target interface{}) bool {
	if err := ctx.ShouldBindQuery(target); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, Response{
			StatusCode: 1,
			StatusMsg:  "参数匹配错误",
		})
		return false
	}
	return true
}

func QueryIDAndValid(ctx *gin.Context, queryName string) uint {
	id, err := strconv.ParseUint(ctx.Query(queryName), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			StatusCode: 1,
			StatusMsg:  queryName + "不是数字",
		})
		return 0
	}
	return uint(id)
}
