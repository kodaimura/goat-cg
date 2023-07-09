package entity


type Column struct {
	ColumnId int `db:"column_id" json:"column_id"`
	TableId int `db:"table_id" json:"table_id"`
	ColumnName string `db:"column_name json:"column_name"`
	ColumnNameLogical string `db:"column_name_logical" json:"column_name_logical"`
	DataTypeCls string `db:"data_type_cls" json:"data_type_cls"`
	Precision int `db:"precision" json:"precision"`
	Scale int `db:"scale" json:"scale"`
	PrimaryKeyFlg int `db:"primary_key_flg" json:"primary_key_flg"`
	NotNullFlg int `db:"not_null_flg" json:"not_null_flg"`
	UniqueFlg int `db:"unique_flg" json:"unique_flg"`
	DefaultValue string `db: "default_value" json:"default_value"`
	Remark string `db:"remark" json:"remark"`
	AlignSeq int `db:"align_seq" json:"align_seq"`
	DelFlg int `db:"del_flg" json:"del_flg"`
	CreateUserId int `db:"create_user_id" json:"create_user_id"`
	UpdateUserId int `db:"update_user_id" json:"update_user_id"`
	CreatedAt string `db:"created_at " json:"created_at "`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
}