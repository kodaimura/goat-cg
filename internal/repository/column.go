package repository

import (
	"database/sql"

	"goat-cg/internal/core/db"
	"goat-cg/internal/model"
)


type ColumnRepository interface {
	GetById(id int) (model.Column, error)
	Insert(c *model.Column) error
	Update(c *model.Column) error
	Delete(c *model.Column) error
	DeleteTx(c *model.Column, tx *sql.Tx) error

	GetByUniqueKey(name string, tableId int) (model.Column, error)
	GetByTableId(tableId int) ([]model.Column, error)
	DeleteByTableIdTx(tableId int, tx *sql.Tx) error
}


type columnRepository struct {
	db *sql.DB
}


func NewColumnRepository() ColumnRepository {
	db := db.GetDB()
	return &columnRepository{db}
}


func (rep *columnRepository) GetById(id int) (model.Column, error) {
	var ret model.Column
	err := rep.db.QueryRow(
		`SELECT 
			column_id,
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
			update_user_id,
			created_at ,
			updated_at
		 FROM
			 column_def
		 WHERE
			 column_id = ?
		 ORDER BY align_seq`,
		 id,
	).Scan(
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


func (rep *columnRepository) Insert(c *model.Column) error {
	_, err := rep.db.Exec(
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
		 ) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
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
	)

	return err
}


func (rep *columnRepository) Update(c *model.Column) error {
	_, err := rep.db.Exec(
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
		 WHERE column_id = ?`,
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
	)

	return err
}


func (rep *columnRepository) Delete(c *model.Column) error {
	_, err := rep.db.Exec(
		`DELETE FROM column_def WHERE column_id = ?`, 
		c.ColumnId,
	)

	return err
}


func (rep *columnRepository) DeleteTx(c *model.Column, tx *sql.Tx) error {
	_, err := tx.Exec(
		`DELETE FROM column_def WHERE column_id = ?`, 
		c.ColumnId,
	)

	return err
}


func (rep *columnRepository) GetByUniqueKey(name string, tableId int) (model.Column, error) {
	var ret model.Column
	err := rep.db.QueryRow(
		`SELECT 
			column_id,
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
			update_user_id,
			created_at ,
			updated_at
		 FROM
			 column_def
		 WHERE
			 table_id = ?
		 AND column_name = ?
		 ORDER BY align_seq`,
		 tableId,
		 name,
	).Scan(
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


func (rep *columnRepository) GetByTableId(tableId int) ([]model.Column, error) {
	rows, err := rep.db.Query(
		`SELECT 
			column_id,
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
			update_user_id,
			created_at ,
			updated_at
		 FROM
			 column_def
		 WHERE
			 table_id = ?
		 ORDER BY align_seq`,
		 tableId,
	)
	defer rows.Close()

	if err != nil {
		return nil, err
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
			return nil, err
		}
		ret = append(ret, c)
	}

	return ret, nil
}


func (rep *columnRepository) DeleteByTableIdTx(tableId int, tx *sql.Tx) error {
	_, err := tx.Exec(
		`DELETE FROM column_def WHERE table_id = ?`, 
		tableId,
	)

	return err
}