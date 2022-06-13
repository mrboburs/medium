package main

import (
	// "fmt"
	// "fmt"
	"mediumuz/configs"
	"mediumuz/package/handler"
	"mediumuz/package/repository"
	"mediumuz/package/service"
	"mediumuz/server"
	"mediumuz/util/logrus"

	_ "github.com/lib/pq"
)

// @title MediumuZ API
// @version 1.0
// @description API Server for MediumuZ Application
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @contact.name   Mr Bobur

func main() {

	logrus := logrus.GetLogger()
	logrus.Info("send email")

	configs, err := configs.InitConfig()
	logrus.Infof("configs %v", configs)
	if err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}
	logrus.Info("successfull checked configs.")
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     configs.DBHost,
		Port:     configs.DBPort,
		Username: configs.DBUsername,
		DBName:   configs.DBName,
		SSLMode:  configs.DBSSLMode,
		Password: configs.DBPassword,
	}, logrus)

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	logrus.Info("successfull connection DB")

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services, logrus, configs)

	server := new(server.Server)
	err = server.Run(configs.HTTPPort, handlers.InitRoutes())

	if err != nil {
		logrus.Fatalf("error occurred while running http server: %s", err.Error())
	}

	defer logrus.Infof("run server port:%v", configs.HTTPPort)
}
