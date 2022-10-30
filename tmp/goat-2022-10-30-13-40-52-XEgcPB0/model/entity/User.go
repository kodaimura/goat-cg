package entity


type User struct {
	UserId int `db:"user_id" json:"user_id"`
	TestColumn string `db:"test_column" json:"test_column"`
	CreateAt string `db:"create_at" json:"create_at"`
	UpdateAt string `db:"update_at" json:"update_at"`
}