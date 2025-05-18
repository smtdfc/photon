package photon

type BaseAdapter interface{
  GetName() string
  Init()
  Start() error
  Listen(port string) error
  Route(method string, path string, handler RouteHandler)
}