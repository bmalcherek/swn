package main

import (
	"github.com/bmalcherek/swn/node"
)

const (
	nodeCount int = 3
)

func main() {
	endChan := make(chan interface{})

	nodes := []*node.Node{}
	for i := 0; i < nodeCount; i++ {
		nodes = append(nodes, &node.Node{})
	}
	for i := 0; i < nodeCount; i++ {
		commChan := make(chan int)
		nodes[i].NextNodeChan = commChan
		nodes[(i+1)%nodeCount].PrevNodeChan = commChan
	}

	for i := 0; i < nodeCount; i++ {
		go nodes[i].Run(i)
	}
	// fmt.Println(nodes[0], nodes[1], nodes[2])

	<-endChan
}
