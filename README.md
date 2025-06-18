# Photon

**Photon** is a lightweight and extensible backend framework for Go. Built with a modular, adapter-based architecture, Photon enables developers to create robust HTTP and WebSocket applications with clarity and flexibility.

---

## ✨ Features

* **Modular Adapters**: Seamlessly switch between HTTP and WebSocket backends (e.g., Echo, Melody).
* **Plug-and-Play Architecture**: Easy integration of custom adapters.
* **Minimalistic Core**: No hidden magic — full transparency and maintainability.

---

## 🚀 Installation

```bash
go get github.com/smtdfc/photon
go install github.com/smtdfc/photon-cli
```

---

## ⚙️ Getting Started

### 1. Create a new Echo-based Photon app

```bash
go mod init example.com/hello
photon gen app hello
```

### 2. Generate a module

```bash
photon gen module hello
```

### 3. Define your module routes

In `internal/hello/init.go`:

```go
package hello

func (m *HelloModule) InitRoute() {
    m.Module.Route("GET", "/hello", m.Handler.Hello)
}
```

In `internal/hello/handler.go`:

```go
package hello

import "github.com/smtdfc/photon"

type HelloModuleHandler struct{}

func (h *HelloModuleHandler) Hello(req photon.Request, res photon.Response) {
    res.JSON(200, map[string]any{
        "message": "Hello World 🚀",
    })
}
```

Register the module in your app:

```go
package app

import (
    "github.com/smtdfc/photon"
    "example.com/hello/internal/hello"
)

func InitModule(app *photon.App) {
    hello.Init(app)
}
```

### 4. Run the app

```bash
go run .
```

Access it at:
[http://127.0.0.1:3000/hello](http://127.0.0.1:3000/hello)

---

## 🔌 Custom Adapter Development

Photon adapters must implement the following interfaces:

### `BaseHTTPAdapter`

```go
type BaseHTTPAdapter interface {
    GetName() string
    Init() error
    Start() error
    Listen(port string) error
    UseSocket(path string, socketAdapter BaseSocketAdapter) error
    Route(method string, path string, handler RouteHandler)
}
```

### `BaseSocketAdapter`

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

## 📆 Ecosystem

* 🔗 [Echo Adapter](https://github.com/smtdfc/photon-echo-adapter)
* 🔗 [Melody Adapter](https://github.com/olahol/photon-melody-adapter)
* 🛠️ [Photon CLI](https://github.com/olahol/photon-cli)

---

## 📜 License

Photon is licensed under the MIT License.
© 2025 [smtdfc](https://github.com/smtdfc)

