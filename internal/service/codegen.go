package service

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"strings"
	"os/exec"

	"goat-cg/internal/shared/constant"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/core/utils"
	"goat-cg/internal/model"
	"goat-cg/internal/repository"
)


type CodegenService interface {
	GenerateGoat(rdbms string, tableIds []int) string
	GenerateDdl(rdbms string, tableIds []int) string
}


type codegenService struct {
	columnRepository repository.ColumnRepository
	tableRepository repository.TableRepository
}


func NewCodegenService() CodegenService {
	columnRepository := repository.NewColumnRepository()
	tableRepository := repository.NewTableRepository()
	return &codegenService{columnRepository, tableRepository}
}


// Generate ddl source and return file path.
// param rdbms: "sqlite3" or "postgresql" 
func (serv *codegenService) GenerateDdl(rdbms string, tableIds []int) string {
	path := "./tmp/ddl-" + time.Now().Format("2006-01-02-15-04-05") + 
		"-" + utils.RandomString(7) + ".sql"

	serv.generateDdlSource(rdbms, tableIds, path)

	return path
}


// Generate goat source and return zip path.
// param rdbms: "sqlite3" or "postgresql" 
func (serv *codegenService) GenerateGoat(rdbms string, tableIds []int) string {
	path := "./tmp/goat-" + time.Now().Format("2006-01-02-15-04-05") + 
		"-" + utils.RandomString(7)

	serv.generateGoatSource(rdbms, tableIds, path)

	if err := exec.Command("zip", "-rm", path + ".zip", path).Run(); err != nil {
		logger.Error(err.Error())
	}

	return path + ".zip"
}


//dataTypeMapSqlite3 map DataTypeCls and sqlite3 data types.
var dataTypeMapSqlite3 = map[string]string {
	constant.DATA_TYPE_CLS_SERIAL: "INTEGER PRIMARY KEY AUTOINCREMENT",
	constant.DATA_TYPE_CLS_TEXT: "TEXT",
	constant.DATA_TYPE_CLS_VARCHAR: "TEXT",
	constant.DATA_TYPE_CLS_CHAR: "TEXT",
	constant.DATA_TYPE_CLS_INTEGER: "INTEGER",
	constant.DATA_TYPE_CLS_NUMERIC: "NUMERIC",
	constant.DATA_TYPE_CLS_TIMESTAMP: "TEXT",
	constant.DATA_TYPE_CLS_DATE: "TEXT",
	constant.DATA_TYPE_CLS_BLOB: "BLOB",
}

//dataTypeMapPostgresql map DataTypeCls and postgresql data types.
var dataTypeMapPostgresql = map[string]string{
	constant.DATA_TYPE_CLS_SERIAL: "SERIAL PRIMARY KEY",
	constant.DATA_TYPE_CLS_TEXT: "TEXT",
	constant.DATA_TYPE_CLS_VARCHAR: "VARCHAR",
	constant.DATA_TYPE_CLS_CHAR: "CHAR",
	constant.DATA_TYPE_CLS_INTEGER: "INTEGER",
	constant.DATA_TYPE_CLS_NUMERIC: "NUMERIC",
	constant.DATA_TYPE_CLS_TIMESTAMP: "TIMESTAMP",
	constant.DATA_TYPE_CLS_DATE: "DATE",
	constant.DATA_TYPE_CLS_BLOB: "BLOB",
}

//dataTypeMapMysql map DataTypeCls and mysql data types.
var dataTypeMapMysql = map[string]string{
	constant.DATA_TYPE_CLS_SERIAL: "INT PRIMARY KEY AUTO_INCREMENT",
	constant.DATA_TYPE_CLS_TEXT: "TEXT",
	constant.DATA_TYPE_CLS_VARCHAR: "VARCHAR",
	constant.DATA_TYPE_CLS_CHAR: "CHAR",
	constant.DATA_TYPE_CLS_INTEGER: "INT",
	constant.DATA_TYPE_CLS_NUMERIC: "NUMERIC",
	constant.DATA_TYPE_CLS_TIMESTAMP: "DATETIME",
	constant.DATA_TYPE_CLS_DATE: "DATE",
	constant.DATA_TYPE_CLS_BLOB: "BLOB",
}

