package module

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"github/guanhg/syncDB-search/cache"
	"github/guanhg/syncDB-search/errorlog"
	"github/guanhg/syncDB-search/schema"
	"log"
)

// 多协程同步更新
func SyncElastic4Mq(numRoutine int){
	rqOptions := cache.MqOptions{Exchange: "db_sync", ExchangeType: "topic", RouteKey: "db_index", Queue: "syncIndex"}

	rq := cache.NewMqContext()
	//rq.DeclareExchangeQueue(rqOptions)
	_ = rq.Qos(100, 0, true)

	forever := make(chan bool)

	for i:=0; i<numRoutine; i++ {
		go func(ser int) {
			defer func() {
				if e:=recover(); e!=nil{
					log.Printf("[Consume Error] Routine %d-> %s\n", ser, e)
				}
			}()
			fmt.Printf("[--> Starting %d] \n", ser)
			msgs, err := rq.Consume(rqOptions.Queue, "", false, false, false, false, nil)
			errorlog.CheckErr(err)
			for msg := range msgs {
				rq.OnMessage(msg, DoConsume)
			}
		}(i)
	}

	<- forever
}

// 消息回调函数
func DoConsume(msg amqp.Delivery)  {
	rowMap := make(map[string]interface{})
	err := json.Unmarshal(msg.Body, &rowMap)
	errorlog.CheckErr(err)

	tableName := rowMap["table"].(string)
	event := int(rowMap["event"].(float64))
	dbName := rowMap["schema"].(string)
	fmt.Printf("============%s.%s Event: %d ===========\n", dbName, tableName, event)

	table := schema.SchemaIndex{Name: rowMap["table"].(string), Context: cache.GetContext("default")}
	if event == 3 {  //删除记录
		data := rowMap["data"].(map[string]interface{})
		err = table.Delete(data["id"].(string))
	}else if event == 2 || event == 1 {  // 插入或更新
		data := rowMap["data"].(map[string]interface{})
		err = table.Upsert(data)
	} else if event == -1 {  // 表更新
		err = table.IndexAll(10)
	} else if event == -2 {  // 删除表
		err = table.DeleteIndexIfExist()
	}

	errorlog.CheckErr(err)
}


