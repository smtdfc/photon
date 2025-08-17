package core

import (
	"sync"
)



type Gateway interface {
	Start() error
}

type HttpContext interface {
	Method() string
	Path() string
	Protocol() string

	Param(key string) string
	Query(key string) string
	QueryDefault(key, def string) string
	Header(key string) string
	Cookie(name string) string
	Body() []byte
	FormValue(key string) string
	FormFile(name string) ([]byte, error)

	Status(code int) HttpContext
	SetHeader(key, value string) HttpContext
	SetCookie(name, value string, options ...any) HttpContext

	Text(code int, data string) HttpContext
	JSON(code int, data any) HttpContext
	HTML(code int, html string) HttpContext
	Blob(code int, contentType string, data []byte) HttpContext
	File(code int, filepath string) HttpContext

	Next() HttpContext
	Abort() HttpContext
	IsAborted() bool

	Set(key string, value any) HttpContext
	Get(key string) any
	MustGet(key string) any
}

type HttpHandler func(HttpContext)

type HttpScope interface {
	Use(mw ...HttpHandler)
	Get(path string, handlers ...HttpHandler)
	Post(path string, handlers ...HttpHandler)
	Put(path string, handlers ...HttpHandler)
	Delete(path string, handlers ...HttpHandler)
	Head(path string, handlers ...HttpHandler)
}

type HttpGateway interface {
	Gateway
	Use(mw ...HttpHandler)
	CreateScope(module *Module, prefix string) HttpScope
}



type WsContext interface{}
type WsHandler func(WsContext)

type WsNamespace interface {
	GetAllRoom() []string
	CreateRoom(name string)
	EmitToRoom(event string, data map[string]any) error
	Emit(event string, data map[string]any) error
	On(event string, handlers ...WsHandler)
}

type WsGateway interface {
	Gateway
	CreateNamespace(module *Module, name string) WsGateway
	OnMessage(event string, handler any)
}



type GatewayManager struct {
	App  *App
	Http HttpGateway
	Ws   WsGateway
}

func (g *GatewayManager) SetGateway(name string, gateway Gateway) {
	switch name {
	case "http":
		g.Http = gateway.(HttpGateway)

	case "ws":
		g.Ws = gateway.(WsGateway)
	}
}

func (g *GatewayManager) StartAll() sync.WaitGroup {
	var wg sync.WaitGroup
	if g.Http != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			g.App.Logger.Info("Starting HTTP Gateway...")
			if err := g.Http.Start(); err != nil {
				g.App.Logger.Error("HTTP Gateway failed: " + err.Error())
			}
			g.App.Logger.Info("HTTP Gateway stopped")
		}()
	}

	if g.Ws != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			g.App.Logger.Info("Starting WS Gateway...")
			if err := g.Ws.Start(); err != nil {
				g.App.Logger.Error("WS Gateway failed: " + err.Error())
			}
			g.App.Logger.Info("WS Gateway stopped")
		}()
	}
	
	return wg
}


func ResolveGateway[T Gateway](module *Module, name string) T {
	var zero T
	var manager = module.App.gatewayManager

	switch name {
	case "http":
		if gw, ok := any(manager.Http).(T); ok {
			return gw
		}
	case "ws":
		if gw, ok := any(manager.Ws).(T); ok {
			return gw
		}
	}

	return zero
}
