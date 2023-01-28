package gee

import (
	"strings"
)

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
	splitedString := strings.Split(path, "/")

}
