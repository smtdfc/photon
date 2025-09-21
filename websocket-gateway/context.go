package websocketGateway

import (
	"github.com/smtdfc/photon/v2/core"
)

type Context struct {
	event     string
	data      any
	Namespace *Namespace
	Gateway   *Gateway
	Client    *Client
}

func (c *Context) CreateRoom(name string) error {
	return c.Gateway.CreateRoom(name)
}

func (ctx *Context) Emit(event string, data any) error {
	return ctx.Gateway.emit(ctx.Client, ctx.Namespace.name+":"+event, data)
}

func (ctx *Context) EmitToRoom(room, event string, data any) error {
	return ctx.Gateway.emitToRoom(ctx.Client, room, ctx.Namespace.name+":"+event, data)
}

func (c *Context) GetAllRoom() []string {
	return c.Gateway.GetAllRoom()
}

func (c *Context) Join(name string) error {
	return c.Gateway.joinRoom(name, c.Client)
}

func (c *Context) HasJoin(name string) bool {
	return c.Gateway.clientInRoom(name, c.Client.Id)
}

func (c *Context) Leave(name string) error {
	return c.Gateway.leaveRoom(name, c.Client.Id)
}

func (c *Context) HasRoom(name string) bool {
	return c.Gateway.HasRoom(name)
}

func (c *Context) GetClientID() string {
	return c.Client.Id
}

func (c *Context) GetData() any {
	return c.data
}

func (c *Context) GetEvent() string {
	return c.event
}

func NewContext(gateway *Gateway, ns *Namespace, client *Client, event string, data any) core.WsContext {
	return &Context{
		event:     event,
		data:      data,
		Gateway:   gateway,
		Namespace: ns,
		Client:    client,
	}
}
