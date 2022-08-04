package dto


import (
	"goat-cg/internal/model/entity"
)

type ServInCreateColumn struct {
	TableId int
	ColumnName string
	ColumnNameLogical string
	DataTypeCls int
	DataByte int
	PrimaryKeyFlg int
	NotNullFlg int
	UniqueFlg int
	Remark string
	CreateUserId int
	UpdateUserId int
}


func (d ServInCreateColumn) ToColumn() entity.Column {
	var c entity.Column

	c.TableId = d.TableId
	c.ColumnName = d.ColumnName
	c.ColumnNameLogical = d.ColumnNameLogical
	c.DataTypeCls = d.DataTypeCls
	c.DataByte = d.DataByte
	c.PrimaryKeyFlg = d.PrimaryKeyFlg
	c.NotNullFlg = d.NotNullFlg
	c.UniqueFlg = d.UniqueFlg
	c.Remark = d.Remark
	c.CreateUserId = d.CreateUserId
	c.UpdateUserId = d.UpdateUserId

	return c
}

