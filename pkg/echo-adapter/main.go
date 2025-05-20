package echoAdapter

import(
  "fmt"
  "errors"
  "github.com/labstack/echo/v4"
  "github.com/smtdfc/photon/pkg/base"
)

type EchoAdapter struct{
  App *photon.App
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

func (a *EchoAdapter) GetInstance() any{
  return a.Instance
}


func (a *EchoAdapter) Init() error{
  return nil
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

func (a *EchoAdapter) SetApp(app *photon.App) error{
  a.App = app
  return nil
}

func (e *EchoAdapter) UseSocket(path string) error {
  socketAdapter:= e.App.SocketAdapter
  
	if socketAdapter == nil {
		return errors.New("socket adapter is nil")
	}

	e.Instance.GET(path, func(c echo.Context) error {
		handler := socketAdapter.HTTPHandler()
		handler(c.Response(), c.Request())
		return nil
	})

	return nil
}


func Init() (photon.BaseAdapter,*echo.Echo){
  ins := echo.New()
  return &EchoAdapter{
    Name: "PhotonEchoAdapter",
    Instance: ins,
  }, ins
}