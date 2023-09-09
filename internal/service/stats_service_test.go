package service_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/client-side96/pi-monitor/internal/mocks"
	"github.com/client-side96/pi-monitor/internal/model"
	"github.com/client-side96/pi-monitor/internal/service"
	"github.com/client-side96/pi-monitor/internal/sub"
)

var subscription *sub.StatsSubscription

var temp = 50.04
var cpu = 22.22
var mem = 15.24

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestService_Stats_GetStats(t *testing.T) {
	subscription = sub.NewStatsSubscription()
	repo := &mocks.MockStatsRepository{
		Temp: temp,
		CPU:  cpu,
		Mem:  mem,
	}
	s := service.NewStatsService(subscription, repo)

	result := s.GetStats()

	if result.Temperature != temp {
		t.Errorf("Expected %f but got %f", temp, result.Temperature)
	}
	if result.CPULoad != cpu {
		t.Errorf("Expected %f but got %f", cpu, result.CPULoad)
	}
	if result.MemoryUsage != mem {
		t.Errorf("Expected %f but got %f", mem, result.MemoryUsage)
	}
}

func TestService_Stats_Connect_SingleConnection(t *testing.T) {
	var result sub.StatsClient
	subscription = sub.NewStatsSubscription()
	repo := &mocks.MockStatsRepository{}
	s := service.NewStatsService(subscription, repo)

	go s.Connect()

	select {
	case result = <-s.Channel:
		if result.ID != 1 {
			t.Errorf("Expected 1 but got %d", result.ID)
		}
		if result.Status != sub.Connected {
			t.Errorf("Expected to be connected but got %s", result.Status)
		}
	}

}

func TestService_Stats_Connect_MultipleConnections(t *testing.T) {
	var connectionCount = 1
	var result sub.StatsClient
	subscription = sub.NewStatsSubscription()
	repo := &mocks.MockStatsRepository{}
	s := service.NewStatsService(subscription, repo)

	go s.Connect()
	go s.Connect()

	for connectionCount < 3 {
		select {
		case result = <-s.Channel:
			if result.ID != connectionCount {
				t.Errorf("Expected %d but got %d", connectionCount, result.ID)
			}
			if result.Status != sub.Connected {
				t.Errorf("Expected to be Connected but got %s", result.Status)
			}
			connectionCount++
		}
	}

}

func TestService_Stats_Disconnect_SingleConnection(t *testing.T) {
	var count = 0
	var result sub.StatsClient
	client := sub.StatsClient{
		Status:  sub.Connected,
		Channel: make(chan *model.StatsModel),
		ID:      1,
	}
	subscription = sub.NewStatsSubscription()
	subscription.SetClients(client)
	repo := &mocks.MockStatsRepository{}
	s := service.NewStatsService(subscription, repo)

	go func() {
		s.Channel <- client
	}()

	for count < 2 {
		select {
		case result = <-s.Channel:
			if count == 0 {
				go func() {
					s.Disconnect(client)
				}()
			}
			count++
		}
	}

	if result.Status != sub.Disconnected {
		t.Errorf("Expected to be Disconnected but got %s", result.Status)
	}
}

func TestService_Stats_PublishToAllClients(t *testing.T) {
	var result *model.StatsModel
	subscription = sub.NewStatsSubscription()
	client := sub.StatsClient{
		Status:  sub.Connected,
		Channel: make(chan *model.StatsModel),
		ID:      1,
	}
	subscription.SetClients(client)
	repo := &mocks.MockStatsRepository{
		Temp: temp,
		CPU:  cpu,
		Mem:  mem,
	}
	s := service.NewStatsService(subscription, repo)

	go s.PublishToAllClients()

	select {
	case result = <-client.Channel:
	}

	if result.Temperature != temp {
		t.Errorf("Expected %f but got %f", temp, result.Temperature)
	}
	if result.CPULoad != cpu {
		t.Errorf("Expected %f but got %f", cpu, result.CPULoad)
	}
	if result.MemoryUsage != mem {
		t.Errorf("Expected %f but got %f", mem, result.MemoryUsage)
	}
}
