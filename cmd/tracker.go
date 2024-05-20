package main

import (
	"fmt"
	"log"

	"github.com/luckycatx/oat-tracker/internal/handler"
	"github.com/luckycatx/oat-tracker/internal/pkg/conf"

	"github.com/gin-gonic/gin"
)

func main() {
	var cfg = conf.Load()
	fmt.Println(cfg)
	var handler = handler.NewHandler(cfg)

	var r = gin.Default()

	// r.Use(gin.Logger())
	// r.Use(gin.Recovery())

	r.GET("/:room/announce", handler.Announce)
	r.GET("/:room/scrape", handler.Scrape)

	log.Fatal(r.Run(cfg.Host + ":" + cfg.Port))
}
