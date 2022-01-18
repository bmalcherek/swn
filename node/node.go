package node

import (
	"fmt"
	"time"

	"github.com/bmalcherek/swn/message"
)

const (
	timeout = time.Second
)

type Node struct {
	NextNodeSendChan  chan<- message.Message
	NextNodeRecvChan  <-chan message.Message
	PrevNodeSendChan  chan<- message.Message
	PrevNodeRecvChan  <-chan message.Message
	NodeId            int
	lastSendId        int
	lastRecievedAckId int
}

func (n *Node) Run() {
	if n.NodeId == 0 {
		n.sendToken(1)
	}

	for {
		select {
		case msg := <-n.PrevNodeRecvChan:
			n.handleTokenMessage(msg)
		case msg := <-n.NextNodeRecvChan:
			n.handleNextNodeMessage(msg)
		}
	}
}

func (n *Node) sendToken(id int) {
	n.NextNodeSendChan <- message.Message{Type: message.Token, Id: id}
	n.lastSendId = id
	go n.checkForRecievedAck()
}

func (n *Node) handleTokenMessage(msg message.Message) {
	fmt.Printf("Node %d recieved msg %v\n", n.NodeId, msg)
	// n.PrevNodeChan <- message.Message{Type: message.Ack, Id: msg.Id}
	// time.Sleep(2 * time.Second)
	msg.Id += 1
	fmt.Println("XDDD")
	// n.NextNodeChan <- msg
	// n.lastSendId = msg.Id
}

func (n *Node) handleNextNodeMessage(msg message.Message) {
	fmt.Printf("Node %d recieved msg %v\n", n.NodeId, msg)
	n.lastRecievedAckId = msg.Id
}

func (n *Node) checkForRecievedAck() {
	time.Sleep(timeout)
	if n.lastSendId > n.lastRecievedAckId {
		fmt.Printf("Node %d didn't recieve ack for %d, %v\n", n.NodeId, n.lastSendId, n.NextNodeSendChan)
		n.sendToken(n.lastSendId)
	}
}
