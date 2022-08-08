package repository

import (
	"database/sql"

	"goat-cg/internal/core/db"
	"goat-cg/internal/model/entity"
)


type ColumnRepository interface {
	Select(id int) (entity.Column, error)
	Insert(c *entity.Column) error
	Update(id int, c *entity.Column) error
	Delete(id int) error

	SelectByNameAndTableId(name string, tableId int) (entity.Column, error)
	SelectByTableId(tableId int) ([]entity.Column, error)
	DeleteByTableId(tableId int) error
}


type columnRepository struct {
	db *sql.DB
}


func NewColumnRepository() ColumnRepository {
	db := db.GetDB()
	return &columnRepository{db}
}


func (rep *columnRepository) Select(id int) (entity.Column, error) {
	var ret entity.Column
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
			remark,
			create_user_id,
			update_user_id,
			del_flg,
			create_at,
			update_at
		 FROM
		 	columns
		 WHERE
		 	column_id = ?`,
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
		&ret.Remark,
		&ret.CreateUserId,
		&ret.CreateUserId,
		&ret.DelFlg,
		&ret.CreateAt,
		&ret.UpdateAt,
	)

	return ret, err
}


func (rep *columnRepository) Insert(c *entity.Column) error {
	_, err := rep.db.Exec(
		`INSERT INTO columns (
			table_id, 
			column_name,
			column_name_logical,
			data_type_cls,
			precision,
			scale,
			primary_key_flg,
			not_null_flg,
			unique_flg,
			remark,
			create_user_id,
			update_user_id,
			del_flg
		 ) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		c.TableId,
		c.ColumnName, 
		c.ColumnNameLogical,
		c.DataTypeCls,
		c.Precision,
		c.Scale,
		c.PrimaryKeyFlg,
		c.NotNullFlg,
		c.UniqueFlg,
		c.Remark,
		c.CreateUserId,
		c.CreateUserId,
		c.DelFlg,
	)
	return err
}


func (rep *columnRepository) Update(id int, c *entity.Column) error {
	_, err := rep.db.Exec(
		`UPDATE columns
		 SET 
			column_name = ?,
			column_name_logical = ?,
			data_type_cls = ?,
			precision = ?,
			scale = ?,
			primary_key_flg = ?,
			not_null_flg = ?,
			unique_flg = ?,
			remark = ?,
			update_user_id = ?,
			del_flg = ?
		 WHERE column_id = ?`,
		c.ColumnName, 
		c.ColumnNameLogical,
		c.DataTypeCls,
		c.Precision,
		c.Scale,
		c.PrimaryKeyFlg,
		c.NotNullFlg,
		c.UniqueFlg,
		c.Remark,
		c.UpdateUserId,
		c.DelFlg,
		id,
	)
	return err
}


func (rep *columnRepository) Delete(id int) error {
	_, err := rep.db.Exec(
		`DELETE FROM columns WHERE column_id = ?`, 
		id,
	)

	return err
}


func (rep *columnRepository) SelectByNameAndTableId(
	name string, 
	tableId int,
) (entity.Column, error) {
	var ret entity.Column
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
			remark,
			create_user_id,
			update_user_id,
			del_flg,
			create_at,
			update_at
		 FROM
		 	columns
		 WHERE
		 	table_id = ?
		 AND column_name = ?`,
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
		&ret.Remark,
		&ret.CreateUserId,
		&ret.CreateUserId,
		&ret.DelFlg,
		&ret.CreateAt,
		&ret.UpdateAt,
	)

	return ret, err
}


func (rep *columnRepository) SelectByTableId(tableId int) ([]entity.Column, error) {
	
	var ret []entity.Column
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
			remark,
			create_user_id,
			update_user_id,
			del_flg,
			create_at,
			update_at
		 FROM
		 	columns
		 WHERE
		 	table_id = ?`,
		 tableId,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		c := entity.Column{}
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
			&c.Remark,
			&c.CreateUserId,
			&c.CreateUserId,
			&c.DelFlg,
			&c.CreateAt,
			&c.UpdateAt,
		)
		if err != nil {
			break
		}
		ret = append(ret, c)
	}

	return ret, err
}


func (rep *columnRepository) DeleteByTableId(tableId int) error {
	_, err := rep.db.Exec(
		`DELETE FROM columns WHERE table_id = ?`, 
		tableId,
	)

	return err
}