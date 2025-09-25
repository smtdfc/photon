package websocketGateway

import (
	"github.com/smtdfc/photon/core"
)

type Listener struct {
	Namespace *Namespace
	Handlers  []core.WsHandler
}
