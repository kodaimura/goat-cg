package dto


import (
	"goat-cg/internal/model"
)

type ServInCreateColumn struct {
	ColumnId int
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
	AlignSeq int
	CreateUserId int
	UpdateUserId int
	DelFlg int
}


func (d ServInCreateColumn) ToColumn() model.Column {
	var c model.Column

	c.ColumnId = d.ColumnId
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
	c.AlignSeq = d.AlignSeq
	c.DelFlg = d.DelFlg
	c.CreateUserId = d.CreateUserId
	c.UpdateUserId = d.UpdateUserId

	return c
}

