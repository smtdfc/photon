package photon

import (
	"fmt"
)

type App struct {
	Adapter       *AdapterManager
	HttpAdapter   BaseHTTPAdapter
	SocketAdapter BaseSocketAdapter
}

func (a *App) Start(port string) error {
	a.Adapter.EsureAdapter("http")

	if err := a.Adapter.HttpAdapter.Start(); err != nil {
		fmt.Printf("[%s]: An error occurred when starting adapter\n", a.HttpAdapter.GetName())
		fmt.Println(err.Error())
		return err
	}

	if err := a.Adapter.HttpAdapter.Listen(port); err != nil {
		fmt.Printf("[%s]: An error occurred when listening\n", a.HttpAdapter.GetName())
		fmt.Println(err.Error())
		return err
	}

	if a.SocketAdapter != nil {
		if err := a.Adapter.SocketAdapter.Start(); err != nil {
			fmt.Printf("[%s]: An error occurred when starting socket adapter\n", a.SocketAdapter.GetName())
			fmt.Println(err.Error())
			return err
		}
	}

	return nil
}

func NewApp() *App {
	app := &App{
		Adapter: &AdapterManager{},
	}

	return app
}
