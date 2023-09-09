package stats

import (
	"log"

	"github.com/client-side96/pi-monitor/internal/util"
)

type OSCommunicator interface {
	ExecuteScript(script string) string
}

type StatsService struct {
	os      OSCommunicator
	sub     Subscriber
	Channel chan Client
}

func NewStatsService(os OSCommunicator) *StatsService {
	return &StatsService{
		os:      os,
		sub:     NewStatsSubscriber(),
		Channel: make(chan Client),
	}
}

func (s *StatsService) EstablishConnection() Client {
	log.Println("Establishing connection...")
	s.sub.SetLatestClientID(s.sub.GetLatestClientID() + 1)
	newClient := Client{
		ID:      s.sub.GetLatestClientID(),
		Channel: make(chan *Stats),
		Status:  Connected,
	}
	s.Channel <- newClient
	log.Printf("Established connection: %d", s.sub.GetLatestClientID())
	return newClient
}

func (s *StatsService) CloseConnection(connection Client) {
	log.Println("Closing connection...")
	for _, c := range s.sub.GetClients() {
		if connection.ID == c.ID {
			c.Status = Disconnected
			s.Channel <- c
			log.Printf("Closed connection: %d", connection.ID)
		}
	}
}

func (s *StatsService) HandleClientSubscription(client Client) {
	status := client.Status
	if status == Connected {
		_, err := s.sub.Subscribe(client)
		if err != nil {
			log.Printf(err.Error())
		}
	}
	if status == Disconnected {
		err := s.sub.Unsubscribe(client.ID)
		if err != nil {
			log.Printf(err.Error())
		}
	}
}

func (s *StatsService) PublishStats() {
	newStats := s.GetStats()
	for _, client := range s.sub.GetClients() {
		client.Channel <- newStats
	}
}

func (s *StatsService) GetStats() *Stats {
	return &Stats{
		Temperature: s.getTemperature(),
		CPULoad:     s.getCPULoad(),
		MemoryUsage: s.getMemoryUsage(),
	}
}

func (s *StatsService) getTemperature() float64 {
	var tempScript = "temperature.sh"
	return util.ToFloat(s.os.ExecuteScript(tempScript))
}

func (s *StatsService) getCPULoad() float64 {
	var cpuScript = "cpu.sh"
	return util.ToFloat(s.os.ExecuteScript(cpuScript))
}

func (s *StatsService) getMemoryUsage() float64 {
	var memoryScript = "memory.sh"
	return util.ToFloat(s.os.ExecuteScript(memoryScript))
}
