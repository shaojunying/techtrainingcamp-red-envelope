#! bin/bash
set -e

username="root" 
password="root"

dbname="red_envelope_rain" #数据库名称


#创建数据库
drop_db_sql="drop database if exists ${dbname};"
create_db_sql="create database if not exists ${dbname};"

mysql -u${username} -p${password} -e "${drop_db_sql} ${create_db_sql}"
echo "数据库创建成功"


#创建表
create_table_sql1="create table if not exists user (
	uid int not null auto_increment,
	amount int default null,
	if_get bool  default false,
	primary key (uid)
)ENGINE=innoDB DEFAULT CHARSET=utf8;"

create_table_sql2="create table if not exists red_envelope (
	envelope_id int not null auto_increment,
	uid int default null,
	opened bool default false,
	value int default null,
	snatch_time timestamp DEFAULT CURRENT_TIMESTAMP,
	primary key(envelope_id)
)ENGINE=innoDB DEFAULT CHARSET=utf8;"

mysql -u${username} -p${password} -D${dbname} -e "${create_table_sql1} ${create_table_sql2}"
echo "表单创建成功"


#插入数据
#大批量导入数据时需要换成copy
#insert_user_1="insert user(amount, if_get) values(5, false);"
#insert_user_2="insert user(amount, if_get) values(20, true);"
#insert_user_3="insert user(amount, if_get) values(15, false);"
#insert_user_4="insert user(amount, if_get) values(40, true);"
#
#insert_envelope_1="insert red_envelope(uid, opened, value) values(1, true, 5);"
#insert_envelope_2="insert red_envelope(opened, value) values(false, 30);"
#insert_envelope_3="insert red_envelope(uid, opened, value) values(2, true, 20);"
#insert_envelope_4="insert red_envelope(uid, opened, value) values(4, true, 15);"

#无需导入数据
#mysql -u${username} -p${password} -D${dbname} -e "${insert_user_1} ${insert_user_2} ${insert_user_3} ${insert_user_4}
#${insert_envelope_1} ${insert_envelope_2} ${insert_envelope_3} ${insert_envelope_4}"
#echo "数据导入成功"