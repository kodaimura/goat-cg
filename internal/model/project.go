package model


type Project struct {
	ProjectId int `db:"project_id" json:"project_id"`
	ProjectName string `db:"project_name" json:"project_name"`
	ProjectMemo string `db:"project_memo" json:"project_memo"`
	UserId string `db:"user_id" json:"user_id"`
	Username string `db:"username" json:"username"`
	CreatedAt string `db:"created_at " json:"created_at "`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}