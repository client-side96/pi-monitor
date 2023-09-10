package mocks

import (
	"github.com/client-side96/pi-monitor/internal/model"
	"github.com/client-side96/pi-monitor/internal/sub"
)

type Called = string

var (
	DisconnectCalled              Called = "DisconnectCalled"
	HandleStatsSubscriptionCalled Called = "HandleStatsSubscriptionCalled"
	PublishToAllClientsCalled     Called = "PublishToAllClientsCalled"
)

type MockStatsService struct {
	Client sub.StatsClient
	Stats  *model.StatsModel
	Called Called
}

func (m *MockStatsService) GetStats() *model.StatsModel {
	return m.Stats
}

func (m *MockStatsService) Connect() sub.StatsClient {
	return m.Client
}

func (m *MockStatsService) Disconnect(_ sub.StatsClient) {
	m.Called = DisconnectCalled
}

func (m *MockStatsService) HandleStatsSubscripition(_ sub.StatsClient) {
	m.Called = HandleStatsSubscriptionCalled
}

func (m *MockStatsService) PublishToAllClients() {
	m.Called = PublishToAllClientsCalled
}
