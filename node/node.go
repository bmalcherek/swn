package node

type Node struct {
	NextNodeChan chan interface{}
	PrevNodeChan chan interface{}
}
