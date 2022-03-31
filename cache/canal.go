package cache

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/withlin/canal-go/client"
	protocol "github.com/withlin/canal-go/protocol"
	"github/guanhg/syncDB-search/config"
	"github/guanhg/syncDB-search/errorlog"
	"log"
	"strconv"
	"time"
)

/*
* 用于数据库的增量同步
*/

type Canal struct {
	Config struct{
		Uri string
		Port int
		UserName string
		Password string
		Dest string
		TimeOut int32
		IdleTimeOut int32
	}

	CanalConn *client.SimpleCanalConnector
}

// 获取canal的dest实例
func GetDestCanal(dest string) *Canal {
	c := new(Canal)

	cfg := config.JsonConfig

	if dest == "default"{
		dest = cfg.Canal.DefaultDest
	}
	c.Config.Dest = dest
	c.Config.Uri = cfg.Canal.Uri
	c.Config.Port = cfg.Canal.Port
	c.Config.UserName = cfg.Canal.Name
	c.Config.Password = cfg.Canal.Password
	c.Config.TimeOut = cfg.Canal.SoTO
	c.Config.IdleTimeOut = cfg.Canal.IdleTO

	c.Connecting()

	return c
}

func GetDefaultCanal() *Canal {
	return GetDestCanal("default")
}

func (c *Canal) Connecting()  {
	connector := client.NewSimpleCanalConnector(
		c.Config.Uri,
		c.Config.Port,
		c.Config.UserName,
		c.Config.Password,
		c.Config.Dest,
		c.Config.TimeOut,
		c.Config.IdleTimeOut)
	err :=connector.Connect()
	if err != nil {
		panic(err)
	}
	log.Println("Connecting To Canal: ", c.Config)
	c.CanalConn = connector
}

func (c *Canal) Get(regex string, size int32) []map[string]interface{} {
	err := c.CanalConn.Subscribe(regex)
	if err != nil {
		panic(err)
	}

	message, err := c.CanalConn.Get(size, nil, nil)
	if err != nil {
		panic(err)
	}

	batchId := message.Id
	if batchId == -1 || len(message.Entries) <= 0 {
		return nil
	}

	var data []map[string]interface{}
	for _, entry := range message.Entries {
		if entry.GetEntryType() == protocol.EntryType_TRANSACTIONBEGIN || entry.GetEntryType() == protocol.EntryType_TRANSACTIONEND {
			continue
		}
		rowChange := new(protocol.RowChange)

		err := proto.Unmarshal(entry.GetStoreValue(), rowChange)
		if err != nil {
			panic(err)
		}
		if rowChange != nil {
			mm := make(map[string]interface{})

			eventType := rowChange.GetEventType()
			header := entry.GetHeader()
			mm["schema"] = header.GetSchemaName()
			mm["table"] = header.GetTableName()

			if eventType == protocol.EventType_ALTER || eventType == protocol.EventType_TRUNCATE ||
				eventType == protocol.EventType_RENAME { // 修改表结构或清除表数据
				mm["event"] = -1
				data = append(data, mm)
			} else if eventType == protocol.EventType_ERASE { // 表删除和新建
				mm["event"] = -2
				data = append(data, mm)
			} else if eventType == protocol.EventType_DELETE || eventType == protocol.EventType_UPDATE ||
				eventType == protocol.EventType_INSERT {  // 表数据更新
				mm["event"] = header.GetEventType()

				var row map[string]interface{}
				for _, rowData := range rowChange.GetRowDatas() {
					if eventType == protocol.EventType_DELETE {
						row = ConvertColumn(rowData.GetBeforeColumns())
					} else if eventType == protocol.EventType_INSERT {
						row = ConvertColumn(rowData.GetAfterColumns())
					} else if eventType == protocol.EventType_UPDATE {
						row = ConvertColumn(rowData.GetAfterColumns())
					}
				}
				mm["data"] = row
				data = append(data, mm)
			}

			fmt.Println(fmt.Sprintf("================> binlog[%s : %d],name[%s,%s], eventType: %s", header.GetLogfileName(), header.GetLogfileOffset(), header.GetSchemaName(), header.GetTableName(), header.GetEventType()))
		}
	}
	return data
}

func ConvertColumn(columns []*protocol.Column) map[string]interface{} {
	row := make(map[string]interface{})
	for _, col := range columns {
		v := getValueTypeOfSqlType(col.GetValue(), col.GetSqlType())
		row[col.GetName()] = v
	}
	return  row
}

// 把Canal返回字段的sqlType，转换成go的常用数据类型
//参考 http://www.docjar.com/html/api/java/sql/Types.java.html

type JDBCSqlType int32
const (
	BIT 			JDBCSqlType = -7
	TINYINT			JDBCSqlType = -6
	SMALLINT		JDBCSqlType = 5
	INTEGER			JDBCSqlType = 4
	BIGINT			JDBCSqlType = -5
	FLOAT			JDBCSqlType = 7
	DOUBLE			JDBCSqlType = 8
	NUMERIC			JDBCSqlType = 2
	DECIMAL			JDBCSqlType = 3
	CHAR			JDBCSqlType = 1
	VARCHAR			JDBCSqlType = 12
	DATE			JDBCSqlType = 91
	TIME			JDBCSqlType = 92
	TIMESTAMP		JDBCSqlType = 93
	BINARY			JDBCSqlType = -2
	VARBINARY		JDBCSqlType = -3
	NULL			JDBCSqlType = 0
	BOOLEAN			JDBCSqlType = 16
)
func getValueTypeOfSqlType(v string, t int32) interface{} {
	if v == ""{
		return nil
	}

	var err error
	var i interface{}
	switch JDBCSqlType(t) {
	case TINYINT:
		i, err = strconv.ParseInt(v, 10, 8)
	case SMALLINT:
		i, err = strconv.ParseInt(v, 10, 8)
	case INTEGER: // int
		i, err = strconv.ParseInt(v, 10, 0)
	case BIGINT: // int64
		i, err = strconv.ParseInt(v, 10, 64)
	case FLOAT:
		i, err = strconv.ParseFloat(v, 32)
	case DOUBLE:
		i, err = strconv.ParseFloat(v, 64)
	case DECIMAL:
		i, err = strconv.ParseInt(v, 10,64)
	case DATE:
		i, err = time.Parse(shortDForm, v)
	case TIME:
		i, err = time.Parse(shortDtForm, v)
	case TIMESTAMP:
		i, err = time.Parse(shortDtForm, v)
	case BOOLEAN:
		i, err = strconv.ParseBool(v)
	default:
		return v
	}

	errorlog.CheckErr(err)
	return i
}
