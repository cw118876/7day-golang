package gee

import (
	"fmt"
	"strings"
)

/*
 * pattern: used as the identifier for handler together with method name
 * path: partial path for traversing the search tree formed with Node
 */
type Node struct {
	part     string
	pattern  string
	isWild   bool
	children []*Node
}

func newNode(part string, wild bool) *Node {
	return &Node{part: part,
		isWild: wild,
	}
}

func parsePath(path string) []string {
	splitedStrs := strings.Split(path, "/")
	parsedStrs := make([]string, 0)
	for _, part := range splitedStrs {
		if part == "" {
			continue
		}
		parsedStrs = append(parsedStrs, part)
		if part[0] == '*' {
			break
		}
	}
	return parsedStrs

}

func (n *Node) String() string {
	return fmt.Sprintf("Node: part: %s, pattern: %s, is Wildcard: %t\n", n.part, n.pattern, n.isWild)

}

func (n *Node) matchSingleChild(part string) *Node {
	for _, c := range n.children {
		if c.part == part || c.isWild {
			return c
		}
	}
	return nil

}

func (n *Node) matchAllChild(part string) []*Node {
	matched := make([]*Node, 0)
	for _, c := range n.children {
		if c.part == part || c.isWild {
			matched = append(matched, c)
		}
	}
	return matched

}

func (n *Node) insert(parts []string, pattern string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	if len(parts) < height {
		return
	}
	p := parts[height]
	c := n.matchSingleChild(p)
	if c == nil {
		c = newNode(p, p[0] == ':' || p[0] == '*')
		n.children = append(n.children, c)
	}
	c.insert(parts, pattern, height+1)
}

func (n *Node) search(parts []string, height int) *Node {
	if len(parts) < height {
		return nil
	}
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern != "" {
			return n
		}
		return nil
	}
	p := parts[height]
	children := n.matchAllChild(p)
	for _, c := range children {
		result := c.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}

func (n *Node) travel(list *([]*Node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, c := range n.children {
		c.travel(list)
	}

}
