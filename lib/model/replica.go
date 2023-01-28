package model

type Replica struct {
	peerID   int
	userName string
	message  string
}

func NewReplica(peerID int, userName, message string) Replica {
	return Replica{
		peerID:   peerID,
		userName: userName,
		message:  message,
	}
}

func (r Replica) PeerID() int {
	return r.peerID
}

func (r Replica) UserName() string {
	return r.userName
}

func (r Replica) Message() string {
	return r.message
}
