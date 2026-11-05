package hamoc

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  int
	Body    string
	JSON    any
	Headers map[string]string
}

func Text(body string) Response { return Response{Body: body} }

func JSONResponse(status int, v any) Response {
	return Response{Status: status, JSON: v}
}

func (r Response) WithHeader(k, v string) Response {
	if r.Headers == nil {
		r.Headers = map[string]string{}
	}
	r.Headers[k] = v
	return r
}

func (r Response) WithStatus(code int) Response {
	r.Status = code
	return r
}

func (r Response) write(w http.ResponseWriter) error {
	if r.Status == 0 {
		r.Status = http.StatusOK
	}
	for k, v := range r.Headers {
		w.Header().Set(k, v)
	}
	if r.JSON != nil {
		b, err := json.Marshal(r.JSON)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(r.Status)
		_, err = w.Write(b)
		return err
	}
	w.WriteHeader(r.Status)
	_, err := w.Write([]byte(r.Body))
	return err
}
