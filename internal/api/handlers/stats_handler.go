package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/client-side96/pi-monitor/internal/service"
	"github.com/client-side96/pi-monitor/internal/sub"
	"github.com/client-side96/pi-monitor/internal/util"
	"github.com/gorilla/websocket"
)

type IStatsHandler interface {
	ConnectStatsWS(http.ResponseWriter, *http.Request)
}

type StatsHandler struct {
	upgrader     websocket.Upgrader
	statsService service.IStatsService
}

func NewStatsHandler(statsService service.IStatsService) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
	}
}

func (h *StatsHandler) ConnectStatsWS(w http.ResponseWriter, r *http.Request) {

	h.upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	ws, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(util.ErrWSUpgrade)
	}

	client := h.statsService.Connect()

	defer cleanupConnection(ws, h.statsService, client)

	go util.WriteWS(ws, func() ([]byte, error) {
		return readStatsChannel(client)
	})

	util.ReadWS(ws)
}

func readStatsChannel(client sub.StatsClient) ([]byte, error) {
	s := <-client.Channel
	return json.Marshal(&s)
}

func cleanupConnection(c *websocket.Conn, s service.IStatsService, cl sub.StatsClient) {
	s.Disconnect(cl)
	c.Close()
}
