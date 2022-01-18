package main

import (
	"fmt"

	"github.com/bmalcherek/swn/node"
)

const (
	nodeCount int = 3
)

func main() {
	nodes := []*node.Node{}
	for i := 0; i < nodeCount; i++ {
		nodes = append(nodes, &node.Node{})
	}
	for i := 0; i < nodeCount; i++ {
		commChan := make(chan interface{})
		nodes[i].NextNodeChan = commChan
		nodes[(i+1)%nodeCount].PrevNodeChan = commChan
	}
	fmt.Println(nodes[0], nodes[1], nodes[2])
}
