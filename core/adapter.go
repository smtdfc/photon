package core

type Adapter interface {
	Init() error
}

func ResolveAdapter[T Adapter](module *Module, name string) T {
	var zero T
	adapters := module.App.Adapters

	if val, ok := adapters[name]; ok && val != nil {
		if adapter, ok := val.(T); ok {
			return adapter
		} else {
			module.App.Logger.Error("Cannot resolve adapter " + name)
			return zero
		}
	} else {
		module.App.Logger.Error("Cannot resolve adapter " + name)
	}

	return zero
}
