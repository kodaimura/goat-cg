package dto


type TableLog struct {
	TableId int
	ProjectId int
	TableName string
	TableNameLogical string
	DelFlg int
	CreateUserId int
	CreateUsername string
	UpdateUserId int
	UpdateUsername string
	CreatedAt string
	UpdatedAt string
}