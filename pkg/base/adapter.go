package photon

import(
	"net/http"
)


type BaseAdapter interface {
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
	Emit(event string, data []byte) error
	Stop() error
	HTTPHandler() func(http.ResponseWriter, *http.Request)
}