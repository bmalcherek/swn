package node

import (
	"fmt"
	"time"

	"github.com/bmalcherek/swn/message"
)

type Node struct {
	NextNodeChan chan message.Message
	PrevNodeChan chan message.Message
	NodeId       int
}

func (n *Node) Run() {
	if n.NodeId == 0 {
		n.NextNodeChan <- message.Message{
			Type: message.Token,
			Id:   0,
		}
	}
	for {
		select {
		case msg := <-n.PrevNodeChan:
			n.handleTokenMessage(msg)
		case msg := <-n.NextNodeChan:
			n.handleNextNodeMessage(msg)
		}
	}
}

func (n *Node) handleTokenMessage(msg message.Message) {
	fmt.Printf("Node %d recieved msg %v\n", n.NodeId, msg)
	n.PrevNodeChan <- message.Message{Type: message.Ack, Id: msg.Id}
	time.Sleep(2 * time.Second)
	msg.Id += 1
	n.NextNodeChan <- msg
}

func (n *Node) handleNextNodeMessage(msg message.Message) {
	fmt.Printf("Node %d recieved msg %v\n", n.NodeId, msg)
}
