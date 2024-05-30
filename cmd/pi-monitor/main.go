package main

import (
	"flag"
	"log"
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
var publishInterval int

func mainLoop(env config.Environment, statsService *service.StatsService) {

	ticker := time.NewTicker(time.Duration(env.PublishInterval) * time.Second)
	done := make(chan bool)

	defer func() {
		ticker.Stop()
		done <- true
	}()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			statsService.PublishToAllClients()
		case client := <-statsService.Channel:
			statsService.HandleStatsSubscripition(client)
		}
	}

}

func main() {
	flag.StringVar(
		&scriptDir,
		"script-dir",
		"/srv/pi-monitor/scripts/",
		"Directory to store bash scripts",
	)
	flag.StringVar(
		&addr,
		"address",
		":8000",
		"Port the webserver runs on",
	)
	flag.IntVar(
		&publishInterval,
		"publish-interval",
		1,
		"Interval in seconds in which the statistics are published to clients",
	)
	flag.Parse()

	env := config.Environment{
		ScriptDir:       scriptDir,
		Addr:            addr,
		PublishInterval: publishInterval,
	}

	// Low level
	linuxCommunicator := os.NewLinuxCommunicator(env)

	// Repositories
	statsRepo := repository.NewStatsRepository(linuxCommunicator)

	// Subscriptions
	statsSub := sub.NewStatsSubscription()

	// Services
	statsService := service.NewStatsService(statsSub, statsRepo)

	// Handlers
	statsHandler := handlers.NewStatsHandler(statsService)

	router := api.NewRouter(statsHandler)

	go mainLoop(env, statsService)

	httpServer := &http.Server{
		Addr:    env.Addr,
		Handler: router.SetupRoutes(),
	}

	err := httpServer.ListenAndServe()
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
