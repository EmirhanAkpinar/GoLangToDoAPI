package models

type Credentials struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Type     int    `json:"type"`
}
