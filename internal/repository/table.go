package repository

import (
	"database/sql"

	"goat-cg/internal/core/db"
	"goat-cg/internal/model"
)


type TableRepository interface {
	Get(t *model.Table) ([]model.Table, error)
	GetOne(t *model.Table) (model.Table, error)
	Insert(t *model.Table, tx *sql.Tx) error
	Update(t *model.Table, tx *sql.Tx) error
	Delete(t *model.Table, tx *sql.Tx) error
}


type tableRepository struct {
	db *sql.DB
}


func NewTableRepository() TableRepository {
	db := db.GetDB()
	return &tableRepository{db}
}


func (rep *tableRepository) Get(t *model.Table) ([]model.Table, error){
	where, binds := db.BuildWhereClause(t)
	query := "SELECT * FROM table_def " + where
	rows, err := rep.db.Query(query, binds...)
	defer rows.Close()

	if err != nil {
		return []model.Table{}, err
	}

	ret := []model.Table{}
	for rows.Next() {
		t := model.Table{}
		err = rows.Scan(
			&t.TableId, 
			&t.ProjectId, 
			&t.TableName,
			&t.TableNameLogical,
			&t.CreateUserId,
			&t.UpdateUserId,
			&t.DelFlg,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return []model.Table{}, err
		}
		ret = append(ret, t)
	}

	return ret, nil
}


func (rep *tableRepository) GetOne(t *model.Table) (model.Table, error){
	var ret model.Table
	where, binds := db.BuildWhereClause(t)
	query := "SELECT * FROM table_def " + where

	err := rep.db.QueryRow(query, binds...).Scan(
		&ret.TableId, 
		&ret.ProjectId, 
		&ret.TableName,
		&ret.TableNameLogical,
		&ret.CreateUserId,
		&ret.UpdateUserId,
		&t.DelFlg,
		&ret.CreatedAt,
		&ret.UpdatedAt,
	)

	return ret, err
}


func (rep *tableRepository) Insert(t *model.Table, tx *sql.Tx) error {
	cmd := 
	`INSERT INTO table_def (
		project_id, 
		table_name,
		table_name_logical,
		del_flg,
		create_user_id,
		update_user_id
	 ) VALUES(?,?,?,?,?,?)`
	binds := []interface{}{
		t.ProjectId, 
		t.TableName,
		t.TableNameLogical,
		t.DelFlg,
		t.CreateUserId,
		t.UpdateUserId,
	}

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
}

func (rep *tableRepository) Update(t *model.Table, tx *sql.Tx) error {
	cmd := 
	`UPDATE table_def
	 SET table_name = ?,
		 table_name_logical = ?,
		 del_flg = ?,
		 update_user_id = ?
	 WHERE table_id= ?`
	binds := []interface{}{
		t.TableName,
		t.TableNameLogical,
		t.DelFlg,
		t.UpdateUserId,
		t.TableId,
	}
	
	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
}


func (rep *tableRepository) Delete(t *model.Table, tx *sql.Tx) error {
	where, binds := db.BuildWhereClause(t)
	cmd := "DELETE FROM table_def " + where

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
}