package json_t

type JSON_payload struct {
	Username string `json:"author"`
	Text     string `json:"text"`
	Time     string `json:"time"`
	Status   string `json:"status"`
}
