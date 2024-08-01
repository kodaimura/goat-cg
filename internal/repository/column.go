package repository

import (
	"database/sql"

	"goat-cg/internal/core/db"
	"goat-cg/internal/model"
)


type ColumnRepository interface {
	Get(c *model.Column) ([]model.Column, error)
	GetOne(c *model.Column) (model.Column, error)
	Insert(c *model.Column, tx *sql.Tx) error
	Update(c *model.Column, tx *sql.Tx) error
	Delete(c *model.Column, tx *sql.Tx) error
}


type columnRepository struct {
	db *sql.DB
}


func NewColumnRepository() ColumnRepository {
	db := db.GetDB()
	return &columnRepository{db}
}


func (rep *columnRepository) Get(c *model.Column) ([]model.Column, error) {
	where, binds := db.BuildWhereClause(c)
	query := "SELECT * FROM column_def " + where

	rows, err := rep.db.Query(query, binds...)
	defer rows.Close()

	if err != nil {
		return []model.Column{}, err
	}

	ret := []model.Column{}
	for rows.Next() {
		c := model.Column{}
		err = rows.Scan(
			&c.ColumnId,
			&c.TableId,
			&c.ColumnName, 
			&c.ColumnNameLogical,
			&c.DataTypeCls,
			&c.Precision,
			&c.Scale,
			&c.PrimaryKeyFlg,
			&c.NotNullFlg,
			&c.UniqueFlg,
			&c.DefaultValue,
			&c.Remark,
			&c.AlignSeq,
			&c.DelFlg,
			&c.CreateUserId,
			&c.UpdateUserId,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return []model.Column{}, err
		}
		ret = append(ret, c)
	}

	return ret, nil
}


func (rep *columnRepository) GetOne(c *model.Column) (model.Column, error) {
	var ret model.Column
	where, binds := db.BuildWhereClause(c)
	query := "SELECT * FROM column_def " + where

	err := rep.db.QueryRow(query, binds...).Scan(
		&ret.ColumnId,
		&ret.TableId,
		&ret.ColumnName, 
		&ret.ColumnNameLogical,
		&ret.DataTypeCls,
		&ret.Precision,
		&ret.Scale,
		&ret.PrimaryKeyFlg,
		&ret.NotNullFlg,
		&ret.UniqueFlg,
		&ret.DefaultValue,
		&ret.Remark,
		&ret.AlignSeq,
		&ret.DelFlg,
		&ret.CreateUserId,
		&ret.UpdateUserId,
		&ret.CreatedAt,
		&ret.UpdatedAt,
	)

	return ret, err
}


func (rep *columnRepository) Insert(c *model.Column, tx *sql.Tx) error {
	cmd := 
	`INSERT INTO column_def (
		table_id, 
		column_name,
		column_name_logical,
		data_type_cls,
		precision,
		scale,
		primary_key_flg,
		not_null_flg,
		unique_flg,
		default_value,
		remark,
		align_seq,
		del_flg,
		create_user_id,
		update_user_id
	 ) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`
	binds := []interface{}{
		c.TableId,
		c.ColumnName, 
		c.ColumnNameLogical,
		c.DataTypeCls,
		c.Precision,
		c.Scale,
		c.PrimaryKeyFlg,
		c.NotNullFlg,
		c.UniqueFlg,
		c.DefaultValue,
		c.Remark,
		c.AlignSeq,
		c.DelFlg,
		c.CreateUserId,
		c.CreateUserId,
	}

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
}


func (rep *columnRepository) Update(c *model.Column, tx *sql.Tx) error {
	cmd := 
	`UPDATE column_def
	 SET 
	    column_name = ?,
	    column_name_logical = ?,
	    data_type_cls = ?,
	    precision = ?,
	    scale = ?,
	    primary_key_flg = ?,
	    not_null_flg = ?,
	    unique_flg = ?,
	    default_value = ?,
	    remark = ?,
	    align_seq = ?,
	    del_flg = ?,
	    update_user_id = ?
	 WHERE column_id = ?`
	binds := []interface{}{
		c.ColumnName, 
		c.ColumnNameLogical,
		c.DataTypeCls,
		c.Precision,
		c.Scale,
		c.PrimaryKeyFlg,
		c.NotNullFlg,
		c.UniqueFlg,
		c.DefaultValue,
		c.Remark,
		c.AlignSeq,
		c.DelFlg,
		c.UpdateUserId,
		c.ColumnId,
	}
	
	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
}


func (rep *columnRepository) Delete(c *model.Column, tx *sql.Tx) error {
	where, binds := db.BuildWhereClause(c)
	cmd := "DELETE FROM column_def " + where

	var err error
	if tx != nil {
        _, err = tx.Exec(cmd, binds...)
    } else {
        _, err = rep.db.Exec(cmd, binds...)
    }
	
	return err
}