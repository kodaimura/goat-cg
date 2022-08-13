package service

import (
	"os"
	"strconv"
	"strings"
	"time"

	"goat-cg/pkg/utils"
	"goat-cg/internal/shared/constant"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/model/entity"
	"goat-cg/internal/model/repository"
)


type CodegenService interface {
	CodeGenerateGoat(projectId int, dbType string, tableIds []int) (string, int)
	CodeGenerateDdl(projectId int, dbType string, tableIds []int) (string, int)
}


type codegenService struct {
	cRep repository.ColumnRepository
	tRep repository.TableRepository
}


func NewCodegenService() CodegenService {
	cRep := repository.NewColumnRepository()
	tRep := repository.NewTableRepository()
	return &codegenService{cRep, tRep}
}


func (serv *codegenService) CodeGenerateGoat(
	projectId int, 
	dbType string,
	tableIds []int, 
) (string, int){
	return "test.txt", 1
}


func (serv *codegenService) CodeGenerateDdl(
	projectId int, 
	dbType string,
	tableIds []int, 
) (string, int){
	ddl := serv.generateDdlCreateTables(dbType, tableIds)
	ddl += serv.generateDdlCreateTriggers(dbType, tableIds)

    fn := "ddl-" + time.Now().Format("2006-01-02-15-04-05") + 
    "-" +utils.RandomString(7) + ".sql"

	f, err := os.Create("tmp/" + fn)
	if err != nil {
		logger.LogError(err.Error())
	}
	if _, err = f.Write([]byte(ddl)); err != nil {
		logger.LogError(err.Error())
	}

	return "/tmp/" + fn, 1
}


func (serv *codegenService) generateDdlCreateTables(dbType string, tableIds []int) string {
	ddl := ""

	for _, tid := range tableIds {
		ddl += serv.generateDdlTable(dbType, tid) + "\n\n"
	}

	return ddl
}


func (serv *codegenService) generateDdlTable(dbType string, tid int) string {
	table, err := serv.tRep.Select(tid)

	if err != nil {
		logger.LogError(err.Error())
	}

	ddl := "CREATE TABLE IF NOT EXISTS " + table.TableName + " (\n"
	ddl += serv.generateDdlColumns(dbType, tid)
	ddl += "\n);"

	return ddl
}


func (serv *codegenService) generateDdlColumns(dbType string, tid int) string {
	columns, err := serv.cRep.SelectByTableId(tid)

	if err != nil {
		logger.LogError(err.Error())
	}

	ddl := ""
	for _, col := range columns {
		ddl += "\t" + serv.generateDdlColumn(dbType, col) + ",\n"
	}

	ddl += serv.generateDdlCommonColumns(dbType)

	pk := serv.generateDdlPrymaryKey(dbType, columns)
	if pk == "" {
		ddl = strings.TrimRight(ddl, ",\n")
	} else {
		ddl += "\t" + pk
	}

	return ddl
}


func (serv *codegenService) generateDdlColumn(dbType string, col entity.Column) string {
	ddl := col.ColumnName
	ddl += " "
	ddl += serv.generateDdlColumnDataType(dbType, col)
	ddl += " "
	ddl += serv.generateDdlColumnConstraints(dbType, col)
	ddl = strings.TrimRight(ddl, " ")
	ddl += " "
	ddl += serv.generateDdlColumnDefault(dbType, col)
	return strings.TrimRight(ddl, " ")
}


var dataTypeMapSqlite3 = map[string]string {
	constant.DATA_TYPE_CLS_SERIAL: "INTEGER AUTOINCREMENT",
	constant.DATA_TYPE_CLS_TEXT: "TEXT",
	constant.DATA_TYPE_CLS_VARCHAR: "TEXT",
	constant.DATA_TYPE_CLS_CHAR: "TEXT",
	constant.DATA_TYPE_CLS_INTEGER: "INTEGER",
	constant.DATA_TYPE_CLS_NUMERIC: "NUMERIC",
	constant.DATA_TYPE_CLS_TIMESTAMP: "TEXT",
	constant.DATA_TYPE_CLS_DATE: "TEXT",
	constant.DATA_TYPE_CLS_BLOB: "BLOB",
}

var dataTypeMapPostgresql = map[string]string{
	constant.DATA_TYPE_CLS_SERIAL: "SERIAL",
	constant.DATA_TYPE_CLS_TEXT: "TEXT",
	constant.DATA_TYPE_CLS_VARCHAR: "VARCHAR",
	constant.DATA_TYPE_CLS_CHAR: "CHAR",
	constant.DATA_TYPE_CLS_INTEGER: "INTEGER",
	constant.DATA_TYPE_CLS_NUMERIC: "NUMERIC",
	constant.DATA_TYPE_CLS_TIMESTAMP: "TIMESTAMP",
	constant.DATA_TYPE_CLS_DATE: "DATE",
	constant.DATA_TYPE_CLS_BLOB: "BLOB",
}

