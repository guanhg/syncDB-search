package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"github/guanhg/syncDB-search/cache"
	"github/guanhg/syncDB-search/errorLog"
	schema "github/guanhg/syncDB-search/schema-index"
)

func main() {

	// sm
	table := schema.SchemaIndex{Name: "sm_record_2017", Context: cache.GetContext("default")}
	err :=table.CreateIndexIfNotExist()
	errorLog.CheckErr(err)
	table.IndexOne(1)
	//fmt.Println(table.BuildFieldMapping())

	// search
	//testSearch()
	//testAggregation()
}

func testAggregation(){
	q := elastic.NewBoolQuery()
	q.Must(elastic.NewTermQuery("sm_id", 1716)).Must(elastic.NewTermQuery("medium_type", 0))
	q.Must(elastic.NewRangeQuery("medium_id").Gt(0))

	sumAgg := elastic.NewSumAggregation().Field("weight")
	disAgg := elastic.NewTermsAggregation().Field("medium_id").SubAggregation("weight", sumAgg).Size(50).OrderByAggregation("weight", false)

	search := schema.Search(q, "sm_record_*")
	search.Size(0).Aggregation("track", disAgg)
	res, _ := search.Do(context.Background())
	aggResult, _ := res.Aggregations["track"].MarshalJSON()

	aggTrack := make(map[string]interface{})
	json.Unmarshal(aggResult, &aggTrack)

	fmt.Println(aggTrack)
}

func testSearch(){
	q := elastic.NewRangeQuery("created_datetime").From("2019-01-01").To("2019-06-01").Boost(3).Relation("within")
	search := schema.Search(q, "sm")
	search.Size(2).From(1)
	res := search.Result()
	fmt.Println(res)
}

