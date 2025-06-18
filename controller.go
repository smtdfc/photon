package photon

type HttpController struct {
	App    *App
	Module *Module
	Routes map[string]*Route
}

func InitHttpController(app *App, module *Module) *HttpController {
	return &HttpController{
		App:    app,
		Module: module,
		Routes: make(map[string]*Route),
	}
}

func (m *HttpController) RouteSocket(path string) error {
	m.App.Adapter.EsureAdapter("http")
	return m.App.Adapter.HttpAdapter.UseSocket(path)
}

func (m *HttpController) Route(method string, path string, handler RouteHandler) *Route {
	m.App.Adapter.EsureAdapter("http")

	route := &Route{
		Module:  m.Module,
		Path:    path,
		Method:  method,
		Handler: handler,
	}

	m.Routes[path] = route
	m.App.Adapter.HttpAdapter.Route(method, path, handler)
	return route
}

type SocketController struct {
	App           *App
	Module        *Module
	SocketAdapter BaseSocketAdapter
}

func InitSocketController(app *App, module *Module) *SocketController {
	return &SocketController{
		App:           app,
		Module:        module,
		SocketAdapter: app.Adapter.SocketAdapter,
	}
}

func (c *SocketController) On(event string, handler SocketEventHandler) {
	c.SocketAdapter.On(event, handler)
}
