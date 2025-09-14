package core

import (
	"errors"
	"github.com/smtdfc/photon/v2/logger"
)

type App struct {
	Logger         *logger.Logger
	Adapters       map[string]any
	GlobalData     map[string]any
	modules        map[string]*Module
	gatewayManager *GatewayManager
}

func (a *App) AddAdapter(name string, adapter Adapter) error {
	if a.Adapters[name] != nil {
		return errors.New("Adapter " + name + " has been added")
	} else {
		a.Adapters[name] = adapter
		return nil
	}
}

func (a *App) Start() {
	a.Logger.Info("Starting application...")

	for name, module := range a.modules {
		a.Logger.Info("Initializing module: " + name)
		if len(module.onStartCallbacks) > 0 {
			module.triggerHook("start")
		}

		a.Logger.Info("Module " + name + " initialized")
	}

	wg := a.gatewayManager.StartAll()
	wg.Wait()
}

func (a *App) SetGateway(name string, gatewayManager Gateway) {
	a.gatewayManager.SetGateway(name, gatewayManager)
}

func NewApp() *App {
	app := &App{
		Logger:     logger.New("stdout"),
		GlobalData: make(map[string]any),
		Adapters:   make(map[string]any),
		modules:    make(map[string]*Module),
	}

	app.gatewayManager = &GatewayManager{
		App: app,
	}

	return app
}
