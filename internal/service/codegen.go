package service

import (
	"os"
	"fmt"
	"time"
	"strconv"
	"strings"
	"os/exec"

	"goat-cg/pkg/utils"
	"goat-cg/internal/shared/constant"
	"goat-cg/internal/core/logger"
	"goat-cg/internal/model"
	"goat-cg/internal/repository"
)


type CodegenService interface {
	CodeGenerateGoat(rdbms string, tableIds []int) string
	CodeGenerateDdl(rdbms string, tableIds []int) string
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


// CodeGenerateDdl generate ddl(create table) source 
// and return file path.
// param rdbms: "sqlite3" or "postgresql" 
func (serv *codegenService) CodeGenerateDdl(rdbms string, tableIds []int) string {
	path := "./tmp/ddl-" + time.Now().Format("2006-01-02-15-04-05") + 
		"-" + utils.RandomString(7) + ".sql"

	serv.cgDdlSource(rdbms, tableIds, path)

	return path
}


// CodeGenerateGoat generate programs(entity, dao for goat) 
// and return zip path .
// param rdbms: "sqlite3" or "postgresql" 
func (serv *codegenService) CodeGenerateGoat(rdbms string, tableIds []int) string {
	path := "./tmp/goat-" + time.Now().Format("2006-01-02-15-04-05") + 
		"-" + utils.RandomString(7)

	serv.cgGoatSource(rdbms, tableIds, path)

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
/// CodeGenerateDdl ///
///////////////////////

// cgDdlSource generate ddl(create table) source.
// main processing of CodeGenerateDdl.
func (serv *codegenService) cgDdlSource(rdbms string, tableIds []int, path string) {
	s := serv.cgDdlCreateTables(rdbms, tableIds) + "\n" +
		serv.cgDdlCreateTriggers(rdbms, tableIds)

	serv.writeFile(path, s)
}


func (serv *codegenService) cgDdlCreateTables(rdbms string, tableIds []int) string {
	s := ""
	for _, tid := range tableIds {
		s += serv.cgDdlCreateTable(rdbms, tid) + "\n\n"
	}

	return s
}


func (serv *codegenService) cgDdlCreateTable(rdbms string, tid int) string {
	s := ""
	table, err := serv.tableRepository.GetById(tid)

	if err != nil {
		logger.Error(err.Error())
		return s
	}

	s += "CREATE TABLE IF NOT EXISTS " + table.TableName + " (\n" +
		serv.cgDdlColumns(rdbms, tid) + "\n);"

	return s
}


func (serv *codegenService) cgDdlColumns(rdbms string, tid int) string {
	s := ""
	cols, err := serv.columnRepository.GetByTableId(tid)

	if err != nil {
		logger.Error(err.Error())
		return s
	}

	for _, col := range cols {
		s += serv.cgDdlColumn(rdbms, col)
	}
	s += serv.cgDdlCommonColumns(rdbms)
	s += serv.cgDdlPrymaryKey(rdbms, cols)

	return strings.TrimRight(s, ",\n")
}


func (serv *codegenService) cgDdlCommonColumns(rdbms string) string {
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


func (serv *codegenService) cgDdlPrymaryKey(rdbms string, cols []model.Column) string {
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


func (serv *codegenService) cgDdlColumn(rdbms string, col model.Column) string {
	s := "\t" + col.ColumnName + " " + serv.cgDdlColumnDataType(rdbms, col)
	if cts := serv.cgDdlColumnConstraints(col); cts != "" {
		s += " " + cts
	}
	if dflt := serv.cgDdlColumnDefault(col); dflt != "" {
		s += " " + dflt
	}

	return s + ",\n"
}


func (serv *codegenService) cgDdlColumnConstraints(col model.Column) string {
	s := ""
	if col.NotNullFlg == constant.FLG_ON {
		s += "NOT NULL "
	}
	if col.UniqueFlg == constant.FLG_ON {
		s += "UNIQUE"
	} 
	
	return strings.TrimRight(s, " ")
}


func (serv *codegenService) cgDdlColumnDefault(col model.Column) string {
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


func (serv *codegenService) cgDdlColumnDataType(rdbms string, col model.Column) string {
	s := ""
	if rdbms == "sqlite3" {
		s = dataTypeMapSqlite3[col.DataTypeCls]

	} else if rdbms == "postgresql" {
		s = serv.cgDdlColumnDataTypePostgresql(col)
	
	} else if rdbms == "mysql" {
		s = serv.cgDdlColumnDataTypeMysql(col)
	}

	return s
}


func (serv *codegenService) cgDdlColumnDataTypePostgresql(col model.Column) string {
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


func (serv *codegenService) cgDdlColumnDataTypeMysql(col model.Column) string {
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


func (serv *codegenService) cgDdlCreateTriggers(rdbms string, tableIds []int) string {
	s := ""
	if rdbms == "postgresql" {
		s += "CREATE FUNCTION set_update_time() returns opaque AS '\n" + 
			"\tBEGIN\n\t\tnew.updated_at := ''now'';\n\t\treturn new;\n\tEND\n" + 
			"' language 'plpgsql';\n\n"
	}

	for _, tid := range tableIds {
		s += serv.cgDdlCreateTrigger(rdbms, tid) + "\n\n"
	}

	return s
}


func (serv *codegenService) cgDdlCreateTrigger(rdbms string, tid int) string {
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
/// CodeGenerateGoat ///
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


// cgGoatSource generate programs(entity, dao for goat).
// main processing of CodeGenerateGoat.
func (serv *codegenService) cgGoatSource(rdbms string, tableIds []int, path string) {
	mePath := path + "/model/entity"
	if err := os.MkdirAll(mePath, 0777); err != nil {
		logger.Error(err.Error())
		return
	}

	mrPath := path + "/model/dao"
	if err := os.MkdirAll(mrPath, 0777); err != nil {
		logger.Error(err.Error())
		return
	}

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

		serv.cgGoatEntitySource(table.TableName, cols, mePath)
		serv.cgGoatDaoSource(rdbms, table.TableName, cols, mrPath)
		//serv.cgGoatController(tn, cPath)
		//serv.cgGoatService(tn, sPath)
	}	
}


/////////////////////////////////
/// CodeGenerateGoat (Entity) ///
/////////////////////////////////

// cgGoatEntitySource generate entity program for goat.
func (serv *codegenService) cgGoatEntitySource(tn string, cols []model.Column, path string) {
	path += "/" + serv.tableNameToFileName(tn)
	s := serv.cgGoatEntity(tn, cols)
	serv.writeFile(path, s)
}


// cgGoatEntity is the main processing of cgGoatEntitySource 
func (serv *codegenService) cgGoatEntity(tn string, cols []model.Column,) string {
	s := "package entity\n\n\n"

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
/// CodeGenerateGoat (Dao) ///
/////////////////////////////////////


// cgGoatDaoSource generate dao program for goat.
func (serv *codegenService) cgGoatDaoSource(rdbms, tn string, cols []model.Column, path string) {
	path += "/" + serv.tableNameToFileName(tn)
	s := serv.cgGoatDao(rdbms, tn, cols)
	serv.writeFile(path, s)
}


// cgGoatDao is the main processing of cgGoatDaoSource 
func (serv *codegenService) cgGoatDao(rdbms, tn string, cols []model.Column) string {
	tnc := SnakeToCamel(tn)
	tnp := SnakeToPascal(tn)

	s := "package dao\n\n\n" +
		"import (\n" + 
		"\t\"database/sql\"\n\n" +
		"\t\"xxxxx/internal/core/db\"\n" +
		"\t\"xxxxx/internal/model/entity\"\n)\n\n\n"

	s += serv.cgGoatDaoInterface(tn, cols)
	
	s += "\n\n\n"
	s += fmt.Sprintf("type %sDao struct {\n\tdb *sql.DB\n}", tnc)
	s += "\n\n\n"
	s += fmt.Sprintf("func New%sDao() *%sDao {\n", tnp, tnc)
	s += fmt.Sprintf("\tdb := db.GetDB()\n\treturn &%sDao{db}\n}", tnc)
	s += "\n\n\n"

	s += serv.cgGoatDaoInsert(rdbms, tn, cols) + "\n\n\n"
	s += serv.cgGoatDaoSelect(rdbms, tn, cols) + "\n\n\n"
	s += serv.cgGoatDaoUpdate(rdbms, tn, cols) + "\n\n\n"
	s += serv.cgGoatDaoDelete(rdbms, tn, cols) + "\n\n\n"
	s += serv.cgGoatDaoSelectAll(rdbms, tn, cols) + "\n"

	return s
}


// cgGoatDaoInterface
// return "type *Dao interface { ... }"
func (serv *codegenService) cgGoatDaoInterface(tn string, cols []model.Column) string {
	s := fmt.Sprintf("type %s interface {\n", SnakeToPascal(tn))

	s += "\t" + serv.cgGoatDaoInterfaceInsert(tn) + "\n"
	s += "\t" + serv.cgGoatDaoInterfaceSelect(tn) + "\n"
	s += "\t" + serv.cgGoatDaoInterfaceUpdate(tn) + "\n"
	s += "\t" + serv.cgGoatDaoInterfaceDelete(tn) + "\n"
	s += "\t" + serv.cgGoatDaoInterfaceSelectAll(tn) + "\n"

	return s + "}"
}


// cgGoatDaoInterfaceInsert
// return "Insert(e entity.Entity) error"
func (serv *codegenService) cgGoatDaoInterfaceInsert(
	tn string,
) string {
	return fmt.Sprintf(
		"Insert(%s *entity.%s) error",
		GetSnakeInitial(tn),
		SnakeToPascal(tn),
	)
}


// cgGoatDaoInterfaceSelect
// return "Select(e *entity.Entity) (entity.Entity, error)"
func (serv *codegenService) cgGoatDaoInterfaceSelect(tn string) string {
	en := SnakeToPascal(tn)

	return fmt.Sprintf(
		"Select(%s *entity.%s) (entity.%s, error)", 
		GetSnakeInitial(tn),
		en,
		en,
	)
}


// cgGoatDaoInterfaceUpdate
// return "Update(x *entity.Entity) error"
func (serv *codegenService) cgGoatDaoInterfaceUpdate(tn string) string {
	return fmt.Sprintf(
		"Update(%s *entity.%s) error",
		GetSnakeInitial(tn),
		SnakeToPascal(tn),
	)
}


// cgGoatDaoInterfaceDelete
// return "Delete(e *entity.Entity) error"
func (serv *codegenService) cgGoatDaoInterfaceDelete(tn string) string {
	return fmt.Sprintf(
		"Delete(%s *entity.%s) error",
		GetSnakeInitial(tn),
		SnakeToPascal(tn), 
	)
}


// cgGoatDaoInterfaceSelectAll
// return "SelectAll() ([]entity.Entity, error)"
func (serv *codegenService) cgGoatDaoInterfaceSelectAll(tn string) string {
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


// cgGoatDaoInsert generate dao function 'Insert'.
// return "func (rep *userDao) Insert(u *entity.User) error {...}"
func (serv *codegenService) cgGoatDaoInsert(rdbms, tn string, cols []model.Column) string {
	s := fmt.Sprintf(
		"func (rep *%sDao) %s {\n", 
		SnakeToCamel(tn), 
		serv.cgGoatDaoInterfaceInsert(tn),
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


// cgGoatDaoSelect generate dao function 'Select'.
// return "func (rep *userDao) Select(u *entity.User) (entity.User, error) {...}"
func (serv *codegenService) cgGoatDaoSelect(rdbms, tn string, cols []model.Column) string {
	bindCount := 0

	s := fmt.Sprintf(
		"func (rep *%sDao) %s {\n",
		SnakeToCamel(tn), 
		serv.cgGoatDaoInterfaceSelect(tn),
	)

	s += fmt.Sprintf("\tvar ret entity.%s\n\n", SnakeToPascal(tn))
	s += "\terr := rep.db.QueryRow(\n\t\t`SELECT\n"
	s += serv.cgGoatDaoSelectSqlColumns(cols)
	s += fmt.Sprintf("\n\t\t FROM %s\n", tn)
	s += serv.cgGoatDaoSqlWhere(rdbms, cols, &bindCount)
	s += "`,\n"
	s += serv.cgGoatDaoSqlWhereBindVals(tn, cols)
	s += "\t).Scan(\n"
	s += serv.cgGoatDaoSelectScanVars(cols, "\t\t&ret.")
	s +=  "\t)\n\n\treturn ret, err\n}"

	return s
}


// cgGoatDaoUpdate generate dao function 'Update'.
// return "func (rep *userDao) Update(u *entity.User) error {...}"
func (serv *codegenService) cgGoatDaoUpdate(rdbms, tn string, cols []model.Column) string {
	bindCount := 0

	s := fmt.Sprintf(
		"func (rep *%sDao) %s {\n",
		SnakeToCamel(tn), 
		serv.cgGoatDaoInterfaceUpdate(tn),
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
	s += serv.cgGoatDaoSqlWhere(rdbms, cols, &bindCount)

	s += "`,\n"
	tni := GetSnakeInitial(tn)
	for _, col := range cols {
		if col.DataTypeCls != constant.DATA_TYPE_CLS_SERIAL && 
		col.PrimaryKeyFlg != constant.FLG_ON {
			s += fmt.Sprintf("\t\t%s.%s,\n", tni, SnakeToPascal(col.ColumnName))
		}
	}
	s += serv.cgGoatDaoSqlWhereBindVals(tn, cols)
	s += "\t)\n\n\treturn err\n}"

	return s
}


// cgGoatDaoDelete generate dao function 'Delete'.
// return "func (rep *userDao) Delete(u *entity.User) error {...}"
func (serv *codegenService) cgGoatDaoDelete(rdbms, tn string, cols []model.Column) string {
	bindCount := 0

	s := fmt.Sprintf(
		"func (rep *%sDao) %s {\n",
		SnakeToCamel(tn), 
		serv.cgGoatDaoInterfaceDelete(tn),
	)
	s += "\t_, err := rep.db.Exec(\n"
	s += fmt.Sprintf("\t\t`DELETE FROM %s\n", tn)
	s += serv.cgGoatDaoSqlWhere(rdbms, cols, &bindCount)
	s += "`,\n"
	s += serv.cgGoatDaoSqlWhereBindVals(tn, cols)
	s += "\t)\n\n\treturn err\n}"

	return s
}


// cgGoatDaoSelectAll generate dao function 'SelectAll'.
// return "func (rep *userDao) SelectAll() ([]entity.User, error) {...}"
func (serv *codegenService) cgGoatDaoSelectAll(rdbms, tn string, cols []model.Column) string {
	tnp := SnakeToPascal(tn)
	tni := GetSnakeInitial(tn)

	s := fmt.Sprintf(
		"func (rep *%sDao) %s {\n",
		 SnakeToCamel(tn),
		 serv.cgGoatDaoInterfaceSelectAll(tn),
	) 
	
	s += fmt.Sprintf("\tvar ret []entity.%s\n\n\trows, err := rep.db.Query(\n", tnp)
	s += "\t\t`SELECT\n"
	s += serv.cgGoatDaoSelectSqlColumns(cols)
	s += fmt.Sprintf("\n\t\t FROM %s`,\n\t)\n\n", tn)
	s += "\tif err != nil {\n\t\treturn nil, err\n\t}\n\n"
	s += "\tfor rows.Next() {\n"
	s += fmt.Sprintf("\t\t%s := entity.%s{}\n", tni, tnp)
	s += "\t\terr = rows.Scan(\n"
	s += serv.cgGoatDaoSelectScanVars(cols, fmt.Sprintf("\t\t\t&%s.", tni))
	s += "\t\t)\n\t\tif err != nil {\n\t\t\tbreak\n\t\t}\n"
	s += fmt.Sprintf("\t\tret = append(ret, %s)\n", tni)
	s +=  "\t}\n\n\treturn ret, err\n}"

	return s
}


func (serv *codegenService) cgGoatDaoSqlWhere(rdbms string, cols []model.Column, bindCount *int) string {
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


func (serv *codegenService) cgGoatDaoSqlWhereBindVals(tn string, cols []model.Column) string {
	s := ""
	tni := GetSnakeInitial(tn)

	for _, col := range cols {
		if col.PrimaryKeyFlg == constant.FLG_ON {
			s += fmt.Sprintf("\t\t%s.%s,\n", tni, SnakeToPascal(col.ColumnName))
		}
	}

	return s
}


func (serv *codegenService) cgGoatDaoSelectSqlColumns(cols []model.Column) string {
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


func (serv *codegenService) cgGoatDaoSelectScanVars(cols []model.Column, prefix string,) string {
	s := ""
	for _, col := range cols {
		s += fmt.Sprintf("%s%s,\n", prefix, SnakeToPascal(col.ColumnName))
	}

	s += fmt.Sprintf("%sCreatedAt,\n", prefix)
	s += fmt.Sprintf("%sUpdatedAt,\n", prefix)

	return s
}
