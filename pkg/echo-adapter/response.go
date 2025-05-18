package echoAdapter

import (
	"net/http"
	"time"
	"github.com/labstack/echo/v4"
	"github.com/smtdfc/photon/pkg/base"
)

type EchoAdapterResponse struct {
	echoCtx echo.Context
}

func (r *EchoAdapterResponse) SetHeader(name string, value string) {
	r.echoCtx.Response().Header().Set(name, value)
}

func (r *EchoAdapterResponse) SetStatus(code int) {
	r.echoCtx.Response().WriteHeader(code)
}

func (r *EchoAdapterResponse) SetCookie(name string, value string, path string, maxAge int) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = value
	cookie.Path = path
	cookie.MaxAge = maxAge
	cookie.Expires = time.Now().Add(time.Duration(maxAge) * time.Second)
	r.echoCtx.SetCookie(cookie)
}

func (r *EchoAdapterResponse) JSON(code int, obj interface{}) error {
	return r.echoCtx.JSON(code, obj)
}

func (r *EchoAdapterResponse) Text(code int, msg string) error {
	return r.echoCtx.String(code, msg)
}

func (r *EchoAdapterResponse) HTML(code int, html string) error {
	return r.echoCtx.HTML(code, html)
}

func (r *EchoAdapterResponse) Blob(code int, contentType string, data []byte) error {
	return r.echoCtx.Blob(code, contentType, data)
}

func (r *EchoAdapterResponse) Redirect(code int, url string) error {
	return r.echoCtx.Redirect(code, url)
}

func (r *EchoAdapterResponse) SendFile(filePath string) error {
	return r.echoCtx.File(filePath)
}

func (r *EchoAdapterResponse) Writer() any {
	return r.echoCtx.Response().Writer
}


func WrapResponse(c echo.Context) photon.Response {
	return &EchoAdapterResponse{
		echoCtx: c,
	}
}