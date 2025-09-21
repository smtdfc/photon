package websocketGateway

import (
	"github.com/smtdfc/photon/v2/core"
)

type Namespace struct {
	Gateway *Gateway
	Module  *core.Module
	name    string
	logger  *core.Logger
}

func (n *Namespace) SetLogger(logger *core.Logger) {
	n.logger = logger
}

func (n *Namespace) On(event string, handlers ...core.WsHandler) {
	n.Gateway.Listeners[n.name+":"+event] = &Listener{
		Namespace: n,
		Handlers:  handlers,
	}
}

func NewNamespace(gateway *Gateway, module *core.Module, name string) *Namespace {
	return &Namespace{
		Gateway: gateway,
		name:    name,
		Module:  module,
	}
}
