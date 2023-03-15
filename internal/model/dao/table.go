package dao

import (
	"database/sql"

	"goat-cg/internal/shared/constant"
	"goat-cg/internal/core/db"
	"goat-cg/internal/model/entity"
)


type TableDao interface {
	Select(id int) (entity.Table, error)
	Insert(t *entity.Table) error
	Update(id int, t *entity.Table) error
	Delete(id int) error

	SelectByProjectId(projectId int) ([]entity.Table, error)
	SelectByNameAndProjectId(name string, projectId int) (entity.Table, error)
	UpdateDelFlg(id, delFlg int) error
}


type tableDao struct {
	db *sql.DB
}


func NewTableDao() TableDao {
	db := db.GetDB()
	return &tableDao{db}
}


func (rep *tableDao) Select(tableId int) (entity.Table, error){
	
	var ret entity.Table
	err := rep.db.QueryRow(
		`SELECT 
			project_id,
			table_id,
			table_name,
			table_name_logical,
			del_flg,
			create_user_id,
			update_user_id
		 FROM 
			 table_def
		 WHERE 
			 table_id = ?`,
		 tableId,
	).Scan(
		&ret.ProjectId, 
		&ret.TableId, 
		&ret.TableName,
		&ret.TableNameLogical,
		&ret.DelFlg,
		&ret.CreateUserId,
		&ret.UpdateUserId,
	)

	return ret, err
}


func (rep *tableDao) Insert(t *entity.Table) error {
	_, err := rep.db.Exec(
		`INSERT INTO table_def (
			project_id, 
			table_name,
			table_name_logical,
			del_flg,
			create_user_id,
			update_user_id
		 ) VALUES(?,?,?,?,?,?)`,
		t.ProjectId, 
		t.TableName,
		t.TableNameLogical,
		constant.FLG_OFF,
		t.CreateUserId,
		t.UpdateUserId,
	)
	return err
}

func (rep *tableDao) Update(id int, t *entity.Table) error {
	_, err := rep.db.Exec(
		`UPDATE table_def
		 SET 
			 table_name = ?,
			 table_name_logical = ?,
			 del_flg = ?,
			 update_user_id = ?
		 WHERE table_id= ?`,
		t.TableName,
		t.TableNameLogical,
		t.DelFlg,
		t.UpdateUserId,
		id,
	)
	return err
}


func (rep *tableDao) Delete(id int) error {
	_, err := rep.db.Exec(
		`DELETE FROM table_def WHERE table_id = ?`, 
		id,
	)

	return err
}


func (rep *tableDao) SelectByNameAndProjectId(
	name string, projectId int,
) (entity.Table, error){
	
	var ret entity.Table
	err := rep.db.QueryRow(
		`SELECT 
			project_id,
			table_id,
			table_name,
			table_name_logical,
			del_flg,
			create_user_id,
			update_user_id
		 FROM 
			 table_def
		 WHERE project_id = ?
		   AND table_name = ?`,
		 projectId,
		 name,
	).Scan(
		&ret.ProjectId, 
		&ret.TableId, 
		&ret.TableName,
		&ret.TableNameLogical,
		&ret.DelFlg,
		&ret.CreateUserId,
		&ret.UpdateUserId,
	)

	return ret, err
}


func (rep *tableDao) SelectByProjectId(projectId int) ([]entity.Table, error){
	
	var ret []entity.Table
	rows, err := rep.db.Query(
		`SELECT 
			table_id,
			table_name,
			table_name_logical,
			del_flg,
			create_user_id,
			update_user_id,
			create_at,
			update_at
		 FROM 
			 table_def
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
			&t.DelFlg,
			&t.CreateUserId,
			&t.UpdateUserId,
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


func (rep *tableDao) UpdateDelFlg(id, delFlg int) error {
	_, err := rep.db.Exec(
		`UPDATE table_def 
		 SET del_flg = ?
		 WHERE table_id = ?`,
		delFlg,
		id,
	)
	return err
}