//dbDataTypeGoTypeMap map DataTypeCls and Golang types.
var dbDataTypeGoTypeMap = map[string]string{
	constant.DATA_TYPE_CLS_SERIAL: "int",
	constant.DATA_TYPE_CLS_TEXT: "string",
	constant.DATA_TYPE_CLS_VARCHAR: "string",
	constant.DATA_TYPE_CLS_CHAR: "string",
	constant.DATA_TYPE_CLS_INTEGER: "int",
	constant.DATA_TYPE_CLS_NUMERIC: "float64",
	constant.DATA_TYPE_CLS_TIMESTAMP: "string",
	constant.DATA_TYPE_CLS_DATE: "string",
	constant.DATA_TYPE_CLS_BLOB: "string",
}


func (serv *codegenService) writeFile(path, content string) {
	f, err := os.Create(path)
	defer f.Close()

	if err != nil {
		logger.Error(err.Error())
	}
	if _, err = f.Write([]byte(content)); err != nil {
		logger.Error(err.Error())
	}
}


func (serv *codegenService) extractPrimaryKeys(cols []model.Column) []model.Column {
	var ret []model.Column

	for _, col := range cols {
		if col.PrimaryKeyFlg == constant.FLG_ON {
			ret = append(ret, col)
		}
	}

	return ret
}


///////////////////////
/// GenerateDdl     ///
///////////////////////

// generateDdlSource generate ddl(create table) source.
// main processing of GenerateDdl.
func (serv *codegenService) generateDdlSource(rdbms string, tableIds []int, path string) {
	s := serv.generateDdlCreateTables(rdbms, tableIds) + "\n" +
		serv.generateDdlCreateTriggers(rdbms, tableIds)

	serv.writeFile(path, s)
}


func (serv *codegenService) generateDdlCreateTables(rdbms string, tableIds []int) string {
	s := ""
	for _, tid := range tableIds {
		s += serv.generateDdlCreateTable(rdbms, tid) + "\n\n"
	}

	return s
}


func (serv *codegenService) generateDdlCreateTable(rdbms string, tid int) string {
	s := ""
	table, err := serv.tableRepository.GetById(tid)

	if err != nil {
		logger.Error(err.Error())
		return s
	}

	s += "CREATE TABLE IF NOT EXISTS " + table.TableName + " (\n" +
		serv.generateDdlColumns(rdbms, tid) + "\n);"

	return s
}


func (serv *codegenService) generateDdlColumns(rdbms string, tid int) string {
	s := ""
	cols, err := serv.columnRepository.GetByTableId(tid)

	if err != nil {
		logger.Error(err.Error())
		return s
	}

	for _, col := range cols {
		s += serv.generateDdlColumn(rdbms, col)
	}
	s += serv.generateDdlCommonColumns(rdbms)
	s += serv.generateDdlPrymaryKey(rdbms, cols)

	return strings.TrimRight(s, ",\n")
}


func (serv *codegenService) generateDdlCommonColumns(rdbms string) string {
	s := ""
	if rdbms == "sqlite3" {
		s = "\tcreated_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),\n" + 
			"\tupdated_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),\n"

	} else if rdbms == "postgresql" {
		s = "\tcreated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" + 
			"\tupdated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,\n"

	} else if rdbms == "mysql" {
		s = "\tcreated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,\n" + 
			"\tupdated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,\n"
	}

	return s
}


func (serv *codegenService) generateDdlPrymaryKey(rdbms string, cols []model.Column) string {
	s := "" 
	pkcols := serv.extractPrimaryKeys(cols)

	for i, col := range pkcols {
		if col.DataTypeCls == constant.DATA_TYPE_CLS_SERIAL {
			return ""
		}

		if i == 0 {
			s += "\tPRIMARY KEY("
		} else {
			s += ", "
		}
		s += col.ColumnName
	}
	
	if s != "" {
		s += "),\n"
	}

	return s
}


func (serv *codegenService) generateDdlColumn(rdbms string, col model.Column) string {
	s := "\t" + col.ColumnName + " " + serv.generateDdlColumnDataType(rdbms, col)
	if cts := serv.generateDdlColumnConstraints(col); cts != "" {
		s += " " + cts
	}
	if dflt := serv.generateDdlColumnDefault(col); dflt != "" {
		s += " " + dflt
	}

	return s + ",\n"
}


