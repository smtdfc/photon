# Photon

**Photon** is a lightweight and extensible backend framework for Go. Built with a modular, adapter-based architecture, Photon enables developers to create robust HTTP applications with clarity and flexibility.

---

## ✨ Features

* **Modular Adapters**: Seamlessly switch between HTTP and WebSocket backends .
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
photon-cli init hello_app
```

### 2. Generate a module

```bash
photon-cli gen module Hello
```

### 3. Define your module routes

In `modules/hello/routes.go`:

```go
package hello

func InitRoutes() {
    self().Router.Get("/hello", handleGet)
}
```

In `modules/hello/handlers.go`:

```go
package hello

import "github.com/smtdfc/photon/http_adapter"

func handleGet(ctx http_adapter.Context) {
    ctx.Res().JSON(http_adapter.JSON{
        "message": "Hello World 🚀",
    })
}
```

Register the module in your app:

```go
package app

import (
    "github.com/smtdfc/photon/core"
    "example.com/modules/hello"
)

func InitModule(app *core.App) {
    hello.Init(app)
}
```

### 4. Run the app

```bash
photon-cli dev
```

Access it at:
[http://127.0.0.1:3000/hello](http://127.0.0.1:3000/hello)


## 📜 License

Photon is licensed under the MIT License.
© 2025 [smtdfc](https://github.com/smtdfc)

