# webDatabaseApp

## web操作


## mysqlclutser数据库操作

- sql0 192.168.50.100 管理节点
- sql1 192.168.50.128 数据节点[11] sql节点
- sql2 192.168.50.129 数据节点[12] sql节点

### mysqlcluster容器创建与连接

```bash
# 创建mysqlBridge网络
sudo docker network create --driver bridge --subnet 192.168.50.0/24 --gateway 192.168.50.1 mysqlBridge


docker pull lihan3238/mysql_ndb_cluster-ubuntu:lihan_ndbd_sql # ndb数据节点和sql节点
docker pull lihan3238/mysql_ndb_cluster-ubuntu:lihan_ndbmgm # mgm管理节点

docker run -di --name sql0 -v /home/lihan/sqlStudy:/home/shareFiles --net mysqlBridge --ip 192.168.50.100 lihan3238/mysql_ndb_cluster-ubuntu:lihan_ndbmgm # 管理节点

docker run -di --name sql1 -v /home/lihan/sqlStudy:/home/shareFiles --net mysqlBridge --ip 192.168.50.128 lihan3238/mysql_ndb_cluster-ubuntu:lihan_ndbd_sql # 数据节点[11] sql节点

docker run -di --name sql2 -v /home/lihan/sqlStudy:/home/shareFiles --net mysqlBridge --ip 192.168.50.129 lihan3238/mysql_ndb_cluster-ubuntu:lihan_ndbd_sql # 数据节点[12] sql节点

# 进入容器
docker exec -it sql0 bash
docker exec -it sql1 bash
docker exec -it sql2 bash

# 启动节点
## 管理节点

ndb_mgmd -f /var/lib/mysql-cluster/config.ini
    
ndb_mgm # 进入管理节点

show # 查看节点状态

## 数据节点

ndbd # 启动数据节点

mysqld --user=root & # 启动sql节点

mysql -u root -p # 进入sql节点，密码：123456

```

### mysqlcluster数据库操作

```sql
-- 创建数据库
create database movie;
use movie;

-- 创建表
create table genomescores(
    movieId int,
    tagId int,
    relevance float,
    primary key(movieId, tagId)
);

create table genometags(
    tagId int,
    tag varchar(255),
    primary key(tagId)
);

create table links(
    movieId int,
    imdbId int,
    tmdbId int,
    primary key(movieId)
);

create table movies(
    movieId int,
    title varchar(255),
    genres varchar(255),
    primary key(movieId)
);

create table ratings(
    userId int,
    movieId int,
    rating float,
    timestamp int,
    primary key(userId, movieId)
);

create table tags(
    userId int,
    movieId int,
    tag varchar(255),
    timestamp int,
    primary key(userId, movieId)
);

create table users(
    userId int,
    gender varchar(255),
    name varchar(255),
    primary key(userId)
);

-- 导入数据 
--  SHOW VARIABLES LIKE 'secure_file_priv'; 检查导入路径
LOAD DATA INFILE '/var/lib/mysql-files/ml-latest/genomescores.csv'
INTO TABLE genomescores
FIELDS TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 ROWS;
  -- 忽略 CSV 文件的首行（标题行）

-- `/var/lib/mysql-files/ml-latest/----.csv` 是你的 CSV 文件的实际路径。
-- `FIELDS TERMINATED BY ','` 指定字段之间的分隔符，这里是逗号。
-- `ENCLOSED BY '"'` 指定字段值的边界符，这里是双引号。
-- `LINES TERMINATED BY '\n'` 指定行的分隔符，这里是换行符。
-- `IGNORE 1 ROWS` 用于忽略 CSV 文件的首行，因为它通常包含列标题。

--  Error 1114 (HY000) The table is full 报错:
-- 检查：show global variables like 'max_heap_table_size';
-- 检查：show global variables like 'tmp_table_size';
--  解决方法1：修改配置文件 my.cnf，增加以下配置：
-- tmp_table_size = 800M // 临时表大小 
-- max_heap_table_size = 800M // 内存表大小 
-- 解决方法2：在数据库中执行以下命令：
-- set global tmp_table_size = 1024 * 1024 * 800*2;
-- set global max_heap_table_size = 1024 * 1024 * 800*2;

USE movie;

LOAD DATA INFILE 'D:/ml-latest/genometags.csv'
INTO TABLE genometags
FIELDS TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 ROWS;

USE movie;

LOAD DATA INFILE 'D:/ml-latest/links.csv'
INTO TABLE links
FIELDS TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 ROWS;

USE movie;

LOAD DATA INFILE 'D:/ml-latest/movies.csv'
INTO TABLE movies
FIELDS TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 ROWS;

USE movie;

LOAD DATA INFILE 'D:/ml-latest/ratings.csv'
INTO TABLE ratings
FIELDS TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 ROWS;

USE movie;

LOAD DATA INFILE 'D:/ml-latest/tags.csv'
INTO TABLE tags
FIELDS TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 ROWS;

USE movie;

LOAD DATA INFILE 'D:/ml-latest/users1.csv'
INTO TABLE users
CHARACTER SET utf8mb4
FIELDS TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\n'
IGNORE 1 ROWS;