func (serv *codegenService) generateDdlColumnConstraints(col model.Column) string {
	s := ""
	if col.NotNullFlg == constant.FLG_ON {
		s += "NOT NULL "
	}
	if col.UniqueFlg == constant.FLG_ON {
		s += "UNIQUE"
	} 
	
	return strings.TrimRight(s, " ")
}


func (serv *codegenService) generateDdlColumnDefault(col model.Column) string {
	s := ""
	if col.DefaultValue != "" {
		if col.DataTypeCls == constant.DATA_TYPE_CLS_NUMERIC ||
		col.DataTypeCls == constant.DATA_TYPE_CLS_INTEGER {
			s = "DEFAULT " + col.DefaultValue
		} else {
			s = "DEFAULT '" + col.DefaultValue + "'"
		}
	}

	return s
}


func (serv *codegenService) generateDdlColumnDataType(rdbms string, col model.Column) string {
	s := ""
	if rdbms == "sqlite3" {
		s = dataTypeMapSqlite3[col.DataTypeCls]

	} else if rdbms == "postgresql" {
		s = serv.generateDdlColumnDataTypePostgresql(col)
	
	} else if rdbms == "mysql" {
		s = serv.generateDdlColumnDataTypeMysql(col)
	}

	return s
}


func (serv *codegenService) generateDdlColumnDataTypePostgresql(col model.Column) string {
	s := dataTypeMapPostgresql[col.DataTypeCls]

	if col.DataTypeCls == constant.DATA_TYPE_CLS_VARCHAR || 
	col.DataTypeCls == constant.DATA_TYPE_CLS_CHAR {
		if col.Precision != 0 {
			s += "(" + strconv.Itoa(col.Precision) + ")"
		}	
	}
	if col.DataTypeCls == constant.DATA_TYPE_CLS_NUMERIC {
		if col.Precision != 0 {
			s += "(" + strconv.Itoa(col.Precision) + "," + strconv.Itoa(col.Scale) + ")"
		}
	}

	return s
}


func (serv *codegenService) generateDdlColumnDataTypeMysql(col model.Column) string {
	s := dataTypeMapMysql[col.DataTypeCls]

	if col.DataTypeCls == constant.DATA_TYPE_CLS_VARCHAR || 
	col.DataTypeCls == constant.DATA_TYPE_CLS_CHAR {
		if col.Precision != 0 {
			s += "(" + strconv.Itoa(col.Precision) + ")"
		}	
	}
	if col.DataTypeCls == constant.DATA_TYPE_CLS_NUMERIC {
		if col.Precision != 0 {
			s += "(" + strconv.Itoa(col.Precision) + "," + strconv.Itoa(col.Scale) + ")"
		}
	}

	return s
}


func (serv *codegenService) generateDdlCreateTriggers(rdbms string, tableIds []int) string {
	s := ""
	if rdbms == "postgresql" {
		s += "CREATE FUNCTION set_update_time() returns opaque AS '\n" + 
			"\tBEGIN\n\t\tnew.updated_at := ''now'';\n\t\treturn new;\n\tEND\n" + 
			"' language 'plpgsql';\n\n"
	}

	for _, tid := range tableIds {
		s += serv.generateDdlCreateTrigger(rdbms, tid) + "\n\n"
	}

	return s
}


func (serv *codegenService) generateDdlCreateTrigger(rdbms string, tid int) string {
	s := ""
	table, err := serv.tableRepository.GetById(tid)

	if err != nil {
		logger.Error(err.Error())
		return s
	}

	if rdbms == "sqlite3" {
		s = "CREATE TRIGGER IF NOT EXISTS trg_" + table.TableName + "_upd " + 
			"AFTER UPDATE ON " + table.TableName + "\n" +
			"BEGIN\n\tUPDATE " + table.TableName + "\n" +
			"\tSET updated_at = DATETIME('now', 'localtime')\n" + 
			"\tWHERE rowid == NEW.rowid;\nEND;"

	} else if rdbms == "postgresql" {
		s = "CREATE TRIGGER trg_" + table.TableName + "_upd " + 
			"BEFORE UPDATE ON " + table.TableName + " FOR EACH ROW\n" + 
			"\texecute procedure set_update_time();" 
	}

	return s
}


