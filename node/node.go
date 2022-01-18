package node

import (
	"fmt"
	"time"
)

type Node struct {
	NextNodeChan chan int
	PrevNodeChan chan int
}

func (n *Node) Run(nodeId int) {
	if nodeId == 0 {
		n.NextNodeChan <- 1
	}
	for {
		msg := <-n.PrevNodeChan
		fmt.Printf("Node %d recieved msg %d\n", nodeId, msg)
		time.Sleep(2 * time.Second)
		msg += 1
		n.NextNodeChan <- msg
	}
}
