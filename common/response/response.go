package response

import (
	"github.com/crazyfrankie/gem/gerrors"
	"github.com/gin-gonic/gin"
	"net/http"
	"zmall/common/constant"
)

type Response struct {
	Code int32
	Msg  string
	Data interface{}
}

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code: constant.Success.Code,
		Msg:  constant.Success.Msg,
		Data: data,
	})
}

func Error(ctx *gin.Context, err error) {
	if bizError, ok := gerrors.FromBizStatusError(err); ok {
		ctx.JSON(http.StatusOK, Response{
			Code: bizError.BizStatusCode(),
			Msg:  bizError.BizMessage(),
		})
	}
}
