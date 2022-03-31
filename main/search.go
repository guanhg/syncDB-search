package main

import (
	"flag"
	"github/guanhg/syncDB-search/module"
	"log"
)

var mod = flag.String("m", "web", "(module)需要调用的执行模块: init/delete/rebuild/mapping/db2mq/es4mq/es4ch/web")
var tb = flag.String("t", "", "(table name)表名，init/delete/rebuild/mapping模块必须要的参数")
var db = flag.String("d", "", "(database name)数据库名，init/delete/rebuild/mapping模块必须要的参数")
var dest = flag.String("dest", "default", "(dest)canal的destination的实例，es4ch/canal2mq模块必须要的参数")
var reg = flag.String("regx", ".*\\..*", "(regex)canal获取表的同步数据，canal2mq模块必须要的参数")
var n = flag.Int("n", 3, "使用n个协程同步mq数据到es，默认3，es4ch/es4mq模块可选参数")
var port = flag.String("p", "8080", "端口号，默认8080，web模块可选参数")

func main() {
	flag.Parse()

	switch *mod {
	case "init":
		module.InitDbIndex(*tb, *db)
	case "delete":
		module.DeleteDbIndex(*tb, *db)
	case "rebuild":
		module.RebuildDbIndex(*tb, *db)
	case "mapping":
		module.MappingIndex(*tb, *db)
	case "db2mq":
		module.SyncCanal2Mq(*dest, *reg)
	case "es4mq":
		module.SyncElastic4Mq(*n)
	case "web":
		module.Application(*port)
	case "es4ch":
		module.SyncElastic4Chan(*n, *reg)
	default:
		log.Printf("Module '%s' 不存在，可选 'init/delete/rebuild/mapping/db2mq/es4mq/es4ch/web'\n", *mod)
	}
}