////////////////////////
/// GenerateGoat     ///
////////////////////////

// tableNameToFileName get file name from tabel name
// user => user.go / USER_TABLE => user_table.go
func (serv *codegenService) tableNameToFileName(tn string) string {
	n := strings.ToLower(tn)
	return n + ".go"
}


//xxx -> Xxx / xxx_yyy -> XxxYyy
func SnakeToPascal(snake string) string {
	ls := strings.Split(strings.ToLower(snake), "_")
	for i, s := range ls {
		ls[i] = strings.ToUpper(s[0:1]) + s[1:]
	}
	return strings.Join(ls, "")
}


//xxx -> xxx / xxx_yyy -> xxxYyy
func SnakeToCamel(snake string) string {
	ls := strings.Split(strings.ToLower(snake), "_")
	for i, s := range ls {
		if i != 0 {
			ls[i] = strings.ToUpper(s[0:1]) + s[1:]
		}
	}
	return strings.Join(ls, "")
}


//xxx -> x / xxx_yyy -> xy
func GetSnakeInitial(snake string) string {
	ls := strings.Split(strings.ToLower(snake), "_")
	ret := ""
	for _, s := range ls {
		ret = s[0:1]
	}
	return ret
}


func (serv *codegenService) generateGoatSource(rdbms string, tableIds []int, rootPath string) {
	path := rootPath + "/internal"
	if err := os.MkdirAll(path, 0777); err != nil {
		logger.Error(err.Error())
		return
	}

	var tables []model.Table

	for _, tid := range tableIds {
		table, err := serv.tableRepository.GetById(tid)
		if err != nil {
			logger.Error(err.Error())
			break
		}

		cols, err := serv.columnRepository.GetByTableId(tid)
		if err != nil {
			logger.Error(err.Error())
			break
		}

		serv.generateGoatModelSource(table.TableName, cols, mePath)
		serv.generateGoatDaoSource(rdbms, table.TableName, cols, mrPath)
		//serv.generateGoatController(tn, cPath)
		//serv.generateGoatService(tn, sPath)
	}
}


/////////////////////////////////
/// GenerateGoat (Model)      ///
/////////////////////////////////

func (serv *codegenService) generateGoatModelSource(tn string, cols []model.Column, path string) {
	path += "/" + serv.tableNameToFileName(tn)
	s := serv.generateGoatModel(tn, cols)
	serv.writeFile(path, s)
}


// generateGoatEntity is the main processing of generateGoatEntitySource 
func (serv *codegenService) generateGoatModel(tn string, cols []model.Column) string {
	s := "package model\n\n\n"

	s += fmt.Sprintf("type %s struct {\n", SnakeToPascal(tn))
	for _, col := range cols {
		s += fmt.Sprintf(
			"\t%s %s `db:\"%s\" json:\"%s\"`\n", 
			SnakeToPascal(col.ColumnName),
			dbDataTypeGoTypeMap[col.DataTypeCls],
			strings.ToLower(col.ColumnName),
			strings.ToLower(col.ColumnName),
		)
	}
	s += "\tCreatedAt string `db:\"created_at \" json:\"created_at \"`\n"
	s += "\tUpdatedAt string `db:\"updated_at\" json:\"updated_at\"`\n"

	return s + "}"
}


/////////////////////////////////////
/// GenerateGoat (Dao)            ///
/////////////////////////////////////


// generateGoatDaoSource generate dao program for goat.
func (serv *codegenService) generateGoatDaoSource(rdbms, tn string, cols []model.Column, path string) {
	path += "/" + serv.tableNameToFileName(tn)
	s := serv.generateGoatDao(rdbms, tn, cols)
	serv.writeFile(path, s)
}


