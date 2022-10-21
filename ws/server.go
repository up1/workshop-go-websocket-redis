package demo

import (
	"encoding/json"
	"fmt"
	"log"
)

type Message struct {
	Data []byte
	Room string
}

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return json
}

type subscription struct {
	conn *connection
	room string
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type Hub struct {
	// Registered connections.
	Rooms map[string]map[*connection]bool

	// Inbound messages from the connections.
	Broadcast chan Message

	// Register requests from the connections.
	Register chan subscription

	// Unregister requests from connections.
	Unregister chan subscription
}

var H = Hub{
	Broadcast:  make(chan Message),
	Register:   make(chan subscription),
	Unregister: make(chan subscription),
	Rooms:      make(map[string]map[*connection]bool),
}

func (h *Hub) Run() {
	go h.subscribeFromRedis()
	for {
		select {
		case s := <-h.Register:
			connections := h.Rooms[s.room]
			if connections == nil {
				connections = make(map[*connection]bool)
				h.Rooms[s.room] = connections
			}
			h.Rooms[s.room][s.conn] = true
		case s := <-h.Unregister:
			connections := h.Rooms[s.room]
			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					delete(connections, s.conn)
					close(s.conn.send)
					if len(connections) == 0 {
						delete(h.Rooms, s.room)
					}
				}
			}
		case m := <-h.Broadcast:
			connections := h.Rooms[m.Room]
			h.publishToRedis(m.encode())

			for c := range connections {
				select {
				case c.send <- m.Data:
				default:
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.Rooms, m.Room)
					}
				}
			}
		}
	}
}

func (h *Hub) subscribeFromRedis() {
	pubsub := redisClient.Subscribe(ctx, "demo-channel")

	ch := pubsub.Channel()

	for msg := range ch {
		var message Message
		if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
			log.Printf("Error on unmarshal JSON message %s", err)
			return
		}
		fmt.Println("Message from redis: ", message)
		h.broadcastToClients(message)
	}
}

func (h *Hub) broadcastToClients(message Message) {
	for client := range h.Rooms[message.Room] {
		client.send <- message.Data
	}
}

func (h *Hub) publishToRedis(message []byte) {
	err := redisClient.Publish(ctx, "demo-channel", message)
	if err != nil {
		log.Println(err)
	}
}


