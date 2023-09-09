package service

import (
	"log"

	"github.com/client-side96/pi-monitor/internal/model"
	"github.com/client-side96/pi-monitor/internal/repository"
	"github.com/client-side96/pi-monitor/internal/sub"
)

type IStatsService interface {
	GetStats() *model.StatsModel
	Connect() sub.StatsClient
	Disconnect(client sub.StatsClient)
	HandleStatsSubscripition(client sub.StatsClient)
	PublishToAllClients()
}

type StatsService struct {
	sub     sub.StatsSubscriber
	repo    repository.IStatsRepository
	Channel chan sub.StatsClient
}

func NewStatsService(subscriber sub.StatsSubscriber, repo repository.IStatsRepository) *StatsService {
	return &StatsService{
		sub:     subscriber,
		repo:    repo,
		Channel: make(chan sub.StatsClient),
	}
}

func (s *StatsService) GetStats() *model.StatsModel {
	return &model.StatsModel{
		Temperature: s.repo.ExecuteTemperatureScript(),
		CPULoad:     s.repo.ExecuteCPUScript(),
		MemoryUsage: s.repo.ExecuteMemoryScript(),
	}
}

func (s *StatsService) Connect() sub.StatsClient {
	log.Println("Establishing connection...")
	s.sub.SetLatestClientID(s.sub.GetLatestClientID() + 1)
	newClient := sub.StatsClient{
		ID:      s.sub.GetLatestClientID(),
		Channel: make(chan *model.StatsModel),
		Status:  sub.Connected,
	}
	s.Channel <- newClient
	log.Printf("Established connection: %d", s.sub.GetLatestClientID())
	return newClient
}

func (s *StatsService) Disconnect(client sub.StatsClient) {
	log.Println("Closing connection...")
	for _, c := range s.sub.GetClients() {
		if client.ID == c.ID {
			c.Status = sub.Disconnected
			s.Channel <- c
			log.Printf("Closed connection: %d", client.ID)
		}
	}
}

func (s *StatsService) HandleStatsSubscripition(client sub.StatsClient) {
	s.sub.HandleSubscription(client)
}

func (s *StatsService) PublishToAllClients() {
	newStats := s.GetStats()
	for _, client := range s.sub.GetClients() {
		client.Channel <- newStats
	}
}
