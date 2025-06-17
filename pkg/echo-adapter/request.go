package echoAdapter

import (
	"github.com/labstack/echo/v4"
	"github.com/smtdfc/photon/pkg/base"
	"io"
	"log"
)

type EchoAdapterRequest struct {
	echoCtx echo.Context
}

func (r *EchoAdapterRequest) GetHeader(name string) string {
	return r.echoCtx.Request().Header.Get(name)
}

func (r *EchoAdapterRequest) GetCookie(name string) (string, error) {
	cookie, err := r.echoCtx.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (r *EchoAdapterRequest) GetIP() string {
	return r.echoCtx.RealIP()
}

func (r *EchoAdapterRequest) GetQuery(key string) string {
	return r.echoCtx.QueryParam(key)
}

func (r *EchoAdapterRequest) GetParam(key string) string {
	return r.echoCtx.Param(key)
}

func (r *EchoAdapterRequest) GetFormValue(key string) string {
	return r.echoCtx.FormValue(key)
}

func (r *EchoAdapterRequest) GetBody() ([]byte, error) {
	body := r.echoCtx.Request().Body
	if body == nil {
		return nil, nil
	}
	defer body.Close()
	return io.ReadAll(body)
}

func (r *EchoAdapterRequest) Method() string {
	return r.echoCtx.Request().Method
}

func (r *EchoAdapterRequest) Path() string {
	return r.echoCtx.Path()
}

func (r *EchoAdapterRequest) Host() string {
	return r.echoCtx.Request().Host
}

func (r *EchoAdapterRequest) URL() string {
	return r.echoCtx.Request().URL.String()
}

func (r *EchoAdapterRequest) ContentType() string {
	return r.echoCtx.Request().Header.Get(echo.HeaderContentType)
}

func (r *EchoAdapterRequest) BindJSON(obj interface{}) error {
	return r.echoCtx.Bind(obj)
}

func (r *EchoAdapterRequest) BindForm(obj interface{}) error {
	return r.echoCtx.Bind(obj)
}

func (r *EchoAdapterRequest) BindQuery(obj interface{}) error {
	log.Println("This method is not implemented in Echo Adapter")
	return r.echoCtx.Bind(obj)
}

func (r *EchoAdapterRequest) FormFile(name string) (any, error) {
	return r.echoCtx.FormFile(name)
}

func (r *EchoAdapterRequest) Context() any {
	return r.echoCtx
}

func WrapRequest(c echo.Context) photon.Request {
	return &EchoAdapterRequest{
		echoCtx: c,
	}
}
