package repository

import (
	"database/sql"

	"goat-cg/internal/shared/constant"
	"goat-cg/internal/core/db"
	"goat-cg/internal/model/entity"
)


type TableRepository interface {
	Insert(t *entity.Table) error
    Update(id int, t *entity.Table) error
    UpdateDelFlg(id, delFlg int) error
}


type tableRepository struct {
	db *sql.DB
}


func NewTableRepository() TableRepository {
	db := db.GetDB()
	return &tableRepository{db}
}


func (rep *tableRepository) Insert(t *entity.Table) error {
	_, err := rep.db.Exec(
		`INSERT INTO tables (
			project_id, 
			table_name,
			table_name_logical,
			create_user_id,
			update_user_id,
			del_flg
		 ) VALUES(?,?,?,?,?,?)`,
		t.ProjectId, 
		t.TableName,
		t.TableNameLogical,
		t.CreateUserId,
		t.CreateUserId,
		constant.FLG_OFF,
	)
	return err
}

func (rep *tableRepository) Update(id int, t *entity.Table) error {
	_, err := rep.db.Exec(
		`UPDATE tables
		 SET 
		 	table_name = ?,
		 	table_name_logical = ?,
		 	update_user_id= ?
		 WHERE table_id= ?`,
		t.TableName,
		t.TableNameLogical,
		t.UpdateUserId,
		id,
	)
	return err
}


func (rep *tableRepository) UpdateDelFlg(id, delFlg int) error {
	_, err := rep.db.Exec(
		`UPDATE tables 
		 SET del_flg = ?
		 WHERE table_id = ?`,
		delFlg,
		id,
	)
	return err
}