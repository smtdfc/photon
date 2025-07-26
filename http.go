package photon

import (
	"net/http"
	"net/url"
)

type HttpAdapterMethodMap = map[string][]HttpHandler

type HttpAdapter struct {
	App    *App
	Logger *Logger
	Config *HttpAdapterConfig
	Routes map[string]HttpAdapterMethodMap
}

type HttpAdapterContext struct {
	Writer    http.ResponseWriter
	Request   *http.Request
	Params    map[string]string
	Query     url.Values
	nextState bool
	_req      *HttpAdapterRequest
}

func (c *HttpAdapterContext) Req() Request {
	if c._req == nil {
		c._req = &HttpAdapterRequest{
			Original: c.Request,
		}
	}
	return c._req
}

func (c *HttpAdapterContext) Next() error {
	c.nextState = true
	return nil
}

func (c *HttpAdapterContext) reset() {
	c.nextState = false
}

func (c *HttpAdapterContext) isNext() bool {
	return c.nextState
}

func CreateContext(w http.ResponseWriter, r *http.Request, match PathMatch) Context {
	return &HttpAdapterContext{
		Writer:    w,
		Request:   r,
		Params:    match.Params,
		Query:     match.Query,
		nextState: false,
	}
}

func (h *HttpAdapter) AddRoute(method string, path string, handlers ...HttpHandler) {
	if h.Routes[path] != nil {
		if h.Routes[path][method] != nil {
			h.Logger.Warn("Path handler overriding " + path)
		}
	} else {
		h.Routes[path] = make(HttpAdapterMethodMap)
	}

	h.Routes[path][method] = handlers
}

func (h *HttpAdapter) handleReq(w http.ResponseWriter, r *http.Request) {
	if h.Config.IsLogging {
		h.Logger.Info("Request: " + r.Method + " " + r.URL.String())
	}

	path := r.URL.String()
	method := r.Method

	for pattern, methodMap := range h.Routes {
		pathMatchRes := parsePathParams(path, pattern)
		if pathMatchRes.Match && methodMap[method] != nil {
			ctx := CreateContext(w, r, pathMatchRes)
			handlers := methodMap[method]

			for _, handler := range handlers {
				ctx.reset()
				handler(ctx)
				if !ctx.isNext() {
					break
				}
			}
		}

	}

}

func (h *HttpAdapter) Start(port string) error {
	http.HandleFunc("/", h.handleReq)

	h.Logger.Info("Server is running at http://localhost:" + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		return err
	}

	return nil
}

func CreateHttpAdapter(config *HttpAdapterConfig) BaseHttpAdapter {
	assertNotNil(config.App, "HttpAdapterConfig.App")

	return &HttpAdapter{
		App:    config.App,
		Logger: CreateLogger("@Adapters/Http"),
		Config: config,
		Routes: make(map[string]HttpAdapterMethodMap),
	}
}
