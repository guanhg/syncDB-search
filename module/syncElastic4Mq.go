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

type CanalMapData []map[string]interface{}

// 多协程同步更新 canal 数据
func SyncElastic4Mq(numRoutine int){
	rqOptions := cache.MqOptions{Exchange: "db_sync", ExchangeType: "topic", RouteKey: "db_index", Queue: "syncIndex"}

	rq := cache.NewMqContext()
	rq.DeclareExchangeQueue(rqOptions)
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

// 不使用mq，而使用chan形式获取canal数据
func SyncElastic4Chan(numRoutine int, regex string)  {
	// 声明RabbitMq队列，用于保存处理chan消息失败后的消息
	rqOptions := cache.MqOptions{Exchange: "db_sync", ExchangeType: "topic", RouteKey: "db_index_dlx", Queue: "syncIndexDlx"}
	rq := cache.NewMqContext()
	rq.DeclareExchangeQueue(rqOptions)

	coroutines := make(chan bool, numRoutine)
	deliveries := make(chan CanalMapData)
	go SyncCanal2Chan(regex, deliveries)

	for {
		dataList := <-deliveries   // 阻塞等待
		coroutines <- true  // 限制创建协程数
		go func() {
			defer func() {
				if e:=recover(); e!=nil {
					log.Printf("[Consume Error] %s\n", e)
					for _, row := range dataList {
						body, _ := json.Marshal(row)
						_ = rq.Publish(rqOptions.Exchange, rqOptions.RouteKey, false, false, amqp.Publishing{Body: body})
					}
				}
			}()

			for _, row := range dataList {
				body, _ := json.Marshal(row)
				msg := amqp.Delivery{Body: body}
				DoConsume(msg)
			}
			<-coroutines
		}()
	}
}

// 消息回调函数
// 注意，json反序列后数值类型都是float64
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
		err = table.Delete(fmt.Sprintf("%v", data["id"]))
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


