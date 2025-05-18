package photon

import(
  "fmt"
)

type App struct{
  Adapter BaseAdapter
}

func (a *App) Start(port string) error{
  err := a.Adapter.Start()
  if err != nil{
    fmt.Println("["+a.Adapter.GetName()+"]: An error occurred when start adapter")
    fmt.Println(err.Error())
    return err
  }
  
  return a.Adapter.Listen(port)
}


func NewApp(adapter BaseAdapter) *App{
  return &App{
    Adapter: adapter,
  }
}