// generateGoatDao is the main processing of generateGoatDaoSource 
func (serv *codegenService) generateGoatDao(rdbms, tn string, cols []model.Column) string {
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)

	s := "package dao\n\n\n" +
		"import (\n" + 
		"\t\"database/sql\"\n\n" +
		"\t\"xxxxx/internal/core/db\"\n" +
		"\t\"xxxxx/internal/model/entity\"\n)\n\n\n"

	s += serv.generateGoatDaoInterface(tn, cols)
	
	s += "\n\n\n"
	s += fmt.Sprintf("type %sDao struct {\n\tdb *sql.DB\n}", tnc)
	s += "\n\n\n"
	s += fmt.Sprintf("func New%sDao() *%sDao {\n", tnp, tnc)
	s += fmt.Sprintf("\tdb := db.GetDB()\n\treturn &%sDao{db}\n}", tnc)
	s += "\n\n\n"

	s += serv.generateGoatDaoInsert(rdbms, tn, cols) + "\n\n\n"
	s += serv.generateGoatDaoSelect(rdbms, tn, cols) + "\n\n\n"
	s += serv.generateGoatDaoUpdate(rdbms, tn, cols) + "\n\n\n"
	s += serv.generateGoatDaoDelete(rdbms, tn, cols) + "\n\n\n"
	s += serv.generateGoatDaoSelectAll(rdbms, tn, cols) + "\n"

	return s
}


// generateGoatDaoInterface
// return "type *Dao interface { ... }"
func (serv *codegenService) generateGoatDaoInterface(tn string, cols []model.Column) string {
	s := fmt.Sprintf("type %s interface {\n", SnakeToPascal(tn))

	s += "\t" + serv.generateGoatDaoInterfaceInsert(tn) + "\n"
	s += "\t" + serv.generateGoatDaoInterfaceSelect(tn) + "\n"
	s += "\t" + serv.generateGoatDaoInterfaceUpdate(tn) + "\n"
	s += "\t" + serv.generateGoatDaoInterfaceDelete(tn) + "\n"
	s += "\t" + serv.generateGoatDaoInterfaceSelectAll(tn) + "\n"

	return s + "}"
}


// generateGoatDaoInterfaceInsert
// return "Insert(e entity.Entity) error"
func (serv *codegenService) generateGoatDaoInterfaceInsert(
	tn string,
) string {
	return fmt.Sprintf(
		"Insert(%s *entity.%s) error",
		GetSnakeInitial(tn),
		SnakeToPascal(tn),
	)
}


// generateGoatDaoInterfaceSelect
// return "Select(e *entity.Entity) (entity.Entity, error)"
func (serv *codegenService) generateGoatDaoInterfaceSelect(tn string) string {
	en := SnakeToPascal(tn)

	return fmt.Sprintf(
		"Select(%s *entity.%s) (entity.%s, error)", 
		GetSnakeInitial(tn),
		en,
		en,
	)
}


// generateGoatDaoInterfaceUpdate
// return "Update(x *entity.Entity) error"
func (serv *codegenService) generateGoatDaoInterfaceUpdate(tn string) string {
	return fmt.Sprintf(
		"Update(%s *entity.%s) error",
		GetSnakeInitial(tn),
		SnakeToPascal(tn),
	)
}


// generateGoatDaoInterfaceDelete
// return "Delete(e *entity.Entity) error"
func (serv *codegenService) generateGoatDaoInterfaceDelete(tn string) string {
	return fmt.Sprintf(
		"Delete(%s *entity.%s) error",
		GetSnakeInitial(tn),
		SnakeToPascal(tn), 
	)
}


// generateGoatDaoInterfaceSelectAll
// return "SelectAll() ([]entity.Entity, error)"
func (serv *codegenService) generateGoatDaoInterfaceSelectAll(tn string) string {
	return fmt.Sprintf(
		"SelectAll() ([]entity.%s, error)",
		SnakeToPascal(tn),
	)
}


func (serv *codegenService) getBindVariable(rdbms string, n int) string {
	if rdbms == "postgresql" {
		return fmt.Sprintf("$%d", n)
	} else {
		return "?"
	}
}


// concatBindVariableWithCommas return ?,?,?,?,... or $1,$2,$3,$4,...
func (serv *codegenService) concatBindVariableWithCommas(rdbms string, bindCount int) string {
	var ls []string
	for i := 1; i <= bindCount; i++ {
		ls = append(ls, serv.getBindVariable(rdbms, i))
	}
	return strings.Join(ls, ",")
}


