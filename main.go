package main

import (
	"context"
	"log"

	"github.com/chandrasitinjak/integrate-pokeapi/config"
	"github.com/chandrasitinjak/integrate-pokeapi/internal/db"
	"github.com/chandrasitinjak/integrate-pokeapi/internal/handler"
	"github.com/chandrasitinjak/integrate-pokeapi/internal/logger"
	"github.com/chandrasitinjak/integrate-pokeapi/internal/repository"
	"github.com/chandrasitinjak/integrate-pokeapi/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	defer logger.Log.Sync()

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Error("failed load config", zap.Error(err))
		return
	}

	mysqlDB, err := db.InitMySQL(cfg)
	if err != nil {
		logger.Log.Error("failed connect db", zap.Error(err))
		return
	}
	defer mysqlDB.Close()

	rdb := db.InitRedis(cfg)
	redisSvc := services.NewRedisService(rdb)

	pokeRepo := repository.NewPokemonRepository(mysqlDB)
	pokeService := services.NewPokemonService(pokeRepo, cfg, redisSvc)
	pokeHandler := handler.NewPokemonHandler(pokeService)

	c := cron.New(cron.WithSeconds())
	// cronExpr := fmt.Sprintf("0 */%d * * * *", 15)
	cronExpr := "*/15 * * * * *"

	_, err = c.AddFunc(cronExpr, func() {
		log.Println("Running background job: refresh Pokémon data")
		if err := pokeService.Sync(context.Background()); err != nil {
			logger.Log.Error("Error running Sync()", zap.Error(err))
		} else {
			logger.Log.Info("Sync completed successfully")
		}
	})
	if err != nil {
		logger.Log.Error("Failed to schedule cron job", zap.Error(err))
	}

	c.Start()
	logger.Log.Info("Cron job started. Waiting...")
	// Router
	r := gin.Default()
	r.GET("/items", pokeHandler.GetAll)
	r.POST("/sync", pokeHandler.Sync)

	r.Run(":" + cfg.Port)
	logger.Log.Info("server is started")

}
