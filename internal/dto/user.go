package dto


type User struct {
	UserId int `json:"user_id"`
	Username string `json:"username"`
	Email string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}