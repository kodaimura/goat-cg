package constant


const APPNAME string = "goat-cg"

//全画面(テンプレート)で共通で設定するパラメータ
var Commons = map[string]string{
	"appname" : APPNAME,
}

//フラグ系(*_FLG)テーブルカラム設定値 
const FLG_ON = 1
const FLG_OFF = 0