package photon

type App struct {
	GlobalData map[string]any
}

func CreateApp() *App {
	return &App{
		GlobalData: make(map[string]any),
	}
}
