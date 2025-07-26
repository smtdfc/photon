package photon

type Module struct {
	Name         string
	provided     any
	injected     map[string]any
	OnModuleInit func()
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
		injected: make(map[string]any),
	}

	app.registerModule(name, mod)
	return mod
}
