package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/client-side96/pi-monitor/internal/api"
	"github.com/client-side96/pi-monitor/internal/api/handlers"
	"github.com/client-side96/pi-monitor/internal/config"
	"github.com/client-side96/pi-monitor/internal/os"
	"github.com/client-side96/pi-monitor/internal/repository"
	"github.com/client-side96/pi-monitor/internal/service"
	"github.com/client-side96/pi-monitor/internal/sub"
)

var scriptDir string
var addr string

func mainLoop(statsService *service.StatsService) {
	ticker := time.NewTicker(1 * time.Second)

	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case client := <-statsService.Channel:
			statsService.HandleStatsSubscripition(client)
		default:
			statsService.PublishToAllClients()

		}
	}
}

func main() {
	flag.StringVar(
		&scriptDir,
		"script-dir",
		"/home/admin/scripts/",
		"Directory to store bash scripts",
	)
	flag.StringVar(
		&addr,
		"address",
		":4000",
		"Port the webserver runs on",
	)
	flag.Parse()

	env := config.Environment{
		ScriptDir: scriptDir,
		Addr:      addr,
	}

	linuxCommunicator := os.NewLinuxCommunicator(env)

	statsRepo := repository.NewStatsRepository(linuxCommunicator)

	statsSub := sub.NewStatsSubscription()

	statsService := service.NewStatsService(statsSub, statsRepo)

	statsHandler := handlers.NewStatsHandler(statsService)

	router := api.NewRouter(statsHandler)

	go mainLoop(statsService)

	httpServer := &http.Server{
		Addr:    env.Addr,
		Handler: router.SetupRoutes(),
	}

	httpServer.ListenAndServe()
}
