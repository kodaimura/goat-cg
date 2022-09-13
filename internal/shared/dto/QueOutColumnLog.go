package dto


type QueOutColumnLog struct {
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
	CreateUserName string
	UpdateUserId int
	UpdateUserName string
	CreateAt string
	UpdateAt string
}