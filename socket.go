package photon

import "errors"

type SocketSession struct {
	ClientID string
	Data     map[string]any
	Instance any
}

type SocketEventMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}


type SocketRoom struct {
	Adapter BaseSocketAdapter
	RoomID  string
	Data    map[string]any
	Clients map[string]*SocketSession
}

func (r *SocketRoom) HasClient(clientID string) bool {
	_, exists := r.Clients[clientID]
	return exists
}

func (r *SocketRoom) Emit(client *SocketSession, msg *SocketEventMessage) error {
	if !r.HasClient(client.ClientID) {
		return errors.New("client not found")
	}
	return r.Adapter.Emit(client, msg)
}



func (r *SocketRoom) EmitAll(msg *SocketEventMessage) error {
	var errs []error
	for _, client := range r.Clients {
		if err := r.Adapter.Emit(client, msg); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.New("some clients failed to receive the message")
	}
	return nil
}