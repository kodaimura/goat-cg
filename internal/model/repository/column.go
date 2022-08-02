package repository

import (
	"database/sql"

	"goat-cg/internal/core/db"
	"goat-cg/internal/model/entity"
)


type ColumnRepository interface {
	Insert(c *entity.Column) error
	Update(id int, c *entity.Column) error
}


type columnRepository struct {
	db *sql.DB
}


func NewColumnRepository() ColumnRepository {
	db := db.GetDB()
	return &columnRepository{db}
}


func (rep *columnRepository) Insert(c *entity.Column) error {
	_, err := rep.db.Exec(
		`INSERT INTO columns (
			table_id, 
			column_name,
			column_name_logical,
			data_type_cls,
			data_byte,
			primary_key_flg,
			not_null_flg,
			unique_flg,
			remark,
			create_user_id,
			update_user_id
		 ) VALUES(?,?,?,?,?,?)`,
		c.TableId,
		c.ColumnName, 
		c.ColumnNameLogical,
		c.DataTypeCls,
		c.DataByte,
		c.PrimaryKeyFlg,
		c.NotNullFlg,
		c.UniqueFlg,
		c.Remark,
		c.CreateUserId,
		c.CreateUserId,
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
			data_byte = ?,
			primary_key_flg = ?,
			not_null_flg = ?,
			unique_flg = ?,
			remark = ?,
			update_user_id = ?
		 WHERE column_id = ?`,
		c.ColumnName, 
		c.ColumnNameLogical,
		c.DataTypeCls,
		c.DataByte,
		c.PrimaryKeyFlg,
		c.NotNullFlg,
		c.UniqueFlg,
		c.Remark,
		c.UpdateUserId,
		c.ColumnId,
	)
	return err
}