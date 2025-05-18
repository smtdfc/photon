package photon

import (
	"fmt"
)

type App struct {
	Adapter       BaseAdapter
	SocketAdapter BaseSocketAdapter 
}

func (a *App) Start(port string) error {
	if err := a.Adapter.Start(); err != nil {
		fmt.Printf("[%s]: An error occurred when starting adapter\n", a.Adapter.GetName())
		fmt.Println(err.Error())
		return err
	}


	if err := a.Adapter.Listen(port); err != nil {
		fmt.Printf("[%s]: An error occurred when listening\n", a.Adapter.GetName())
		fmt.Println(err.Error())
		return err
	}

	if a.SocketAdapter != nil {
		if err := a.SocketAdapter.Start(); err != nil {
			fmt.Printf("[%s]: An error occurred when starting socket adapter\n", a.SocketAdapter.GetName())
			fmt.Println(err.Error())
			return err
		}
	}

	return nil
}

func NewApp(adapter BaseAdapter, socketAdapter ...BaseSocketAdapter) *App {
	app := &App{
		Adapter: adapter,
	}
	
	if len(socketAdapter) > 0 {
		app.SocketAdapter = socketAdapter[0]
	}
	
	adapter.SetApp(app)
	return app
}