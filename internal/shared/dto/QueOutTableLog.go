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
	CreateAt string
	UpdateAt string
}