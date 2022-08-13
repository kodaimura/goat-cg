package repository

import (
	"database/sql"

	"goat-cg/internal/shared/constant"
	"goat-cg/internal/core/db"
	"goat-cg/internal/model/entity"
)


type TableRepository interface {
	Select(id int) (entity.Table, error)
	Insert(t *entity.Table) error
    Update(id int, t *entity.Table) error
    Delete(id int) error

    SelectByProjectId(projectId int) ([]entity.Table, error)
    SelectByNameAndProjectId(name string, projectId int) (entity.Table, error)
    UpdateDelFlg(id, delFlg int) error
    UpdateLastLog(id int, msg string) error
}


type tableRepository struct {
	db *sql.DB
}


func NewTableRepository() TableRepository {
	db := db.GetDB()
	return &tableRepository{db}
}


func (rep *tableRepository) Select(tableId int) (entity.Table, error){
	
	var ret entity.Table
	err := rep.db.QueryRow(
		`SELECT 
			project_id,
			table_id,
			table_name,
			table_name_logical,
			create_user_id,
			update_user_id,
			del_flg
		 FROM 
		 	tables
		 WHERE 
		 	table_id = ?`,
		 tableId,
	).Scan(
		&ret.ProjectId, 
		&ret.TableId, 
		&ret.TableName,
		&ret.TableNameLogical,
		&ret.CreateUserId,
		&ret.UpdateUserId,
		&ret.DelFlg,
	)

	return ret, err
}


func (rep *tableRepository) Insert(t *entity.Table) error {
	_, err := rep.db.Exec(
		`INSERT INTO tables (
			project_id, 
			table_name,
			table_name_logical,
			create_user_id,
			update_user_id,
			last_log,
			del_flg
		 ) VALUES(?,?,?,?,?,?,?)`,
		t.ProjectId, 
		t.TableName,
		t.TableNameLogical,
		t.CreateUserId,
		t.CreateUserId,
		"create",
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
		 	update_user_id = ?,
		 	last_log = ?,
		 	del_flg = ?
		 WHERE table_id= ?`,
		t.TableName,
		t.TableNameLogical,
		t.UpdateUserId,
		"update",
		t.DelFlg,
		id,
	)
	return err
}


func (rep *tableRepository) Delete(id int) error {
	_, err := rep.db.Exec(
		`DELETE FROM tables WHERE table_id = ?`, 
		id,
	)

	return err
}


func (rep *tableRepository) SelectByNameAndProjectId(
	name string, projectId int,
) (entity.Table, error){
	
	var ret entity.Table
	err := rep.db.QueryRow(
		`SELECT 
			project_id,
			table_id,
			table_name,
			table_name_logical,
			create_user_id,
			update_user_id,
			del_flg
		 FROM 
		 	tables
		 WHERE project_id = ?
		   AND table_name = ?`,
		 projectId,
		 name,
	).Scan(
		&ret.ProjectId, 
		&ret.TableId, 
		&ret.TableName,
		&ret.TableNameLogical,
		&ret.CreateUserId,
		&ret.UpdateUserId,
		&ret.DelFlg,
	)

	return ret, err
}


func (rep *tableRepository) SelectByProjectId(projectId int) ([]entity.Table, error){
	
	var ret []entity.Table
	rows, err := rep.db.Query(
		`SELECT 
			table_id,
			table_name,
			table_name_logical,
			create_user_id,
			update_user_id,
			del_flg,
			last_log,
			create_at,
			update_at
		 FROM 
		 	tables
		 WHERE 
		 	project_id = ?`, 
		 projectId,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		t := entity.Table{}
		err = rows.Scan(
			&t.TableId, 
			&t.TableName,
			&t.TableNameLogical,
			&t.CreateUserId,
			&t.UpdateUserId,
			&t.DelFlg,
			&t.LastLog,
			&t.CreateAt,
			&t.UpdateAt,
		)
		if err != nil {
			break
		}
		ret = append(ret, t)
	}

	return ret, err
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


func (rep *tableRepository) UpdateLastLog(id int, msg string) error {
	_, err := rep.db.Exec(
		`UPDATE tables 
		 SET last_log = ?
		 WHERE table_id = ?`,
		msg,
		id,
	)
	return err
}