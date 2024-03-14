package server

import (
	"context"
	"log"

	"mmf/config"
	"mmf/pkg/redis"
	"mmf/pkg/server/api"
	"mmf/wires"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		config: config,
	}
}

func (server *Server) Start() {
	//initCrawler(server.config.redis)

	r := gin.Default()
	r.Use(gin.Logger())

	wires.Init()
	redis.Init(server.config, context.Background())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	api.RegisterVersion(r, context.Background())

	err := r.Run(":" + server.config.Server.Port)

	if err != nil {
		log.Fatal("Could not start the server" + err.Error())
		return
	}

	println("Starting server on port: " + server.config.Server.Port)
}