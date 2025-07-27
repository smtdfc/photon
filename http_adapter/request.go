package http_adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

type HttpAdapterRequest struct {
	Original *http.Request
}

func (r *HttpAdapterRequest) GetHeader(name string) string {
	return r.Original.Header.Get(name)
}

func (r *HttpAdapterRequest) GetMethod() string {
	return r.Original.Method
}

func (r *HttpAdapterRequest) GetPath() string {
	return r.Original.URL.String()
}

func (r *HttpAdapterRequest) GetFullURL() string {
	scheme := "http"
	if r.Original.TLS != nil {
		scheme = "https"
	}

	return scheme + "://" + r.Original.Host + r.Original.RequestURI
}

func (r *HttpAdapterRequest) GetBody() (map[string]any, error) {
	var body map[string]any

	data, err := io.ReadAll(r.Original.Body)
	if err != nil {
		return nil, err
	}

	defer r.Original.Body.Close()

	err = json.Unmarshal(data, &body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (r *HttpAdapterRequest) GetBodyReusable() (map[string]any, error) {
	var body map[string]any

	data, err := io.ReadAll(r.Original.Body)
	if err != nil {
		return nil, err
	}
	r.Original.Body.Close()

	r.Original.Body = io.NopCloser(bytes.NewReader(data))

	err = json.Unmarshal(data, &body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (r *HttpAdapterRequest) GetRawBody() ([]byte, error) {
	data, err := io.ReadAll(r.Original.Body)
	if err != nil {
		return nil, err
	}
	defer r.Original.Body.Close()
	return data, nil
}

func (r *HttpAdapterRequest) BindBody(target any) error {
	body, err := io.ReadAll(r.Original.Body)
	if err != nil {
		return err
	}

	defer r.Original.Body.Close()

	return json.Unmarshal(body, target)
}

func (r *HttpAdapterRequest) GetQuery(name string) string {
	return r.Original.URL.Query().Get(name)
}

func (r *HttpAdapterRequest) GetAllQuery() map[string][]string {
	return r.Original.URL.Query()
}

func (r *HttpAdapterRequest) GetCookie(name string) (string, error) {
	cookie, err := r.Original.Cookie(name)
	if err != nil {
		return "", err
	} else {
		return cookie.Value, nil
	}
}

func (r *HttpAdapterRequest) GetFormValue(name string) string {
	return r.Original.FormValue(name)
}

func (r *HttpAdapterRequest) GetPostFormValue(name string) (string, error) {
	err := r.Original.ParseForm()
	if err != nil {
		return "", err
	}
	return r.Original.PostFormValue(name), nil
}

func (r *HttpAdapterRequest) GetFile(name string, sizeMB int) (*File, error) {
	if r.Original.MultipartForm == nil {
		if err := r.Original.ParseMultipartForm(int64(sizeMB) << 20); err != nil {
			return nil, err
		}
	}

	file, fileHeader, err := r.Original.FormFile(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if fileHeader.Size > int64(sizeMB)<<20 {
		return nil, fmt.Errorf("file %s exceeds max size", fileHeader.Filename)
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = mime.TypeByExtension(filepath.Ext(fileHeader.Filename))
	}

	return &File{
		Filename: fileHeader.Filename,
		Size:     fileHeader.Size,
		MIME:     contentType,
		Content:  content,
		Reader:   bytes.NewReader(content),
	}, nil
}

func (r *HttpAdapterRequest) GetAllFiles() (map[string][]File, error) {
	if r.Original.MultipartForm == nil {
		err := r.Original.ParseMultipartForm(32 << 20)
		if err != nil {
			return nil, err
		}
	}

	result := make(map[string][]File)

	for field, fileHeaders := range r.Original.MultipartForm.File {
		for _, fh := range fileHeaders {
			f, err := fh.Open()
			if err != nil {
				return nil, fmt.Errorf("cannot open file %s: %w", fh.Filename, err)
			}

			content, err := io.ReadAll(f)
			f.Close()
			if err != nil {
				return nil, fmt.Errorf("cannot read file %s: %w", fh.Filename, err)
			}

			reader := bytes.NewReader(content)
			mime := fh.Header.Get("Content-Type")

			result[field] = append(result[field], File{
				Filename: fh.Filename,
				Size:     fh.Size,
				MIME:     mime,
				Content:  content,
				Reader:   reader,
			})
		}
	}

	return result, nil
}

func (r *HttpAdapterRequest) GetAllForm() (map[string]string, error) {
	err := r.Original.ParseForm()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for key, values := range r.Original.Form {
		if len(values) > 0 {
			result[key] = values[0]
		}
	}
	return result, nil
}

func (r *HttpAdapterRequest) GetAllHeader() (map[string]string, error) {
	headers := make(map[string]string)
	for headerName, headerValue := range r.Original.Header {
		headers[headerName] = strings.Join(headerValue, ", ")
	}

	return headers, nil
}
