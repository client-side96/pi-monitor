package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/client-side96/pi-monitor/internal/api/handlers"
	"github.com/client-side96/pi-monitor/internal/mocks"
	"github.com/client-side96/pi-monitor/internal/model"
	"github.com/client-side96/pi-monitor/internal/sub"
	"github.com/gorilla/websocket"
)

func TestHandler_Stats_ConnectStatsWS_StatsFromChannel(t *testing.T) {
	var result *model.StatsModel
	s := &mocks.MockStatsService{
		Client: sub.StatsClient{
			Status:  sub.Connected,
			Channel: make(chan *model.StatsModel),
			ID:      1,
		},
	}
	stats := &model.StatsModel{
		Temperature: 10,
		CPULoad:     20,
		MemoryUsage: 30,
	}
	h := handlers.NewStatsHandler(s)

	srv := httptest.NewServer(http.HandlerFunc(h.ConnectStatsWS))
	url := "ws" + strings.TrimPrefix(srv.URL, "http")

	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	go func() {
		s.Client.Channel <- stats
	}()

	_, bytes, err := ws.ReadMessage()
	if err != nil {
		t.Errorf("%v", err)
	}
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		t.Errorf("%v", err)
	}

	if result.Temperature != stats.Temperature {
		t.Errorf("Expected %f but got %f", stats.Temperature, result.Temperature)
	}
	if result.CPULoad != stats.CPULoad {
		t.Errorf("Expected %f but got %f", stats.CPULoad, result.CPULoad)
	}
	if result.MemoryUsage != stats.MemoryUsage {
		t.Errorf("Expected %f but got %f", stats.MemoryUsage, result.MemoryUsage)
	}
}
