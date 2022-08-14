package entity


type Project struct {
	ProjectId int `db:"project_id" json:"project_id"`
	ProjectCd string `db:"project_cd" json:"project_cd"`
	ProjectName string `db:"project_name" json:"project_name"`
	CreateAt string `db:"create_at" json:"create_at"`
	UpdateAt string `db:"update_at" json:"update_at"`
}