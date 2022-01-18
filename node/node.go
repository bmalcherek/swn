package node

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bmalcherek/swn/message"
)

const (
	timeout        = time.Second
	ackFaultChance = 0.3
)

type Node struct {
	NextNodeSendChan    chan<- message.Message
	NextNodeRecvChan    <-chan message.Message
	PrevNodeSendChan    chan<- message.Message
	PrevNodeRecvChan    <-chan message.Message
	NodeId              int
	lastSendId          int
	lastRecievedAckId   int
	lastRecievedTokenId int
}

func (n *Node) Run() {
	if n.NodeId == 0 {
		n.sendToken(1)
	}

	for {
		select {
		case msg := <-n.PrevNodeRecvChan:
			n.handlePrevNodeMessage(msg)
		case msg := <-n.NextNodeRecvChan:
			n.handleNextNodeMessage(msg)
		}
	}
}

func (n *Node) sendToken(id int) {
	msg := message.Message{Type: message.Token, Id: id}
	n.NextNodeSendChan <- msg
	fmt.Printf("Node %d send token %v\n", n.NodeId, msg)
	n.lastSendId = id
	go n.checkForRecievedAck()
}

func (n *Node) handlePrevNodeMessage(msg message.Message) {
	fmt.Printf("Node %d recieved msg %v\n", n.NodeId, msg)
	if n.lastRecievedTokenId == msg.Id {
		return
	}
	n.lastRecievedTokenId = msg.Id
	if rand.Float32() < ackFaultChance {
		fmt.Printf("Node %d ack channel failed\n", n.NodeId)
		time.Sleep(6 * time.Second)
	}
	n.PrevNodeSendChan <- message.Message{Type: message.Ack, Id: msg.Id}
	time.Sleep(2 * time.Second)
	n.sendToken(msg.Id + 1)
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
