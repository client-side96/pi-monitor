package sub

import (
	"errors"
	"log"
)

type ClientConnectionStatus = string

const (
	Connected    ClientConnectionStatus = "connected"
	Subscribed   ClientConnectionStatus = "subscribed"
	Disconnected ClientConnectionStatus = "disconnected"
)

type client[T any] struct {
	ID      int
	Channel chan *T
	Status  ClientConnectionStatus
}

type Subscriber[T any] interface {
	HandleSubscription(client client[T])
	Subscribe(client client[T]) (int, error)
	Unsubscribe(clientID int) error
	GetLatestClientID() int
	SetLatestClientID(clientID int)
	GetClients() []client[T]
	SetClients(client client[T])
}

type Subscription[T any] struct {
	latestClientID int
	clients        []client[T]
}

func newSubscription[T any]() *Subscription[T] {
	return &Subscription[T]{
		latestClientID: 0,
		clients:        []client[T]{},
	}
}

func (sub *Subscription[T]) handleSubscription(client client[T]) {
	status := client.Status
	if status == Connected {
		_, err := sub.subscribe(client)
		if err != nil {
			log.Println(err.Error())
		}
	}
	if status == Disconnected {
		err := sub.unsubscribe(client.ID)
		if err != nil {
			log.Println(err.Error())
		}
	}
}

func (sub *Subscription[T]) subscribe(client client[T]) (int, error) {
	if client.Status != Connected {
		return client.ID, errors.New("Subscription error: Not connected")
	}
	client.Status = Subscribed
	sub.clients = append(sub.clients, client)
	return client.ID, nil
}

func (sub *Subscription[T]) unsubscribe(clientID int) error {
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

func (sub *Subscription[T]) getLatestClientID() int {
	return sub.latestClientID
}

func (sub *Subscription[T]) setLatestClientID(clientID int) {
	sub.latestClientID = clientID
}

func (sub *Subscription[T]) getClients() []client[T] {
	return sub.clients
}

func (sub *Subscription[T]) setClients(client client[T]) {
	sub.clients = append(sub.clients, client)
}
