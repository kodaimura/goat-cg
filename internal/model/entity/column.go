package entity


type Column struct {
	ColumnId int `db:"column_id" json:"column_id"`
	TableId int `db:"table_id" json:"table_id"`
	ColumnName string `db:"column_name json:"column_name"`
	ColumnNameLogical string `db:"column_name_logical" json:"column_name_logical"`
	DataTypeCls string `db:"data_type_cls" json:"data_type_cls"`
	DataByte int `db:"data_byte" json:"data_byte"`
	DecimalByte int `db:"decimal_byte" json:"decimal_byte"`
	PrimaryKeyFlg int `db:"primary_key_flg" json:"primary_key_flg"`
	NotNullFlg int `db:"not_null_flg" json:"not_null_flg"`
	UniqueFlg int `db:"unique_flg" json:"unique_flg"`
	Remark string `db:"remark" json:"remark"`
	CreateUserId int `db:"create_user_id" json:"create_user_id"`
	UpdateUserId int `db:"update_user_id" json:"update_user_id"`
	CreateAt string `db:"create_at" json:"create_at"`
	UpdateAt string `db:"update_at" json:"update_at"`
}