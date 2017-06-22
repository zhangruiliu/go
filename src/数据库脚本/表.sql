-- 样本对象表
CREATE TABLE `personinfo` (
	`id`	TEXT NOT NULL UNIQUE,
	`person_name`	TEXT,
	`gender`	TEXT,
	`race`	TEXT,
	`birth_datetime`	TEXT,
	`native_place_region`	TEXT,
	`native_place_addr`	TEXT,
	`id_card_no`	TEXT,
	`sample_lab_no`	TEXT NOT NULL,
	`collect_user`	TEXT,
	`collect_datetime`	TEXT,
	`collect_agency_code`	TEXT,
	`collect_agency_name`	TEXT,
	`server`	TEXT,
	`remark`	TEXT,
	`Person_id`	TEXT,
	`create_user`	TEXT,
	`create_datetime`	TEXT,
	`update_user`	TEXT,
	`update_datetime`	TEXT,
	`sample_type`	TEXT,
	`flag`	INTEGER NOT NULL,
	`foot`	BLOB,
	PRIMARY KEY(`id`)
);
-- 用户表
CREATE TABLE `sys_user` (
	`id`	TEXT NOT NULL,
	`init_server_no`	TEXT NOT NULL,
	`lab_id`	TEXT,
	`user_type`	TEXT,
	`user_name`	TEXT,
	`password`	TEXT,
	`full_name`	TEXT,
	`phone`	TEXT,
	`certificate_type`	TEXT,
	`certificate_no`	TEXT,
	`email`	TEXT,
	`organization_regionalism`	TEXT,
	`organization_name`	TEXT,
	`login_ip_range`	TEXT,
	`active_flag`	INTEGER,
	`is_lab_employee`	INTEGER,
	`lab_employee_id`	TEXT,
	`is_transfer_user`	INTEGER,
	`remark`	TEXT,
	`data_source`	TEXT,
	`delete_flag`	INTEGER NOT NULL,
	`user_level`	INTEGER,
	`create_user`	TEXT NOT NULL,
	`create_datetime`	timestamp,
	`update_user`	TEXT,
	`update_datetime`	timestamp,
	`va`	TEXT,
	`inspection_personnel_id`	TEXT,
	PRIMARY KEY(`ID`)
);
-- 查重表
CREATE TABLE check_person (
    id            TEXT (20)  UNIQUE
                             PRIMARY KEY,
    id_card_no    TEXT (18),
    server        TEXT (12),
    status        TEXT (1),
    name          TEXT (256),
    sample_lab_no TEXT (64) 
);
CREATE TABLE `idcard` (
	`id_card_no`	TEXT(18) NOT NULL,
	`person_name`	TEXT(256),
	`gender`	TEXT(12),
	`race`	TEXT(64),
	`address`	TEXT,
	PRIMARY KEY(`id_card_no`)
);
--字典表
CREATE TABLE `sys_dict` (
	`id`	TEXT NOT NULL,
	`dict_category`	TEXT,
	`dict_sub_category`	TEXT,
	`root_flag`	INTEGER,
	`dict_key`	TEXT,
	`dict_national_key`	TEXT,
	`dict_value1`	TEXT,
	`dict_value2`	TEXT,
	`dict_value3`	TEXT,
	`dict_alias`	TEXT,
	`parent_key`	TEXT,
	`download_flag`	INTEGER,
	`readonly_flag`	INTEGER,
	`ord`	INTEGER,
	`dict_py`	TEXT,
	`active_flag`	INTEGER,
	`remark`	TEXT,
	`create_user`	TEXT,
	`create_datetime`	timestamp,
	`update_user`	TEXT,
	`update_datetime`	timestamp,
	PRIMARY KEY(`id`)
);
CREATE TABLE `server` (
	`server`	TEXT(6) NOT NULL,
	`name`	TEXT(250),
	`urlmain`	TEXT(256),
	`dsnmainmysql`	TEXT(256) NOT NULL
);