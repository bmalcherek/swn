package node

import (
	"log"
	"math/rand"
	"time"

	"github.com/bmalcherek/swn/message"
)

const (
	timeout                     = 1500 * time.Millisecond
	prevNodeRecvChanFaultChance = 0.2
	nextNodeRecvChanFaultChance = 0.2
	nodeCount                   = 3
)

type Node struct {
	NextNodeSendChan          chan<- message.Message
	NextNodeRecvChan          <-chan message.Message
	PrevNodeSendChan          chan<- message.Message
	PrevNodeRecvChan          <-chan message.Message
	NodeId                    int
	lastSentId                int
	lastRecievedAckId         int
	lastRecievedTokenId       int
	lastSentRecoveryTokenId   int
	lastRecievedRecoveryAckId int
}

func (n *Node) Run() {
	if n.NodeId == 0 {
		n.sendToken(1)
	}

	for {
		select {
		case msg := <-n.PrevNodeRecvChan:
			if rand.Float32() < prevNodeRecvChanFaultChance {
				log.Printf("NODE %d PREV NODE RECV CHAN FAULT\n", n.NodeId)
				time.Sleep(time.Duration(rand.Intn(10)+1) * time.Second)
			}
			n.handlePrevNodeMessage(msg)
		case msg := <-n.NextNodeRecvChan:
			if rand.Float32() < nextNodeRecvChanFaultChance {
				log.Printf("NODE %d NEXT NODE RECV CHAN FAULT\n", n.NodeId)
				time.Sleep(time.Duration(rand.Intn(10)+1) * time.Second)
			}
			n.handleNextNodeMessage(msg)
		}
	}
}

func (n *Node) sendToken(id int) {
	msg := message.Message{Type: message.Token, Id: id}
	n.NextNodeSendChan <- msg
	log.Printf("Node %d sent token %v\n", n.NodeId, msg)
	n.lastSentId = id
	go n.checkForRecievedAck()
}

func (n *Node) sendAck(id int) {
	msg := message.Message{Type: message.Ack, Id: id}
	n.PrevNodeSendChan <- msg
	log.Printf("Node %d sent ack %v\n", n.NodeId, msg)
}

func (n *Node) sendRecoveryToken(id, target int) {
	msg := message.Message{Type: message.RecoveryToken, Id: id, Target: target}
	n.PrevNodeSendChan <- msg
	n.lastSentRecoveryTokenId = id
	log.Printf("Node %d sent recovery token %v\n", n.NodeId, msg)
	go n.checkForRecievedRecoveryAck(target)
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
		log.Printf("Node %d recieved token %v\n", n.NodeId, msg)
		n.sendAck(msg.Id)
		if n.lastRecievedTokenId < msg.Id {
			n.lastRecievedTokenId = msg.Id
			n.enterCriticalRegion(msg.Id)
		}
	} else if msg.Type == message.RecoveryAck {
		log.Printf("Node %d recieved recovery ack %v\n", n.NodeId, msg)
		n.lastRecievedRecoveryAckId = msg.Id
	}
}

func (n *Node) handleNextNodeMessage(msg message.Message) {
	log.Printf("Node %d recieved msg %v\n", n.NodeId, msg)
	if msg.Type == message.Ack {
		n.lastRecievedAckId = msg.Id
	} else if msg.Type == message.RecoveryToken {
		n.sendRecoveryAck(msg.Id)
		if msg.Target != n.NodeId {
			n.sendRecoveryToken(msg.Id, msg.Target)
			return
		}
		if n.lastRecievedTokenId < msg.Id {
			n.lastRecievedTokenId = msg.Id
			n.enterCriticalRegion(msg.Id)
		}
	}
}

func (n *Node) enterCriticalRegion(id int) {
	log.Printf("Node %d entered critical region\n", n.NodeId)
	time.Sleep(2 * time.Second)
	n.sendToken(id + 1)
}

func (n *Node) checkForRecievedAck() {
	time.Sleep(timeout)
	if n.lastSentId > n.lastRecievedAckId && n.lastSentId > n.lastRecievedTokenId {
		log.Printf("Node %d didn't recieve ack for %d\n", n.NodeId, n.lastSentId)
		n.sendToken(n.lastSentId)
		if n.lastSentId > n.lastRecievedRecoveryAckId {
			n.sendRecoveryToken(n.lastSentId, (n.NodeId+1)%nodeCount)
		}
	}
}

func (n *Node) checkForRecievedRecoveryAck(target int) {
	time.Sleep(timeout)
	if n.lastSentRecoveryTokenId > n.lastRecievedRecoveryAckId {
		log.Printf("Node %d didn't recieve recovery ack for %d\n", n.NodeId, n.lastSentId)
		n.sendRecoveryToken(n.lastSentRecoveryTokenId, target)
	}
}
