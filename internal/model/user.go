package model


type User struct {
	UserId int `db:"user_id" json:"user_id"`
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
	CreatedAt string `db:"created_at " json:"created_at "`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}