package models

type URLInfo struct {
	UUID      string `json:"id,omitempty"`
	UserID    string `json:"userid,omitempty"`
	LongURL   string `json:"longurl,omitempty"`
	ShortURL  string `json:"shorturl,omitempty"`
	IsDeleted int    `json:"is_deleted,omitempty"`
}
