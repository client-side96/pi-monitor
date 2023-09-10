package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/client-side96/pi-monitor/internal/api/handlers"
	"github.com/client-side96/pi-monitor/internal/model"
	"github.com/client-side96/pi-monitor/internal/sub"
	"github.com/client-side96/pi-monitor/internal/testutil"
	"github.com/client-side96/pi-monitor/internal/testutil/mocks"
	"github.com/gorilla/websocket"
)

func readMessage(t *testing.T, ws *websocket.Conn) []byte {
	_, bytes, err := ws.ReadMessage()
	if err != nil {
		t.Errorf("%v", err)
	}
	return bytes
}

func readJSONMessage[T any](t *testing.T, ws *websocket.Conn) *T {
	var result T
	bytes := readMessage(t, ws)
	err := json.Unmarshal(bytes, &result)
	if err != nil {
		t.Errorf("%v", err)
	}
	return &result
}

func createWS(t *testing.T, hf http.HandlerFunc) *websocket.Conn {
	srv := httptest.NewServer(http.HandlerFunc(hf))
	url := "ws" + strings.TrimPrefix(srv.URL, "http")

	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	return ws
}

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

	ws := createWS(t, h.ConnectStatsWS)
	defer ws.Close()

	go func() {
		s.Client.Channel <- stats
	}()

	result = readJSONMessage[model.StatsModel](t, ws)

	testutil.AssertEqual(t, stats.Temperature, result.Temperature)
	testutil.AssertEqual(t, stats.CPULoad, result.CPULoad)
	testutil.AssertEqual(t, stats.MemoryUsage, result.MemoryUsage)
}

func TestHandler_Stats_ConnectStatsWS_DisconnectChannel(t *testing.T) {
	s := &mocks.MockStatsService{
		Client: sub.StatsClient{
			Status:  sub.Connected,
			Channel: make(chan *model.StatsModel),
			ID:      1,
		},
		Called: "Not called",
	}
	stats := &model.StatsModel{
		Temperature: 10,
		CPULoad:     20,
		MemoryUsage: 30,
	}
	h := handlers.NewStatsHandler(s)

	ws := createWS(t, h.ConnectStatsWS)
	defer ws.Close()

	go func() {
		s.Client.Channel <- stats
		err := ws.WriteMessage(websocket.CloseMessage, []byte{})
		if err != nil {
			t.Error(err)
		}
	}()

	readMessage(t, ws)

	testutil.AssertEqual(t, s.Called, mocks.DisconnectCalled)
}
