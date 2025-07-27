package photon

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"
)

type HttpAdapterResponse struct {
	Writer http.ResponseWriter
}

func (r *HttpAdapterResponse) RawWriter() http.ResponseWriter {
	return r.Writer
}

func (r *HttpAdapterResponse) SetHeader(name, value string) Response {
	r.Writer.Header().Set(name, value)
	return r
}

func (r *HttpAdapterResponse) SetContentType(ct string) Response {
	r.Writer.Header().Set("Content-Type", ct)
	return r
}

func (r *HttpAdapterResponse) SetCookie(name, value string, opts CookieOption) Response {
	expires := time.Now().Add(time.Duration(opts.MaxAge) * time.Second)
	http.SetCookie(r.Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     opts.Path,
		HttpOnly: opts.HttpOnly,
		Secure:   opts.Secure,
		Expires:  expires,
	})
	return r
}

func (r *HttpAdapterResponse) Code(code int) Response {
	r.Writer.WriteHeader(code)
	return r
}

func (r *HttpAdapterResponse) JSON(data JSON) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(r.Writer, "Failed to marshal", http.StatusInternalServerError)
		return err
	}
	r.Writer.Header().Set("Content-Type", "application/json")
	r.Writer.Write(jsonBytes)
	return nil
}

func (r *HttpAdapterResponse) Text(data string) error {
	r.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err := r.Writer.Write([]byte(data))
	return err
}

func (r *HttpAdapterResponse) HTML(html string) error {
	r.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := r.Writer.Write([]byte(html))
	return err
}

func (r *HttpAdapterResponse) File(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, _ := f.Read(buf)
	mimeType := http.DetectContentType(buf[:n])

	r.Writer.Header().Set("Content-Type", mimeType)

	f.Seek(0, io.SeekStart)

	_, err = io.Copy(r.Writer, f)
	return err
}

func (r *HttpAdapterResponse) Stream(reader io.Reader) error {
	_, err := io.Copy(r.Writer, reader)
	return err
}

func (r *HttpAdapterResponse) Redirect(status int, location string) error {
	http.Redirect(r.Writer, nil, location, status)
	return nil
}
