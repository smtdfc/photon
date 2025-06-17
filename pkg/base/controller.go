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
	}
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
