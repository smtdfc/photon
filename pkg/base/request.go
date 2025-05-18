package photon


type Request interface{
  GetHeader(name string) string
  GetCookie(name string) (string, error)
  GetIP() string
}
