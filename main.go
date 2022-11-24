package main

import (
	"log"
	_ "net/http"

	_ "github.com/AlexandrBurak/TaxiApp/docs"
	"github.com/AlexandrBurak/TaxiApp/internal/cache"
	"github.com/AlexandrBurak/TaxiApp/internal/config"
	"github.com/AlexandrBurak/TaxiApp/internal/handlers"
	"github.com/AlexandrBurak/TaxiApp/internal/logger"
	"github.com/AlexandrBurak/TaxiApp/internal/middleware"
	_ "github.com/AlexandrBurak/TaxiApp/internal/migrations"
	"github.com/AlexandrBurak/TaxiApp/internal/repository"
	"github.com/AlexandrBurak/TaxiApp/internal/service/AuthService"
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//@title UserService
//@version 1.0
//description  user service for taxi app

//@host localhost:8080
//@BasePath /

// @securityDefinitions.apikey SignIn
// @in header
// @name Authorization
func main() {
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatal(err)
	}
	repos := getRepos()
	cache, err := cache.NewCache()
	if err != nil {
		log.Fatal(err)
	}
	service := AuthService.NewService(repos, logger, *cache)
	handler := handlers.NewHandler(service)

	startServer(handler)
}

func getRepos() repository.Repository {
	cfg, err := config.GetDbCfg()
	if err != nil {
		log.Fatal(err)
	}
	repos, err := repository.NewRepository(cfg)
	if err != nil {
		log.Fatal(err)
	}
	return *repos

}

func startServer(handler handlers.Handler) {
	r := gin.Default()
	r.Use(ginprom.PromMiddleware(nil))
	pprof.Register(r)
	r.Use(middleware.SetUserStatus())

	r.GET("/metrics", ginprom.PromHandler(promhttp.Handler()))

	r.POST("/register", middleware.EnsureNotLoggedIn(), handler.Register)

	r.POST("/login", middleware.EnsureNotLoggedIn(), handler.Login)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.GET("/logout", middleware.EnsureLoggedIn(), handler.Logout)
	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
