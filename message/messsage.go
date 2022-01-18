package message

const (
	Token         string = "TOKEN"
	RecoveryToken string = "RECOVERY_TOKEN"
	Ack           string = "ACK"
)

type Message struct {
	Type string
	Id   int
}
