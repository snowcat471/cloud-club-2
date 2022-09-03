package server

import (
	"strconv"
	"test-app/controller"

	"github.com/gin-gonic/gin"
)

type GinServer struct {
	port int32
}

func (s *GinServer) Run() {
	r := gin.Default()
	s.addRoute(r)

	r.Run(":" + strconv.Itoa(int(s.port)))
}

func (s *GinServer) addRoute(r *gin.Engine) {
	r.GET("/", controller.Hello)
}

func NewGinServer(port int32) Server {
	return &GinServer{
		port: port,
	}
}
