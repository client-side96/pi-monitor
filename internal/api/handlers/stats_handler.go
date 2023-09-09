package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/client-side96/pi-monitor/internal/service"
	"github.com/client-side96/pi-monitor/internal/sub"
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
		log.Println(err)
	}

	client := h.statsService.Connect()

	go writeWsMessages(ws, h.statsService, client)

	readWsMessages(ws, h.statsService, client)
}

func readWsMessages(
	conn *websocket.Conn,
	statsService service.IStatsService,
	client sub.StatsClient,
) {
	defer func() {
		statsService.Disconnect(client)
		conn.Close()
	}()
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func writeWsMessages(
	conn *websocket.Conn,
	statsService service.IStatsService,
	client sub.StatsClient,
) {

	defer func() {
		statsService.Disconnect(client)
		conn.Close()
	}()

	for {
		s := <-client.Channel
		data, err := json.Marshal(&s)
		if err != nil {
			log.Println(err)
			return
		}
		err = conn.WriteMessage(1, data)
		if err != nil {
			log.Printf("WS ERROR: %s", err)
			return
		}

	}
}
