package model


type Table struct {
	TableId int `db:"table_id" json:"table_id"`
	ProjectId int `db:"project_id" json:"project_id"`
	TableName string `db:"table_name" json:"table_name"`
	TableNameLogical string `db:"table_name_logical" json:"table_name_logical"`
	DelFlg int `db:"del_flg" json:"del_flg"`
	CreateUserId int `db:"create_user_id" json:"create_user_id"`
	UpdateUserId int `db:"update_user_id" json:"update_user_id"`
	CreatedAt string `db:"created_at " json:"created_at "`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}