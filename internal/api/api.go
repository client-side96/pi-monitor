package api

import (
	"net/http"

	"github.com/client-side96/pi-monitor/internal/api/handlers"
)

type Router struct {
	statsHandler handlers.IStatsHandler
}

func NewRouter(statsHandler handlers.IStatsHandler) *Router {
	return &Router{
		statsHandler: statsHandler,
	}
}

func (r *Router) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/monitor/api/stats", r.statsHandler.ConnectStatsWS)
	return mux
}
