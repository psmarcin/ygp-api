package main

import (
	"github.com/psmarcin/youtubegoespodcast/internal/adapters"
	application "github.com/psmarcin/youtubegoespodcast/internal/app"
	"github.com/psmarcin/youtubegoespodcast/internal/config"
	"github.com/psmarcin/youtubegoespodcast/internal/logger"
	"github.com/psmarcin/youtubegoespodcast/internal/ports"
	"github.com/sirupsen/logrus"
)

var l = logrus.WithField("source", "cmd")

func main() {
	// Config
	config.Init()
	traceFlusher := config.InitTracer(config.Cfg)
	defer traceFlusher()
	// Logger
	logger.Setup()

	// Cache
	cacheRepository, err := adapters.NewCacheRepository()
	if err != nil {
		l.WithError(err).Fatalf("can't create cache repository")
	}
	cacheService := application.NewCacheService(&cacheRepository)

	// YouTube
	yt, err := adapters.NewYouTube()
	if err != nil {
		l.WithError(err).Fatalf("can't connect to youtube service")
	}

	youTubeAPIRepository := adapters.NewYouTubeAPIRepository(yt)
	youTubeRepository, err := adapters.NewYouTubeRepository()
	if err != nil {
		l.WithError(err).Fatalf("can't create youtube request repository")
	}

	youTubeService := application.NewYouTubeService(youTubeRepository, youTubeAPIRepository, cacheService)

	// YTDL
	//ytdlRepository := adapters.NewYTDLRepository()
	//ytdlService := application.NewYTDLService(ytdlRepository)

	// Video
	fileService := application.NewFileService()

	// API
	fiberServer := ports.CreateHTTPServer()
	h := ports.NewHttpServer(fiberServer, youTubeService, fileService)
	app := h.Serve()

	logrus.Fatal(app.Listen(":" + config.Cfg.Port))
}
