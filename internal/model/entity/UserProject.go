package entity


type UserProject struct {
	UserId int `db:"user_id" json:"user_id"`
	ProjectId int `db:"project_id" json:"project_id"`
	StateCls string `db:"state_cls" json:"state_cls"`
	RoleCls string `db:"role_cls" json:"role_cls"`
	CreateAt string `db:"create_at" json:"create_at"`
	UpdateAt string `db:"create_at" json:"create_at"`
}