#! bin/bash
set -ex

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
	id int not null auto_increment,
	amount int default null,
	if_get bool  default false,
	primary key (id)
)ENGINE=innoDB DEFAULT CHARSET=utf8;"

create_table_sql2="create table if not exists red_envelope (
	id int not null auto_increment,
	user_id int default null,
	if_open bool default false,
	money int default null,
	open_time timestamp default null,
	primary key(id)
)ENGINE=innoDB DEFAULT CHARSET=utf8;"

mysql -u${username} -p${password} -D${dbname} -e "${create_table_sql1} ${create_table_sql2}"
echo "表单创建成功"


#插入数据
#大批量导入数据时需要换成copy
insert_user_1="insert user(amount, if_get) values(5, false);"
insert_user_2="insert user(amount, if_get) values(20, true);"
insert_user_3="insert user(amount, if_get) values(15, false);"
insert_user_4="insert user(amount, if_get) values(40, true);"

insert_envelope_1="insert red_envelope(user_id, if_open, money, open_time) values(1, true, 5, '2021-09-12 13:40:33');"
insert_envelope_2="insert red_envelope(if_open, money) values(false, 30);"
insert_envelope_3="insert red_envelope(user_id, if_open, money, open_time) values(2, true, 20, '2021-10-12 15:34:36');"
insert_envelope_4="insert red_envelope(user_id, if_open, money, open_time) values(4, true, 15, '2021-10-30 08:23:47');"


mysql -u${username} -p${password} -D${dbname} -e "${insert_user_1} ${insert_user_2} ${insert_user_3} ${insert_user_4}
${insert_envelope_1} ${insert_envelope_2} ${insert_envelope_3} ${insert_envelope_4}"
echo "数据导入成功"