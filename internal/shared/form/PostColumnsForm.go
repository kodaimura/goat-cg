package form

import (
	"goat-cg/internal/shared/dto"
)


type PostColumnsForm struct {
	ColumnName string `form:"column_name" binding:"required,max=50,min=1"`
	ColumnNameLogical string `form:"column_name_logical" binding:"required,max=50,min=1"`
	DataTypeCls string `form:"data_type_cls" binding:"required"`
	DataByte int `form:"data_byte"`
	DecimalByte int `form:"decimal_byte"`
	PrimaryKeyFlg int `form:"primary_key_flg"`
	NotNullFlg int `form:"not_null_flg"`
	UniqueFlg int `form:"unique_flg"`
	Remark string `form:"remark"`
}


func (f PostColumnsForm) ToServInCreateColumn(
	tableId int, 
	userId int,
) dto.ServInCreateColumn {
	var ret dto.ServInCreateColumn

	ret.TableId = tableId
	ret.ColumnName = f.ColumnName
	ret.ColumnNameLogical = f.ColumnNameLogical
	ret.DataTypeCls = f.DataTypeCls
	ret.DataByte = f.DataByte
	ret.DecimalByte = f.DecimalByte
	ret.PrimaryKeyFlg = f.PrimaryKeyFlg
	ret.NotNullFlg = f.NotNullFlg
	ret.UniqueFlg = f.UniqueFlg
	ret.Remark = f.Remark
	ret.CreateUserId = userId
	ret.UpdateUserId = userId

	return ret
}

