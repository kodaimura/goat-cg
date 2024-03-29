package repository

import (
	"database/sql"

	"goat-cg/internal/shared/constant"
	"goat-cg/internal/core/db"
	"goat-cg/internal/model"
)


type TableRepository interface {
	GetById(id int) (model.Table, error)
	Insert(t *model.Table) (int, error)
	Update(t *model.Table) error
	Delete(t *model.Table) error
	DeleteTx(t *model.Table, tx *sql.Tx) error

	GetByProjectId(projectId int) ([]model.Table, error)
	GetByUniqueKey(name string, projectId int) (model.Table, error)
	DeleteByProjectIdTx(projectId int, tx *sql.Tx) error
}


type tableRepository struct {
	db *sql.DB
}


func NewTableRepository() TableRepository {
	db := db.GetDB()
	return &tableRepository{db}
}


func (rep *tableRepository) GetById(id int) (model.Table, error){
	var ret model.Table
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
		 id,
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


func (rep *tableRepository) Insert(t *model.Table) (int, error) {
	var tableId int

	err := rep.db.QueryRow(
		`INSERT INTO table_def (
			project_id, 
			table_name,
			table_name_logical,
			del_flg,
			create_user_id,
			update_user_id
		 ) VALUES(?,?,?,?,?,?)
		 RETURNING table_id`,
		t.ProjectId, 
		t.TableName,
		t.TableNameLogical,
		constant.FLG_OFF,
		t.CreateUserId,
		t.UpdateUserId,
	).Scan(
		&tableId,
	)

	return tableId, err
}

func (rep *tableRepository) Update(t *model.Table) error {
	_, err := rep.db.Exec(
		`UPDATE table_def
		 SET table_name = ?,
			 table_name_logical = ?,
			 del_flg = ?,
			 update_user_id = ?
		 WHERE table_id= ?`,
		t.TableName,
		t.TableNameLogical,
		t.DelFlg,
		t.UpdateUserId,
		t.TableId,
	)

	return err
}


func (rep *tableRepository) Delete(t *model.Table) error {
	_, err := rep.db.Exec(
		`DELETE FROM table_def WHERE table_id = ?`, 
		t.TableId,
	)

	return err
}


func (rep *tableRepository) DeleteTx(t *model.Table, tx *sql.Tx) error {
	_, err := tx.Exec(
		`DELETE FROM table_def WHERE table_id = ?`, 
		t.TableId,
	)

	return err
}


func (rep *tableRepository) GetByUniqueKey(name string, projectId int) (model.Table, error){
	var ret model.Table
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


func (rep *tableRepository) GetByProjectId(projectId int) ([]model.Table, error){
	rows, err := rep.db.Query(
		`SELECT 
			table_id,
			table_name,
			table_name_logical,
			del_flg,
			create_user_id,
			update_user_id,
			created_at ,
			updated_at
		 FROM 
			 table_def
		 WHERE 
			 project_id = ?`, 
		 projectId,
	)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	ret := []model.Table{}
	for rows.Next() {
		t := model.Table{}
		err = rows.Scan(
			&t.TableId, 
			&t.TableName,
			&t.TableNameLogical,
			&t.DelFlg,
			&t.CreateUserId,
			&t.UpdateUserId,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		ret = append(ret, t)
	}

	return ret, nil
}


func (rep *tableRepository) DeleteByProjectIdTx(projectId int, tx *sql.Tx) error {
	_, err := tx.Exec(
		`DELETE FROM table_def WHERE project_id = ?`, 
		projectId,
	)

	return err
}