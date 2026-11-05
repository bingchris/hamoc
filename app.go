package hamoc

import "net/http"

type App struct {
	router      router
	middlewares []Middleware
}

func New() *App {
	app := &App{}
	app.Use(Logger())
	app.Use(Recoverer())
	return app
}

func (a *App) Use(mw Middleware)          { a.middlewares = append(a.middlewares, mw) }
func (a *App) Get(p string, h Handler)    { a.router.add("GET", p, h) }
func (a *App) Post(p string, h Handler)   { a.router.add("POST", p, h) }
func (a *App) Put(p string, h Handler)    { a.router.add("PUT", p, h) }
func (a *App) Delete(p string, h Handler) { a.router.add("DELETE", p, h) }
func (a *App) Patch(p string, h Handler)  { a.router.add("PATCH", p, h) }

func (a *App) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h, params := a.router.match(req.Method, req.URL.Path)
	if h == nil {
		http.NotFound(w, req)
		return
	}
	handler := h
	for i := len(a.middlewares) - 1; i >= 0; i-- {
		handler = a.middlewares[i](handler)
	}
	r := &Request{Request: req, params: params}
	if err := handler(r).write(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *App) Listen(addr string) error {
	return http.ListenAndServe(addr, a)
}
