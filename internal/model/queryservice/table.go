package queryservice

import (
	"database/sql"

	"goat-cg/internal/shared/constant"
	"goat-cg/internal/core/db"
	"goat-cg/internal/model/entity"
)


type TableQueryService interface {
	QueryTables(projectId int) ([]entity.Table, error)
	QueryTable(projectId int, tableId int) (entity.Table, error)
}

type tableQueryService struct {
	db *sql.DB
}

func NewTableQueryService() TableQueryService {
	db := db.GetDB()
	return &tableQueryService{db}
}


func (qs *tableQueryService) QueryTables(
	projectId int,
) ([]entity.Table, error){
	
	var ret []entity.Table
	rows, err := qs.db.Query(
		`SELECT 
			table_id,
			table_name,
			table_name_logical,
			create_user_id,
			update_user_id
		 FROM 
		 	tables
		 WHERE 
		 	project_id = ?
		 AND del_flg = ?`, 
		 projectId,
		 constant.FLG_OFF,
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
		)
		if err != nil {
			break
		}
		ret = append(ret, t)
	}

	return ret, err
}


func (qs *tableQueryService) QueryTable(
	projectId int,
	tableId int,
) (entity.Table, error){
	
	var ret entity.Table
	err := qs.db.QueryRow(
		`SELECT 
			table_id,
			table_name,
			table_name_logical,
			create_user_id,
			update_user_id
		 FROM 
		 	tables
		 WHERE 
		 	project_id = ?
		 	table_id = ?
		 AND del_flg = ?`, 
		 projectId,
		 tableId,
		 constant.FLG_OFF,
	).Scan(
		&ret.TableId, 
		&ret.TableName,
		&ret.TableNameLogical,
		&ret.CreateUserId,
		&ret.UpdateUserId,
	)

	return ret, err
}