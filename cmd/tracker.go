package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"gioui.org/app"
	"github.com/gin-gonic/gin"

	"github.com/luckycatx/oat-tracker/internal/handler"
	"github.com/luckycatx/oat-tracker/internal/pkg/conf"
	"github.com/luckycatx/oat-tracker/internal/view"
)

func main() {
	go func() {
		var window = new(app.Window)
		window.Option(app.Title("Oat-Tracker"))
		if err := view.Run(window, tracker); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func tracker(out io.Writer) {
	var cfg = conf.Load()
	var handler = handler.New(cfg)
	go handler.Cleanup()

	gin.DefaultWriter, gin.DefaultErrorWriter = out, out
	var r = gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.GET("/:room/announce", pausable(handler.Announce))
	r.GET("/:room/scrape", pausable(handler.Scrape))

	log.Fatal(r.Run(cfg.Host + ":" + cfg.Port))
}

func pausable(f gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if view.Paused {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "Tracker is paused",
			})
			return
		}
		f(c)
	}
}
