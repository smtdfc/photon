# Photon

**Photon** is a lightweight and extensible backend framework for Go. Built with a modular, adapter-based architecture, Photon enables developers to create robust HTTP applications with clarity and flexibility.

---

## ✨ Features

- **Gateway**: Seamlessly switch between HTTP and WebSocket backends .
- **Plug-and-Play Architecture**: Easy integration of custom adapters.
- **Minimalistic Core**: No hidden magic — full transparency and maintainability.

---

## 🚀 Installation

```bash
go get github.com/smtdfc/photon
go get github.com/smtdfc/photon-http-gateway
go install github.com/smtdfc/photon/cli
```

---

## ⚙️ Getting Started

### 1. Create a new Echo-based Photon app

```bash
go mod init example.com/hello
photon/cli init hello_app
```

### 2. Generate a module

```bash
photon/cli gen module Hello
```

### 3. Define your module routes

In `modules/hello/http.go`:

```go
package hello

import "github.com/smtdfc/photon/core"
import "github.com/smtdfc/photon-http-gateway"
import "github.com/smtdfc/photon-http-gateway/middlewares/cors"

var (
  ROUTE_PREFIX = ""
)

func initRoute(){
  logger := self().Logger
  gateway := core.ResolveGateway[*httpGateway.Gateway](self(),"http")
  http:= gateway.CreateScope(self(),ROUTE_PREFIX)
  http.SetLogger(logger)
  http.Use(cors.Middleware(cors.Config{}))

  http.Get(
    "/hello",
    func(ctx core.HttpContext){
      ctx.Text(200,"Hello world ")
    },
  )
}

```

In `modules/hello/lifecycle.go`:

```go
package hello

// This function will be called when the module is initialized.
func onModuleInit(){
  initRoute() // Call when module start
}

//This function will be called when the module starts.
func onModuleStart(){}
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
photon/cli dev
```

Access it at:
[http://127.0.0.1:3000/hello](http://127.0.0.1:3000/hello)

## 📜 License

Photon is licensed under the MIT License.
© 2025 [smtdfc](https://github.com/smtdfc)