// generateGoatDaoInsert generate dao function 'Insert'.
// return "func (rep *userDao) Insert(u *entity.User) error {...}"
func (serv *codegenService) generateGoatDaoInsert(rdbms, tn string, cols []model.Column) string {
	s := fmt.Sprintf(
		"func (rep *%sDao) %s {\n", 
		SnakeToCamel(tn), 
		serv.generateGoatDaoInterfaceInsert(tn),
	)

	s += "\t_, err := rep.db.Exec(\n"
	s += fmt.Sprintf("\t\t`INSERT INTO %s (\n", tn)

	bindCount := 0
	for _, col := range cols {
		if col.DataTypeCls != constant.DATA_TYPE_CLS_SERIAL {
			bindCount += 1
			if bindCount == 1 {
				s += fmt.Sprintf("\t\t\t%s", col.ColumnName)
			} else {
				s += fmt.Sprintf("\n\t\t\t,%s", col.ColumnName)
			}
		}	
	}
	s += fmt.Sprintf("\n\t\t ) VALUES(%s)`,\n", serv.concatBindVariableWithCommas(rdbms, bindCount))

	tni := GetSnakeInitial(tn)
	for _, col := range cols {
		if col.DataTypeCls != constant.DATA_TYPE_CLS_SERIAL {
			s += fmt.Sprintf("\t\t%s.%s,\n", tni, SnakeToPascal(col.ColumnName))
		}
	}
	s += "\t)\n\n\treturn err\n}"

	return s
}


// generateGoatDaoSelect generate dao function 'Select'.
// return "func (rep *userDao) Select(u *entity.User) (entity.User, error) {...}"
func (serv *codegenService) generateGoatDaoSelect(rdbms, tn string, cols []model.Column) string {
	bindCount := 0

	s := fmt.Sprintf(
		"func (rep *%sDao) %s {\n",
		SnakeToCamel(tn), 
		serv.generateGoatDaoInterfaceSelect(tn),
	)

	s += fmt.Sprintf("\tvar ret entity.%s\n\n", SnakeToPascal(tn))
	s += "\terr := rep.db.QueryRow(\n\t\t`SELECT\n"
	s += serv.generateGoatDaoSelectSqlColumns(cols)
	s += fmt.Sprintf("\n\t\t FROM %s\n", tn)
	s += serv.generateGoatDaoSqlWhere(rdbms, cols, &bindCount)
	s += "`,\n"
	s += serv.generateGoatDaoSqlWhereBindVals(tn, cols)
	s += "\t).Scan(\n"
	s += serv.generateGoatDaoSelectScanVars(cols, "\t\t&ret.")
	s +=  "\t)\n\n\treturn ret, err\n}"

	return s
}


// generateGoatDaoUpdate generate dao function 'Update'.
// return "func (rep *userDao) Update(u *entity.User) error {...}"
func (serv *codegenService) generateGoatDaoUpdate(rdbms, tn string, cols []model.Column) string {
	bindCount := 0

	s := fmt.Sprintf(
		"func (rep *%sDao) %s {\n",
		SnakeToCamel(tn), 
		serv.generateGoatDaoInterfaceUpdate(tn),
	)
	
	s += "\t_, err := rep.db.Exec(\n"
	s += fmt.Sprintf("\t\t`UPDATE %s\n\t\t SET\n", tn)
	for _, col := range cols {
		if col.DataTypeCls != constant.DATA_TYPE_CLS_SERIAL && 
		col.PrimaryKeyFlg != constant.FLG_ON {
			bindCount += 1
			if bindCount == 1 {
				s += fmt.Sprintf("\t\t\t%s = %s", col.ColumnName, serv.getBindVariable(rdbms, bindCount))
			} else {
				s += fmt.Sprintf("\n\t\t\t,%s = %s", col.ColumnName, serv.getBindVariable(rdbms, bindCount))
			}
		}
	}

	s += "\n"
	s += serv.generateGoatDaoSqlWhere(rdbms, cols, &bindCount)

	s += "`,\n"
	tni := GetSnakeInitial(tn)
	for _, col := range cols {
		if col.DataTypeCls != constant.DATA_TYPE_CLS_SERIAL && 
		col.PrimaryKeyFlg != constant.FLG_ON {
			s += fmt.Sprintf("\t\t%s.%s,\n", tni, SnakeToPascal(col.ColumnName))
		}
	}
	s += serv.generateGoatDaoSqlWhereBindVals(tn, cols)
	s += "\t)\n\n\treturn err\n}"

	return s
}


