package handler

import (
	"fmt"
	"mediumuz/package/service"
	"mediumuz/util/logrus"

	"mediumuz/configs"
	_ "mediumuz/docs"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type Handler struct {
	services *service.Service
	logrus   *logrus.Logger
	config   *configs.Configs
}

func NewHandler(services *service.Service, logrus *logrus.Logger, config *configs.Configs) *Handler {
	return &Handler{services: services, logrus: logrus, config: config}
}

func (handler *Handler) InitRoutes() *gin.Engine {
	config := handler.config
	fmt.Println(config)
	// docs.SwaggerInfo_swagger.Title = config.AppName
	// docs.SwaggerInfo_swagger.Version = config.Version
	// docs.SwaggerInfo_swagger.Host = config.ServiceHost + config.HTTPPort
	// docs.SwaggerInfo_swagger.Schemes = []string{ "https"}
	router := gin.New()
	router.Static("/public", "./public/")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	auth := router.Group("/auth")

	api := router.Group("/api", handler.userIdentity)
	{
		auth.POST("/sign-up", handler.signUp)
		auth.POST("/sign-in", handler.signIn)
		api.POST("/verify", handler.ConfirmByEmail)
		//recoveryPassword
		// auth.GET("/recovery")
		// auth.GET("/verify-code")
		// auth.GET("/recovery-password")
	}
	{
		account := api.Group("/account")
		{
			account.PATCH("/upload-image", handler.uploadAccountImage)
			account.POST("/update", handler.UpdateProfile)
			account.GET("/get", handler.GetUserById)
			account.DELETE("/delete", handler.DeleteUser)
			account.GET("/getUsers", handler.GetAllUsers)

		}
		post := api.Group("/post")
		{
			post.POST("/create", handler.createPost)
			router.Group("/post").GET("/get-comments", handler.getComments)
			post.GET("/count-like/:id", handler.ClickLike)
			post.POST("/comment", handler.commentPost)
			post.GET("/view", handler.viewPost)
			post.POST("/update/:id", handler.UpdatePost)
			router.Group("/post").GET("/get/:id", handler.getPostID)
			router.Group("/post").GET("/get-all", handler.GetAllPosts)
			post.DELETE("/delete/:id", handler.PostDelete)
			post.PATCH("/upload/:id", handler.uploadPostImg)
		}
	}
	return router
}
