package photon

type App struct {
	GlobalData map[string]any
	Logger *Logger
}

func (a *App) Start(port string) error{
	a.Logger.Info("App started successful")
	return nil
}

func CreateApp() *App {
	return &App{
		GlobalData: make(map[string]any),
		Logger: CreateLogger("@App"),
	}
}
