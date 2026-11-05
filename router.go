package hamoc

import "strings"

// Handler returns a Response instead of writing to a ResponseWriter.
type Handler func(r *Request) Response

// Middleware wraps a Handler.
type Middleware func(next Handler) Handler

type route struct {
	method  string
	pattern string
	handler Handler
}

// router matches paths and extracts :params.
type router struct {
	routes []route
}

func (rt *router) add(method, pattern string, h Handler) {
	rt.routes = append(rt.routes, route{method, pattern, h})
}

func (rt *router) match(method, path string) (Handler, map[string]string) {
	sp := splitPath(path)
	for _, r := range rt.routes {
		if r.method != method {
			continue
		}
		pp := splitPath(r.pattern)
		if len(pp) != len(sp) {
			continue
		}
		params := map[string]string{}
		ok := true
		for i := range pp {
			if strings.HasPrefix(pp[i], ":") {
				params[pp[i][1:]] = sp[i]
			} else if pp[i] != sp[i] {
				ok = false
				break
			}
		}
		if ok {
			return r.handler, params
		}
	}
	return nil, nil
}

func splitPath(p string) []string {
	p = strings.Trim(p, "/")
	if p == "" {
		return []string{}
	}
	return strings.Split(p, "/")
}
