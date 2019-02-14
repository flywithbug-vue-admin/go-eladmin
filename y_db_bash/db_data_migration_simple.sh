#!/usr/bin/env bash


#备份数据到当前目录下文件夹
mongodump -d  doc_manager -o ./db_backup/

#拷贝数据库到远端文
scp  -r ./db_backup/doc_manager name@host:/root/dump

#远程恢复数据库内容   数据库地址和端口           本地数据库目录            数据库用户名 密码     数据库名    覆盖远端数据库
mongorestore -h host:port ./db_backup/doc_manager -u name -p pass -d docmanager  --drop