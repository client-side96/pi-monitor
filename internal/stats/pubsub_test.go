package stats

import "testing"

func TestStats_PubsubSubscribe(t *testing.T) {
	client := Client{
		ID:      1,
		Channel: make(chan *Stats),
		Status:  Connected,
	}
	sub := NewStatsSubscriber()

	id, err := sub.Subscribe(client)

	if sub.clients[0].ID != client.ID {
		t.Errorf("Client %d should be added to list upon subscription", client.ID)
	}
	if sub.clients[0].Status != Subscribed {
		t.Errorf("Client connection status should be 'Subscribed', but received %s", client.Status)
	}
	if id != client.ID {
		t.Errorf("Expected client id %d, but got %d", client.ID, id)
	}
	if err != nil {
		t.Errorf("Unexpected error was thrown: %s", err.Error())
	}
}

func TestStats_PubsubSubscribeAlreadySubscribed(t *testing.T) {
	client := Client{
		ID:      1,
		Channel: make(chan *Stats),
		Status:  Subscribed,
	}
	sub := NewStatsSubscriber()

	_, err := sub.Subscribe(client)

	if err == nil {
		t.Errorf("Failed: Expected subscription error")
	}
}

func TestStats_PubsubUnsubscribe(t *testing.T) {
	client := Client{
		ID:      1,
		Channel: make(chan *Stats),
		Status:  Subscribed,
	}
	sub := NewStatsSubscriber()
	sub.clients = []Client{client}

	err := sub.Unsubscribe(client.ID)

	if err != nil {
		t.Errorf("Unexpected error was thrown: %s", err.Error())
	}
	if len(sub.clients) > 0 {
		t.Errorf("Expected client to be removed from clients")
	}
}

func TestStats_PubsubUnsubscribeNoConnection(t *testing.T) {
	client := Client{
		ID:      1,
		Channel: make(chan *Stats),
		Status:  Subscribed,
	}
	sub := NewStatsSubscriber()

	err := sub.Unsubscribe(client.ID)

	if err == nil {
		t.Errorf("Failed: Expected subscription error")
	}
}
