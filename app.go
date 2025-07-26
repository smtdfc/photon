package photon

import (
	"sync"
)

type AdapterManager struct {
	Http BaseHttpAdapter
}

type App struct {
	GlobalData map[string]any
	Logger     *Logger
	Adapter    *AdapterManager
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
		Logger:     CreateLogger("@App"),
		Adapter:    &AdapterManager{},
	}
}
