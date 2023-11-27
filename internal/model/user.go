package model


type User struct {
	UserId int `db:"user_id" json:"user_id"`
	UserName string `db:"user_name" json:"user_name"`
	Password string `db:"password" json:"password"`
	CreatedAt string `db:"created_at " json:"created_at "`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}