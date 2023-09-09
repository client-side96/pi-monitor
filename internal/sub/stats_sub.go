package sub

import "github.com/client-side96/pi-monitor/internal/model"

type StatsSubscriber interface {
	Subscriber[model.StatsModel]
}

type StatsSubscription struct {
	*Subscription[model.StatsModel]
}

type StatsClient = client[model.StatsModel]

func NewStatsSubscription() *StatsSubscription {
	return &StatsSubscription{
		Subscription: newSubscription[model.StatsModel](),
	}
}

func (sub *StatsSubscription) HandleSubscription(client StatsClient) {
	sub.handleSubscription(client)
}

func (sub *StatsSubscription) Subscribe(client StatsClient) (int, error) {
	return sub.subscribe(client)
}

func (sub *StatsSubscription) Unsubscribe(clientID int) error {
	return sub.unsubscribe(clientID)
}

func (sub *StatsSubscription) GetLatestClientID() int {
	return sub.getLatestClientID()
}

func (sub *StatsSubscription) SetLatestClientID(clientID int) {
	sub.setLatestClientID(clientID)
}

func (sub *StatsSubscription) GetClients() []StatsClient {
	return sub.getClients()
}

func (sub *StatsSubscription) SetClients(client StatsClient) {
	sub.setClients(client)
}
