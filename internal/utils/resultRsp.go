package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// ResultRsp 类似于Java中的Result工具类类
type ResultRsp struct {
	Code int         `json:"code"`
	Msg  string      `json:"success"`
	Data interface{} `json:"data"`
}

// 下面是Result类的方法

// Success 正确状态处理，状态码已经包含在StatusOK里面了
func Success(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  msg,
		"data": data,
	})
}

// Error 错误状态处理，状态码已经包含在StatusBadRequest里面了
func Error(c *gin.Context, msg string) {
	// 这里的JSON类似于Java的@RestController
	c.JSON(http.StatusBadRequest,
		gin.H{
			"error": msg,
		})
}
