package types

type Sentence struct {
	ID         uint   `json:"id"`
	UUID       string `json:"uuid"`
	Hitokoto   string `json:"hitokoto"`
	Type       string `json:"type"`
	From       string `json:"from"`
	FromWho    string `json:"from_who"`
	Creator    string `json:"creator"`
	CreatorUID uint   `json:"creator_uid"`
	Reviewer   uint   `json:"reviewer"`
	CommitFrom string `json:"commit_from"`
	CreatedAt  string `json:"created_at"` // Will be converted to time.Time
	Length     uint   `json:"length"`
}
