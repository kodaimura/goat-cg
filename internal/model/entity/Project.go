package entity


type Project struct {
	ProjectId int `db:"project_id" json:"project_id"`
	ProjectCd string `db:"project_cd" json:"project_cd"`
	ProjectName string `db:"project_name" json:"project_name"`
	CreatedAt string `db:"created_at " json:"created_at "`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}