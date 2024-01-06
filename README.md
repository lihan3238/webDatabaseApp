# webDatabaseApp

## web操作

- gin+html/css/js

- 前端根据用户输入向后端传入参数：用户ID、关键词、年份、标签、任务
- 后端执行数据库查询并以json返回查询结果到前端输出

```golang
// main.go
package main

import (
    // 导入所需的包

func main() {
    // 设置数据库连接信息
    // ...

    // 设置静态文件目录
    // ...

    // 设置路由规则
    // ...

    // 运行服务器
    // ...
}

// 结构体定义
type ResultA struct {
    Movie  string   `json:"movie"`
    Rating string   `json:"rating"`
    Tag    []string `json:"tag"`
}

// executeQuery 函数
func executeQuery(db *sql.DB, searchUserID string, searchKeyword string, searchYear string, searchTag string, task string) ([]ResultA, error) {
    // 根据不同任务执行相应的数据库查询
    switch task {
        case "task_a":
            // 任务A的查询逻辑
            // ...
        case "task_b":
            // 任务B的查询逻辑
            // ...
        case "task_c":
            // 任务C的查询逻辑
            // ...
        case "task_d":
            // 任务D的查询逻辑
            // ...
        case "task_e":
            // 任务E的查询逻辑
            // ...
        default:
            return nil, fmt.Errorf("Invalid task specified" + searchUserID)
    }

    return nil, nil
}

// ginHtml 函数
func ginHtml(c *gin.Context) {
    // 处理HTML请求的逻辑
    // ...
}

```

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
```
### 查询例子简要展示
此处展示任务a的查询结果

![Alt text](assets/imgs/image.png)

当输入用户id时展示查询结果

![Alt text](assets/imgs/image-1.png)

![Alt text](assets/imgs/image-2.png)








