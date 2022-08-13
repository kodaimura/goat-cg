package dto


import (
	"goat-cg/internal/model/entity"
)

type ServInCreateColumn struct {
	TableId int
	ColumnName string
	ColumnNameLogical string
	DataTypeCls string
	Precision int
	Scale int
	PrimaryKeyFlg int
	NotNullFlg int
	UniqueFlg int
	DefaultValue string
	Remark string
	CreateUserId int
	UpdateUserId int
	DelFlg int
}


func (d ServInCreateColumn) ToColumn() entity.Column {
	var c entity.Column

	c.TableId = d.TableId
	c.ColumnName = d.ColumnName
	c.ColumnNameLogical = d.ColumnNameLogical
	c.DataTypeCls = d.DataTypeCls
	c.Precision = d.Precision
	c.Scale = d.Scale
	c.PrimaryKeyFlg = d.PrimaryKeyFlg
	c.NotNullFlg = d.NotNullFlg
	c.UniqueFlg = d.UniqueFlg
	c.DefaultValue = d.DefaultValue
	c.Remark = d.Remark
	c.CreateUserId = d.CreateUserId
	c.UpdateUserId = d.UpdateUserId
	c.DelFlg = d.DelFlg

	return c
}

