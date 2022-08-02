package entity


type Table struct {
	TableId int `db:"table_id" json:"table_id"`
	ProjectId int `db:"project_id" json:"project_id"`
	TableName string `db:"table_name" json:"table_name"`
	TableNameLogical string `db:"table_name_logical" json:"table_name_logical"`
	CreateUserId int `db:"create_user_id" json:"create_user_id"`
	UpdateUserId int `db:"create_user_id" json:"create_user_id"`
	DelFlg int `db:"del_flg" json:"del_flg"`
	CreateAt string `db:"create_at" json:"create_at"`
	UpdateAt string `db:"update_at" json:"update_at"`
}