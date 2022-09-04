package server

import (
	"strconv"
	"test-app/controller"

	"github.com/gin-gonic/gin"
)

type GinServer struct {
	port   int32
	router *gin.Engine
}

func (s *GinServer) Run() {
	s.router.Run(":" + strconv.Itoa(int(s.port)))
}

func NewGinServer(port int32) Server {
	r := gin.Default()
	addRoutes(r)

	return &GinServer{
		port:   port,
		router: r,
	}
}

func addRoutes(r *gin.Engine) {
	r.GET("/", controller.Hello)
}
