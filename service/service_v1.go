package service

import (
	"github.com/soyvural/simple-device-api/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	ginSwagger "github.com/swaggo/gin-swagger"
)

type Service struct {
	router    *gin.Engine
	deviceSvc *device
}

func New(r *gin.Engine) *Service {
	r.Use(requestIDMiddleware())
	return &Service{
		router:    r,
		deviceSvc: newDeviceService(store.NewCache()),
	}
}

func (s *Service) SetRoute_v1() {
	v1 := s.router.Group("/api/v1")
	deviceRouter := v1.Group("/devices")

	deviceRouter.POST("", s.deviceSvc.CreateDevice)
	deviceRouter.GET(":id", s.deviceSvc.GetDevice)
	deviceRouter.DELETE(":id", s.deviceSvc.DeleteDevice)

	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Request-Id", uuid.NewString())
		c.Next()
	}
}
