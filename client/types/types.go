package types

type SessionState int

type JSON_payload struct {
	Username string `json:"author"`
	Text     string `json:"text"`
	Time     string `json:"time"`
	Room     string `json:"room"`
}

const (
	state SessionState = iota
	JoinRoom
	CreateRoom
	ChatRoom
	CancelCreate
)

type RecMsg struct {
	Username string `json:"author"`
	Time     string `json:"time"`
	Text     string `json:"text"`
}

type StateChangeMsg struct {
	Msg SessionState
}
