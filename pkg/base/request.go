package photon

type Request interface {
    GetHeader(name string) string
    GetCookie(name string) (string, error)
    GetIP() string

    GetQuery(key string) string
    GetParam(key string) string
    GetFormValue(key string) string
    GetBody() ([]byte, error)

    Method() string
    Path() string
    Host() string
    URL() string
    ContentType() string

    BindJSON(obj interface{}) error
    BindForm(obj interface{}) error
    BindQuery(obj interface{}) error

    Context() any
    FormFile(name string) (any, error)
}