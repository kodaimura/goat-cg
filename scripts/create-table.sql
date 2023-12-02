CREATE TABLE IF NOT EXISTS users (
	user_id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	created_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	updated_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime'))
);

CREATE TRIGGER IF NOT EXISTS trg_users_upd AFTER UPDATE ON users
BEGIN
    UPDATE users 
    	SET updated_at = DATETIME('now', 'localtime') 
    	WHERE rowid == NEW.rowid;
END;


CREATE TABLE IF NOT EXISTS project (
	project_id INTEGER PRIMARY KEY AUTOINCREMENT,
	project_name TEXT NOT NULL,
	project_memo TEXT,
	user_id INTEGER NOT NULL,
	username TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	updated_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	UNIQUE(project_name, username)
);

CREATE TRIGGER IF NOT EXISTS trg_project_upd AFTER UPDATE ON project
BEGIN
    UPDATE project
    	SET updated_at = DATETIME('now', 'localtime') 
    	WHERE rowid == NEW.rowid;
END;

CREATE TABLE IF NOT EXISTS project_member (
	project_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	user_status TEXT NOT NULL,
	user_role TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	updated_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	PRIMARY KEY(project_id, user_id)
);

CREATE TRIGGER IF NOT EXISTS trg_project_member_upd AFTER UPDATE ON project_member
BEGIN
    UPDATE project_member
    	SET updated_at = DATETIME('now', 'localtime') 
    	WHERE rowid == NEW.rowid;
END;


CREATE TABLE IF NOT EXISTS table_def (
	table_id INTEGER PRIMARY KEY AUTOINCREMENT,
	project_id INTEGER NOT NULL,
	table_name TEXT NOT NULL,
	table_name_logical TEXT,
	create_user_id INTEGER,
	update_user_id INTEGER,
	del_flg INTEGER NOT NULL DEFAULT 0,
	created_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	updated_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime'))
);

CREATE TRIGGER IF NOT EXISTS trg_table_def_upd AFTER UPDATE ON table_def
BEGIN
    UPDATE table_def
    	SET updated_at = DATETIME('now', 'localtime') 
    	WHERE rowid == NEW.rowid;

    INSERT INTO table_def_log
	SELECT * FROM table_def WHERE table_id == NEW.table_id;
END;

CREATE TRIGGER IF NOT EXISTS trg_table_def_ins AFTER INSERT ON table_def
BEGIN
	INSERT INTO table_def_log
	SELECT * FROM table_def WHERE table_id == NEW.table_id;
END;

CREATE TRIGGER IF NOT EXISTS trg_table_def_del AFTER DELETE ON table_def
BEGIN
	DELETE FROM table_def_log
	WHERE table_id == OLD.table_id;
END;


CREATE TABLE IF NOT EXISTS table_def_log (
	table_id INTEGER,
	project_id INTEGER,
	table_name TEXT,
	table_name_logical TEXT,
	create_user_id INTEGER,
	update_user_id INTEGER,
	del_flg INTEGER,
	created_at TEXT,
	updated_at TEXT
);


CREATE TABLE IF NOT EXISTS column_def (
	column_id INTEGER PRIMARY KEY AUTOINCREMENT,
	table_id INTEGER NOT NULL,
	column_name TEXT NOT NULL,
	column_name_logical TEXT,
	data_type_cls TEXT,
	precision INTEGER DEFAULT 0,
	scale INTEGER DEFAULT 0,
	primary_key_flg INTEGER DEFAULT 0,
	not_null_flg INTEGER DEFAULT 0,
	unique_flg INTEGER DEFAULT 0,
	default_value TEXT,
	remark TEXT,
	align_seq INTEGER,
	del_flg INTEGER NOT NULL DEFAULT 0,
	create_user_id INTEGER,
	update_user_id INTEGER,
	created_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	updated_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime'))
);

CREATE TRIGGER IF NOT EXISTS trg_column_def_upd AFTER UPDATE ON column_def
BEGIN
    UPDATE column_def
    SET updated_at = DATETIME('now', 'localtime') 
    WHERE rowid == NEW.rowid;

    INSERT INTO column_def_log
	SELECT * FROM column_def WHERE column_id == NEW.column_id;
END;

CREATE TRIGGER IF NOT EXISTS trg_column_def_ins AFTER INSERT ON column_def
BEGIN
	INSERT INTO column_def_log
	SELECT * FROM column_def WHERE column_id == NEW.column_id;
END;

CREATE TRIGGER IF NOT EXISTS trg_column_def_del AFTER DELETE ON column_def
BEGIN
	DELETE FROM column_def_log
	WHERE column_id == OLD.column_id;
END;

CREATE TABLE IF NOT EXISTS column_def_log (
	column_id INTEGER,
	table_id INTEGER,
	column_name TEXT,
	column_name_logical TEXT,
	data_type_cls TEXT,
	precision INTEGER,
	scale INTEGER,
	primary_key_flg INTEGER,
	not_null_flg INTEGER,
	unique_flg INTEGER,
	default_value TEXT,
	remark TEXT,
	align_seq INTEGER,
	del_flg INTEGER,
	create_user_id INTEGER,
	update_user_id INTEGER,
	created_at TEXT,
	updated_at TEXT
);


CREATE TABLE IF NOT EXISTS general (
	class TEXT,
	key1 TEXT,
	value1 TEXT,
	key2 TEXT,
	value2 TEXT,
	remark TEXT,
	created_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	updated_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime'))
);

CREATE TRIGGER IF NOT EXISTS trg_general_upd AFTER UPDATE ON general
BEGIN
    UPDATE general
    	SET updated_at = DATETIME('now', 'localtime') 
    	WHERE rowid == NEW.rowid;
END;