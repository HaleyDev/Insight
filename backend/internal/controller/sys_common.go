package controller

import (
	"insight/internal/pkg/errors"
	r "insight/internal/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Api struct {
	errors.Error
}

// Success 业务成功响应
func (api *Api) Success(c *gin.Context, data ...any) {
	response := r.Resp()
	if len(data) > 0 {
		response.WithDataSuccess(c, data[0])
		return
	}
	response.Success(c)
}

// FailCode 业务失败响应
func (api *Api) FailCode(c *gin.Context, code int, data ...any) {
	response := r.Resp()
	if len(data) > 0 {
		response.WithData(data[0]).FailCode(c, code)
		return
	}
	response.FailCode(c, code)
}

// Fail 业务失败响应
func (api *Api) Fail(c *gin.Context, code int, message string, data ...any) {
	response := r.Resp()
	if len(data) > 0 {
		response.WithData(data[0]).FailCode(c, code, message)
		return
	}
	response.FailCode(c, code, message)
}

func (api *Api) Err(c *gin.Context, e error) {
	businessError, err := api.AsBusinessError(e)
	if err != nil {
		api.FailCode(c, http.StatusInternalServerError, err.Error())
		return
	}

	api.Fail(c, businessError.GetCode(), businessError.GetMessage())
}
