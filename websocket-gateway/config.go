package websocketGateway

import "net/http"


type Config struct {
	Port string
	CheckOrigin func(r *http.Request) bool
}
