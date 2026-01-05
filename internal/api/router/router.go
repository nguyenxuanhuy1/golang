package router

import (
	"traingolang/internal/api/handler"
	"traingolang/internal/auth"
	"traingolang/internal/config"
	"traingolang/internal/repository"

	// "traingolang/internal/websocket"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// r.MaxMultipartMemory = 1 << 20
	postRepo := repository.NewPostRepo(config.DB)
	imageRepo := repository.NewImageRepository(config.DB)
	// PUBLIC ROUTES
	r.POST("/api/user/register", handler.Register)
	r.POST("/api/user/login", handler.Login)
	r.POST("/api/search/post", handler.SearchPostsHandler(postRepo))
	// PROTECTED ROUTES
	api := r.Group("/api")
	api.Use(auth.Middleware())
	{
		api.POST("/match/create", handler.CreateMatch)
		api.POST("/match/join", handler.JoinMatch)
		api.GET("/profile", handler.Profile)
		api.POST("/upload", handler.UploadHandler)

		api.POST(
			"/create/post",
			auth.AdminOnly(),
			auth.LimitUploadSize(1<<20),
			handler.CreatePost(postRepo, imageRepo),
		)
		api.POST("/update/post/:id",
			auth.AdminOnly(),
			auth.LimitUploadSize(1<<20),
			handler.UpdatePost(postRepo, imageRepo),
		)

		api.POST("/delete/post/:id",
			auth.AdminOnly(),
			handler.DeletePost(postRepo, imageRepo),
		)

	}

	return r
}
