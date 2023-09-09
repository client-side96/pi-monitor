package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/client-side96/pi-monitor/internal/stats"
	"github.com/gorilla/websocket"
)

type Server struct {
	upgrader     websocket.Upgrader
	statsService *stats.StatsService
}

func NewServer(statsService *stats.StatsService) *Server {
	return &Server{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		statsService: statsService,
	}
}

func (s *Server) GetRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/stats", s.wsStats)
	return mux
}

func (s *Server) wsStats(w http.ResponseWriter, r *http.Request) {
	s.upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	statsConnection := s.statsService.EstablishConnection()

	go writeWsMessages(ws, s.statsService, statsConnection)

	readWsMessages(ws, s.statsService, statsConnection)
}

func readWsMessages(
	conn *websocket.Conn,
	statsService *stats.StatsService,
	statsConnection stats.Client,
) {
	defer func() {
		statsService.CloseConnection(statsConnection)
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
	statsService *stats.StatsService,
	statsConnection stats.Client,
) {

	defer func() {
		statsService.CloseConnection(statsConnection)
		conn.Close()
	}()

	for {
		select {
		case s := <-statsConnection.Channel:
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
}
