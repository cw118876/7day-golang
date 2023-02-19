package gee

import (
	"fmt"
	"net/http"
)

type (
	HandleFunc func(c *Context)
	Engine     struct {
		*RouterGroup
		router *Router
		groups []*RouterGroup // store all groups
	}
	RouterGroup struct {
		prefix      string
		middlewares []HandleFunc // support middleware
		parent      *RouterGroup // support nesting
		engine      *Engine      // all groups share an engine instance
	}
)

func NewEngine() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (g *RouterGroup) Group(prefix string) *RouterGroup {
	engine := g.engine
	newGroup := &RouterGroup{
		prefix: engine.prefix + prefix,
		parent: g,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (g *RouterGroup) addRouter(method, path string, handler HandleFunc) {
	pattern := g.prefix + path
	fmt.Printf("method: %4s, path: %s\n", method, pattern)
	g.engine.router.addRouter(method, pattern, handler)

}

func (g *RouterGroup) POST(path string, handler HandleFunc) {
	g.addRouter("POST", path, handler)
}

func (g *RouterGroup) GET(path string, handler HandleFunc) {
	g.addRouter("GET", path, handler)
}

func (e *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, e)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	e.router.handle(c)
}
