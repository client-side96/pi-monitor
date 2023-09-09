package main

import (
	"flag"
	"net/http"
	"time"

	"github.com/client-side96/pi-monitor/internal/config"
	"github.com/client-side96/pi-monitor/internal/os"
	"github.com/client-side96/pi-monitor/internal/stats"
	"github.com/client-side96/pi-monitor/internal/web"
)

var scriptDir string
var addr string

func mainLoop(statsService *stats.StatsService) {
	ticker := time.NewTicker(1 * time.Second)

	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case client := <-statsService.Channel:
			statsService.HandleClientSubscription(client)
		default:
			statsService.PublishStats()

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
	statsService := stats.NewStatsService(linuxCommunicator)

	go mainLoop(statsService)

	server := web.NewServer(statsService)
	httpServer := &http.Server{
		Addr:    env.Addr,
		Handler: server.GetRoutes(),
	}

	httpServer.ListenAndServe()
}
