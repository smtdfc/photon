package photon


type Route struct {
  Module *Module
  Path string
  Method string
  Handler RouteHandler
}