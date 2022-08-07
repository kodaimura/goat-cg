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

	SelectByNameAndTableId(name string, tableId int) (entity.Column, error)
	SelectByTableId(tableId int) ([]entity.Column, error)
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
			data_byte,
			decimal_byte,
			primary_key_flg,
			not_null_flg,
			unique_flg,
			remark,
			create_user_id,
			update_user_id,
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
		&ret.DataByte,
		&ret.DecimalByte,
		&ret.PrimaryKeyFlg,
		&ret.NotNullFlg,
		&ret.UniqueFlg,
		&ret.Remark,
		&ret.CreateUserId,
		&ret.CreateUserId,
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
			data_byte,
			decimal_byte,
			primary_key_flg,
			not_null_flg,
			unique_flg,
			remark,
			create_user_id,
			update_user_id
		 ) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)`,
		c.TableId,
		c.ColumnName, 
		c.ColumnNameLogical,
		c.DataTypeCls,
		c.DataByte,
		c.DecimalByte,
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
			decimal_byte = ?,
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
		c.DecimalByte,
		c.PrimaryKeyFlg,
		c.NotNullFlg,
		c.UniqueFlg,
		c.Remark,
		c.UpdateUserId,
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
			data_byte,
			decimal_byte,
			primary_key_flg,
			not_null_flg,
			unique_flg,
			remark,
			create_user_id,
			update_user_id,
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
		&ret.DataByte,
		&ret.DecimalByte,
		&ret.PrimaryKeyFlg,
		&ret.NotNullFlg,
		&ret.UniqueFlg,
		&ret.Remark,
		&ret.CreateUserId,
		&ret.CreateUserId,
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
			data_byte,
			decimal_byte,
			primary_key_flg,
			not_null_flg,
			unique_flg,
			remark,
			create_user_id,
			update_user_id,
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
			&c.DecimalByte,
			&c.DataByte,
			&c.PrimaryKeyFlg,
			&c.NotNullFlg,
			&c.UniqueFlg,
			&c.Remark,
			&c.CreateUserId,
			&c.CreateUserId,
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