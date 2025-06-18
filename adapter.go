package photon

import (
	"net/http"
)

type BaseHTTPAdapter interface {
	GetInstance() any
	SetApp(app *App) error
	GetName() string
	Init() error
	Start() error
	Listen(port string) error
	UseSocket(path string) error
	Route(method string, path string, handler RouteHandler)
}

type BaseSocketAdapter interface {
	GetName() string
	Init() error
	Start() error
	On(event string, handler SocketEventHandler)
	Emit(client *SocketSession, message *SocketEventMessage) error
	Stop() error
	HTTPHandler() func(http.ResponseWriter, *http.Request)
}

type AdapterManager struct {
	App           *App
	HttpAdapter   BaseHTTPAdapter
	SocketAdapter BaseSocketAdapter
}

func (m *AdapterManager) UseHttpAdapter(adapter BaseHTTPAdapter) {
	m.HttpAdapter = adapter
	adapter.SetApp(m.App)
}

func (m *AdapterManager) UseSocketAdapter(adapter BaseSocketAdapter) {
	m.SocketAdapter = adapter
}

func (m *AdapterManager) EsureAdapter(adapterType string) {
	if adapterType == "http" && m.HttpAdapter == nil {
		panic("Cannot find Http Adapter !")
	}

	if adapterType == "socket" && m.SocketAdapter == nil {
		panic("Cannot find Socket Adapter !")
	}
}
