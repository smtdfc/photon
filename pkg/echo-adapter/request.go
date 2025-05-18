package echoAdapter

import(
  "github.com/labstack/echo/v4"
  "github.com/smtdfc/photon/pkg/base"
)

type EchoAdapterRequest struct{
  Context echo.Context
}

func(r *EchoAdapterRequest) GetHeader(name string) string {
  return r.Context.Request().Header.Get(name)
}

func(r *EchoAdapterRequest) GetCookie(name string) (string, error){
  _, err := r.Context.Cookie(name)
  return "",err
}

func(r *EchoAdapterRequest) GetIP() string {
  return r.Context.RealIP()
}

func WrapRequest(c echo.Context) photon.Request{
  return &EchoAdapterRequest{
    Context: c,
  }
}