--
-- 由SQLiteStudio v3.1.1 产生的文件 周三 6月 7 17:51:21 2017
--
-- 文本编码：System
--
PRAGMA foreign_keys = off;
BEGIN TRANSACTION;

-- 表：server
CREATE TABLE server (server TEXT (6) NOT NULL, name TEXT (250), urlmain TEXT (256), dsnmainmysql TEXT (256) NOT NULL);
INSERT INTO server (server, name, urlmain, dsnmainmysql) VALUES ('650001', '650001', '650001', 'root:root@tcp(192.168.1.28:3306)/limsxjhetian?charset=utf8&parseTime=true');

COMMIT TRANSACTION;
PRAGMA foreign_keys = on;
