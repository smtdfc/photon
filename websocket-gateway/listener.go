package websocketGateway

import (
	"github.com/smtdfc/photon/v2/core"
)

type Listener struct {
	Namespace *Namespace
	Handlers  []core.WsHandler
}