func (serv *codegenService) generateDdlColumnDataType(dbType string, col entity.Column) string {
	ddl := ""

	if dbType == "sqlite3" {
		ddl = dataTypeMapSqlite3[col.DataTypeCls]

	} else if dbType == "postgresql" {
		ddl = dataTypeMapPostgresql[col.DataTypeCls]

		if col.DataTypeCls == constant.DATA_TYPE_CLS_VARCHAR || 
		col.DataTypeCls == constant.DATA_TYPE_CLS_CHAR {

			if col.Precision != 0 {
				ddl += "(" + strconv.Itoa(col.Precision) + ")"
			}	
		}

		if col.DataTypeCls == constant.DATA_TYPE_CLS_NUMERIC {
			if col.Precision != 0 {
				ddl += "(" + strconv.Itoa(col.Precision)
				ddl += "," + strconv.Itoa(col.Scale) + ")"
			}
		}
	}
	return ddl
}


func (serv *codegenService) generateDdlColumnConstraints(dbType string, col entity.Column) string {
	ddl := ""

	if col.NotNullFlg == constant.FLG_ON {
		ddl += "NOT NULL"
	}

	if col.UniqueFlg == constant.FLG_ON {
		if ddl != "" {
			ddl += " "
		} 
		ddl += "UNIQUE"
	} 
	
	return ddl
}


func (serv *codegenService) generateDdlPrymaryKey(dbType string, columns []entity.Column) string {
	ddl := "PRIMARY KEY("
	pkc := 0

	for _, col := range columns {
		if col.PrimaryKeyFlg == constant.FLG_ON {
			if pkc != 0 {
				ddl += ", "
			}
			ddl += col.ColumnName
			pkc++
		}
	}
	
	if pkc == 0 {
		ddl = ""
	} else {
		ddl += ")"
	}

	return ddl
}


func (serv *codegenService) generateDdlColumnDefault(dbType string, col entity.Column) string {
	if col.DefaultValue != "" {
		if col.DataTypeCls == constant.DATA_TYPE_CLS_NUMERIC ||
		col.DataTypeCls == constant.DATA_TYPE_CLS_INTEGER {
			return "DEFAULT " + col.DefaultValue
		} else {
			return "DEFAULT '" + col.DefaultValue + "'"
		}
	}

	return ""
}


func (serv *codegenService) generateDdlCommonColumns(dbType string) string {
	if dbType == "sqlite3" {
		return "\tcreate_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),\n" + 
		"\tupdate_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),\n"

	} else if dbType == "postgresql" {
		return "\tcreate_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" + 
		"\tupdate_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,\n"
	}

	return ""
}



func (serv *codegenService) generateDdlCreateTriggers(dbType string, tableIds []int) string {
	ddl := ""

	ddl += serv.generateDdlFunction(dbType)

	for _, tid := range tableIds {
		ddl += serv.generateDdlTrigger(dbType, tid) + "\n\n"
	}

	return ddl
}


func (serv *codegenService) generateDdlFunction(dbType string) string {
	if dbType == "postgresql" {
		return "CREATE FUNCTION set_update_time() returns opaque AS '\n" + 
		"\tBEGIN\n\t\tnew.updated_at := ''now'';\n\t\treturn new;\n\tEND\n" + 
		"' language 'plpgsql';\n\n"
	}

	return ""
}


func (serv *codegenService) generateDdlTrigger(dbType string, tid int) string {
	table, err := serv.tRep.Select(tid)

	if err != nil {
		logger.LogError(err.Error())
	}

	if dbType == "sqlite3" {
		return "CREATE TRIGGER IF NOT EXISTS " + table.TableName + "_update_trg " + 
		"AFTER UPDATE ON " + table.TableName + 
		"\nBEGIN\n\tUPDATE " + table.TableName + 
		"\n\tSET update_at = DATETIME('now', 'localtime')\n" + 
		"\n\tWHERE rowid == NEW.rowid;\nEND;"

	} else if dbType == "postgresql" {
		return "CREATE TRIGGER " + table.TableName + "_update_trg " + 
		"AFTER UPDATE ON " + table.TableName + " FOR EACH ROW" + 
		"\n\texecute procedure set_update_time()" 
	}

	return ""
}
