package photon

type RouteHandler func(req Request, res Response)
type SocketEventHandler func(client *SocketSession, msg *SocketEventMessage) error
