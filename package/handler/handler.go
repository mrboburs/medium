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
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", handler.signUp)
		auth.POST("/sign-in", handler.signIn)
		auth.POST("/verify", handler.ConfirmEmail)
		//recoveryPassword
		// auth.GET("/recovery")
		// auth.GET("/verify-code")
		// auth.GET("/recovery-password")
	}
	api := router.Group("/api", handler.userIdentity)
	{
		account := api.Group("/account")
		{
			account.PATCH("/upload-image", handler.uploadAccountImage)
			account.POST("/update", handler.UpdateProfile)

			account.GET("/get", handler.getUser)
			account.GET("/resend", handler.resendCodeToEmail)
			account.GET("/search", handler.getUser)
		}
	}
	return router
}
