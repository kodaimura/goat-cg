package dto


type QueOutTableLog struct {
	TableId int
	ProjectId int
	TableName string
	TableNameLogical string
	DelFlg int
	CreateUserId int
	CreateUserName string
	UpdateUserId int
	UpdateUserName string
	CreatedAt string
	UpdatedAt string
}