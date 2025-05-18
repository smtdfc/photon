package photon

type RouteHandler func(req Request, res Response)
type SocketEventHandler func(data []byte) error