// generateGoatDaoDelete generate dao function 'Delete'.
// return "func (rep *userDao) Delete(u *entity.User) error {...}"
func (serv *codegenService) generateGoatDaoDelete(rdbms, tn string, cols []model.Column) string {
	bindCount := 0

	s := fmt.Sprintf(
		"func (rep *%sDao) %s {\n",
		SnakeToCamel(tn), 
		serv.generateGoatDaoInterfaceDelete(tn),
	)
	s += "\t_, err := rep.db.Exec(\n"
	s += fmt.Sprintf("\t\t`DELETE FROM %s\n", tn)
	s += serv.generateGoatDaoSqlWhere(rdbms, cols, &bindCount)
	s += "`,\n"
	s += serv.generateGoatDaoSqlWhereBindVals(tn, cols)
	s += "\t)\n\n\treturn err\n}"

	return s
}


// generateGoatDaoSelectAll generate dao function 'SelectAll'.
// return "func (rep *userDao) SelectAll() ([]entity.User, error) {...}"
func (serv *codegenService) generateGoatDaoSelectAll(rdbms, tn string, cols []model.Column) string {
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	s := fmt.Sprintf(
		"func (rep *%sDao) %s {\n",
		 SnakeToCamel(tn),
		 serv.generateGoatDaoInterfaceSelectAll(tn),
	) 
	
	s += fmt.Sprintf("\tvar ret []entity.%s\n\n\trows, err := rep.db.Query(\n", tnp)
	s += "\t\t`SELECT\n"
	s += serv.generateGoatDaoSelectSqlColumns(cols)
	s += fmt.Sprintf("\n\t\t FROM %s`,\n\t)\n\n", tn)
	s += "\tif err != nil {\n\t\treturn nil, err\n\t}\n\n"
	s += "\tfor rows.Next() {\n"
	s += fmt.Sprintf("\t\t%s := entity.%s{}\n", tni, tnp)
	s += "\t\terr = rows.Scan(\n"
	s += serv.generateGoatDaoSelectScanVars(cols, fmt.Sprintf("\t\t\t&%s.", tni))
	s += "\t\t)\n\t\tif err != nil {\n\t\t\tbreak\n\t\t}\n"
	s += fmt.Sprintf("\t\tret = append(ret, %s)\n", tni)
	s +=  "\t}\n\n\treturn ret, err\n}"

	return s
}


func (serv *codegenService) generateGoatDaoSqlWhere(rdbms string, cols []model.Column, bindCount *int) string {
	s := "\t\t WHERE "

	isFirst := true
	for _, col := range cols {
		if col.PrimaryKeyFlg == constant.FLG_ON {
			*bindCount += 1
			if isFirst {
				s += fmt.Sprintf("%s = %s", col.ColumnName, serv.getBindVariable(rdbms, *bindCount))
				isFirst = false
			} else {
				s += fmt.Sprintf("\n\t\t   AND %s = %s", col.ColumnName, serv.getBindVariable(rdbms, *bindCount))
			}
		}
	}

	return s
}


func (serv *codegenService) generateGoatDaoSqlWhereBindVals(tn string, cols []model.Column) string {
	s := ""
	tni := GetSnakeInitial(tn)

	for _, col := range cols {
		if col.PrimaryKeyFlg == constant.FLG_ON {
			s += fmt.Sprintf("\t\t%s.%s,\n", tni, SnakeToPascal(col.ColumnName))
		}
	}

	return s
}


func (serv *codegenService) generateGoatDaoSelectSqlColumns(cols []model.Column) string {
	s := ""

	for i, col := range cols {
		if i == 0 {
			s += fmt.Sprintf("\t\t\t%s", col.ColumnName)
		} else {
			s += fmt.Sprintf("\n\t\t\t,%s", col.ColumnName)
		}
	}
	s += "\n\t\t\t,created_at"
	s += "\n\t\t\t,updated_at"

	return s
}


func (serv *codegenService) generateGoatDaoSelectScanVars(cols []model.Column, prefix string,) string {
	s := ""
	for _, col := range cols {
		s += fmt.Sprintf("%s%s,\n", prefix, SnakeToPascal(col.ColumnName))
	}

	s += fmt.Sprintf("%sCreatedAt,\n", prefix)
	s += fmt.Sprintf("%sUpdatedAt,\n", prefix)

	return s
}
