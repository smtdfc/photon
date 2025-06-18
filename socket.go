package photon

type SocketSession struct {
	ClientID string
	Data     map[string]any
	Instance any
}

type SocketEventMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}
