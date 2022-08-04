package form

import (
	"goat-cg/internal/shared/dto"
)


type PostColumnsForm struct {
	ColumnName string `json:"column_name" binding:"required,max=50,min=1"`
	ColumnNameLogical string `json:"column_name_logical" binding:"required,max=50,min=1"`
	DataTypeCls int `json:"data_type_cls" binding:"required"`
	DataByte int `json:"data_byte" binding:"required"`
	PrimaryKeyFlg int `json:"primary_key_flg"`
	NotNullFlg int `json:"not_null_flg"`
	UniqueFlg int `json:"unique_flg"`
	Remark string `json:"remark"`
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
	ret.PrimaryKeyFlg = f.PrimaryKeyFlg
	ret.NotNullFlg = f.NotNullFlg
	ret.UniqueFlg = f.UniqueFlg
	ret.Remark = f.Remark
	ret.CreateUserId = userId
	ret.UpdateUserId = userId

	return ret
}

