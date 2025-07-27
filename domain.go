package photon

import (
	//	"context"
	"io"
	"net/http"
)

type File struct {
	Filename string
	Size     int64
	MIME     string
	Content  []byte
	Reader   io.Reader
}

type CookieOption struct {
	Path     string
	Domain   string
	MaxAge   int
	Secure   bool
	HttpOnly bool
	SameSite http.SameSite
}

type Request interface {
	GetHeader(name string) string
	GetAllHeader() (map[string]string, error)
	/*
	   GetMethod() string
	   GetPath() string
	   GetFullURL() string

	   GetBody() (map[string]any, error)
	   GetRawBody() ([]byte, error)
	   BindBody(target any) error

	   GetQuery(name string) string
	   GetAllQuery() map[string]string

	   GetParam(name string) string
	   GetAllParam() map[string]string

	   GetCookie(name string) (string, error)

	   GetForm(name string) string
	   GetAllForm() (map[string]string, error)

	   GetFile(name string) (File, error)
	   GetAllFiles() (map[string][]File, error)

	   GetContext() context.Context
	   SetContext(ctx context.Context)
	*/
}

type Response interface {
	SetHeader(name, value string) Response
	SetContentType(ct string) Response
	SetCookie(name, value string, opts ...CookieOption) Response

	Code(code int) Response

	JSON(data any) error
	Text(data string) error
	HTML(html string) error
	File(filepath string) error
	Stream(reader io.Reader) error

	Redirect(status int, location string) error

	RawWriter() http.ResponseWriter
}

type JSON = map[string]any
type Context interface {
	Req() Request
	// Res() Response
	// Set(key string, val any)
	// Get(key string) (any, bool)
	Params()map[string]string
	reset()
	isNext() bool
	Next() error
}

type HttpHandler func(ctx Context)

type BaseHttpAdapter interface {
	AddRoute(method string, path string, handlers ...HttpHandler)
	Start(port string) error
}

type HttpAdapterConfig struct {
	App       *App
	IsLogging bool
}
