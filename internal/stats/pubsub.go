package stats

import (
	"errors"
)

type ClientConnectionStatus = string

const (
	Connected    ClientConnectionStatus = "connected"
	Subscribed   ClientConnectionStatus = "subscribed"
	Disconnected ClientConnectionStatus = "disconnected"
)

type Client struct {
	ID      int
	Channel chan *Stats
	Status  ClientConnectionStatus
}

type Subscriber interface {
	Subscribe(client Client) (int, error)
	Unsubscribe(clientID int) error
	GetLatestClientID() int
	SetLatestClientID(clientID int)
	GetClients() []Client
}

type StatsSubscriber struct {
	latestClientID int
	clients        []Client
}

func NewStatsSubscriber() *StatsSubscriber {
	return &StatsSubscriber{
		latestClientID: 0,
		clients:        []Client{},
	}
}

func (sub *StatsSubscriber) Subscribe(client Client) (int, error) {
	if client.Status != Connected {
		return client.ID, errors.New("Subscription error: Not connected")
	}
	client.Status = Subscribed
	sub.clients = append(sub.clients, client)
	return client.ID, nil
}

func (sub *StatsSubscriber) Unsubscribe(clientID int) error {
	index := -1
	for i, client := range sub.clients {
		if client.ID == clientID {
			index = i
		}
	}
	if index == -1 {
		return errors.New("Subscription error: Connection for client does not exist")
	}
	// remove the client from the list
	sub.clients[index] = sub.clients[len(sub.clients)-1]
	sub.clients = sub.clients[:len(sub.clients)-1]
	return nil
}

func (sub *StatsSubscriber) GetLatestClientID() int {
	return sub.latestClientID
}

func (sub *StatsSubscriber) SetLatestClientID(clientID int) {
	sub.latestClientID = clientID
}

func (sub *StatsSubscriber) GetClients() []Client {
	return sub.clients
}
