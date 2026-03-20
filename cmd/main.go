package main

import (
	"forum/internal/config"
	database "forum/internal/db"
	"forum/internal/repository"
	"forum/internal/server"
	"forum/internal/service"
	"log"
)

func main() {
	cfg := config.GetConfig()
	db, err := database.InitDatabase(cfg.DbCfg)
	if err != nil {
		log.Fatal(err)

	}
	repo := repository.NewRepository(db)
	srvc := service.NewService(repo)
	srv := new(server.Server)
	if err := srv.Run(cfg.ServerCfg, cfg.WebCfg, srvc); err != nil {
		log.Fatal("Cannot Start the Server")
	}

}
