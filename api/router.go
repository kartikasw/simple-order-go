package api

import (
	"simple-order-go/internal/handler"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router       *gin.Engine
	orderHandler handler.OrderHandler
}

func NewServer(orderHandler handler.OrderHandler) *Server {
	server := &Server{orderHandler: orderHandler}
	server.setupRouter()
	return server
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/orders", server.orderHandler.CreateOrder)
	router.GET("/orders", server.orderHandler.GetAllOrders)
	router.GET("/orders/:id", server.orderHandler.GetOrderByID)
	router.PUT("/orders/:id", server.orderHandler.UpdateOrder)
	router.DELETE("/orders/:id", server.orderHandler.DeleteOrder)

	server.router = router
}

func (server *Server) Start(port string) error {
	return server.router.Run(port)
}
