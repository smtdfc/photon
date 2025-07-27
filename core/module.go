package core

import (
	http "github.com/smtdfc/photon/http_adapter"
	"github.com/smtdfc/photon/internal/helpers"
)

type Module struct {
	Name         string
	App          *App
	Adapter      *AdapterManager
	provided     any
	injected     map[string]any
	Router       *ModuleRouter
	OnModuleInit func()
}

type ModuleRouter struct {
	Module *Module
	Prefix string
}

func (r *ModuleRouter) Get(pattern string, handlers ...http.HttpHandler) {
	helpers.AssertNotNil(r.Module.Adapter.Http, "HttpAdapter")
	r.Module.Adapter.Http.AddRoute("GET", r.Prefix+pattern, handlers...)
}

func (r *ModuleRouter) Post(pattern string, handlers ...http.HttpHandler) {
	helpers.AssertNotNil(r.Module.Adapter.Http, "HttpAdapter")
	r.Module.Adapter.Http.AddRoute("POST", r.Prefix+pattern, handlers...)
}

func (r *ModuleRouter) Put(pattern string, handlers ...http.HttpHandler) {
	helpers.AssertNotNil(r.Module.Adapter.Http, "HttpAdapter")
	r.Module.Adapter.Http.AddRoute("PUT", r.Prefix+pattern, handlers...)
}

func (r *ModuleRouter) Head(pattern string, handlers ...http.HttpHandler) {
	helpers.AssertNotNil(r.Module.Adapter.Http, "HttpAdapter")
	r.Module.Adapter.Http.AddRoute("HEAD", r.Prefix+pattern, handlers...)
}

func (r *ModuleRouter) Option(pattern string, handlers ...http.HttpHandler) {
	helpers.AssertNotNil(r.Module.Adapter.Http, "HttpAdapter")
	r.Module.Adapter.Http.AddRoute("OPTION", r.Prefix+pattern, handlers...)
}

func (m *Module) GetInject(edge string) any {
	return m.injected[edge]
}

func (m *Module) GetProvide() any {
	return m.provided
}

func (m *Module) Inject(edge string, provider Provider) {
	m.injected[edge] = provider.GetProvide()
}

func (m *Module) Provide(value any) {
	m.provided = value
}

func CreateModule(app *App, name string) *Module {
	mod := &Module{
		Name:     name,
		App:      app,
		Adapter:  app.Adapter,
		injected: make(map[string]any),
	}

	mod.Router = &ModuleRouter{
		Module: mod,
		Prefix: "",
	}

	app.registerModule(name, mod)
	return mod
}
