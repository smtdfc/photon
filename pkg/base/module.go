package photon

type Module struct{
  Name string
  Prefix string
  App *App
  Adapter *AdapterManager
}


func NewModule(name string, app *App) *Module{
  return &Module{
    Name:name,
    Prefix:"",
    App:app,
    Adapter: app.Adapter,
  }
}
