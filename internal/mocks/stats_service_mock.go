package mocks

import (
	"github.com/client-side96/pi-monitor/internal/model"
	"github.com/client-side96/pi-monitor/internal/sub"
)

type MockStatsService struct {
	Client sub.StatsClient
	Stats  *model.StatsModel
}

func (m *MockStatsService) GetStats() *model.StatsModel {
	return m.Stats
}

func (m *MockStatsService) Connect() sub.StatsClient {
	return m.Client
}

func (m *MockStatsService) Disconnect(client sub.StatsClient) {}

func (m *MockStatsService) HandleStatsSubscripition(client sub.StatsClient) {}

func (m *MockStatsService) PublishToAllClients() {}
