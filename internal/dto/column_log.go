package dto


type ColumnLog struct {
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
	DelFlg int
	CreateUserId int
	CreateUsername string
	UpdateUserId int
	UpdateUsername string
	CreatedAt string
	UpdatedAt string
}