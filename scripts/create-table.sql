CREATE TABLE IF NOT EXISTS users (
	user_id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_name TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	create_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	update_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime'))
);

CREATE TRIGGER IF NOT EXISTS trigger_users_updated_at AFTER UPDATE ON users
BEGIN
    UPDATE users 
    	SET update_at = DATETIME('now', 'localtime') 
    	WHERE rowid == NEW.rowid;
END;


CREATE TABLE IF NOT EXISTS projects (
	project_id INTEGER PRIMARY KEY AUTOINCREMENT,
	project_cd TEXT NOT NULL UNIQUE,
	project_name TEXT,
	create_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	update_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime'))
);

CREATE TRIGGER IF NOT EXISTS trigger_projects_updated_at AFTER UPDATE ON projects
BEGIN
    UPDATE projects
    	SET update_at = DATETIME('now', 'localtime') 
    	WHERE rowid == NEW.rowid;
END;


CREATE TABLE IF NOT EXISTS users_projects (
	user_id INTEGER NOT NULL,
	project_id INTEGER NOT NULL,
	state_cls TEXT NOT NULL,
	role_cls TEXT,
	create_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	update_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	PRIMARY KEY(user_id, project_id)
);

CREATE TRIGGER IF NOT EXISTS trigger_users_projects_updated_at AFTER UPDATE ON users_projects
BEGIN
    UPDATE users_projects
    	SET update_at = DATETIME('now', 'localtime') 
    	WHERE rowid == NEW.rowid;
END;


CREATE TABLE IF NOT EXISTS tables (
	table_id INTEGER PRIMARY KEY AUTOINCREMENT,
	project_id INTEGER NOT NULL,
	table_name TEXT NOT NULL,
	table_name_logical TEXT,
	create_user_id INTEGER,
	update_user_id INTEGER,
	del_flg INTEGER NOT NULL DEFAULT 0,
	create_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	update_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime'))
);

CREATE TRIGGER IF NOT EXISTS trigger_tables_updated_at AFTER UPDATE ON tables
BEGIN
    UPDATE tables
    	SET update_at = DATETIME('now', 'localtime') 
    	WHERE rowid == NEW.rowid;
END;


CREATE TABLE IF NOT EXISTS columns (
	column_id INTEGER PRIMARY KEY AUTOINCREMENT,
	table_id INTEGER NOT NULL,
	column_name TEXT NOT NULL,
	column_name_logical TEXT,
	data_type_cls TEXT,
	precision INTEGER,
	scale INTEGER,
	primary_key_flg INTEGER DEFAULT 0,
	not_null_flg INTEGER DEFAULT 0,
	unique_flg INTEGER DEFAULT 0,
	remark TEXT,
	create_user_id INTEGER,
	update_user_id INTEGER,
	del_flg INTEGER NOT NULL DEFAULT 0,
	create_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	update_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime'))
);

CREATE TRIGGER IF NOT EXISTS trigger_columns_updated_at AFTER UPDATE ON columns
BEGIN
    UPDATE columns
    	SET update_at = DATETIME('now', 'localtime') 
    	WHERE rowid == NEW.rowid;
END;


CREATE TABLE IF NOT EXISTS generals (
	class TEXT,
	key1 TEXT,
	value1 TEXT,
	key2 TEXT,
	value2 TEXT,
	remark TEXT,
	create_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime')),
	update_at TEXT NOT NULL DEFAULT (DATETIME('now', 'localtime'))
);

CREATE TRIGGER IF NOT EXISTS trigger_generals_updated_at AFTER UPDATE ON generals
BEGIN
    UPDATE generals
    	SET update_at = DATETIME('now', 'localtime') 
    	WHERE rowid == NEW.rowid;
END;