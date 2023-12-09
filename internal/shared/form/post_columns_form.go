package form

import (
	"goat-cg/internal/dto"
)


type PostColumnsForm struct {
	ColumnId int `form:"column_id"`
	ColumnName string `form:"column_name" binding:"required,max=50,min=1"`
	ColumnNameLogical string `form:"column_name_logical"`
	DataTypeCls string `form:"data_type_cls" binding:"required"`
	Precision int `form:"precision"`
	Scale int `form:"scale"`
	PrimaryKeyFlg int `form:"primary_key_flg"`
	NotNullFlg int `form:"not_null_flg"`
	UniqueFlg int `form:"unique_flg"`
	DefaultValue string `form:"default_value"`
	Remark string `form:"remark"`
	AlignSeq int `form:"align_seq"`
	DelFlg int `form:"del_flg"`
}


func (f PostColumnsForm) ToCreateColumn(
	tableId int, 
	userId int,
) dto.CreateColumn {
	var ret dto.CreateColumn

	ret.ColumnId = f.ColumnId
	ret.TableId = tableId
	ret.ColumnName = f.ColumnName
	ret.ColumnNameLogical = f.ColumnNameLogical
	ret.DataTypeCls = f.DataTypeCls
	ret.Precision = f.Precision
	ret.Scale = f.Scale
	ret.PrimaryKeyFlg = f.PrimaryKeyFlg
	ret.NotNullFlg = f.NotNullFlg
	ret.UniqueFlg = f.UniqueFlg
	ret.DefaultValue = f.DefaultValue
	ret.Remark = f.Remark
	ret.AlignSeq = f.AlignSeq
	ret.DelFlg = f.DelFlg
	ret.CreateUserId = userId
	ret.UpdateUserId = userId

	return ret
}

