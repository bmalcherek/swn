package node

import (
	"log"
	"math/rand"
	"time"

	"github.com/bmalcherek/swn/message"
)

const (
	timeout        = 1500 * time.Millisecond
	ackFaultChance = 0.5
	nodeCount      = 3
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
			if rand.Float32() > 0.9999 {
				n.handlePrevNodeMessage(msg)
			}
		case msg := <-n.NextNodeRecvChan:
			n.handleNextNodeMessage(msg)
		}
	}
}

func (n *Node) sendToken(id int) {
	msg := message.Message{Type: message.Token, Id: id}
	n.NextNodeSendChan <- msg
	log.Printf("Node %d sent token %v\n", n.NodeId, msg)
	n.lastSendId = id
	go n.checkForRecievedAck()
}

func (n *Node) sendAck(id int) {
	if rand.Float32() < ackFaultChance {
		log.Printf("Node %d ack channel failed\n", n.NodeId)
		time.Sleep(6 * time.Second)
	}
	msg := message.Message{Type: message.Ack, Id: id}
	n.PrevNodeSendChan <- msg
	log.Printf("Node %d sent ack %v\n", n.NodeId, msg)
}

func (n *Node) sendRecoveryToken(id int) {
	msg := message.Message{Type: message.RecoveryToken, Id: id, Target: (n.NodeId + 1) % nodeCount}
	n.PrevNodeSendChan <- msg
	log.Printf("Node %d sent recovery token %v\n", n.NodeId, msg)
}

func (n *Node) sendRecoveryAck(id int) {
	msg := message.Message{Type: message.RecoveryAck, Id: id}
	n.NextNodeSendChan <- msg
	log.Printf("Node %d sent recovery ack %v\n", n.NodeId, msg)
}

func (n *Node) handlePrevNodeMessage(msg message.Message) {
	if msg.Type == message.Token {
		if n.lastRecievedTokenId == msg.Id {
			return
		}
		log.Printf("Node %d recieved msg %v\n", n.NodeId, msg)
		n.lastRecievedTokenId = msg.Id
		go n.sendAck(msg.Id)
		time.Sleep(2 * time.Second)
		n.sendToken(msg.Id + 1)
	}
}

func (n *Node) handleNextNodeMessage(msg message.Message) {
	log.Printf("Node %d recieved msg %v\n", n.NodeId, msg)
	if msg.Type == message.Ack {
		n.lastRecievedAckId = msg.Id
	} else if msg.Type == message.RecoveryToken {
		n.sendRecoveryAck(msg.Id)
	}
}

func (n *Node) checkForRecievedAck() {
	time.Sleep(timeout)
	if n.lastSendId > n.lastRecievedAckId && n.lastSendId > n.lastRecievedTokenId {
		log.Printf("Node %d didn't recieve ack for %d, %v\n", n.NodeId, n.lastSendId, n.NextNodeSendChan)
		n.sendToken(n.lastSendId)
		n.sendRecoveryToken(n.lastSendId)
	}
}
