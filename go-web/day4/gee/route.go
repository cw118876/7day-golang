package gee

import (
	"net/http"
	"strings"
)

type Router struct {
	roots    map[string]*Node
	handlers map[string]HandleFunc
}

func newRouter() *Router {
	return &Router{roots: make(map[string]*Node, 0),
		handlers: make(map[string]HandleFunc, 0),
	}
}

func (r *Router) addRouter(method, pattern string, handler HandleFunc) {
	key := method + "-" + pattern
	r.handlers[key] = handler
	parts := parsePath(pattern)
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = newNode("", false)
	}
	r.roots[method].insert(parts, pattern, 0)
}

func (r *Router) getRouter(method, path string) (*Node, map[string]string) {
	searchParts := parsePath(path)
	params := make(map[string]string, 0)
	_, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	node := r.roots[method].search(searchParts, 0)
	if node != nil {
		parts := parsePath(node.pattern)
		for index, p := range parts {
			if p[0] == ':' {
				params[p[1:]] = searchParts[index]
			}
			if p[0] == '*' && len(p) > 1 {
				params[p[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return node, params

	}
	return nil, nil

}

func (r *Router) getRoutes(method string) []*Node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*Node, 0)
	root.travel(&nodes)
	return nodes
}

func (r *Router) handle(c *Context) {
	n, params := r.getRouter(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)

	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
