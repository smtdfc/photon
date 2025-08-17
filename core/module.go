package core

import (
	"errors"
	"github.com/smtdfc/photon/logger"
)

type Module struct {
	App      *App
	Logger   *logger.Logger
	Name     string
	provided any
	injected map[string]any
	onStartCallbacks  []func()
}

func (m *Module) triggerHook(name string){
	if(name == "start"){
		for _,cb := range m.onStartCallbacks{
			cb()
		}
	}
}

func (m *Module) OnStart(callback func()){
	m.onStartCallbacks = append(m.onStartCallbacks,callback)
}

func (m *Module) inject(edge string, provider Provider) error {
	if m.injected[edge] != nil {
		return errors.New("Edge " + edge + " used in module " + m.Name)
	}
	m.injected[edge] = provider.provide()
	return nil
}

func (m *Module) provide() any {
	return m.provided
}

func NewModule(app *App, name string) *Module {
	return &Module{
		Name:     name,
		App:      app,
		Logger:   logger.New("stdout"),
		injected: make(map[string]any),
	}
}
