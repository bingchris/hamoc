package hamoc

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// Request wraps *http.Request with Flask-like conveniences.
type Request struct {
	*http.Request
	params map[string]string
	body   []byte
	read   bool
}

// Param returns a path parameter by name.
func (r *Request) Param(name string) string { return r.params[name] }

// Params returns all path parameters.
func (r *Request) Params() map[string]string { return r.params }

// Query returns a query string value (first value).
func (r *Request) Query(name string) string { return r.URL.Query().Get(name) }

// QueryAll returns all values for a query key.
func (r *Request) QueryAll(name string) []string { return r.URL.Query()[name] }

// Cookies returns all cookies as a name->value map.
func (r *Request) Cookies() map[string]string {
	out := map[string]string{}
	for _, c := range r.Request.Cookies() {
		out[c.Name] = c.Value
	}
	return out
}

// Cookie returns a single cookie value by name.
func (r *Request) Cookie(name string) string { return r.Cookies()[name] }

// Headers returns all headers as a name->value map.
func (r *Request) Headers() map[string]string {
	out := map[string]string{}
	for k, v := range r.Request.Header {
		out[k] = strings.Join(v, ", ")
	}
	return out
}

// Header returns a single header value by name.
func (r *Request) Header(name string) string { return r.Request.Header.Get(name) }

// BodyBytes reads and caches the request body.
func (r *Request) BodyBytes() ([]byte, error) {
	if r.read {
		return r.body, nil
	}
	b, err := io.ReadAll(r.Request.Body)
	if err != nil {
		return nil, err
	}
	r.body = b
	r.read = true
	return b, nil
}

// JSON decodes the request body as JSON into v.
func (r *Request) JSON(v any) error {
	b, err := r.BodyBytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}
