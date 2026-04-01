package middleware

import (
	"homework04/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 全局异常处理
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err != nil {
			if bizErr, ok := err.Err.(*utils.BizError); ok {
				// 记录业务错误
				utils.ErrorLogger.Printf("BizError: code=%d, msg=%s\n", bizErr.Code, bizErr.Msg)
				// 返回 HTTP 状态码
				httpStatus := 400
				if bizErr.Code >= 500 {
					httpStatus = 500
				} else if bizErr.Code == 401 {
					httpStatus = 401
				} else if bizErr.Code == 403 {
					httpStatus = 403
				} else if bizErr.Code == 404 {
					httpStatus = 404
				}

				c.JSON(httpStatus, utils.Response{
					Code:    bizErr.Code,
					Message: bizErr.Msg,
				})
			} else {
				// 未知错误
				utils.ErrorLogger.Printf("UnknownError: %v\n", err.Err)
				c.JSON(http.StatusInternalServerError, utils.Response{
					Code:    500,
					Message: "internal server error",
				})
			}
			c.Abort()
		}
	}
}
