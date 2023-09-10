package service_test

import (
	"io"
	"log"
	"os"
	"testing"

	"github.com/client-side96/pi-monitor/internal/model"
	"github.com/client-side96/pi-monitor/internal/service"
	"github.com/client-side96/pi-monitor/internal/sub"
	"github.com/client-side96/pi-monitor/internal/testutil"
	"github.com/client-side96/pi-monitor/internal/testutil/mocks"
)

var subscription *sub.StatsSubscription

var temp = 50.04
var cpu = 22.22
var mem = 15.24

func TestMain(m *testing.M) {
	log.SetOutput(io.Discard)
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

	testutil.AssertEqual(t, result.Temperature, temp)
	testutil.AssertEqual(t, result.CPULoad, cpu)
	testutil.AssertEqual(t, result.MemoryUsage, mem)
}

func TestService_Stats_Connect_SingleConnection(t *testing.T) {
	var result sub.StatsClient
	subscription = sub.NewStatsSubscription()
	repo := &mocks.MockStatsRepository{}
	s := service.NewStatsService(subscription, repo)

	go s.Connect()

	result = <-s.Channel

	testutil.AssertEqual(t, 1, result.ID)
	testutil.AssertEqual(t, sub.Connected, result.Status)

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
		result = <-s.Channel

		testutil.AssertEqual(t, connectionCount, result.ID)
		testutil.AssertEqual(t, sub.Connected, result.Status)

		connectionCount++
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
		result = <-s.Channel
		if count == 0 {
			go func() {
				s.Disconnect(client)
			}()
		}
		count++
	}

	testutil.AssertEqual(t, sub.Disconnected, result.Status)
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

	result = <-client.Channel

	testutil.AssertEqual(t, result.Temperature, temp)
	testutil.AssertEqual(t, result.CPULoad, cpu)
	testutil.AssertEqual(t, result.MemoryUsage, mem)
}
