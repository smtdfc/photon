package core

import (
	"github.com/smtdfc/photon/http_adapter"
	"github.com/smtdfc/photon/logger"
	"sync"
)

type AdapterManager struct {
	Http http_adapter.BaseHttpAdapter
}

type App struct {
	GlobalData map[string]any
	Logger     *logger.Logger
	Adapter    *AdapterManager
	modules    map[string]*Module
}

func (a *App) registerModule(name string, mod *Module) {
	if a.modules[name] != nil {
		a.Logger.Warn("Module " + name + " was registered multiple times")
	}
	a.modules[name] = mod
}

func (a *App) Start(port string) error {
	var wg sync.WaitGroup

	if a.Adapter.Http != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := a.Adapter.Http.Start(port); err != nil {
				a.Logger.Error("HTTP Adapter error: " + err.Error())
			}
		}()
	}
	/*
		if a.Adapter.Socket != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := a.Adapter.Socket.StartSocket("9090"); err != nil {
					a.Logger.Error("Socket Adapter error: " + err.Error())
				}
			}()
		}
	*/

	a.Logger.Info("App started all active adapters")
	wg.Wait()
	return nil
}

func CreateApp() *App {
	return &App{
		GlobalData: make(map[string]any),
		Logger:     logger.CreateLogger("@App"),
		Adapter:    &AdapterManager{},
		modules:    make(map[string]*Module),
	}
}
