package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"github/guanhg/syncDB-search/cache"
	"github/guanhg/syncDB-search/errorlog"
	"log"
	"time"
)

func main() {
	defer func() {
		if e:=recover(); e!=nil{
			log.Printf("[Publishing Error] Routine %v\n", e)
		}
	}()

	canal := cache.GetDefaultCanal()
	rq := cache.NewMqContext()
	log.Printf("=======[Start Sync DB]======\n")
	rqDeadOptions := cache.MqOptions{Exchange: "db_sync", ExchangeType: "topic", RouteKey: "db_index_dlx", Queue: "syncIndexDlx"}

	rqOptions := cache.MqOptions{Exchange: "db_sync", ExchangeType: "topic", RouteKey: "db_index", Queue: "syncIndex"}
	// 设置死信队列，在消息消费时出错会进入死信队列'rqDeadOptions'
	rqOptions.QueueArgs = map[string]interface{}{"x-dead-letter-exchange": rqDeadOptions.Exchange, "x-dead-letter-routing-key": rqDeadOptions.RouteKey}

	rq.DeclareExchangeQueue(rqOptions)
	rq.DeclareExchangeQueue(rqDeadOptions)

	for  {
		rows := canal.Get(".*\\..*", 100)
		if len(rows) <= 0 {
			time.Sleep(3 * time.Second)
			continue
		}
		for i :=range rows{  // 发送到mq
			fmt.Println(rows[i])
			body, err := json.Marshal(rows[i])
			errorlog.CheckErr(err)
			err = rq.Publish(rqOptions.Exchange, rqOptions.RouteKey, false, false, amqp.Publishing{Body: body})
			errorlog.CheckErr(err)
		}
	}
}
