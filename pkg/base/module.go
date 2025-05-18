package photon

type Module struct{
  Name string
  Prefix string
  App *App
  Adapter BaseAdapter
  Routes map[string] *Route
}


func (m *Module) Route(method string,path string, handler RouteHandler) *Route{
  route := &Route{
    Module:m,
    Path:path,
    Method:method,
    Handler:handler,
  }
  
  m.Routes[path] = route
  m.Adapter.Route(method,path,handler)
  return route
}


func NewModule(name string, app *App) *Module{
  return &Module{
    Name:name,
    Prefix:"",
    App:app,
    Adapter: app.Adapter,
    Routes: make(map[string] *Route),
  }
}
