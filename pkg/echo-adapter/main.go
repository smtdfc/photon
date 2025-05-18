package echoAdapter

import(
  "fmt"
  "github.com/labstack/echo/v4"
  "github.com/smtdfc/photon/pkg/base"
)

type EchoAdapter struct{
  Name string
  Instance *echo.Echo
}

func wrapHandler(handler photon.RouteHandler) echo.HandlerFunc{
  return func(c echo.Context) error{
    handler(
      WrapRequest(c),
      WrapResponse(c),
    )
    
    return nil
  }
}


func (a *EchoAdapter) Init(){
  
}

func (a *EchoAdapter) Listen(port string) error{
  return a.Instance.Start(":" + port)
}

func (a *EchoAdapter) GetName() string{
  return a.Name
}

func (a *EchoAdapter) Route(method string, path string, handler photon.RouteHandler ){
  switch(method){
    case "GET":
    a.Instance.GET(path,wrapHandler(handler))
  }
}


func (a *EchoAdapter) Start() error{
  fmt.Println("Starting Photon Echo Adapter....")
  return nil
}

func Init() photon.BaseAdapter{
  return &EchoAdapter{
    Name: "PhotonEchoAdapter",
    Instance: echo.New(),
  }
}