package pubsub

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

const (
	PUBLISH     = "publish"
	SUBSCRIBE   = "subscribe"
	UNSUBSCRIBE = "unsubscribe"
)

type PubSub struct {
	Clients       []Client
	Subscriptions []Subscription
}

type Client struct {
	Id         string
	Connection *websocket.Conn
}

type Message struct {
	Action  string          `json:"action"`
	Topic   string          `json:"topic"`
	Message json.RawMessage `json:"message"`
}

type Subscription struct {
	Topic  string
	Client *Client
}

func (ps *PubSub) AddClient(client Client) *PubSub {
	ps.Clients = append(ps.Clients, client)
	payload := []byte("Hello Client ID:" + client.Id)
	client.Connection.WriteMessage(1, payload)
	return ps
}

func (ps *PubSub) RemoveClient(client Client) *PubSub {
	for index, sub := range ps.Subscriptions {
		if client.Id == sub.Client.Id {
			ps.Subscriptions = append(ps.Subscriptions[:index], ps.Subscriptions[index+1:]...)
		}
	}

	for index, c := range ps.Clients {
		if c.Id == client.Id {
			ps.Clients = append(ps.Clients[:index], ps.Clients[index+1:]...)
		}
	}

	return ps
}

func (ps *PubSub) GetSubscriptions(topic string, client *Client) []Subscription {
	var subscriptionList []Subscription
	for _, subscription := range ps.Subscriptions {
		if client != nil {
			if subscription.Client.Id == client.Id && subscription.Topic == topic {
				subscriptionList = append(subscriptionList, subscription)
			}
		} else {
			if subscription.Topic == topic {
				subscriptionList = append(subscriptionList, subscription)
			}
		}
	}
	return subscriptionList
}

func (ps *PubSub) Subscribe(client *Client, topic string) *PubSub {
	clientSubs := ps.GetSubscriptions(topic, client)
	if len(clientSubs) > 0 {
		return ps
	}
	newSubscription := Subscription{
		Topic:  topic,
		Client: client,
	}
	ps.Subscriptions = append(ps.Subscriptions, newSubscription)
	return ps
}

func (ps *PubSub) Publish(topic string, message []byte, excludeClient *Client) {
	subscriptions := ps.GetSubscriptions(topic, nil)
	for _, sub := range subscriptions {
		fmt.Printf("Sending to client id %s message is %s \n", sub.Client.Id, message)
		sub.Client.Send(message)
	}
}

func (client *Client) Send(message []byte) error {
	return client.Connection.WriteMessage(1, message)
}

func (ps *PubSub) Unsubscribe(client *Client, topic string) *PubSub {
	for index, sub := range ps.Subscriptions {
		if sub.Client.Id == client.Id && sub.Topic == topic {
			ps.Subscriptions = append(ps.Subscriptions[:index], ps.Subscriptions[index+1:]...)
		}
	}
	return ps
}

func (ps *PubSub) HandleReceiveMessage(client Client, messageType int, payload []byte) *PubSub {
	m := Message{}

	err := json.Unmarshal(payload, &m)
	if err != nil {
		fmt.Println("This is not a correct message payload")
		return ps
	}

	switch m.Action {
	case PUBLISH:
		fmt.Println("This is a publish message")
		ps.Publish(m.Topic, m.Message, nil)

	case SUBSCRIBE:
		ps.Subscribe(&client, m.Topic)
		fmt.Println("new subscriber to topic", m.Topic, len(ps.Subscriptions), client.Id)

	case UNSUBSCRIBE:
		fmt.Println("Client want to unsubscribe the topic", m.Topic, client.Id)
		ps.Unsubscribe(&client, m.Topic)

	default:
		break
	}
	return ps
}

var upgrader = websocket.Upgrader{}

func autoId() string {
	return uuid.Must(uuid.NewV4(), nil).String()
}

var ps = &PubSub{}

func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := Client{
		Id:         autoId(),
		Connection: conn,
	}

	ps.AddClient(client)
	fmt.Println("New Client is connected, total: ", len(ps.Clients))

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Something went wrong", err)

			ps.RemoveClient(client)
			log.Println("total clients and subscriptions ", len(ps.Clients), len(ps.Subscriptions))

			return
		}
		ps.HandleReceiveMessage(client, messageType, p)
	}
}

// func main() {
// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		http.ServeFile(w, r, "static")
// 	})
// 	http.HandleFunc("/ws", WebsocketHandler)
// 	fmt.Println("Server is running: http://localhost:3000")
// 	http.ListenAndServe(":3000", nil)
// }
