package photon

type RouteHandler func(req Request, res Response)

type Route struct {
  Module *Module
  Path string
  Method string
  Handler RouteHandler
}