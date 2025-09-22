package chat

import (
	"github.com/VitalyCone/websocket-messenger/internal/app/model"
	"github.com/VitalyCone/websocket-messenger/internal/app/service"
	"github.com/sirupsen/logrus"
)

// Hub is a struct that holds all the clients and the messages that are sent to them
type Hub struct {
	service *service.Service
	// Registered clients.
	clients map[string]map[*Client]bool
	//Unregistered clients.
	unregister chan *Client
	// Register requests from the clients.
	register chan *Client
	// Inbound messages from the clients.
	broadcast chan model.MessageWS
}

func NewHub(service *service.Service) *Hub {
	return &Hub{
		service: service,
		clients:    make(map[string]map[*Client]bool),
		unregister: make(chan *Client),
		register:   make(chan *Client),
		broadcast:  make(chan model.MessageWS),
	}
}

// Core function to run the hub
func (h *Hub) Run() {
	for {
		select {
		// Register a client.
		case client := <-h.register:
			h.RegisterNewClient(client)
			// Unregister a client.
		case client := <-h.unregister:
			h.RemoveClient(client)
			// Broadcast a message to all clients.
		case message := <-h.broadcast:
			//Check if the message is a type of "message"
			h.HandleMessage(message)

		}
	}
}

// function check if room exists and if not create it and add client to it
func (h *Hub) RegisterNewClient(client *Client) {
	connections := h.clients[client.Username]
	if connections == nil {
		connections = make(map[*Client]bool)
		h.clients[client.Username] = connections
	}
	h.clients[client.Username][client] = true

	logrus.Println("Size of clients: ", len(h.clients[client.Username]))
}

// function to remvoe client from room
func (h *Hub) RemoveClient(client *Client) {
	if _, ok := h.clients[client.Username]; ok {
		delete(h.clients[client.Username], client)
		close(client.send)
		logrus.Println("Removed client")
	}
}

// function to handle message based on type of message
func (h *Hub) HandleMessage(message model.MessageWS) {
	modelChat, err := h.service.Chat.GetChat(message.ChatID)
	if err != nil {
		logrus.Errorf("failed to get chat for %d : %v", message.ChatID, err)
		return
	}

	modelMessage := message.ToModel()
	err = h.service.Message.CreateMessage(&modelMessage)
	if err != nil {
		logrus.Errorf("failed to create message for %d : %v", message.ChatID, err)
	}
	
	for _, user := range modelChat.Users {
		if user.Username != message.Sender {
			message.Recipients = append(message.Recipients, user.Username)
		}
	}
	//Check if the message is a type of "message"
	if message.Type == "message" {
		for _, recipient := range message.Recipients   {
			clients := h.clients[recipient]
			for client := range clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients[message.Sender], client)
				}
			}
		}
	}

	//Check if the message is a type of "notification"
	if message.Type == "notification" {
		for _, recipient := range message.Recipients {
			logrus.Println("Notification: ", message.Content)
			clients := h.clients[recipient]
			for client := range clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients[recipient], client)
				}
			}
		}
	}

}

// if message.Type == "message" {
// 	for _, recipient := range message.Recipients{
// 		var clients map[*Client]bool
// 		clientsByte, _ := h.clients.Get(context.Background(), recipient).Result()
// 		json.Unmarshal([]byte(clientsByte), &clients)
// 		for client := range clients {
// 			select {
// 			case client.send <- message:
// 			default:
// 				close(client.send)
// 				// delete(h.clients[message.ID], client)
// 			}
// 		}
// 	}
// }

// //Check if the message is a type of "notification"
// if message.Type == "notification" {
// 	logrus.Println("Notification: ", message.Content)
// 	for _, recipient := range message.Recipients{
// 		var clients map[*Client]bool
// 		clientsByte, _ := h.clients.Get(context.Background(), recipient).Result()
// 		json.Unmarshal([]byte(clientsByte), &clients)
// 		for client := range clients {
// 			select {
// 			case client.send <- message:
// 			default:
// 				close(client.send)
// 				// delete(h.clients[recipient], client)
// 			}
// 		}
// 	}
// }
