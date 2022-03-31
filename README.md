### 同步Mysql表数据到Elastic

​使用ali的Canal数据库同步工具，把Mysql数据的增删改同步到RabbitMQ，然后从MQ中拿取消息同步到Elastic中

 - Mysql --> Canal --> RabbitMq --> Elastic

1. 依赖安装 [Canal](https://github.com/alibaba/canal) 、RabbitMq、MySql、ElasticSearch 7.5

2. 编辑config/config.json配置文件

3. 主体文件main/search.go

	go run main/search.go -m module [paramters]

	目录module中的文件函数是主要执行体

	包括初始化索引，删除索引，重构索引，获取索引mapping，同步canal数据到mq，同步mq数据到es，web应用

```
-m string
      (module)需要调用的执行模块: init/delete/rebuild/mapping/db2mq/es4mq/web
-t string
      (table name)表名，init/delete/rebuild/mapping模块必选参数
-d string
      (database name)数据库名，init/delete/rebuild/mapping模块必选参数
-n int
      使用n个协程同步mq数据到es，默认10，es4mq模块可选参数 (default 10)
-p string
      端口号，默认8080，web模块可选参数 (default "8080")
-r string
      (regex)canal获取表的同步数据，canal2mq模块必选参数 (default ".*\\..*")
```

```
如 初始化|删除|重构|mapping 数据库statement的表track索引
go run search.go -m init|delete|rebuild|mapping -t table_name -d db_name
如 同步数据statement中所有表更新数据到mq
go run search.go -m db2mq -r "example\\..*"
如 10个协程同步mq数据到es
go run search.go -m es4mq -n 10
如 启动search web应用,端口8080
go run search.go -m web -p 8080
```

