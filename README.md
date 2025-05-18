# Photon

**Photon** is a lightweight, extensible backend framework for Go, designed with an adapter-based architecture. It supports both HTTP and WebSocket protocols through pluggable adapters such as Echo and Melody, enabling developers to build robust backend applications with simplicity and flexibility.

---

## Key Features

- Modular HTTP and WebSocket adapters  
- Plug-and-play adapter architecture  
- Minimalistic core with zero magic, ensuring full transparency and maintainability  

---

## Installation

```bash
go get github.com/smtdfc/photon
```

---

## Getting Started

### 1. Initialize a basic Echo application

```bash
go mod init example.com/hello
photon gen app hello
```

### 2. Generate a module

```bash
photon gen module hello
```

### 3. Configure the Hello module

Define the routing logic by implementing the `InitRoute` method in the `hello` module:

```go
package hello

func (m *HelloModule) InitRoute() {
    m.Module.Route("GET", "/hello", m.Handler.Hello)
}
```

Implement the request handler within the `HelloModuleHandler` struct:

```go
package hello

import (
    "github.com/smtdfc/photon/pkg/base"
)

type HelloModuleHandler struct{}

func (h *HelloModuleHandler) Hello(req photon.Request, res photon.Response) {
    res.JSON(200, map[string]any{
        "message": "Hello World ....",
    })
}
```

Register the module in the application by invoking its initialization function:

```go
package app

import (
    "github.com/smtdfc/photon/pkg/base"
    "example.com/hello/internal/hello"
)

func InitModule(app *photon.App) {
    hello.Init(app)
}
```

### 4. Run the application

```bash
go run .
```

Access the endpoint in your browser:  
`http://127.127.127:3000/hello`

---

## Custom Adapters

You can implement custom adapters by conforming to the following interfaces:

### `BaseAdapter` Interface

```go
type BaseAdapter interface {
    GetName() string
    Init() error
    Start() error
    Listen(port string) error
    UseSocket(path string, socketAdapter BaseSocketAdapter) error
    Route(method string, path string, handler RouteHandler)
}
```

### `BaseSocketAdapter` Interface

```go
type BaseSocketAdapter interface {
    GetName() string
    Init() error
    Start() error
    On(event string, handler SocketEventHandler)
    Emit(event string, data []byte) error
    Stop() error
    HTTPHandler() func(http.ResponseWriter, *http.Request)
}
```

---

## License

Photon is licensed under the MIT License © 2025 [smtdfc](https://github.com/smtdfc)

---

## Acknowledgments

Photon leverages the following outstanding open-source libraries:

- [Echo](https://github.com/labstack/echo)  
- [Melody](https://github.com/olahol/melody)
