package routes

import (
	"github.com/ahmedmalmoselhy/unione_go/internal/handlers"
	"github.com/ahmedmalmoselhy/unione_go/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/register", handlers.Register)
			auth.GET("/me", middlewares.AuthMiddleware(), handlers.GetMe)
		}
	}

	return r
}
