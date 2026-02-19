package incident

type Incident struct {
	ID       string                  `json:"id"`
	TeamID   string                  `json:"team_id"`
	Content  *map[string]interface{} `json:"content"`
	Resolved bool
}
