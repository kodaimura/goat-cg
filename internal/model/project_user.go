package model


type ProjectUser struct {
	UserId int `db:"user_id" json:"user_id"`
	ProjectId int `db:"project_id" json:"project_id"`
	StateCls string `db:"state_cls" json:"state_cls"`
	RoleCls string `db:"role_cls" json:"role_cls"`
	CreatedAt string `db:"created_at " json:"created_at "`
	UpdatedAt string `db:"created_at " json:"created_at "`
}