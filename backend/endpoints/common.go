package endpoints

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func parsePaginationQueryParams(ctx *gin.Context) (int, int, string) {
	page := 1
	limit := 10
	sortBy := "created_at desc"
	if ctx.Query("page") != "" {
		page, _ = strconv.Atoi(ctx.Query("page"))
	}
	if ctx.Query("limit") != "" {
		limit, _ = strconv.Atoi(ctx.Query("limit"))
	}
	if ctx.Query("sortBy") != "" {
		sortBy = ctx.Query("sortBy")
	}
	return page, limit, sortBy
}

type PaginationResponse struct {
	Total int64       `json:"total"`
	Data  interface{} `json:"data"`
}
