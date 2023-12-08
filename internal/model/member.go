package model


type Member struct {
	ProjectId int `db:"project_id" json:"project_id"`
	UserId int `db:"user_id" json:"user_id"`
	UserStatus string `db:"user_status" json:"user_status"`
	UserRole string `db:"user_role" json:"user_role"`
	CreatedAt string `db:"created_at " json:"created_at "`
	UpdatedAt string `db:"created_at " json:"created_at "`
}