package constant


/*
区分系(*_CLS)カラム区分値

const *_CLS_XXXX = "00"
const *_CLS_YYYY = "01"
const *_CLS_ZZZZ = "99"
*/

//USER_PROJECTS.STATE_CLS
const (
	STATE_CLS_NOMAL = "00"
	STATE_CLS_JOIN = "01"
	STATE_CLS_REQUEST = "02"
	STATE_CLS_BLOCK = "03"
)

//USER_PROJECTS.ROLE_CLS
const (
	ROLE_CLS_NOMAL = "00"
	ROLE_CLS_ADMIN = "88"
	ROLE_CLS_OWNER = "99"
)


//COLUMNS.DATE_TYPE_CLS
const (
	DATA_TYPE_CLS_SERIAL = "01"
	DATA_TYPE_CLS_TEXT = "10"
	DATA_TYPE_CLS_VARCHAR = "11"
	DATA_TYPE_CLS_CHAR = "12"
	DATA_TYPE_CLS_INTEGER = "20"
	DATA_TYPE_CLS_NUMERIC = "30"
	DATA_TYPE_CLS_TIMESTAMP = "40"
	DATA_TYPE_CLS_DATE = "41"
	DATA_TYPE_CLS_BLOB = "50"
)