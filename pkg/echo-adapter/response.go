package echoAdapter

import(
  "github.com/labstack/echo/v4"
  "github.com/smtdfc/photon/pkg/base"
)


type EchoAdapterResponse struct{
  Context echo.Context
}

func WrapResponse(c echo.Context) photon.Response{
  return &EchoAdapterResponse{
    Context: c,
  }
}