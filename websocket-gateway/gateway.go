package websocketGateway

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/smtdfc/photon/core"
)

func generateClientID() string {
	return uuid.New().String()
}

type Client struct {
	Id   string
	Conn *websocket.Conn
}

type Members = []*Client

type Gateway struct {
	Listeners map[string]*Listener
	Upgrader  websocket.Upgrader
	Rooms     map[string]Members
	Clients   map[string]*Client
	config    Config
	mu        sync.RWMutex
}

type Message struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

func (g *Gateway) addClient(c *Client) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Clients[c.Id] = c
}

func (g *Gateway) removeClient(id string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.Clients, id)

	// remove client from all rooms
	for room, members := range g.Rooms {
		newMembers := make(Members, 0, len(members))
		for _, member := range members {
			if member.Id != id {
				newMembers = append(newMembers, member)
			}
		}
		g.Rooms[room] = newMembers
	}
}

func (g *Gateway) Broadcast(event string, data any) {
	msg := Message{Event: event, Data: data}
	raw, _ := json.Marshal(msg)

	g.mu.RLock()
	defer g.mu.RUnlock()
	for _, c := range g.Clients {
		_ = c.Conn.WriteMessage(websocket.TextMessage, raw)
	}
}

func (g *Gateway) dispatch(client *Client, event string, data any) {
	if g.Listeners[event] != nil {
		for _, handler := range g.Listeners[event].Handlers {
			logger := g.Listeners[event].Namespace.logger
			if logger != nil {
				d, _ := json.Marshal(data)
				logger.Info("(Websocket) " + event + ": " + string(d))
			}

			ctx := NewContext(
				g,
				g.Listeners[event].Namespace,
				client,
				event,
				data,
			)
			handler(ctx)
		}
	}
}

func (g *Gateway) handler(w http.ResponseWriter, r *http.Request) {
	conn, err := g.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	client := &Client{
		Id:   generateClientID(),
		Conn: conn,
	}
	g.addClient(client)

	log.Println("Client connected:", client.Id)

	defer func() {
		conn.Close()
		g.removeClient(client.Id)
		log.Println("Client disconnected:", client.Id)
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("Invalid JSON:", err)
			continue
		}

		g.dispatch(client, msg.Event, msg.Data)
	}
}

func (g *Gateway) Start() error {
	port := g.config.Port
	if port == "" {
		port = "8000"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", g.handler)

	fmt.Println("Websocket server is running http://localhost:" + port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Println("Error when starting Websocket server:", err)
		return err
	}
	return nil
}

func (g *Gateway) HasRoom(name string) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	_, ok := g.Rooms[name]
	return ok
}

func (g *Gateway) CreateRoom(name string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, ok := g.Rooms[name]; ok {
		return errors.New("room already exists")
	}
	g.Rooms[name] = []*Client{}
	return nil
}

func (g *Gateway) GetAllRoom() []string {
	g.mu.RLock()
	defer g.mu.RUnlock()

	rooms := make([]string, 0, len(g.Rooms))
	for name := range g.Rooms {
		rooms = append(rooms, name)
	}
	return rooms
}

func (g *Gateway) clientInRoom(room string, clientID string) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	members, ok := g.Rooms[room]
	if !ok {
		return false
	}
	for _, member := range members {
		if member.Id == clientID {
			return true
		}
	}
	return false
}

func (g *Gateway) emit(client *Client, event string, data any) error {
	msg := Message{Event: event, Data: data}
	raw, _ := json.Marshal(msg)

	g.mu.Lock()
	defer g.mu.Unlock()
	_ = client.Conn.WriteMessage(websocket.TextMessage, raw)
	return nil
}

func (g *Gateway) emitToRoom(sender *Client, room string, event string, data any) error {
	msg := Message{Event: event, Data: data}
	raw, _ := json.Marshal(msg)

	g.mu.Lock()
	defer g.mu.Unlock()

	members, ok := g.Rooms[room]
	if !ok {
		return errors.New("Client not joining room")
	}
	for _, c := range members {
		if c.Id == sender.Id {
			continue
		}
		_ = c.Conn.WriteMessage(websocket.TextMessage, raw)
	}

	return nil
}

func (g *Gateway) joinRoom(room string, client *Client) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	members, ok := g.Rooms[room]
	if !ok {
		return errors.New("room does not exist")
	}
	for _, member := range members {
		if member.Id == client.Id {
			return errors.New("client already in room")
		}
	}
	g.Rooms[room] = append(members, client)
	return nil
}

func (g *Gateway) leaveRoom(room string, clientID string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	members, ok := g.Rooms[room]
	if !ok {
		return errors.New("room not found")
	}

	newMembers := make(Members, 0, len(members))
	for _, member := range members {
		if member.Id != clientID {
			newMembers = append(newMembers, member)
		}
	}
	g.Rooms[room] = newMembers
	return nil
}

func (g *Gateway) CreateNamespace(module *core.Module, name string) core.WsNamespace {
	return NewNamespace(g, module, name)
}

func New(config Config) core.WsGateway {
	co := config.CheckOrigin
	if co == nil {
		co = func(r *http.Request) bool {
			return true
		}
	}

	return &Gateway{
		config:    config,
		Listeners: make(map[string]*Listener),
		Rooms:     make(map[string]Members),
		Clients:   make(map[string]*Client),
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     co,
		},
	}
}
