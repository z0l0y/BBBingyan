package main

import (
	"BBBingyan/internal/log"
	"BBBingyan/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	log.InitFile("logs")
	router := gin.New()
	router.Use(middleware.LogMiddleware())
	middleware.RedisInit()
	router.POST("/send-verification-code", middleware.SendVerificationCode)
	router.POST("/verify-verification-code", middleware.VerifyVerificationCode)
	err := router.Run(":8080")
	if err != nil {
		return
	}
}
