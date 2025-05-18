package photon

type Response interface {
	SetHeader(name string, value string)
	SetStatus(code int)
	SetCookie(name string, value string, path string, maxAge int)

	JSON(code int, obj interface{}) error
	Text(code int, msg string) error
	HTML(code int, html string) error
	Blob(code int, contentType string, data []byte) error

	Redirect(code int, url string) error
	SendFile(filePath string) error

	Writer() any
}