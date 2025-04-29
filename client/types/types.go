package types

type SessionState int

type JSON_payload struct {
	Username string `json:"author"`
	Text     string `json:"text"`
	Time     string `json:"time"`
	Status   string `json:"status"`
}

const (
	state SessionState = iota
	JoinRoom
	CreateRoom
	LeaveRoom
	ListRooms
	ChatRoom
	CancelCreate
	CancelJoin
	CancelList
)

type StateChangeMsg struct {
	Msg SessionState
}
